// Package main demonstrates variables, constants, and data types in Go
package main

import (
	"fmt"
)

// Global constants (package level)
const (
	// Constants can be typed or untyped
	Pi           = 3.14159
	AppName      = "Go Learning App"
	VersionMajor = 1
	VersionMinor = 0
)

// Global variables (package level)
var (
	// Declared but not initialized - will have zero values
	globalCounter int
	globalMessage string
	globalFlag    bool

	// Declared and initialized
	globalRate  = 0.05
	globalAdmin = "admin@example.com"
)

func main() {
	// SECTION 1: Variable Declaration
	fmt.Println("\n=== Variable Declaration ===\n")

	// Method 1: Declare then assign (separate steps)
	var age int
	age = 30
	fmt.Println("Age:", age)

	// Method 2: Declare and initialize in one step
	var name string = "John Doe"
	fmt.Println("Name:", name)

	// Method 3: Type inference (compiler determines type)
	var salary = 75000.50 // float64 inferred
	fmt.Println("Salary:", salary)

	// Method 4: Short declaration (most common inside functions)
	city := "New York"
	fmt.Println("City:", city)

	// Multiple variables at once
	var x, y, z int = 10, 20, 30
	fmt.Println("x, y, z:", x, y, z)

	a, b, c := 5, "hello", true
	fmt.Println("a, b, c:", a, b, c)

	// SECTION 2: Data Types
	fmt.Println("\n=== Data Types ===\n")

	// Numeric types
	var (
		integerNum int        = 42
		floatNum   float64    = 3.14
		complexNum complex128 = 1 + 2i
	)

	fmt.Println("Integer:", integerNum)
	fmt.Println("Float:", floatNum)
	fmt.Println("Complex:", complexNum)

	// Boolean type
	isActive := true
	fmt.Println("Is Active:", isActive)

	// String type
	message := "Hello, Go!"
	fmt.Println("Message:", message)
	fmt.Println("Length:", len(message))
	fmt.Println("First character:", string(message[0]))

	// SECTION 3: Type Conversion
	fmt.Println("\n=== Type Conversion ===\n")

	// Go requires explicit type conversion
	var i int = 42
	var f float64 = float64(i)
	var u uint = uint(f)

	fmt.Println("int to float64 to uint:", i, f, u)

	// String conversion
	var ascii int = 65
	var char = string(rune(ascii)) // Convert to rune first for proper Unicode handling
	fmt.Println("ASCII 65 to string:", char)

	// SECTION 4: Constants
	fmt.Println("\n=== Constants ===\n")

	// Using the global constants
	fmt.Println("Pi:", Pi)
	fmt.Println("App:", AppName)
	fmt.Println("Version:", VersionMajor, ".", VersionMinor)

	// Local constants
	const (
		MaxUsers = 100
		Timeout  = 30 // seconds
	)

	fmt.Println("Max Users:", MaxUsers)
	fmt.Println("Timeout:", Timeout, "seconds")

	// Constant expressions
	const (
		Byte = 1
		KB   = 1024 * Byte
		MB   = 1024 * KB
		GB   = 1024 * MB
	)

	fmt.Println("Kilobyte:", KB, "bytes")
	fmt.Println("Megabyte:", MB, "bytes")
	fmt.Println("Gigabyte:", GB, "bytes")

	// SECTION 5: Zero Values
	fmt.Println("\n=== Zero Values ===\n")

	// Variables declared without initialization get zero values
	var (
		zeroInt    int
		zeroFloat  float64
		zeroString string
		zeroBool   bool
	)

	fmt.Println("Zero int:", zeroInt)       // 0
	fmt.Println("Zero float:", zeroFloat)   // 0.0
	fmt.Println("Zero string:", zeroString) // "" (empty string)
	fmt.Println("Zero bool:", zeroBool)     // false

	// Global variables with zero values
	fmt.Println("Global counter:", globalCounter) // 0
	fmt.Println("Global message:", globalMessage) // "" (empty string)
	fmt.Println("Global flag:", globalFlag)       // false

	// Global variables with initialized values
	fmt.Println("Global rate:", globalRate)   // 0.05
	fmt.Println("Global admin:", globalAdmin) // "admin@example.com"
}
