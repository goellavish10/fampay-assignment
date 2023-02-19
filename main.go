package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/goellavish10/fampay-assignment/configs"
	"github.com/goellavish10/fampay-assignment/models"
	"github.com/goellavish10/fampay-assignment/responses"
	"github.com/goellavish10/fampay-assignment/routes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Connect to database
	configs.ConnectDB()

	// Routes
	routes.SearchRoutes(router)
	var apiKey string = configs.GetEnv("GOOGLE_CONSOLE_API_KEY_0")
	go callYouTubeApi(apiKey)

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func callYouTubeApi(apiKey string) {
	var apiKeysUsed = []string{}
	var count int = 0
	recallApiTicker := time.NewTicker(20 * time.Second)
	var validate = validator.New()
	var videoCollection *mongo.Collection = configs.GetMongoCollection(configs.DB, "videos")
	for msg := range recallApiTicker.C {
		fmt.Println("Calling YouTube API at: ", msg)
		// Calling the YouTube API
		client := &http.Client{}

		req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/search?key="+apiKey+"&part=snippet&maxResults=50&order=date&publishedAfter=2018-01-01T00:00:00Z&q=asmr&type=video", nil)

		if err != nil {
			fmt.Println(err.Error())
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Println(resp.StatusCode)
			if resp.StatusCode == 403 {

				fmt.Println("API Key has exceeded the quota")

			USED_API_KEY_COUNTER:
				apiKeysUsed = append(apiKeysUsed, strconv.Itoa(count))
				count++
				var newApiKeyCount string = strconv.Itoa(count)
				// Check if the newApiCount is already present in the apiKeysUsed array
				var idx int = searchInApiArray(apiKeysUsed, newApiKeyCount)
				if idx == -1 {
					apiKey = configs.GetEnv("GOOGLE_CONSOLE_API_KEY_" + newApiKeyCount)
					fmt.Println("New Api Key Used: ", apiKey)
				} else {
					goto USED_API_KEY_COUNTER
				}

				fmt.Println("New Api Key Used: ", count)
				continue
			}
			return
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println(err.Error())
		}

		var responseObject responses.Response
		json.Unmarshal(bodyBytes, &responseObject)
		// fmt.Println(responseObject.Items)

		for k := 0; k < len(responseObject.Items); k++ {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			var video models.Video
			defer cancel()

			if validationErr := validate.Struct(&video); validationErr != nil {
				fmt.Println("Validation Error: ", validationErr)
				continue
			}

			// Check if the video already exists in the database
			// objId, _ := primitive.ObjectIDFromHex(responseObject.Items[k].ID.VideoID)
			var filter = bson.M{"videoid": responseObject.Items[k].ID.VideoID}
			err := videoCollection.FindOne(ctx, filter).Decode(&video)

			fmt.Println(err)

			if err != nil {
				// Video does not exist in the database
				newVideo := models.Video{
					Id:          primitive.NewObjectID(),
					Title:       responseObject.Items[k].Snippet.Title,
					Description: responseObject.Items[k].Snippet.Description,
					PublishedAt: responseObject.Items[k].Snippet.PublishedAt,
					VideoId:     responseObject.Items[k].ID.VideoID,
					Thumbnail: models.Thumbnail{
						Default: models.ThumbnailProperties{
							Url:    responseObject.Items[k].Snippet.Thumbnails.Default.Url,
							Width:  responseObject.Items[k].Snippet.Thumbnails.Default.Width,
							Height: responseObject.Items[k].Snippet.Thumbnails.Default.Height,
						},
						Medium: models.ThumbnailProperties{
							Url:    responseObject.Items[k].Snippet.Thumbnails.Medium.Url,
							Width:  responseObject.Items[k].Snippet.Thumbnails.Medium.Width,
							Height: responseObject.Items[k].Snippet.Thumbnails.Medium.Height,
						},
						High: models.ThumbnailProperties{
							Url:    responseObject.Items[k].Snippet.Thumbnails.High.Url,
							Width:  responseObject.Items[k].Snippet.Thumbnails.High.Width,
							Height: responseObject.Items[k].Snippet.Thumbnails.High.Height,
						},
					},
				}

				// Inserting the video into the database
				result, err := videoCollection.InsertOne(ctx, newVideo)

				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("Inserted a single document: ", result.InsertedID)
				continue
			}

			fmt.Println("Video already exists in the database: ", video.Title)

		}

		resp.Body.Close()
		// fmt.Println(responseObject)
	}

}

func searchInApiArray(apiArray []string, newApiKeyCount string) int {

	idx := slices.IndexFunc(apiArray, func(c string) bool { return c == newApiKeyCount })

	fmt.Println(idx)

	return idx
}
