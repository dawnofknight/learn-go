// Package main demonstrates maps in Go
package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	// SECTION 1: Map Declaration and Initialization
	fmt.Println("\n=== Map Declaration and Initialization ===\n")

	// Method 1: Using map literal
	colors := map[string]string{
		"red":   "#FF0000",
		"green": "#00FF00",
		"blue":  "#0000FF",
	}
	fmt.Println("Colors map:", colors)

	// Method 2: Using make function
	scores := make(map[string]int)
	fmt.Println("Empty scores map:", scores)

	// Method 3: Empty map (nil map)
	var empty map[string]int
	fmt.Println("Nil map:", empty)
	fmt.Println("Is nil?", empty == nil)

	// SECTION 2: Map Operations
	fmt.Println("\n=== Map Operations ===\n")

	// Adding elements
	scores["Alice"] = 95
	scores["Bob"] = 80
	scores["Charlie"] = 90
	fmt.Println("Scores after adding elements:", scores)

	// Accessing elements
	fmt.Println("Bob's score:", scores["Bob"])

	// Accessing a non-existent key returns the zero value
	fmt.Println("Dave's score (non-existent):", scores["Dave"])

	// Checking if a key exists
	score, exists := scores["Eve"]
	if exists {
		fmt.Println("Eve's score:", score)
	} else {
		fmt.Println("Eve is not in the map")
	}

	// Updating elements
	scores["Bob"] = 85
	fmt.Println("Scores after updating Bob's score:", scores)

	// Deleting elements
	delete(scores, "Charlie")
	fmt.Println("Scores after deleting Charlie:", scores)

	// SECTION 3: Iterating Over Maps
	fmt.Println("\n=== Iterating Over Maps ===\n")

	// Using range
	fmt.Println("Iterating over colors map:")
	for key, value := range colors {
		fmt.Printf("%s -> %s\n", key, value)
	}

	// Iterating in a specific order (maps are unordered by default)
	fmt.Println("\nIterating over scores map in sorted order:")

	// Get all keys
	var names []string
	for name := range scores {
		names = append(names, name)
	}

	// Sort the keys
	sort.Strings(names)

	// Iterate in sorted order
	for _, name := range names {
		fmt.Printf("%s: %d\n", name, scores[name])
	}

	// SECTION 4: Map of Maps
	fmt.Println("\n=== Map of Maps ===\n")

	// Creating a nested map
	users := map[string]map[string]string{
		"alice": {
			"email":   "alice@example.com",
			"phone":   "555-1234",
			"address": "123 Main St",
		},
		"bob": {
			"email":   "bob@example.com",
			"phone":   "555-5678",
			"address": "456 Oak Ave",
		},
	}

	fmt.Println("Users map:")
	for user, details := range users {
		fmt.Printf("%s:\n", user)
		for key, value := range details {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	// Adding a new nested map
	users["charlie"] = map[string]string{
		"email":   "charlie@example.com",
		"phone":   "555-9012",
		"address": "789 Pine Rd",
	}

	// Accessing nested map values
	fmt.Println("\nBob's email:", users["bob"]["email"])

	// SECTION 5: Maps as Sets
	fmt.Println("\n=== Maps as Sets ===\n")

	// Using a map as a set (only care about keys, not values)
	fruitSet := map[string]bool{
		"apple":  true,
		"banana": true,
		"cherry": true,
	}

	// Check if an element is in the set
	fmt.Println("Is 'apple' in the set?", fruitSet["apple"])

	// Add to the set
	fruitSet["date"] = true

	// Remove from the set
	delete(fruitSet, "banana")

	// List all elements in the set
	fmt.Println("Fruits in the set:")
	for fruit := range fruitSet {
		fmt.Println(fruit)
	}

	// A more memory-efficient set using empty struct
	type void struct{}
	var member void

	vegetableSet := make(map[string]void)
	vegetableSet["carrot"] = member
	vegetableSet["broccoli"] = member
	vegetableSet["spinach"] = member

	fmt.Println("\nVegetables in the set:")
	for vegetable := range vegetableSet {
		fmt.Println(vegetable)
	}

	// SECTION 6: Map Capacity and Performance
	fmt.Println("\n=== Map Capacity and Performance ===\n")

	// Creating a map with initial capacity
	citiesWithCapacity := make(map[string]string, 10)
	fmt.Printf("Created map with initial capacity for %d elements\n", 10)

	// Adding elements
	citiesWithCapacity["USA"] = "Washington D.C."
	citiesWithCapacity["UK"] = "London"
	citiesWithCapacity["France"] = "Paris"

	fmt.Printf("Map now has %d elements\n", len(citiesWithCapacity))

	// SECTION 7: Practical Examples
	fmt.Println("\n=== Practical Examples ===\n")

	// Example 1: Word frequency counter
	fmt.Println("Word frequency counter:")
	text := "the quick brown fox jumps over the lazy dog the fox was quick"
	words := strings.Fields(text) // Split by whitespace

	frequency := make(map[string]int)
	for _, word := range words {
		frequency[word]++
	}

	for word, count := range frequency {
		fmt.Printf("%s: %d\n", word, count)
	}

	// Example 2: Grouping data
	fmt.Println("\nGrouping students by grade:")
	studentScores := map[string]int{
		"Alice":   92,
		"Bob":     78,
		"Charlie": 85,
		"Dave":    95,
		"Eve":     88,
		"Frank":   72,
	}

	// Group students by grade (A, B, C, D, F)
	gradeGroups := map[string][]string{
		"A": {},
		"B": {},
		"C": {},
		"D": {},
		"F": {},
	}

	for student, score := range studentScores {
		switch {
		case score >= 90:
			gradeGroups["A"] = append(gradeGroups["A"], student)
		case score >= 80:
			gradeGroups["B"] = append(gradeGroups["B"], student)
		case score >= 70:
			gradeGroups["C"] = append(gradeGroups["C"], student)
		case score >= 60:
			gradeGroups["D"] = append(gradeGroups["D"], student)
		default:
			gradeGroups["F"] = append(gradeGroups["F"], student)
		}
	}

	for grade, students := range gradeGroups {
		if len(students) > 0 {
			fmt.Printf("Grade %s: %v\n", grade, students)
		}
	}

	// Example 3: Simple cache
	fmt.Println("\nSimple cache implementation:")

	// Create a function that's expensive to compute
	expensiveOperation := func(key string) string {
		fmt.Printf("Computing expensive result for %s\n", key)
		return "Result for " + key
	}

	// Create a cache
	cache := make(map[string]string)

	// Function that uses the cache
	getResult := func(key string) string {
		// Check if result is in cache
		if result, found := cache[key]; found {
			fmt.Printf("Cache hit for %s\n", key)
			return result
		}

		// Not in cache, compute it
		result := expensiveOperation(key)

		// Store in cache for next time
		cache[key] = result
		return result
	}

	// First call - should compute
	fmt.Println(getResult("A"))

	// Second call - should use cache
	fmt.Println(getResult("A"))

	// Different key - should compute
	fmt.Println(getResult("B"))
}
