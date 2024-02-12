package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to PostgreSQL database
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Create database tables if they don't exist
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

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Running on AWS Lambda
		lambda.Start(handler)
	} else {
		// Running locally
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Retrieve environment variables if needed
	databaseURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")

	// Handle requests...
}

func setupDB() error {
	// Database setup function...
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	// Register user handler...
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	// Login user handler...
}

func getAllLibraries(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	// Fetch libraries concurrently using goroutines and channels
	ch := make(chan []Library)
	go fetchLibraries(limit, ch)

	// Collect results from channel
	libraries := <-ch

	// Return response
	json.NewEncoder(w).Encode(libraries)
}

func fetchLibraries(limit int, ch chan<- []Library) {
	rows, err := db.Query("SELECT id, name FROM libraries ORDER BY name LIMIT $1", limit)
	if err != nil {
		log.Printf("Error fetching libraries: %v", err)
		ch <- nil
		return
	}
	defer rows.Close()

	var libraries []Library
	for rows.Next() {
		var library Library
		if err := rows.Scan(&library.ID, &library.Name); err != nil {
			log.Printf("Error scanning library row: %v", err)
			ch <- nil
			return
		}
		libraries = append(libraries, library)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over library rows: %v", err)
		ch <- nil
		return
	}

	ch <- libraries
}

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	// Similar implementation as getAllLibraries, but for books
}

func addLibrary(w http.ResponseWriter, r *http.Request) {
	// Add library handler...
}

func addBook(w http.ResponseWriter, r *http.Request) {
	// Add book handler...
}