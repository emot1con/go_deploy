# Simple Go API

A very simple REST API built with Go and Gorilla Mux that manages users.

## Features

- ✅ RESTful endpoints for user management
- ✅ JSON responses
- ✅ In-memory data storage
- ✅ Health check endpoint
- ✅ Proper HTTP status codes
- ✅ Input validation

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Home page with API info |
| GET | `/api/v1/health` | Health check |
| GET | `/api/v1/users` | Get all users |
| POST | `/api/v1/users` | Create a new user |
| GET | `/api/v1/users/{id}` | Get user by ID |
| PUT | `/api/v1/users/{id}` | Update user by ID |
| DELETE | `/api/v1/users/{id}` | Delete user by ID |

## Running the API

1. Make sure you have Go installed
2. Navigate to the project directory:
   ```powershell
   cd "c:\Users\arqan\Documents\golang\go_deploy"
   ```

3. Install dependencies:
   ```powershell
   go mod tidy
   ```

4. Run the API:
   ```powershell
   go run cmd/main.go
   ```

5. The API will start on `http://localhost:8080`

## Usage Examples

### Get all users
```bash
curl http://localhost:8080/api/v1/users
```

### Create a new user
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Johnson", "email": "alice@example.com"}'
```

### Get a specific user
```bash
curl http://localhost:8080/api/v1/users/1
```

### Update a user
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "John Updated", "email": "john.updated@example.com"}'
```

### Delete a user
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

### Health check
```bash
curl http://localhost:8080/api/v1/health
```

## Sample Response

When you visit `http://localhost:8080/api/v1/users`, you'll get:

```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created": "2025-05-29T10:30:00Z"
  },
  {
    "id": 2,
    "name": "Jane Smith",
    "email": "jane@example.com",
    "created": "2025-05-29T10:30:00Z"
  }
]
```

## Project Structure

```
go_deploy/
├── cmd/
│   └── main.go     # Main application file
├── go.mod          # Go module file
└── README.md       # This file
```

## Notes

- This API uses in-memory storage, so data will be lost when the server restarts
- For production use, consider adding a database layer
- Add authentication and authorization as needed
- Consider adding logging middleware
- Add rate limiting for production use
