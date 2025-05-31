// Package main is the entry point for our API application
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// SECTION 1: DATA STRUCTURES
//===============================

// APIResponse defines the structure of our API responses
// This ensures consistent response format across all endpoints
type APIResponse struct {
	Status     string      `json:"status"`         // "success" or "error"
	Message    string      `json:"message"`        // Human-readable message
	Data       interface{} `json:"data,omitempty"` // Optional data payload
	Timestamp  string      `json:"timestamp"`      // When the response was generated
	StatusCode int         `json:"-"`              // HTTP status code (not shown in JSON)
}

// User defines the structure of our user data
// The `json` tags tell the encoder how to convert to/from JSON
type User struct {
	ID    string `json:"id"`    // Unique identifier
	Name  string `json:"name"`  // User's name
	Email string `json:"email"` // User's email
}

// SECTION 2: HANDLER STRUCTURE
//===============================

// UserHandler manages user-related operations
// It contains a map to store users in memory
type UserHandler struct {
	users map[string]User // Map of user ID to User struct
}

// NewUserHandler creates a new UserHandler instance
// This is a constructor function
func NewUserHandler() *UserHandler {
	return &UserHandler{
		users: make(map[string]User), // Initialize empty map
	}
}

// SECTION 3: HELPER FUNCTIONS
//===============================

// sendResponse is a helper function to send standardized JSON responses
// It handles setting headers, status code, and logging
func sendResponse(w http.ResponseWriter, statusCode int, status, message string, data interface{}) {
	// Create the response structure
	response := APIResponse{
		Status:     status,
		Message:    message,
		Data:       data,
		Timestamp:  time.Now().Format(time.RFC3339),
		StatusCode: statusCode,
	}

	// Set response headers and status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Encode and send the response as JSON
	json.NewEncoder(w).Encode(response)

	// Log the response for debugging
	log.Printf("[Response] Status: %d, Message: %s", statusCode, message)
}

// SECTION 4: HTTP HANDLERS
//===============================

// handleGet processes GET requests
// It either returns all users or a specific user by ID
func (h *UserHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GET] Request received: %s", r.URL.Path)

	// Extract user ID from URL path
	id := strings.TrimPrefix(r.URL.Path, "/users/")

	// If no ID provided, return all users
	if id == "" {
		users := make([]User, 0, len(h.users))
		for _, user := range h.users {
			users = append(users, user)
		}
		sendResponse(w, http.StatusOK, "success", "Users retrieved successfully", users)
		log.Printf("[GET] Retrieved %d users", len(users))
		return
	}

	// If ID provided, return specific user
	if user, exists := h.users[id]; exists {
		sendResponse(w, http.StatusOK, "success",
			fmt.Sprintf("User %s retrieved successfully", id), user)
		log.Printf("[GET] Retrieved user: %s", id)
		return
	}

	// User not found
	sendResponse(w, http.StatusNotFound, "error",
		fmt.Sprintf("User %s not found", id), nil)
	log.Printf("[GET] User not found: %s", id)
}

// handlePost processes POST requests to create new users
func (h *UserHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	log.Printf("[POST] Request received")

	// Decode JSON request body into User struct
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendResponse(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		log.Printf("[POST] Error decoding request: %v", err)
		return
	}

	// Check if user already exists
	if _, exists := h.users[user.ID]; exists {
		sendResponse(w, http.StatusConflict, "error",
			fmt.Sprintf("User %s already exists", user.ID), nil)
		log.Printf("[POST] User already exists: %s", user.ID)
		return
	}

	// Save new user
	h.users[user.ID] = user
	sendResponse(w, http.StatusCreated, "success", "User created successfully", user)
	log.Printf("[POST] Created user: %s", user.ID)
}

// PUT handler with logging
func (h *UserHandler) handlePut(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL path
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	log.Printf("[PUT] Request received for user: %s", id)

	// Decode the request body into User struct
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendResponse(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		log.Printf("[PUT] Error decoding request: %v", err)
		return
	}

	// Set the ID from the URL path
	user.ID = id

	// Update or create the user
	h.users[id] = user

	// Send success response
	sendResponse(w, http.StatusOK, "success", fmt.Sprintf("User %s updated successfully", id), user)
	log.Printf("[PUT] Updated user: %s", id)
}

// PATCH handler with logging
func (h *UserHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	log.Printf("[PATCH] Request received for user: %s", id)

	// Check if user exists
	existingUser, exists := h.users[id]
	if !exists {
		sendResponse(w, http.StatusNotFound, "error", fmt.Sprintf("User %s not found", id), nil)
		log.Printf("[PATCH] User not found: %s", id)
		return
	}

	// Decode the updates from request body
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		sendResponse(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		log.Printf("[PATCH] Error decoding request: %v", err)
		return
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		existingUser.Name = name
	}
	if email, ok := updates["email"].(string); ok {
		existingUser.Email = email
	}

	// Save updated user
	h.users[id] = existingUser
	sendResponse(w, http.StatusOK, "success", fmt.Sprintf("User %s patched successfully", id), existingUser)
	log.Printf("[PATCH] Patched user: %s", id)
}

func (h *UserHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	log.Printf("[DELETE] Request received for user: %s", id)

	// Check if user exists
	_, exists := h.users[id]
	if !exists {
		sendResponse(w, http.StatusNotFound, "error", fmt.Sprintf("User %s not found", id), nil)
		log.Printf("[DELETE] User not found: %s", id)
		return
	}

	// Delete the user
	delete(h.users, id)
	sendResponse(w, http.StatusOK, "success", fmt.Sprintf("User %s deleted successfully", id), nil)
	log.Printf("[DELETE] Deleted user: %s", id)
}

// SECTION 5: MAIN ROUTER
//===============================

// ServeHTTP is the main router that handles all incoming requests
func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log all incoming requests
	log.Printf("[Request] %s %s", r.Method, r.URL.Path)

	// Route to appropriate handler based on HTTP method
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodPut:
		h.handlePut(w, r)
	case http.MethodPatch:
		h.handlePatch(w, r)
	case http.MethodDelete:
		h.handleDelete(w, r)
	default:
		sendResponse(w, http.StatusMethodNotAllowed, "error",
			"Method not allowed", nil)
		log.Printf("[Error] Method not allowed: %s", r.Method)
	}
}

// SECTION 6: STRUCT EXAMPLES
//===============================

// demonstrateStructUsage shows various ways to use structs in Go
// This function demonstrates creating, accessing, and manipulating structs
func demonstrateStructUsage() {
	fmt.Println("\n=== STRUCT USAGE EXAMPLES ===")

	// 1. Creating and Using Struct Variables
	fmt.Println("\n1. Creating and Using Struct Variables:")

	// Create an instance of APIResponse
	response := APIResponse{
		Status:     "success",
		Message:    "User retrieved successfully",
		Data:       nil, // Using the interface{} field
		Timestamp:  time.Now().Format(time.RFC3339),
		StatusCode: 200,
	}

	// Create a User struct
	user := User{
		ID:    "u123",
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Assign the user to the Data field of response
	response.Data = user

	// Access struct fields
	fmt.Println("Response status:", response.Status)
	fmt.Println("User name:", user.Name)

	// 2. Using Structs in Conditionals
	fmt.Println("\n2. Using Structs in Conditionals:")

	// Check response status
	if response.Status == "success" {
		fmt.Println("Operation was successful")
	} else {
		fmt.Println("Operation failed with message:", response.Message)
	}

	// Check HTTP status code
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		fmt.Println("Successful HTTP response")
	} else if response.StatusCode >= 400 && response.StatusCode < 500 {
		fmt.Println("Client error")
	} else {
		fmt.Println("Server error or other status")
	}

	// Check if Data field is empty
	if response.Data == nil {
		fmt.Println("No data provided")
	} else {
		fmt.Println("Data is provided")
	}

	// 3. Using Structs in Loops (with a slice of structs)
	fmt.Println("\n3. Using Structs in Loops:")

	// Create a slice of User structs
	users := []User{
		{ID: "u123", Name: "John Doe", Email: "john@example.com"},
		{ID: "u456", Name: "Jane Smith", Email: "jane@example.com"},
		{ID: "u789", Name: "Bob Johnson", Email: "bob@example.com"},
	}

	// Iterate through the slice using a for loop
	fmt.Println("All users:")
	for i, user := range users {
		fmt.Printf("%d: %s (ID: %s, Email: %s)\n", i+1, user.Name, user.ID, user.Email)
	}

	// Filter users based on a condition
	fmt.Println("\nUsers with 'J' names:")
	for _, user := range users {
		if strings.HasPrefix(user.Name, "J") {
			fmt.Printf("%s starts with J\n", user.Name)
		}
	}

	// 4. Using Structs with Maps
	fmt.Println("\n4. Using Structs with Maps:")

	// Create a map with user IDs as keys and User structs as values
	userMap := make(map[string]User)

	// Add users to the map
	for _, user := range users {
		userMap[user.ID] = user
	}

	// Look up a specific user
	if user, exists := userMap["u456"]; exists {
		fmt.Printf("Found user: %s\n", user.Name)
	} else {
		fmt.Println("User not found")
	}

	// Iterate through the map
	fmt.Println("\nAll users in map:")
	for id, user := range userMap {
		fmt.Printf("ID %s: %s\n", id, user.Name)
	}

	// 5. Using the `interface{}` Field
	fmt.Println("\n5. Using the interface{} Field:")

	// Create different types of responses with various data types
	responses := []APIResponse{
		{
			Status:    "success",
			Message:   "Retrieved user",
			Data:      User{ID: "u123", Name: "John Doe", Email: "john@example.com"},
			Timestamp: time.Now().Format(time.RFC3339),
		},
		{
			Status:    "success",
			Message:   "Retrieved count",
			Data:      42, // integer data
			Timestamp: time.Now().Format(time.RFC3339),
		},
		{
			Status:    "success",
			Message:   "Retrieved names",
			Data:      []string{"Alice", "Bob", "Charlie"}, // slice data
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}

	// Process different data types using type assertions
	for i, resp := range responses {
		fmt.Printf("\nResponse %d: %s\n", i+1, resp.Message)

		// Use type assertion to handle different data types
		switch data := resp.Data.(type) {
		case User:
			fmt.Printf("User data: %s (ID: %s)\n", data.Name, data.ID)
		case int:
			fmt.Printf("Integer data: %d\n", data)
		case []string:
			fmt.Println("String slice data:")
			for _, name := range data {
				fmt.Printf("  - %s\n", name)
			}
		default:
			fmt.Println("Unknown data type")
		}
	}
}

// SECTION 7: MAIN FUNCTION
//===============================

func main() {
	// Configure logging to include date, time, and file information
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Demonstrate struct usage examples
	demonstrateStructUsage()

	// Create new handler instance
	userHandler := NewUserHandler()

	// Register handler for /users/ path
	http.Handle("/users/", userHandler)

	// Start the server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
