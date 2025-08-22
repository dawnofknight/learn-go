// Package main is the entry point for executable Go programs
package main

// Import the fmt package for formatted I/O operations
// import "fmt"

// addNumbers adds two integers and returns the result
func addNumbers(a, b int) int {
	return a + b
}

// getPointerValue returns the value pointed to by an integer pointer
func getPointerValue(ptr *int) int {
	return *ptr
}

func getNewData(data int) int {
	return data
}
