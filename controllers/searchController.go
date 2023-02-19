package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/goellavish10/fampay-assignment/configs"
	"go.mongodb.org/mongo-driver/mongo"
)

var videosCollection *mongo.Collection = configs.GetMongoCollection(configs.ConnectDB(), "videos")

var validate = validator.New()

func SearchController() string {
	return "Search"
}
