package utils

import (
	"net/http"
	"encoding/json"
)


type GlobalMap map[string]any

func getJSONMessage(jsonType, message string) GlobalMap {
	switch jsonType {
		case "e":
			return GlobalMap{"error": message}
		default:
			return GlobalMap{"message": message}
	}
}

func SendResponse(res http.ResponseWriter, responseMessage, status string, responseData ...interface{}) {
	switch status {
		case "b":
			res.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(res).Encode(getJSONMessage("e", responseMessage))
		case "e":
			res.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(res).Encode(getJSONMessage("e", responseMessage))
		case "u":
			res.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(res).Encode(getJSONMessage("e", responseMessage))
		default:
			if len(responseData) > 0 {
				data := getJSONMessage("", responseMessage)
				data["data"] = responseData[0]
				res.WriteHeader(http.StatusOK)
				json.NewEncoder(res).Encode(data)
			} else {
				res.WriteHeader(http.StatusOK)
				json.NewEncoder(res).Encode(getJSONMessage("", responseMessage))
			}
	}
}
