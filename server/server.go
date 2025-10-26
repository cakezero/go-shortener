package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cakezero/go-shortener/src/routes"
	"github.com/cakezero/go-shortener/src/utils"
	"github.com/rs/cors"
)

func init() {
	loadEnvErr := utils.LoadEnv()

	if loadEnvErr != nil {
		panic(loadEnvErr)
	}

	dbError := utils.DB()
	if dbError != nil {
		panic(dbError)
	}

	loggerError := utils.StartLogger()

	if loggerError != nil {
		panic(loggerError)
	}

	utils.Logger.Info("DB Connected!")
}

func main() {
	router := routes.Routes()

	CORS := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173", "https://go-u-sh.vercel.app"},
		AllowedMethods: []string{"POST", "GET", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	corsHandler := CORS.Handler(router)

	port := utils.PORT

	serverRunningMsg := fmt.Sprintf("Server is running on port %s\n", port)

	utils.Logger.Info(serverRunningMsg)

	err := http.ListenAndServe(port, corsHandler)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
