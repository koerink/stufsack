package gnutils

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ClientMongo() *mongo.Client {
	mongoTestClient, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI("mongodb://application:2024##appsAdmin@103.82.92.6:26030/backendAPI?readPreference=primaryPreferred"))
	if err != nil {
		log.Fatal("Error while connecting to DB", err)
	}
	log.Println("Connection Successfully")
	err = mongoTestClient.Ping(context.Background(), readpref.PrimaryPreferred())
	if err != nil {
		log.Fatal("Ping Error")
	}
	log.Println("Ping Success")

	return mongoTestClient
}
