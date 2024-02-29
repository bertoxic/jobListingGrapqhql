package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


var connectionString = "mongodb://localhost:27017/"

type DB struct {
	client *mongo.Client
}

func Connect()*DB {
mongoContext, cancel := context.WithTimeout(context.Background(), 20*time.Second)
defer cancel()
mongoClient, err:= mongo.Connect(mongoContext, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	err = mongoClient.Ping(mongoContext,readpref.Primary())
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Println("CONNECTED TO DATABASE")
	return &DB{
		client: mongoClient,
	}
}