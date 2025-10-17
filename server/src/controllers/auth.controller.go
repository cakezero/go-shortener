package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/cakezero/go-shortener/src/models"
	"github.com/cakezero/go-shortener/src/utils"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/kamva/mgm/v3"
	"golang.org/x/crypto/bcrypt"
)

func Register(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var user models.User

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		utils.Logger.Error(err.Error())
		utils.SendResponse(res, "Request data received was not appropriate", "b")
		return
	}

	if validateError := NewValidator.Struct(user); validateError != nil {
		utils.Logger.Error(validateError.Error())
		utils.SendResponse(res, "All fields are required", "b")
		return
	}

	if length := len(user.Password); length < 8 {
		utils.SendResponse(res, "Password must be gte 8", "b")
		return
	}

	userModel := &models.User{}
	checkUserErr := mgm.Coll(userModel).First(bson.M{"email": user.Email}, userModel)

	if checkUserErr == nil { // this means email exists cause an error wasn't thrown
		utils.SendResponse(res, "Email exists", "b")
		return
	}

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if hashErr != nil {
		utils.Logger.Error(hashErr.Error())
		utils.SendResponse(res, "Internal server error", "e")
		return
	}

	user.Password = string(hashedPassword)
	saveErr := mgm.Coll(&user).Create(&user)

	if saveErr != nil {
		utils.Logger.Error(saveErr.Error())
		utils.SendResponse(res, "Internal server error", "e")
		return
	}

	accessToken, refreshToken, tokenFetchErr := utils.GenerateJWTs(user.ID.String())

	if tokenFetchErr != nil {
		utils.Logger.Error(tokenFetchErr.Error())
		utils.SendResponse(res, "Internal server error", "e")
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
	})

	data := utils.GlobalMap{
		"user":  user,
		"token": accessToken,
	}

	utils.SendResponse(res, "User saved", "", data)
}

func Login(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var user models.User

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		utils.SendResponse(res, "Inappropriate data recieved from the request", "b")
		return
	}

	if user.Password == "" || user.Email == "" {
		utils.SendResponse(res, "Both password and email are required", "b")
		return
	}

	fetchedUser := mgm.Coll(&models.User{}).FindOne(context.Background(), bson.M{"email": user.Email})

	if fetchedUser.Err() != nil {
		utils.SendResponse(res, "email or password is invalid", "b")
		return
	}

	var decodedUser models.User

	if decodeErr := fetchedUser.Decode(&decodedUser); decodeErr != nil {
		utils.SendResponse(res, "Internal server error", "e")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(decodedUser.Password), []byte(user.Password)); err != nil {
		utils.SendResponse(res, "email or password is invalid", "b")
		return
	}

	accessToken, refreshToken, tokenFetchErr := utils.GenerateJWTs(user.ID.String())
	if tokenFetchErr != nil {
		utils.SendResponse(res, "Error creating access/refresh token, try again", "e")
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
	})

	loginData := utils.GlobalMap{
		"user":        decodedUser,
		"accessToken": accessToken,
	}

	utils.SendResponse(res, "User logged In", "", loginData)
}

func Logout(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		utils.SendResponse(res, "Invalid auth header format", "u")
		return
	}

	accessToken := parts[1]
	_, claims, decodeErr := utils.DecodeJWT(res, accessToken)

	if decodeErr != nil {
		utils.Logger.Error(decodeErr.Error())
		return
	}

	exp := int64(0)
	if expFloat, ok := claims["exp"].(float64); ok {
		exp = int64(expFloat)
	}

	if exp > 0 {
		ttl := time.Until(time.Unix(exp, 0))
		err := utils.GetRedisClient().Set(Ctx, accessToken, "revoked", ttl).Err()

		if err != nil {
			utils.Logger.Error(err.Error())
			utils.SendResponse(res, "Error revoking token", "e")
			return
		}
	}

	http.SetCookie(res, &http.Cookie{
		Name: "refresh_token",
		Value: "",
		MaxAge: -1,
		Path: "/",
		HttpOnly: true,
		Secure: true,
		Expires: time.Unix(0, 0),
	})

	utils.SendResponse(res, "User logged out", "")

}
