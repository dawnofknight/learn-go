// Package main demonstrates functions in Go
package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
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

// SECTION 8: OOP-Style Programming with Interfaces (Like Java)

// Interface definition - similar to Java interfaces
// Any type that implements these methods satisfies this interface
type Shape interface {
	Area() float64
	Perimeter() float64
	GetName() string
}

// Interface for objects that can be drawn
type Drawable interface {
	Draw()
}

// Interface composition - combining interfaces (like extending in Java)
type DrawableShape interface {
	Shape
	Drawable
}

// Base struct with common fields (like a base class in Java)
type BaseShape struct {
	name  string
	color string
}

// "Constructor" function for BaseShape (Go doesn't have constructors)
func NewBaseShape(name, color string) BaseShape {
	return BaseShape{
		name:  name,
		color: color,
	}
}

// Getter methods (encapsulation like Java)
func (b BaseShape) GetName() string {
	return b.name
}

func (b BaseShape) GetColor() string {
	return b.color
}

// Setter methods with pointer receiver (like Java setters)
func (b *BaseShape) SetName(name string) {
	b.name = name
}

func (b *BaseShape) SetColor(color string) {
	b.color = color
}

// Rectangle struct "extends" BaseShape (composition over inheritance)
type Rectangle struct {
	BaseShape // Embedded struct - similar to inheritance
	Width     float64
	Height    float64
}

// "Constructor" for Rectangle
func NewRectangle(name, color string, width, height float64) *Rectangle {
	return &Rectangle{
		BaseShape: NewBaseShape(name, color),
		Width:     width,
		Height:    height,
	}
}

// Rectangle implements Shape interface
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Rectangle implements Drawable interface
func (r Rectangle) Draw() {
	fmt.Printf("Drawing a %s rectangle (%s) with dimensions %.1fx%.1f\n",
		r.GetColor(), r.GetName(), r.Width, r.Height)
}

// Method with pointer receiver (can modify the receiver)
func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

// Circle struct also "extends" BaseShape
type Circle struct {
	BaseShape
	Radius float64
}

// "Constructor" for Circle
func NewCircle(name, color string, radius float64) *Circle {
	return &Circle{
		BaseShape: NewBaseShape(name, color),
		Radius:    radius,
	}
}

// Circle implements Shape interface
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Circle implements Drawable interface
func (c Circle) Draw() {
	fmt.Printf("Drawing a %s circle (%s) with radius %.1f\n",
		c.GetColor(), c.GetName(), c.Radius)
}

// Method to scale circle
func (c *Circle) Scale(factor float64) {
	c.Radius *= factor
}

// SECTION 9: Polymorphism and Interface Usage (Like Java)

// Function that works with any Shape (polymorphism)
func printShapeInfo(s Shape) {
	fmt.Printf("Shape: %s, Area: %.2f, Perimeter: %.2f\n",
		s.GetName(), s.Area(), s.Perimeter())
}

// Function that works with any Drawable object
func drawShape(d Drawable) {
	d.Draw()
}

// Function demonstrating interface composition
func processDrawableShape(ds DrawableShape) {
	fmt.Printf("Processing %s:\n", ds.GetName())
	ds.Draw()
	fmt.Printf("Area: %.2f, Perimeter: %.2f\n", ds.Area(), ds.Perimeter())
}

// SECTION 10: Factory Pattern (Common in Java)

// Factory function to create shapes based on type
func CreateShape(shapeType, name, color string, dimensions ...float64) Shape {
	switch strings.ToLower(shapeType) {
	case "rectangle":
		if len(dimensions) >= 2 {
			return NewRectangle(name, color, dimensions[0], dimensions[1])
		}
	case "circle":
		if len(dimensions) >= 1 {
			return NewCircle(name, color, dimensions[0])
		}
	}
	return nil
}

// SECTION 11: Collection Processing (Like Java Streams)

// Function to process a collection of shapes
func processShapes(shapes []Shape) {
	totalArea := 0.0
	for _, shape := range shapes {
		totalArea += shape.Area()
		printShapeInfo(shape)
	}
	fmt.Printf("Total area of all shapes: %.2f\n", totalArea)
}

// Higher-order function similar to Java's Stream.map()
func mapShapes(shapes []Shape, mapper func(Shape) float64) []float64 {
	result := make([]float64, len(shapes))
	for i, shape := range shapes {
		result[i] = mapper(shape)
	}
	return result
}

// Filter function similar to Java's Stream.filter()
func filterShapes(shapes []Shape, predicate func(Shape) bool) []Shape {
	var result []Shape
	for _, shape := range shapes {
		if predicate(shape) {
			result = append(result, shape)
		}
	}
	return result
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

	// SECTION 8: OOP-Style Programming with Structs and Methods
	fmt.Println("\n=== OOP-Style Programming ===\n")

	// Creating instances using "constructors"
	rect := NewRectangle("MyRectangle", "blue", 5.0, 3.0)
	circle := NewCircle("MyCircle", "red", 4.0)

	fmt.Printf("Created rectangle: %s (%s) - %.1fx%.1f\n",
		rect.GetName(), rect.GetColor(), rect.Width, rect.Height)
	fmt.Printf("Created circle: %s (%s) - radius %.1f\n",
		circle.GetName(), circle.GetColor(), circle.Radius)

	// SECTION 9: Interface Usage and Polymorphism
	fmt.Println("\n=== Interface Usage and Polymorphism ===\n")

	// Using interfaces for polymorphism
	var shapes []Shape = []Shape{rect, circle}

	fmt.Println("Processing shapes polymorphically:")
	for _, shape := range shapes {
		printShapeInfo(shape)
	}

	// Using Drawable interface
	fmt.Println("\nDrawing shapes:")
	drawShape(rect)
	drawShape(circle)

	// Using composed interface
	fmt.Println("\nUsing composed interface:")
	processDrawableShape(rect)
	fmt.Println()
	processDrawableShape(circle)

	// SECTION 10: Factory Pattern
	fmt.Println("\n=== Factory Pattern ===\n")

	// Creating shapes using factory
	factoryRect := CreateShape("rectangle", "FactoryRect", "green", 6.0, 4.0)
	factoryCircle := CreateShape("circle", "FactoryCircle", "yellow", 3.0)

	if factoryRect != nil {
		printShapeInfo(factoryRect)
	}
	if factoryCircle != nil {
		printShapeInfo(factoryCircle)
	}

	// SECTION 11: Collection Processing (Java Streams-like)
	fmt.Println("\n=== Collection Processing ===\n")

	allShapes := []Shape{rect, circle, factoryRect, factoryCircle}

	fmt.Println("All shapes:")
	processShapes(allShapes)

	// Map operation - get all areas
	areas := mapShapes(allShapes, func(s Shape) float64 {
		return s.Area()
	})
	fmt.Println("\nAreas:", areas)

	// Filter operation - get shapes with area > 20
	largeShapes := filterShapes(allShapes, func(s Shape) bool {
		return s.Area() > 20
	})
	fmt.Println("\nShapes with area > 20:")
	for _, shape := range largeShapes {
		printShapeInfo(shape)
	}

	// SECTION 12: Encapsulation and Getters/Setters
	fmt.Println("\n=== Encapsulation (Getters/Setters) ===\n")

	// Using getters and setters
	fmt.Printf("Rectangle name before: %s\n", rect.GetName())
	rect.SetName("UpdatedRectangle")
	fmt.Printf("Rectangle name after: %s\n", rect.GetName())

	// Scaling objects (modifying through methods)
	fmt.Printf("Rectangle area before scaling: %.2f\n", rect.Area())
	rect.Scale(2.0)
	fmt.Printf("Rectangle area after scaling by 2: %.2f\n", rect.Area())

	// BONUS: Higher-order function example
	fmt.Println("\n=== Higher-Order Function Example ===\n")

	// Map is a higher-order function that applies a function to each element in a slice
	strings := []string{"hello", "world", "go", "functions"}
	//uppercased := mapStrings(strings, append(strings), "!")
	fmt.Println("Original:", strings)
	//fmt.Println("Uppercased:", uppercased)
}

// Higher-order function that applies a function to each string in a slice
func mapStrings(strs []string, f func(string) string) []string {
	result := make([]string, len(strs))
	for i, str := range strs {
		result[i] = f(str)
	}
	return result
}
