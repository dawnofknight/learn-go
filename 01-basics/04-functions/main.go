// Package main demonstrates functions in Go
package main

import (
	"errors"
	"fmt"
	"math"
)

// SECTION 1: Basic Function Declaration

// Simple function with no parameters and no return value
func sayHello() {
	fmt.Println("Hello, Go Functions!")
}

// Function with parameters
func greet(name string) {
	fmt.Println("Hello,", name)
}

// Function with return value
func add(a, b int) int {
	return a + b
}

// SECTION 2: Multiple Return Values

// Function returning multiple values
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Function with named return values
func calculateCircle(radius float64) (area, circumference float64) {
	area = math.Pi * radius * radius
	circumference = 2 * math.Pi * radius
	// Naked return - returns the named return values
	return
}

// SECTION 3: Variadic Functions

// Variadic function (variable number of arguments)
func sum(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

// SECTION 4: Anonymous Functions and Closures

// Function returning a function (closure)
func makeMultiplier(factor int) func(int) int {
	// This anonymous function is a closure because it "closes over" the factor variable
	return func(x int) int {
		return x * factor
	}
}

// SECTION 5: Function as a Type

// Defining a function type
type MathFunc func(int, int) int

// Function that takes a function as an argument
func calculate(a, b int, operation MathFunc) int {
	return operation(a, b)
}

// Functions that match the MathFunc type
func multiply(a, b int) int {
	return a * b
}

func subtract(a, b int) int {
	return a - b
}

// SECTION 6: Defer Statement

// Function demonstrating defer
func processFile(filename string) {
	fmt.Println("Opening file:", filename)

	// Defer executes when the surrounding function returns
	defer fmt.Println("Closing file:", filename)

	// This code runs before the deferred statement
	fmt.Println("Processing file contents...")
}

// Function demonstrating multiple defers (LIFO order)
func deferDemo() {
	fmt.Println("Start")
	defer fmt.Println("First defer")
	defer fmt.Println("Second defer")
	defer fmt.Println("Third defer")
	fmt.Println("End")
}

// SECTION 7: Recursion

// Recursive function to calculate factorial
func factorial(n uint) uint {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}

// Recursive function to calculate Fibonacci numbers
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// SECTION 8: Methods

// Define a struct type
type Rectangle struct {
	Width  float64
	Height float64
}

// Method with a receiver (makes this function a method of Rectangle)
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Method with a pointer receiver (can modify the receiver)
func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

func main() {
	// SECTION 1: Basic Function Calls
	fmt.Println("\n=== Basic Function Calls ===\n")

	sayHello()
	greet("Gopher")
	result := add(5, 3)
	fmt.Println("5 + 3 =", result)

	// SECTION 2: Multiple Return Values
	fmt.Println("\n=== Multiple Return Values ===\n")

	// Using multiple return values
	quotient, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("10 / 2 =", quotient)
	}

	// Error handling
	quotient, err = divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("10 / 0 =", quotient) // This won't execute
	}

	// Named return values
	area, circumference := calculateCircle(5)
	fmt.Printf("Circle with radius 5: Area = %.2f, Circumference = %.2f\n", area, circumference)

	// SECTION 3: Variadic Functions
	fmt.Println("\n=== Variadic Functions ===\n")

	// Calling variadic function with different numbers of arguments
	fmt.Println("Sum of no numbers:", sum())
	fmt.Println("Sum of one number:", sum(5))
	fmt.Println("Sum of three numbers:", sum(1, 2, 3))

	// Passing a slice to a variadic function
	numbers := []int{1, 2, 3, 4, 5}
	fmt.Println("Sum of numbers slice:", sum(numbers...)) // Note the ... to unpack the slice

	// SECTION 4: Anonymous Functions and Closures
	fmt.Println("\n=== Anonymous Functions and Closures ===\n")

	// Anonymous function defined and called immediately
	func() {
		fmt.Println("This is an anonymous function")
	}()

	// Anonymous function with parameters
	func(name string) {
		fmt.Println("Hello,", name)
	}("Anonymous Gopher")

	// Using a closure
	doubler := makeMultiplier(2)
	tripler := makeMultiplier(3)

	fmt.Println("Double 5:", doubler(5))
	fmt.Println("Triple 5:", tripler(5))

	// SECTION 5: Function as a Type
	fmt.Println("\n=== Function as a Type ===\n")

	// Using functions as arguments
	fmt.Println("5 * 3 =", calculate(5, 3, multiply))
	fmt.Println("5 - 3 =", calculate(5, 3, subtract))

	// Using an anonymous function as an argument
	result = calculate(5, 3, func(a, b int) int {
		return a + b
	})
	fmt.Println("5 + 3 =", result)

	// SECTION 6: Defer Statement
	fmt.Println("\n=== Defer Statement ===\n")

	processFile("example.txt")
	fmt.Println()
	deferDemo()

	// SECTION 7: Recursion
	fmt.Println("\n=== Recursion ===\n")

	fmt.Println("Factorial of 5:", factorial(5))

	fmt.Println("Fibonacci sequence:")
	for i := 0; i < 10; i++ {
		fmt.Print(fibonacci(i), " ")
	}
	fmt.Println()

	// SECTION 8: Methods
	fmt.Println("\n=== Methods ===\n")

	rect := Rectangle{Width: 5, Height: 3}
	fmt.Printf("Rectangle: %.1f x %.1f\n", rect.Width, rect.Height)
	fmt.Printf("Area: %.1f\n", rect.Area())

	// Using a pointer receiver method
	rect.Scale(2)
	fmt.Printf("After scaling: %.1f x %.1f\n", rect.Width, rect.Height)
	fmt.Printf("New area: %.1f\n", rect.Area())

	// BONUS: Higher-order function example
	fmt.Println("\n=== Higher-Order Function Example ===\n")

	// Map is a higher-order function that applies a function to each element in a slice
	strings := []string{"hello", "world", "go", "functions"}
	uppercased := mapStrings(strings, strings.ToUpper)
	fmt.Println("Original:", strings)
	fmt.Println("Uppercased:", uppercased)
}

// Higher-order function that applies a function to each string in a slice
func mapStrings(strs []string, f func(string) string) []string {
	result := make([]string, len(strs))
	for i, str := range strs {
		result[i] = f(str)
	}
	return result
}
