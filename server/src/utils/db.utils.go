package utils

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client;

func DB() error {
	dbConnctErr := mgm.SetDefaultConfig(nil, "go-shortener", options.Client().ApplyURI(DB_URI))

	return dbConnctErr;
}

func GetRedisClient() *redis.Client {
	if redisClient != nil {
		return redisClient
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr: REDIS_URI,
		Username: REDIS_USERNAME,
		Password: REDIS_PASSWORD,
		DB: 0,
	})

	return redisClient
}
