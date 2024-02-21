package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type Library struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Book struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	LibraryID uuid.UUID `json:"library_id"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err := setupDB(); err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}
}

func main() {
	http.HandleFunc("/register", registerUser)
	http.HandleFunc("/login", loginUser)
	http.HandleFunc("/libraries", getAllLibraries)
	http.HandleFunc("/books", getAllBooks)
	http.HandleFunc("/add-library", addLibrary)
	http.HandleFunc("/add-book", addBook)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func setupDB() error {
	// Database setup logic remains the same
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	// User registration handler logic remains the same
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	// User login handler logic remains the same
}

func getAllLibraries(w http.ResponseWriter, r *http.Request) {
	// Fetch libraries logic remains the same

	// Adjust endpoint URL for AWS Lambda
	lambdaEndpoint := os.Getenv("LAMBDA_ENDPOINT")
	if lambdaEndpoint != "" {
		// Modify the endpoint URL for Lambda
		lambdaEndpoint += "/books"
	}

	resp, err := http.Get(lambdaEndpoint)
	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}
	log.Printf("Response body: %s", body)
}

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	// Fetch books logic remains the same

	// Adjust endpoint URL for AWS Lambda
	lambdaEndpoint := os.Getenv("LAMBDA_ENDPOINT")
	if lambdaEndpoint != "" {
		// Modify the endpoint URL for Lambda
		lambdaEndpoint += "/login"
	}

	requestBody := map[string]string{
		"username": "user1",
		"password": "password1",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshaling request body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(lambdaEndpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Error making HTTP POST request to add library: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}
	log.Printf("Response body: %s", body)
}

func addLibrary(w http.ResponseWriter, r *http.Request) {
	// Add library handler logic remains the same
}

func addBook(w http.ResponseWriter, r *http.Request) {
	// Add book handler logic remains the same
}
