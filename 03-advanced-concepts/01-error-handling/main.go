package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// Basic error handling
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Custom error type
type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Msg)
}

func validateUsername(username string) error {
	if len(username) < 3 {
		return &ValidationError{"username", "too short"}
	}
	if len(username) > 20 {
		return &ValidationError{"username", "too long"}
	}
	return nil
}

// Error wrapping (Go 1.13+)
func openConfigFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return data, nil
}

// Sentinel errors
var (
	ErrNotFound     = errors.New("item not found")
	ErrUnauthorized = errors.New("unauthorized access")
)

func findItem(id string) (string, error) {
	// Simulate a database lookup
	if id == "404" {
		return "", ErrNotFound
	}
	if id == "401" {
		return "", ErrUnauthorized
	}
	return "Item " + id, nil
}

// Error handling with multiple return values
func processItem(id string) (string, int, error) {
	item, err := findItem(id)
	if err != nil {
		return "", 0, err
	}
	return item, len(item), nil
}

// Panic and recover
func recoverExample() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	fmt.Println("Calling function that may panic...")
	panickingFunction()
	fmt.Println("This line will not be executed")
}

func panickingFunction() {
	fmt.Println("About to panic...")
	panic("something went terribly wrong")
}

// Defer statement examples
func deferExample() string {
	file, err := os.CreateTemp("", "example")
	if err != nil {
		return fmt.Sprintf("Error creating temp file: %v", err)
	}
	defer file.Close()
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte("Hello, World!")); err != nil {
		return fmt.Sprintf("Error writing to file: %v", err)
	}

	return fmt.Sprintf("Successfully wrote to temp file: %s", file.Name())
}

// Multiple defer statements (executed in LIFO order)
func multipleDeferExample() {
	fmt.Println("Start")
	defer fmt.Println("First defer")
	defer fmt.Println("Second defer")
	defer fmt.Println("Third defer")
	fmt.Println("End")
}

// Practical example: database transaction
func simulateTransaction() error {
	// In a real application, this would be a database transaction
	fmt.Println("Starting transaction...")

	// Simulate beginning a transaction
	defer fmt.Println("Transaction ended")

	// Simulate operations that might fail
	if err := simulateOperation("Operation 1"); err != nil {
		fmt.Println("Rolling back transaction due to error")
		return err
	}

	if err := simulateOperation("Operation 2"); err != nil {
		fmt.Println("Rolling back transaction due to error")
		return err
	}

	fmt.Println("Committing transaction")
	return nil
}

func simulateOperation(name string) error {
	fmt.Printf("Executing %s...\n", name)
	// Simulate a random failure
	if time.Now().UnixNano()%10 == 0 {
		return fmt.Errorf("%s failed", name)
	}
	return nil
}

// Error handling with type assertions
func handleDifferentErrors() {
	errs := []error{
		errors.New("generic error"),
		&ValidationError{"email", "invalid format"},
		fmt.Errorf("wrapped error: %w", ErrNotFound),
	}

	for _, err := range errs {
		switch e := err.(type) {
		case *ValidationError:
			fmt.Printf("Validation error on field %s: %s\n", e.Field, e.Msg)
		default:
			// Check for wrapped errors
			if errors.Is(err, ErrNotFound) {
				fmt.Println("Item not found error")
			} else {
				fmt.Printf("Unknown error: %v\n", err)
			}
		}
	}
}

// Implementing a retry mechanism with errors
func retryOperation(operation func() error, maxRetries int) error {
	var err error

	for i := 0; i < maxRetries; i++ {
		err = operation()
		if err == nil {
			return nil // Success
		}

		fmt.Printf("Attempt %d failed: %v. Retrying...\n", i+1, err)
		time.Sleep(time.Duration(i*100) * time.Millisecond) // Exponential backoff
	}

	return fmt.Errorf("operation failed after %d attempts: %w", maxRetries, err)
}

func unreliableOperation() error {
	// Simulate an operation that sometimes fails
	if time.Now().UnixNano()%3 == 0 {
		return errors.New("temporary failure")
	}
	return nil
}

// Error handling with cleanup
func processFile(filename string) error {
	file, err := os.CreateTemp("", filename)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	// Ensure the file is removed when we're done
	defer func() {
		file.Close()
		os.Remove(file.Name())
		fmt.Println("Cleanup complete")
	}()

	// Process the file
	if _, err := file.Write([]byte("Processing data")); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println("File processed successfully")
	return nil
}

func main() {
	fmt.Println("=== Basic Error Handling ===")
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}

	result, err = divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}

	fmt.Println("\n=== Custom Error Types ===")
	err = validateUsername("jo")
	if err != nil {
		fmt.Println("Error:", err)
		// Type assertion to access the custom error fields
		if validationErr, ok := err.(*ValidationError); ok {
			fmt.Printf("Field: %s, Message: %s\n", validationErr.Field, validationErr.Msg)
		}
	}

	fmt.Println("\n=== Error Wrapping ===")
	_, err = openConfigFile("nonexistent.conf")
	if err != nil {
		fmt.Println("Error:", err)

		// Unwrap to get the underlying error
		var pathError *os.PathError
		if errors.As(err, &pathError) {
			fmt.Printf("Path error: op=%s path=%s\n", pathError.Op, pathError.Path)
		}
	}

	fmt.Println("\n=== Sentinel Errors ===")
	item, err := findItem("123")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Found:", item)
	}

	_, err = findItem("404")
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			fmt.Println("Item not found error occurred")
		} else {
			fmt.Println("Other error:", err)
		}
	}

	fmt.Println("\n=== Panic and Recover ===")
	recoverExample()

	fmt.Println("\n=== Defer Statements ===")
	fmt.Println(deferExample())

	fmt.Println("\n=== Multiple Defer Statements ===")
	multipleDeferExample()

	fmt.Println("\n=== Transaction Example ===")
	err = simulateTransaction()
	if err != nil {
		fmt.Println("Transaction failed:", err)
	}

	fmt.Println("\n=== Error Type Handling ===")
	handleDifferentErrors()

	fmt.Println("\n=== Retry Mechanism ===")
	err = retryOperation(unreliableOperation, 5)
	if err != nil {
		fmt.Println("Final error:", err)
	} else {
		fmt.Println("Operation succeeded")
	}

	fmt.Println("\n=== Error Handling with Cleanup ===")
	err = processFile("example.txt")
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("\nProgram completed successfully")
}
