package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// User represents a simple user structure
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Created string `json:"created"`
}

// In-memory storage for users
var users []User
var nextID = 1

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	// Initialize with some sample data
	users = []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com", Created: time.Now().Format(time.RFC3339)},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Created: time.Now().Format(time.RFC3339)},
	}
	nextID = 3
	// Create router
	r := mux.NewRouter()

	// Add CORS middleware
	r.Use(corsMiddleware)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/users", getUsers).Methods("GET")
	api.HandleFunc("/users", createUser).Methods("POST")
	api.HandleFunc("/users/{id}", getUser).Methods("GET")
	api.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	api.HandleFunc("/health", healthCheck).Methods("GET")
	api.HandleFunc("/home", homeHandler).Methods("GET")

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	// Root endpoint - serve the frontend
	r.HandleFunc("/", frontendHandler).Methods("GET")

	// Get port from environment variable or use default
	portEnv := os.Getenv("TEST_ENV")
	if portEnv == "" {
		portEnv = "8080" // Default port if TEST_ENV is not set
	}
	port := ":" + portEnv
	fmt.Printf("🚀 Server starting on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET    /              - Frontend (User Management)")
	fmt.Println("  GET    /api/v1/home   - API Home page")
	fmt.Println("  GET    /api/v1/health - Health check")
	fmt.Println("  GET    /api/v1/users  - Get all users")
	fmt.Println("  POST   /api/v1/users  - Create new user")
	fmt.Println("  GET    /api/v1/users/{id} - Get user by ID")
	fmt.Println("  PUT    /api/v1/users/{id} - Update user by ID")
	fmt.Println("  DELETE /api/v1/users/{id} - Delete user by ID")
	fmt.Printf("\n🌐 Frontend available at: http://0.0.0.0%s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0"+port, r))
}

// corsMiddleware adds CORS headers to all responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// frontendHandler serves the frontend HTML file
func frontendHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

// homeHandler handles the API home endpoint
func homeHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Welcome to the Simple API",
		"version": "1.0.0",
		"status":  "running",
		"endpoints": map[string]string{
			"health": "/api/v1/health",
			"users":  "/api/v1/users",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// healthCheck handles health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    "running",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getUsers returns all users
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// getUser returns a specific user by ID
func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	for _, user := range users {
		if user.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

// createUser creates a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if user.Name == "" || user.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	// Set ID and creation time
	user.ID = nextID
	nextID++
	user.Created = time.Now().Format(time.RFC3339)

	// Add to users slice
	users = append(users, user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// updateUser updates an existing user
func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Find and update user
	for i, user := range users {
		if user.ID == id {
			updatedUser.ID = id
			updatedUser.Created = user.Created // Keep original creation time
			users[i] = updatedUser

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedUser)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

// deleteUser deletes a user by ID
func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Find and remove user
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}
