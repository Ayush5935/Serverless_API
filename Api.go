package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	// Make HTTP request to another API
	resp, err := http.Get("https://api.example.com/endpoint")
	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Read response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		log.Printf("Error decoding response body: %v", err)
		return
	}

	// Print response body
	log.Println("Response:", responseBody)
}
￼Enter
