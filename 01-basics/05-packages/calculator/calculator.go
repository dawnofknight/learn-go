// Package calculator provides basic and advanced mathematical operations
package calculator

import (
	"errors"
	"fmt"
)

// Version is an exported package variable
var Version = "1.0.0"

// unexportedMaxValue is a package-private variable (not accessible outside)
var unexportedMaxValue = 1000000

// init function runs when the package is initialized
func init() {
	fmt.Println("Calculator package initialized with version", Version)
}

// Add returns the sum of two integers
func Add(a, b int) int {
	return a + b
}

// Subtract returns the difference between two integers
func Subtract(a, b int) int {
	return a - b
}

// Multiply returns the product of two integers
func Multiply(a, b int) int {
	return a * b
}

// Divide returns the quotient of two numbers
// Returns an error if dividing by zero
func Divide(a, b int) float64 {
	// We're ignoring the error here for simplicity
	// In a real application, you would handle this error
	result, _ := safeDivide(a, b)
	return result
}

// safeDivide is an unexported function that handles division with error checking
// It's not accessible outside this package
func safeDivide(a, b int) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return float64(a) / float64(b), nil
}

// checkRange is an unexported helper function
func checkRange(value int) bool {
	return value >= -unexportedMaxValue && value <= unexportedMaxValue
}
