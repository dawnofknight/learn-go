// Package main demonstrates interfaces in Go
package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// SECTION 1: Basic Interface Declaration and Implementation

// Speaker is a simple interface with a single method
type Speaker interface {
	Speak() string
}

// Dog implements the Speaker interface
type Dog struct {
	Name string
	Age  int
}

// Speak returns the sound a dog makes
func (d Dog) Speak() string {
	return fmt.Sprintf("%s says: Woof!", d.Name)
}

// Cat implements the Speaker interface
type Cat struct {
	Name string
	Age  int
}

// Speak returns the sound a cat makes
func (c Cat) Speak() string {
	return fmt.Sprintf("%s says: Meow!", c.Name)
}

// SECTION 2: Interface Composition

// Mover is an interface for things that can move
type Mover interface {
	Move() string
}

// Sleeper is an interface for things that can sleep
type Sleeper interface {
	Sleep() string
}

// Animal combines multiple interfaces
type Animal interface {
	Speaker // Embedded interface
	Mover
	Sleeper
}

// Implement the additional methods for Dog
func (d Dog) Move() string {
	return fmt.Sprintf("%s runs around", d.Name)
}

func (d Dog) Sleep() string {
	return fmt.Sprintf("%s is sleeping", d.Name)
}

// Implement the additional methods for Cat
func (c Cat) Move() string {
	return fmt.Sprintf("%s prowls silently", c.Name)
}

func (c Cat) Sleep() string {
	return fmt.Sprintf("%s is taking a nap", c.Name)
}

// SECTION 3: Empty Interface

// PrintAny prints any value using the empty interface
func PrintAny(v interface{}) {
	fmt.Printf("Value: %v, Type: %T\n", v, v)
}

// SECTION 4: Type Assertions and Type Switches

// GetSound extracts sound from a Speaker using type assertion
func GetSound(s Speaker) string {
	// Type assertion to check if the Speaker is a Dog
	if dog, ok := s.(Dog); ok {
		return fmt.Sprintf("Dog %s (age %d) says Woof!", dog.Name, dog.Age)
	}

	// Type assertion to check if the Speaker is a Cat
	if cat, ok := s.(Cat); ok {
		return fmt.Sprintf("Cat %s (age %d) says Meow!", cat.Name, cat.Age)
	}

	// Default case
	return s.Speak()
}

// DescribeAnimal uses a type switch to handle different animal types
func DescribeAnimal(a Animal) string {
	switch v := a.(type) {
	case Dog:
		return fmt.Sprintf("This is a dog named %s, age %d", v.Name, v.Age)
	case Cat:
		return fmt.Sprintf("This is a cat named %s, age %d", v.Name, v.Age)
	default:
		return "Unknown animal type"
	}
}

// SECTION 5: Interfaces with Standard Library

// Implementing the sort.Interface for custom sorting
type Person struct {
	Name string
	Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

type ByName []Person

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// SECTION 6: Implementing the Stringer Interface

// Implementing fmt.Stringer for custom string representation
func (p Person) String() string {
	return fmt.Sprintf("%s (%d years)", p.Name, p.Age)
}

// SECTION 7: Interface Values and Nil

// NoiseMaker is an interface for things that make noise
type NoiseMaker interface {
	MakeNoise() string
}

// Robot implements NoiseMaker
type Robot struct {
	Name string
}

func (r *Robot) MakeNoise() string {
	if r == nil {
		return "..."
	}
	return fmt.Sprintf("%s says: Beep Boop!", r.Name)
}

// SECTION 8: Practical Examples

// Example 1: Geometry shapes with interfaces

// Shape is an interface for geometric shapes
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle implements Shape
type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle implements Shape
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Triangle implements Shape
type Triangle struct {
	A, B, C float64 // sides
}

func (t Triangle) Area() float64 {
	// Heron's formula
	s := (t.A + t.B + t.C) / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

// TotalArea calculates the total area of multiple shapes
func TotalArea(shapes ...Shape) float64 {
	var total float64
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}

// Example 2: IO interfaces

// SimpleWriter implements a basic io.Writer-like interface
type SimpleWriter interface {
	Write([]byte) (int, error)
}

// ConsoleWriter writes to console
type ConsoleWriter struct{}

func (cw ConsoleWriter) Write(data []byte) (int, error) {
	n, err := fmt.Print(string(data))
	return n, err
}

// UpperWriter converts to uppercase before writing
type UpperWriter struct {
	Writer SimpleWriter
}

func (uw UpperWriter) Write(data []byte) (int, error) {
	upperData := []byte(strings.ToUpper(string(data)))
	return uw.Writer.Write(upperData)
}

func main() {
	// SECTION 1: Basic Interface Declaration and Implementation
	fmt.Println("=== Basic Interface Declaration and Implementation ===")

	// Create instances that implement the Speaker interface
	dog := Dog{Name: "Rex", Age: 3}
	cat := Cat{Name: "Whiskers", Age: 2}

	// Use the interface
	fmt.Println(dog.Speak())
	fmt.Println(cat.Speak())

	// Store interface values in a slice
	speakers := []Speaker{dog, cat}
	fmt.Println("\nAll speakers:")
	for _, s := range speakers {
		fmt.Println(s.Speak())
	}

	// SECTION 2: Interface Composition
	fmt.Println("\n=== Interface Composition ===")

	// Use the composed Animal interface
	animals := []Animal{dog, cat}
	fmt.Println("\nAll animals:")
	for _, a := range animals {
		fmt.Println(a.Speak())
		fmt.Println(a.Move())
		fmt.Println(a.Sleep())
		fmt.Println() // Empty line for readability
	}

	// SECTION 3: Empty Interface
	fmt.Println("=== Empty Interface ===")

	// The empty interface can hold values of any type
	PrintAny(42)
	PrintAny("Hello, Go!")
	PrintAny(true)
	PrintAny(3.14159)
	PrintAny(dog)
	PrintAny([]string{"apple", "banana", "cherry"})

	// SECTION 4: Type Assertions and Type Switches
	fmt.Println("\n=== Type Assertions and Type Switches ===")

	// Using type assertions
	fmt.Println("\nType assertions:")
	fmt.Println(GetSound(dog))
	fmt.Println(GetSound(cat))

	// Using type switches
	fmt.Println("\nType switches:")
	fmt.Println(DescribeAnimal(dog))
	fmt.Println(DescribeAnimal(cat))

	// SECTION 5: Interfaces with Standard Library
	fmt.Println("\n=== Interfaces with Standard Library ===")

	// Using sort.Interface for custom sorting
	people := []Person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 20},
		{"Diana", 35},
	}

	// Sort by age
	sort.Sort(ByAge(people))
	fmt.Println("\nSorted by age:")
	for _, p := range people {
		fmt.Println(p)
	}

	// Sort by name
	sort.Sort(ByName(people))
	fmt.Println("\nSorted by name:")
	for _, p := range people {
		fmt.Println(p)
	}

	// SECTION 6: Implementing the Stringer Interface
	fmt.Println("\n=== Implementing the Stringer Interface ===")

	// The String() method is automatically used by fmt.Println
	person := Person{"John", 28}
	fmt.Println("Person:", person) // Uses the String() method

	// SECTION 7: Interface Values and Nil
	fmt.Println("\n=== Interface Values and Nil ===")

	// A nil pointer value can still be used with interfaces
	var r1 *Robot = nil
	var noisemaker NoiseMaker = r1

	// The interface value itself is not nil
	fmt.Println("Is interface nil?", noisemaker == nil) // false
	fmt.Println("Noise from nil robot:", noisemaker.MakeNoise())

	// Create a non-nil Robot
	r2 := &Robot{Name: "R2D2"}
	noisemaker = r2
	fmt.Println("Noise from robot:", noisemaker.MakeNoise())

	// SECTION 8: Practical Examples
	fmt.Println("\n=== Practical Examples ===")

	// Example 1: Geometry shapes with interfaces
	fmt.Println("\nGeometry shapes example:")

	rect := Rectangle{Width: 5, Height: 3}
	circ := Circle{Radius: 2}
	tri := Triangle{A: 3, B: 4, C: 5}

	// Use the Shape interface
	shapes := []Shape{rect, circ, tri}

	fmt.Println("Individual shapes:")
	for i, shape := range shapes {
		fmt.Printf("Shape %d - Area: %.2f, Perimeter: %.2f\n",
			i+1, shape.Area(), shape.Perimeter())
	}

	// Calculate total area
	total := TotalArea(rect, circ, tri)
	fmt.Printf("Total area: %.2f\n", total)

	// Example 2: IO interfaces
	fmt.Println("\nIO interfaces example:")

	// Basic writer
	consoleWriter := ConsoleWriter{}
	fmt.Print("Direct output: ")
	consoleWriter.Write([]byte("hello world\n"))

	// Decorator pattern with interfaces
	upperWriter := UpperWriter{Writer: consoleWriter}
	fmt.Print("Uppercase output: ")
	upperWriter.Write([]byte("hello world\n"))
}
