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
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK) // Return 200 OK for preflight request
		return
	}

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
		jsonData, err := json.Marshal(RequestBody{Message: word})
		if err != nil {
			http.Error(w, "Error encoding message", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s\n", jsonData)
		log.Printf("Sent message: %s", jsonData) // Log the sent message
		flusher.Flush()
		time.Sleep(300 * time.Millisecond)

		if index == len(words)-1 && sendImage {
			imgMessage, _ := json.Marshal(RequestBody{Message: "Image URL: https://source.unsplash.com/random"})
			fmt.Fprintf(w, "%s\n", imgMessage)
			log.Printf("Sent message: %s", imgMessage) // Log the sent image message
			flusher.Flush()
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
