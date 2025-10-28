package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cakezero/go-shortener/src/models"
	"github.com/cakezero/go-shortener/src/utils"
	"github.com/jpillora/go-tld"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/go-playground/validator/v10"
	"github.com/kamva/mgm/v3"
)

var NewValidator = validator.New()
var Ctx = context.Background()

func Shorten(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	type requestBody struct {
		Url string
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

	longUrl := reqBody.Url

	if !strings.HasPrefix(longUrl, "https://") && !strings.HasPrefix(longUrl, "http://") {
		utils.SendResponse(res, "URL must start with https:// or http://", "")
		return
	}

	_, err := tld.Parse(longUrl)
	if err != nil {
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
			"shorturl": utils.Domain + shortUrl,
			"longurl":  longUrl,
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

	utils.SendResponse(res, "URL shortened", "")
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

	shortUrl := shortenedLinkReqBody.ShortUrl

	if shortUrl == "" {
		utils.SendResponse(res, "shortUrl is required", "b")
		return
	}

	var urlFound models.Url

	urlNotFound := mgm.Coll(&models.Url{}).First(bson.M{"shorturl": shortUrl}, &urlFound)

	if urlNotFound != nil {
		utils.SendResponse(res, "shortUrl doesn't exist or is invalid", "e")
		return
	}

	utils.SendResponse(res, "redirecting...", "", urlFound.LongUrl)
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

func FetchUrls(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "applicaton/json")

	id := req.URL.Query().Get("id")

	if id == "" {
		utils.SendResponse(res, "id is required", "")
		return
	}

	var urls []models.Url 

	noUrlsFound := mgm.Coll(&models.Url{}).SimpleFind(&urls, bson.M{"user": id})

	if noUrlsFound != nil {
		utils.SendResponse(res, "No urls shortened", "")
		return
	}

	utils.SendResponse(res, "urls fetched!", "", urls)
}

func DeleteAllUrls(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	id := req.URL.Query().Get("id")

	_, deleteManyErr := mgm.Coll(&models.Url{}).DeleteMany(Ctx, bson.M{"user": id})

	if deleteManyErr != nil {
		utils.Logger.Error(deleteManyErr.Error())
		utils.SendResponse(res, "Error deleting urls", "e")
		return
	}

	utils.SendResponse(res, "Urls deleted successfully", "")
}
func DeleteSelectedUrls(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	type selectedUrls struct {
		Ids []string
	}

	var selectedUrlIds selectedUrls

	if decodedErr := json.NewDecoder(req.Body).Decode(&selectedUrlIds); decodedErr != nil {
		utils.SendResponse(res, "Data received from request is inappropriate", "")
		return
	}

	for _, id := range selectedUrlIds.Ids {
		deletedUrl := mgm.Coll(&models.Url{}).FindOneAndDelete(Ctx, bson.M{"_id": id})

		if deletedUrl.Err() != nil {
			utils.Logger.Error(deletedUrl.Err().Error())
			utils.SendResponse(res, "Error deleting url", "e")
			return
		}
	}

	utils.SendResponse(res, "Selected urls deleted", "")
}
