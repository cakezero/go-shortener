package utils

import (
	"time"
	"math/rand"
	"go.uber.org/zap"
)

var Logger *zap.Logger;

const baseCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";

func GenerateShortUrl(urlText string) string {
	rand.New(rand.NewSource(time.Now().UnixNano()));

	b := make([]byte, 5);

	charset := baseCharset + urlText;
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))];
	}

	return string(b);
}

func StartLogger () error {
	var logger *zap.Logger;
	var err error;

	if ENVIRONMENT == "production" {
		logger, err = zap.NewProduction();
	} else {
		logger, err = zap.NewDevelopment();
	}

	defer logger.Sync();

	Logger = logger;

	return err;
}
