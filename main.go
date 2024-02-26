package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type RequestBody struct {
	Message string `json:"message"`
}

var (
	loremIpsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
	words      = strings.Split(loremIpsum, " ")
)

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	sendImage := strings.Contains(strings.ToLower(requestBody.Message), "image")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for index, word := range words {
		fmt.Fprintf(w, "data: %s\n\n", word)
		flusher.Flush()
		time.Sleep(300 * time.Millisecond)

		if index == len(words)-1 {
			if sendImage {
				fmt.Fprintf(w, "data: Image URL: https://source.unsplash.com/random\n\n")
				flusher.Flush()
			}
			return
		}
	}
}

func main() {
	http.HandleFunc("/events", sseHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
