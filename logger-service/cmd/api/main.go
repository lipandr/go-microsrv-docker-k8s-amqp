package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	//webPort  = "80"
	//rpcPort  = "5001"
	//gRpcPort = "50001"
	mongoURL = "mongodb://mongo:27001"
)

var client *mongo.Client

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient
}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}
	return c, nil
}