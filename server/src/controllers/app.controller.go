package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/cakezero/go-shortener/src/models"
	"github.com/cakezero/go-shortener/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"
	"github.com/kamva/mgm/v3"
)

var NewValidator = validator.New()
var Ctx = context.Background()

func Home(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	id := req.URL.Query().Get("id")

	var shortenUrlHistory []models.Url

	urlHistory, findErr := mgm.Coll(&models.Url{}).Find(Ctx, bson.M{"user": id})

	if findErr != nil {
		utils.SendResponse(res, "No url has been shortened", "e")
		return
	}

	if fetchErr := urlHistory.All(Ctx, &shortenUrlHistory); fetchErr != nil {
		utils.SendResponse(res, "Error fetching url history", "e")
		return
	}

	utils.SendResponse(res, "Url fetched", "", shortenUrlHistory)
}

func Shorten(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	type requestBody struct {
		url string
		Id  string
	}

	var reqBody requestBody
	var urlModel models.Url

	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		utils.Logger.Error(err.Error())
		utils.SendResponse(res, "Request data received was not appropriate", "b")
		return
	}

	validatorErr := NewValidator.Struct(reqBody)

	if validatorErr != nil {
		utils.Logger.Error(validatorErr.Error())
		utils.SendResponse(res, "All fields are required", "b")
		return
	}

	longUrl := reqBody.url

	utils.Logger.Info("Url from request", zap.String("longUrl", longUrl))

	// if _, notURLErr := url.ParseRequestURI(longUrl); notURLErr != nil {
	if _, notURLErr := url.Parse(longUrl); notURLErr != nil {
		utils.Logger.Error(notURLErr.Error())
		utils.SendResponse(res, "Invalid URL passed, send urls like so: https://example.com or http://example.com", "b")
		return
	}

	urlBody := strings.Split(strings.Split(longUrl, "://")[1], ".")[0]

	shortUrl := utils.GenerateShortUrl(urlBody)

	id := reqBody.Id

	urlModel.ShortUrl = shortUrl
	urlModel.LongUrl = longUrl

	if id == "" {
		saveErr := mgm.Coll(&urlModel).Create(&urlModel)

		if saveErr != nil {
			utils.Logger.Error(saveErr.Error())
			utils.SendResponse(res, "Internal server error", "e")
			return
		}

		data := utils.GlobalMap{
			"shortUrl": shortUrl,
			"longUrl":  longUrl,
		}

		utils.SendResponse(res, "URL shortened", "", data)
		return
	}

	urlModel.User = &id

	saveErr := mgm.Coll(&urlModel).Create(&urlModel)

	if saveErr != nil {
		utils.Logger.Error(saveErr.Error())
		utils.SendResponse(res, "Internal server error", "e")
		return
	}

	utils.SendResponse(res, "URL shortened", "", utils.Domain + shortUrl)
}


func VisitLongUrl(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	type shortenedLinkParams struct {
		ShortUrl string
	}

	var shortenedLinkReqBody shortenedLinkParams

	decodeError := json.NewDecoder(req.Body).Decode(&shortenedLinkReqBody)
	if decodeError != nil {
		utils.Logger.Error(decodeError.Error())
		utils.SendResponse(res, "Data received from request is not appropriate", "e")
		return
	}

	if shortenedLinkReqBody.ShortUrl == "" {
		utils.SendResponse(res, "shortUrl is required", "b")
		return
	}

	urlFound := mgm.Coll(&models.Url{}).FindOne(Ctx, bson.M{"shortUrl": shortenedLinkReqBody.ShortUrl});
	if urlFound.Err() != nil {
		utils.Logger.Error(urlFound.Err().Error())
		utils.SendResponse(res, "shortUrl doesn't exist or is invalid", "e")
		return
	}

	var urlProp models.Url

	urlDecodeErr := urlFound.Decode(&urlProp)
	if urlDecodeErr != nil {
		utils.Logger.Error(urlDecodeErr.Error())
		utils.SendResponse(res, "Internal server error", "e")
		return
	}

	utils.SendResponse(res, "redirecting...", "", urlProp.LongUrl)
}

func DeleteUrl(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	
	id := req.URL.Query().Get("id")

	if id == "" {
		utils.SendResponse(res, "id is required", "b")
		return
	}


	deletedUrl := mgm.Coll(&models.Url{}).FindOneAndDelete(Ctx, bson.M{"_id": id})

	if deletedUrl.Err() != nil {
		utils.Logger.Error(deletedUrl.Err().Error())
		utils.SendResponse(res, "Error deleting url", "e")
		return
	}

	utils.SendResponse(res, "Url deleted successfully", "")
}