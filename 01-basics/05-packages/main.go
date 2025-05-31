// Package main is the entry point for the executable program
package main

import (
	"fmt"

	// Import our custom packages
	"github.com/learn-go/01-basics/05-packages/calculator"
	"github.com/learn-go/01-basics/05-packages/geometry"
)

func main() {
	fmt.Println("=== Go Packages Demo ===\n")

	// Using the calculator package
	fmt.Println("Calculator Package:")
	fmt.Printf("Addition: %d + %d = %d\n", 10, 5, calculator.Add(10, 5))
	fmt.Printf("Subtraction: %d - %d = %d\n", 10, 5, calculator.Subtract(10, 5))
	fmt.Printf("Multiplication: %d * %d = %d\n", 10, 5, calculator.Multiply(10, 5))
	fmt.Printf("Division: %d / %d = %.2f\n", 10, 5, calculator.Divide(10, 5))

	// Using advanced calculator functions
	fmt.Printf("Power: %d^%d = %d\n", 2, 3, calculator.Power(2, 3))
	fmt.Printf("Square Root of %d = %.2f\n", 16, calculator.SquareRoot(16))

	// Using the geometry package
	fmt.Println("\nGeometry Package:")

	// Create a circle
	circle := geometry.NewCircle(5)
	fmt.Printf("Circle with radius %.1f:\n", circle.Radius)
	fmt.Printf("  Area: %.2f\n", circle.Area())
	fmt.Printf("  Perimeter: %.2f\n", circle.Perimeter())

	// Create a rectangle
	rect := geometry.NewRectangle(4, 6)
	fmt.Printf("\nRectangle with width %.1f and height %.1f:\n", rect.Width, rect.Height)
	fmt.Printf("  Area: %.2f\n", rect.Area())
	fmt.Printf("  Perimeter: %.2f\n", rect.Perimeter())

	// Using geometry utilities
	fmt.Println("\nGeometry Utilities:")
	fmt.Printf("Distance between points (1,2) and (4,6): %.2f\n",
		geometry.Distance(1, 2, 4, 6))
	fmt.Printf("Is point (2,3) inside rectangle (1,1,5,5): %t\n",
		geometry.IsPointInRectangle(2, 3, 1, 1, 5, 5))

	// Accessing package variables
	fmt.Println("\nPackage Variables:")
	fmt.Printf("Calculator Version: %s\n", calculator.Version)
	fmt.Printf("Geometry Pi Value: %.5f\n", geometry.Pi)

	// We can't access unexported (private) members
	// This would cause a compilation error:
	// fmt.Println(calculator.multiply(5, 5))
	// fmt.Println(geometry.calculateArea(5, 5))
}
