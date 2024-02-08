package redis

import (
	redis "github.com/redis/go-redis/v9"
)

var client *redis.Client

func NewClient() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func GetClient() *redis.Client {
	if client == nil {
		NewClient()
	}

	return client
}
