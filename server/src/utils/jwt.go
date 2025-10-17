package utils

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTs(uid string) (accessToken, refreshToken string, err error) {
	accessClaims := jwt.MapClaims{
		"id": uid,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	}

	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessJWT.SignedString(JWT_SECRET)
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.MapClaims{
		"id": uid,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshJWT.SignedString(REFRESH_SECRET)
	return 
}

func DecodeJWT(res http.ResponseWriter, jwt_value string) (string, jwt.MapClaims, error) {
	parsedToken, parseErr := jwt.Parse(jwt_value, func (token *jwt.Token) (interface{}, error) {
		return REFRESH_SECRET, nil
	})

	if parseErr != nil || !parsedToken.Valid {
		SendResponse(res, "Refresh token is invalid", "u")
		return "", nil, parseErr
	}

	claims := parsedToken.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	return id, claims, nil
}