package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Running on AWS Lambda
		lambda.Start(handler)
	} else {
		// Running locally
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Server listening on port %s...", port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "GET":
		switch request.Path {
		case "/libraries":
			libraries, err := getAllLibraries()
			if err != nil {
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			}
			body, err := json.Marshal(libraries)
			if err != nil {
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			}
			return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(body)}, nil
		case "/books":
			books, err := getAllBooks()
			if err != nil {
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			}
			body, err := json.Marshal(books)
			if err != nil {
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			}
			return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(body)}, nil
		default:
			return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, nil
		}
	case "POST":
		switch request.Path {
		case "/add-library":
			var library Library
			if err := json.Unmarshal([]byte(request.Body), &library); err != nil {
				return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
			}
			if err := addLibrary(library); err != nil {
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			}
			return events.APIGatewayProxyResponse{StatusCode: http.StatusCreated}, nil
		case "/add-book":
			var book Book
			if err := json.Unmarshal([]byte(request.Body), &book); err != nil {
				return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
			}
			if err := addBook(book); err != nil {
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			}
			return events.APIGatewayProxyResponse{StatusCode: http.StatusCreated}, nil
		default:
			return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, nil
		}
	default:
		return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed}, nil
	}
}

func setupDB() error {
	// Database setup logic...
}

func getAllLibraries() ([]Library, error) {
	// Retrieve all libraries from the database
	// and return them as a slice of Library structs
}

func getAllBooks() ([]Book, error) {
	// Retrieve all books from the database
	// and return them as a slice of Book structs
}

func addLibrary(library Library) error {
	// Add a new library to the database
}

func addBook(book Book) error {
	// Add a new book to the database
}
