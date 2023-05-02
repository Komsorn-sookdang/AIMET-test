package databases

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var redisClient *redis.Client

func getRedisURI() (string, error) {
	URI := viper.GetString("redis.uri")
	if URI == "" {
		return "", fmt.Errorf("Redis URI is empty")
	}
	return URI, nil
}

func InitRedisClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	redisURI, err := getRedisURI()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisURI,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	redisClient = client
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func CloseRedisClient() {
	redisClient.Close()
}
