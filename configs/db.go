package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoURIString()))

	if err != nil {
		log.Fatal(err)
	}

	// Defining timeout context for connecting to database
	ctx, _ := context.WithTimeout(context.Background(), 12*time.Second)

	// checking for errors and connecting to MongoDB Client
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB....")
	return client
}

var DB *mongo.Client = ConnectDB()

func GetMongoCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("fampay-assignment").Collection(collectionName)
	// mod := mongo.IndexModel{Keys: bson.D{{Key: "title", Value: "text"}}}
	// name, err := collection.Indexes().CreateOne(context.TODO(), mod)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Created index: ", name)

	return collection
}
