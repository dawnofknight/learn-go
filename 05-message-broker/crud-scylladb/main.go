package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

// User represents the user data structure
type User struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}

// UserTable metadata for ScyllaDB operations
var userMetadata = table.Metadata{
	Name:    "users",
	Columns: []string{"id", "name", "email", "created_at"},
	PartKey: []string{"id"},
}

var userTable = table.New(userMetadata)

// Database configuration
const (
	KeyspaceName = "example"
	TableName    = "users"
	ServerPort   = ":8080"
)

// Global session variable for HTTP handlers
var globalSession gocqlx.Session

// API Response structures
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// initializeDatabase creates keyspace and table if they don't exist
func initializeDatabase(session gocqlx.Session) error {
	// Create keyspace
	keyspaceQuery := fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s 
		WITH replication = {
			'class': 'SimpleStrategy',
			'replication_factor': 1
		}
	`, KeyspaceName)
	
	if err := session.ExecStmt(keyspaceQuery); err != nil {
		return fmt.Errorf("failed to create keyspace: %w", err)
	}
	
	// Create table in the keyspace (fully qualified name)
	tableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.%s (
			id text PRIMARY KEY,
			name text,
			email text,
			created_at timestamp
		)
	`, KeyspaceName, TableName)
	
	if err := session.ExecStmt(tableQuery); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	
	return nil
}

// createUser inserts a new user into the database
func createUser(session gocqlx.Session, user User) error {
	q := session.Query(userTable.Insert()).BindStruct(user)
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// getUserByID retrieves a user by ID
func getUserByID(session gocqlx.Session, id string) (*User, error) {
	var user User
	q := session.Query(userTable.Get()).BindMap(qb.M{"id": id})
	if err := q.GetRelease(&user); err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// updateUser updates an existing user
func updateUser(session gocqlx.Session, user User) error {
	q := session.Query(userTable.Update("name", "email")).BindStruct(user)
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// deleteUser removes a user by ID
func deleteUser(session gocqlx.Session, id string) error {
	q := session.Query(userTable.Delete()).BindMap(qb.M{"id": id})
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// getAllUsers retrieves all users from the database
func getAllUsers(session gocqlx.Session) ([]User, error) {
	var users []User
	q := session.Query(userTable.SelectAll())
	if err := q.SelectRelease(&users); err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	return users, nil
}

// HTTP Handlers

// createUserHandler handles POST /users
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Validate required fields
	if req.Name == "" || req.Email == "" {
		response := APIResponse{
			Success: false,
			Message: "Name and email are required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Create user
	user := User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}
	
	if err := createUser(globalSession, user); err != nil {
		response := APIResponse{
			Success: false,
			Message: "Failed to create user",
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	response := APIResponse{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// getUserHandler handles GET /users/{id}
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	vars := mux.Vars(r)
	userID := vars["id"]
	
	user, err := getUserByID(globalSession, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		
		response := APIResponse{
			Success: false,
			Message: "Failed to get user",
			Error:   err.Error(),
		}
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	response := APIResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	}
	json.NewEncoder(w).Encode(response)
}

// getAllUsersHandler handles GET /users
func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	users, err := getAllUsers(globalSession)
	if err != nil {
		response := APIResponse{
			Success: false,
			Message: "Failed to get users",
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	response := APIResponse{
		Success: true,
		Message: fmt.Sprintf("Retrieved %d users", len(users)),
		Data:    users,
	}
	json.NewEncoder(w).Encode(response)
}

// updateUserHandler handles PUT /users/{id}
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	vars := mux.Vars(r)
	userID := vars["id"]
	
	// Get existing user
	existingUser, err := getUserByID(globalSession, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		
		response := APIResponse{
			Success: false,
			Message: "User not found",
			Error:   err.Error(),
		}
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Update fields if provided
	if req.Name != "" {
		existingUser.Name = req.Name
	}
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	
	if err := updateUser(globalSession, *existingUser); err != nil {
		response := APIResponse{
			Success: false,
			Message: "Failed to update user",
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	response := APIResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    existingUser,
	}
	json.NewEncoder(w).Encode(response)
}

// deleteUserHandler handles DELETE /users/{id}
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	vars := mux.Vars(r)
	userID := vars["id"]
	
	// Check if user exists
	_, err := getUserByID(globalSession, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		
		response := APIResponse{
			Success: false,
			Message: "User not found",
			Error:   err.Error(),
		}
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	if err := deleteUser(globalSession, userID); err != nil {
		response := APIResponse{
			Success: false,
			Message: "Failed to delete user",
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	response := APIResponse{
		Success: true,
		Message: "User deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// healthHandler handles GET /health
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	response := APIResponse{
		Success: true,
		Message: "API is healthy",
		Data: map[string]interface{}{
			"timestamp": time.Now(),
			"version":   "1.0.0",
			"database":  "ScyllaDB",
		},
	}
	json.NewEncoder(w).Encode(response)
}

// setupRoutes configures all API routes
func setupRoutes() *mux.Router {
	r := mux.NewRouter()
	
	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/health", healthHandler).Methods("GET")
	api.HandleFunc("/users", createUserHandler).Methods("POST")
	api.HandleFunc("/users", getAllUsersHandler).Methods("GET")
	api.HandleFunc("/users/{id}", getUserHandler).Methods("GET")
	api.HandleFunc("/users/{id}", updateUserHandler).Methods("PUT")
	api.HandleFunc("/users/{id}", deleteUserHandler).Methods("DELETE")
	
	return r
}

// runDemo runs the original CRUD demo
func runDemo(session gocqlx.Session) {
	// Generate a unique ID for the user
	userID := uuid.New().String()
	
	// Create a new user
	user := User{
		ID:        userID,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
	}
	
	// Demonstrate CRUD operations
	fmt.Println("\n=== CRUD Operations Demo ===")
	
	// CREATE
	fmt.Println("\n1. Creating user...")
	if err := createUser(session, user); err != nil {
		log.Fatalf("Create operation failed: %v", err)
	}
	fmt.Printf("✓ User created successfully with ID: %s\n", userID)
	
	// READ
	fmt.Println("\n2. Reading user...")
	fetchedUser, err := getUserByID(session, userID)
	if err != nil {
		log.Fatalf("Read operation failed: %v", err)
	}
	fmt.Printf("✓ Found user: %+v\n", *fetchedUser)
	
	// UPDATE
	fmt.Println("\n3. Updating user...")
	fetchedUser.Name = "John Smith"
	fetchedUser.Email = "johnsmith@example.com"
	if err := updateUser(session, *fetchedUser); err != nil {
		log.Fatalf("Update operation failed: %v", err)
	}
	fmt.Println("✓ User updated successfully")
	
	// READ again to verify update
	updatedUser, err := getUserByID(session, userID)
	if err != nil {
		log.Fatalf("Read after update failed: %v", err)
	}
	fmt.Printf("✓ Updated user: %+v\n", *updatedUser)
	
	// LIST ALL
	fmt.Println("\n4. Listing all users...")
	allUsers, err := getAllUsers(session)
	if err != nil {
		log.Fatalf("List operation failed: %v", err)
	}
	fmt.Printf("✓ Found %d users:\n", len(allUsers))
	for i, u := range allUsers {
		fmt.Printf("   %d. %s (%s) - %s\n", i+1, u.Name, u.Email, u.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	
	// DELETE
	fmt.Println("\n5. Deleting user...")
	if err := deleteUser(session, userID); err != nil {
		log.Fatalf("Delete operation failed: %v", err)
	}
	fmt.Println("✓ User deleted successfully")
	
	// Verify deletion
	_, err = getUserByID(session, userID)
	if err != nil {
		fmt.Println("✓ Confirmed: User no longer exists")
	} else {
		fmt.Println("⚠ Warning: User still exists after deletion")
	}
	
	fmt.Println("\n=== CRUD Operations Demo Completed ===")
}

func main() {
	// Initialize ScyllaDB cluster
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Consistency = gocql.LocalQuorum
	cluster.ConnectTimeout = time.Second * 10
	cluster.Timeout = time.Second * 10
	
	// Create session for initialization
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatalf("Failed to connect to ScyllaDB: %v", err)
	}
	
	fmt.Println("Connected to ScyllaDB successfully!")
	
	// Initialize database (create keyspace and table)
	if err := initializeDatabase(session); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	
	fmt.Println("Database initialized successfully!")
	
	// Close the initial session
	session.Close()
	
	// Create a new session connected to the keyspace
	cluster.Keyspace = KeyspaceName
	keyspaceSession, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatalf("Failed to connect to keyspace: %v", err)
	}
	defer keyspaceSession.Close()
	
	// Set global session for HTTP handlers
	globalSession = keyspaceSession
	
	// Run demo if requested
	if len(os.Args) > 1 && os.Args[1] == "demo" {
		runDemo(session)
		return
	}
	
	// Setup HTTP routes
	router := setupRoutes()
	
	// Start HTTP server
	fmt.Printf("🚀 Starting REST API server on http://localhost%s\n", ServerPort)
	fmt.Println("📚 API Documentation:")
	fmt.Println("   GET    /api/v1/health          - Health check")
	fmt.Println("   GET    /api/v1/users           - Get all users")
	fmt.Println("   POST   /api/v1/users           - Create user")
	fmt.Println("   GET    /api/v1/users/{id}      - Get user by ID")
	fmt.Println("   PUT    /api/v1/users/{id}      - Update user")
	fmt.Println("   DELETE /api/v1/users/{id}      - Delete user")
	fmt.Println("\n💡 Run with 'go run main.go demo' to see CRUD demo")
	
	log.Fatal(http.ListenAndServe(ServerPort, router))
}