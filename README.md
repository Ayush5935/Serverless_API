# Serverless_API


Sure, here are the endpoints along with example JSON payloads for each:

1. **Register User Endpoint** (`POST /register`):
   - Example JSON Payload:
     ```json
     {
         "username": "example_user",
         "password": "example_password"
     }
     ```

2. **Login User Endpoint** (`POST /login`):
   - Example JSON Payload:
     ```json
     {
         "username": "example_user",
         "password": "example_password"
     }
     ```

3. **Get All Libraries Endpoint** (`GET /libraries`):
   - No JSON payload required. This endpoint retrieves all libraries from the database.

4. **Get All Books Endpoint** (`GET /books`):
   - No JSON payload required. This endpoint retrieves all books from the database.

5. **Add Library Endpoint** (`POST /add-library`):
   - Example JSON Payload:
     ```json
     {
         "name": "Example Library"
     }
     ```

6. **Add Book Endpoint** (`POST /add-book`):
   - Example JSON Payload:
     ```json
     {
         "title": "Example Book",
         "author": "Example Author",
         "library_id": "Library_ID"
     }
     ```





     go get github.com/lib/pq
go get github.com/google/uuid
go get github.com/joho/godotenv


Replace `"Library_ID"` in the example JSON payload for the `add-book` endpoint with the actual UUID of the library to which you want to add the book.

You can use these example JSON payloads in Postman to test each endpoint. Ensure that you set the correct request method (`POST` for endpoints that require data creation or modification, `GET` for endpoints that retrieve data) and that the server URL matches the server you are running locally (e.g., `http://localhost:8080`).





DATABASE_URL=postgres://username:password@localhost:5432/library_system
PORT=8080

You can add key-value pairs for your environment variables. For example:
Key: DATABASE_URL
Value: postgres://username:password@localhost:5432/library_system
Key: PORT
Value: 8080

