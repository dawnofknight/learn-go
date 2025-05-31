package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

func main() {
	fmt.Println("This file contains functions that will be tested in the test files.")
	fmt.Println("Run 'go test' to execute the tests.")

	// Example usage of the functions
	fmt.Println("\nExample usage:")

	// Calculator functions
	fmt.Printf("Add(5, 3) = %d\n", Add(5, 3))
	fmt.Printf("Subtract(5, 3) = %d\n", Subtract(5, 3))
	fmt.Printf("Multiply(5, 3) = %d\n", Multiply(5, 3))

	result, err := Divide(10, 2)
	fmt.Printf("Divide(10, 2) = %.2f, error: %v\n", result, err)

	result, err = Divide(10, 0)
	fmt.Printf("Divide(10, 0) = %.2f, error: %v\n", result, err)

	// String functions
	fmt.Printf("Reverse(\"hello\") = %s\n", Reverse("hello"))
	fmt.Printf("IsPalindrome(\"radar\") = %t\n", IsPalindrome("radar"))
	fmt.Printf("IsPalindrome(\"hello\") = %t\n", IsPalindrome("hello"))

	// User functions
	user := NewUser("john_doe", "john@example.com", 25)
	fmt.Printf("User: %+v\n", user)
	fmt.Printf("IsAdult: %t\n", user.IsAdult())

	err = user.UpdateEmail("invalid")
	fmt.Printf("UpdateEmail(\"invalid\"): %v\n", err)

	err = user.UpdateEmail("john.doe@example.com")
	fmt.Printf("UpdateEmail(\"john.doe@example.com\"): %v\n", err)
	fmt.Printf("Updated user: %+v\n", user)
}

// Calculator functions

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
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// String functions

// Reverse returns the reversed string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome checks if a string is a palindrome
func IsPalindrome(s string) bool {
	// Convert to lowercase and remove non-alphanumeric characters
	s = strings.ToLower(s)
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		}
	}
	s = result.String()

	// Check if the string is equal to its reverse
	return s == Reverse(s)
}

// User type and methods

// User represents a user in the system
type User struct {
	Username string
	Email    string
	Age      int
}

// NewUser creates a new user
func NewUser(username, email string, age int) *User {
	return &User{
		Username: username,
		Email:    email,
		Age:      age,
	}
}

// IsAdult checks if the user is an adult (18 or older)
func (u *User) IsAdult() bool {
	return u.Age >= 18
}

// UpdateEmail updates the user's email if it's valid
func (u *User) UpdateEmail(email string) error {
	if !strings.Contains(email, "@") {
		return errors.New("invalid email address")
	}
	u.Email = email
	return nil
}

// Geometry types and methods for testing interfaces

// Shape is an interface for geometric shapes
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle represents a rectangle
type Rectangle struct {
	Width  float64
	Height float64
}

// Area returns the area of the rectangle
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter returns the perimeter of the rectangle
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle represents a circle
type Circle struct {
	Radius float64
}

// Area returns the area of the circle
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter returns the circumference of the circle
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Database interface for mocking in tests

// Database is an interface for database operations
type Database interface {
	GetUser(id string) (*User, error)
	SaveUser(user *User) error
}

// UserService provides user-related functionality
type UserService struct {
	DB Database
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id string) (*User, error) {
	user, err := s.DB.GetUser(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// SaveUser saves a user to the database
func (s *UserService) SaveUser(user *User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if user.Username == "" {
		return errors.New("username cannot be empty")
	}

	if err := s.DB.SaveUser(user); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}
