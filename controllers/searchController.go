package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/goellavish10/fampay-assignment/configs"
	"github.com/goellavish10/fampay-assignment/models"
	"github.com/goellavish10/fampay-assignment/responses"

	// "github.com/hashicorp/hcl/hcl/strconv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var validate = validator.New()
var videosCollection *mongo.Collection = configs.GetMongoCollection(configs.ConnectDB(), "videos")

func SearchController() http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		var searchQuery string = r.URL.Query().Get("searchQuery")
		var searchType string = r.URL.Query().Get("searchType")
		// var searchQuery string = params["searchQuery"]
		// var video models.Video
		defer cancel()

		var filter bson.M = bson.M{}

		if searchType == "title" {
			filter = bson.M{"title": bson.M{"$regex": primitive.Regex{Pattern: "^" + searchQuery + "*", Options: "i"}}}
		} else if searchType == "description" {
			filter = bson.M{"description": bson.M{"$regex": primitive.Regex{Pattern: "^" + searchQuery + "*", Options: "i"}}}
		}

		cursor, err := videosCollection.Find(ctx, filter)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		var results []models.Video
		if err = cursor.All(ctx, &results); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		for _, result := range results {
			res, _ := json.Marshal(result)
			fmt.Println(string(res))
		}

		fmt.Println(results)

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": results}}
		json.NewEncoder(rw).Encode(response)

	}
}

func GetAllStoredVideos() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		var page string = r.URL.Query().Get("page")
		if page == "" {
			page = "1"
		}
		defer cancel()
		toSkip, err := strconv.Atoi(page)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		opts := options.Find().SetSort(bson.D{{"publishedat", -1}}).SetSkip(int64(toSkip - 1)).SetLimit(10)

		cursor, err := videosCollection.Find(ctx, bson.M{}, opts)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		var results []models.Video
		if err = cursor.All(ctx, &results); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		for _, result := range results {
			res, _ := json.Marshal(result)
			fmt.Println(string(res))
		}

		fmt.Println(results)

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": results}}
		json.NewEncoder(rw).Encode(response)

	}
}
