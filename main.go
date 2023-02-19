package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/goellavish10/fampay-assignment/configs"
	"github.com/goellavish10/fampay-assignment/responses"
	"github.com/goellavish10/fampay-assignment/routes"
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
	recallApiTicker := time.NewTicker(3 * time.Second)

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
		resp.Body.Close()
		// fmt.Println(responseObject)
	}

}

func searchInApiArray(apiArray []string, newApiKeyCount string) int {

	idx := slices.IndexFunc(apiArray, func(c string) bool { return c == newApiKeyCount })

	fmt.Println(idx)

	return idx
}
