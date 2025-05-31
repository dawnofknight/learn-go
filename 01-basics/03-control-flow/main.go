// Package main demonstrates control flow in Go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// SECTION 1: If-Else Statements
	fmt.Println("\n=== If-Else Statements ===\n")

	// Basic if statement
	age := 19
	if age >= 18 {
		fmt.Println("You are an adult")
	} else {
		fmt.Println("You are a minor")
	}

	// If with a short statement
	if score := rand.Intn(100); score >= 70 {
		fmt.Println("Score:", score, "- You passed!")
	} else if score >= 50 {
		fmt.Println("Score:", score, "- You need improvement")
	} else {
		fmt.Println("Score:", score, "- You failed")
	}

	// Multiple conditions
	temperature := 25
	humidity := 60

	if temperature > 30 && humidity > 80 {
		fmt.Println("It's hot and humid")
	} else if temperature > 30 || humidity > 80 {
		fmt.Println("It's either hot or humid")
	} else {
		fmt.Println("Weather is pleasant")
	}

	// SECTION 2: Switch Statements
	fmt.Println("\n=== Switch Statements ===\n")

	// Basic switch
	day := time.Now().Weekday()
	switch day {
	case time.Saturday, time.Sunday:
		fmt.Println("It's the weekend!")
	default:
		fmt.Println("It's a weekday")
	}

	// Switch with no expression (like if-else chain)
	hour := time.Now().Hour()
	switch {
	case hour < 12:
		fmt.Println("Good morning!")
	case hour < 17:
		fmt.Println("Good afternoon!")
	default:
		fmt.Println("Good evening!")
	}

	// Switch with fallthrough
	n := rand.Intn(10)
	fmt.Println("Number:", n)
	switch n {
	case 0:
		fmt.Println("n is zero")
	case 1:
		fmt.Println("n is one")
		fallthrough // Executes the next case regardless of its condition
	case 2:
		fmt.Println("n is <= two")
	case 3, 4, 5:
		fmt.Println("n is between 3 and 5")
	default:
		fmt.Println("n is greater than 5")
	}

	// Type switch
	var x interface{} = "hello"
	switch v := x.(type) {
	case nil:
		fmt.Println("x is nil")
	case int:
		fmt.Println("x is an int:", v)
	case string:
		fmt.Println("x is a string:", v)
	case bool:
		fmt.Println("x is a bool:", v)
	default:
		fmt.Printf("x is some other type: %T\n", v)
	}

	// SECTION 3: For Loops
	fmt.Println("\n=== For Loops ===\n")

	// Basic for loop (like C/Java)
	fmt.Println("Basic for loop:")
	for i := 0; i < 5; i++ {
		fmt.Print(i, " ")
	}
	fmt.Println()

	// For as a while loop
	fmt.Println("\nFor as a while loop:")
	counter := 0
	for counter < 5 {
		fmt.Print(counter, " ")
		counter++
	}
	fmt.Println()

	// Infinite loop with break
	fmt.Println("\nInfinite loop with break:")
	sum := 0
	for {
		sum++
		if sum > 5 {
			break
		}
		fmt.Print(sum, " ")
	}
	fmt.Println()

	// For loop with continue
	fmt.Println("\nFor loop with continue (skip even numbers):")
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue // Skip even numbers
		}
		fmt.Print(i, " ")
	}
	fmt.Println()

	// For loop with range (for arrays/slices)
	fmt.Println("\nFor loop with range (array):")
	numbers := [5]int{10, 20, 30, 40, 50}
	for index, value := range numbers {
		fmt.Printf("Index: %d, Value: %d\n", index, value)
	}

	// For loop with range (for maps)
	fmt.Println("\nFor loop with range (map):")
	capitals := map[string]string{
		"USA":    "Washington D.C.",
		"France": "Paris",
		"Japan":  "Tokyo",
	}
	for country, capital := range capitals {
		fmt.Printf("The capital of %s is %s\n", country, capital)
	}

	// For loop with range (for strings)
	fmt.Println("\nFor loop with range (string):")
	word := "Go语言" // A string with both ASCII and Unicode characters
	for i, char := range word {
		fmt.Printf("Character '%c' starts at byte position %d\n", char, i)
	}

	// SECTION 4: Nested Loops and Labels
	fmt.Println("\n=== Nested Loops and Labels ===\n")

	// Nested loops
	fmt.Println("Multiplication table (3x3):")
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			fmt.Printf("%d x %d = %d\t", i, j, i*j)
		}
		fmt.Println() // New line after each row
	}

	// Using labels with break
	fmt.Println("\nUsing labels with break:")
OuterLoop:
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if i*j > 5 {
				fmt.Printf("Breaking outer loop at i=%d, j=%d\n", i, j)
				break OuterLoop // Break out of the outer loop
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
		fmt.Println()
	}
}
