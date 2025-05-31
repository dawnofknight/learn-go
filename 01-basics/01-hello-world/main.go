// Package main is the entry point for executable Go programs
package main

// Import the fmt package for formatted I/O operations
import "fmt"

// The main function is the entry point for the program
// It will be executed automatically when the program runs
func main() {
	// Print a message to the console
	fmt.Println("Hello, World!")

	// Print a formatted message
	name := "Gopher"
	fmt.Printf("Hello, %s! Welcome to Go programming.\n", name)

	// Print multiple values
	fmt.Println("Go is", "awesome", "and", "fun!")

	// Print with newlines for each item
	fmt.Println("Learning Go:\n1. Basics\n2. Data Structures\n3. Advanced Concepts")
}
