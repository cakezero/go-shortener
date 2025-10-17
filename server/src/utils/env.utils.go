package utils

import (
	"os"
	"github.com/joho/godotenv"
)

var (
	DB_URI, Domain string
	ENVIRONMENT string
	JWT_SECRET, REFRESH_SECRET []byte
	PORT, REDIS_PASSWORD string
	REDIS_URI, REDIS_USERNAME string
)

func LoadEnv() (err error) {
	err = godotenv.Load()

  DB_URI = os.Getenv("DB_URI")

	ENVIRONMENT = os.Getenv("ENVIRONMENT")

	if ENVIRONMENT == "production" {
		Domain = "https://go-u-sh.vercel.app/"
	} else {
		Domain = "http://localhost:5173/"
	}

	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
	
	if port := os.Getenv("PORT"); port != "" {
		PORT = port
	} else {
		PORT = ":3030"
	}
	
	REDIS_URI = os.Getenv("REDIS_URI")

	REDIS_USERNAME = os.Getenv("REDIS_USERNAME")

  REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")

	REFRESH_SECRET = []byte(os.Getenv("REFRESH_SECRET"))

	return
}
