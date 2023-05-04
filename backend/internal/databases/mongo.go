package databases

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client

func getMongoURI() (string, error) {
	URI := viper.GetString("mongo.uri")
	if URI == "" {
		return "", fmt.Errorf("Mongo URI is empty")
	}
	return URI, nil
}

func InitMongoClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI, err := getMongoURI()
	if err != nil {
		panic(err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Println(mongoURI)
		panic(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	mongoClient = client
}

func GetMongoClient() *mongo.Client {
	return mongoClient
}

func CloseMongoClient() {
	if err := mongoClient.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}
