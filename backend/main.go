package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	//Check for persisted data or create a new one
	file, err := os.OpenFile("reviews.json", os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	content, _ := io.ReadAll(file)

	var previousReviews Response
	json.Unmarshal(content, &previousReviews)
	mostRecentTimestamp := time.Time{}

	//Validate what the most recent timestamp from the previous reviews is
	if len(previousReviews.Feed.Entry) > 0 {
		mostRecentTimestamp, _ = time.Parse(time.RFC3339, previousReviews.Feed.Entry[0].Updated.Label)
	}

	//Fetch reviews
	newResponses := getReviewsFromURL()

	var newEntries []Entry

	//If there were previous entries, check if the latest API responses occurred after the last timestamp
	if len(previousReviews.Feed.Entry) > 0 {
		for i := 0; i < len(newResponses.Feed.Entry); i++ {
			currentTime, _ := time.Parse(time.RFC3339, newResponses.Feed.Entry[i].Updated.Label)
			if currentTime.After(mostRecentTimestamp) {
				newEntries = append(newEntries, newResponses.Feed.Entry[i])
			}
		}
	}

	//Add any new entries to the persisted file
	allReviews := append(newEntries, newResponses.Feed.Entry...)
	newResponses.Feed.Entry = allReviews
	updatedFile, err := json.Marshal(newResponses)
	if err != nil {
		fmt.Println(err)
	}

	os.WriteFile("reviews.json", updatedFile, 0644)

	data, err := os.ReadFile("reviews.json")
	if err != nil {
		fmt.Println("File reading error", err)
	}

	port := ":8080"

	//Expose file to endpoint for client
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	fmt.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(port, nil))

}

// Handle fetching a few pages of customer reviews. 5 was an arbitrary number of pages but I used a more popular app for many of my tests
func getReviewsFromURL() Response {
	var responses Response
	for i := 1; i < 5; i++ {
		url := "https://itunes.apple.com/us/rss/customerreviews/id=447188370/sortBy=mostRecent/page=" + strconv.Itoa(i) + "/json"
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var newResponse Response
		json.Unmarshal(body, &newResponse)
		if err != nil {
			fmt.Println(err)
		}
		responses.Feed.Entry = append(responses.Feed.Entry, newResponse.Feed.Entry...)
	}
	return responses
}

// Handle CORS for local development
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// Structs for parsing iTunes url json
type Response struct {
	Feed Feed `json:"feed"`
}

type Feed struct {
	Author Author  `json:"author"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	Author  Author `json:"author"`
	Updated struct {
		Label string `json:"label"`
	} `json:"updated"`
	Rating struct {
		Label string `json:"label"`
	} `json:"im:rating"`
	Content struct {
		Label string `json:"label"`
	} `json:"content"`
}

type Author struct {
	URI struct {
		Label string `json:"label"`
	} `json:"uri"`
	Name struct {
		Label string `json:"label"`
	} `json:"name"`
}
