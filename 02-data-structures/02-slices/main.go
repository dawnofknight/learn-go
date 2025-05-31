// Package main demonstrates slices in Go
package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	// SECTION 1: Slice Declaration and Initialization
	fmt.Println("\n=== Slice Declaration and Initialization ===\n")

	// Method 1: Using slice literal
	names := []string{"Alice", "Bob", "Charlie"}
	fmt.Println("Names slice:", names)

	// Method 2: Using make function
	// make([]T, length, capacity)
	numbers := make([]int, 5) // length=5, capacity=5
	fmt.Println("Numbers slice:", numbers)

	numbersWithCap := make([]int, 3, 10) // length=3, capacity=10
	fmt.Println("Numbers with capacity:", numbersWithCap)
	fmt.Printf("Length: %d, Capacity: %d\n", len(numbersWithCap), cap(numbersWithCap))

	// Method 3: Empty slice (nil slice)
	var empty []int
	fmt.Println("Empty slice:", empty)
	fmt.Println("Is nil?", empty == nil)

	// Method 4: Creating a slice from an array
	array := [5]int{10, 20, 30, 40, 50}
	sliceFromArray := array[1:4] // elements 1, 2, 3 (indices are inclusive:exclusive)
	fmt.Println("Slice from array:", sliceFromArray)

	// SECTION 2: Slice Operations
	fmt.Println("\n=== Slice Operations ===\n")

	// Accessing elements
	fmt.Println("First name:", names[0])
	fmt.Println("Last name:", names[len(names)-1])

	// Modifying elements
	numbers[2] = 42
	fmt.Println("Modified numbers slice:", numbers)

	// Slicing a slice
	fruit := []string{"apple", "banana", "cherry", "date", "elderberry"}
	fmt.Println("Fruit slice:", fruit)

	// Slice syntax: slice[low:high:max]
	// - low: starting index (inclusive)
	// - high: ending index (exclusive)
	// - max: limiting capacity

	// Basic slicing
	fmt.Println("fruit[1:3]:", fruit[1:3]) // [banana cherry]
	fmt.Println("fruit[:2]:", fruit[:2])   // [apple banana]
	fmt.Println("fruit[2:]:", fruit[2:])   // [cherry date elderberry]
	fmt.Println("fruit[:]:", fruit[:])     // [apple banana cherry date elderberry]

	// Slicing with capacity control
	limitedCap := fruit[1:3:4] // elements 1,2 with capacity limited to index 4
	fmt.Println("Limited capacity slice:", limitedCap)
	fmt.Printf("Length: %d, Capacity: %d\n", len(limitedCap), cap(limitedCap))

	// SECTION 3: Appending to Slices
	fmt.Println("\n=== Appending to Slices ===\n")

	// Append elements
	colors := []string{"red", "green"}
	fmt.Println("Original colors:", colors)

	colors = append(colors, "blue")
	fmt.Println("After appending 'blue':", colors)

	// Append multiple elements
	colors = append(colors, "yellow", "purple")
	fmt.Println("After appending multiple colors:", colors)

	// Append one slice to another
	moreColors := []string{"orange", "pink"}
	colors = append(colors, moreColors...)
	fmt.Println("After appending another slice:", colors)

	// Demonstrating capacity growth
	fmt.Println("\nCapacity growth demonstration:")
	s := make([]int, 0)
	fmt.Printf("Initial slice - Length: %d, Capacity: %d\n", len(s), cap(s))

	for i := 0; i < 10; i++ {
		s = append(s, i)
		fmt.Printf("After appending %d - Length: %d, Capacity: %d\n", i, len(s), cap(s))
	}

	// SECTION 4: Copying Slices
	fmt.Println("\n=== Copying Slices ===\n")

	// Using copy function
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, len(src))
	copied := copy(dst, src)

	fmt.Println("Source slice:", src)
	fmt.Println("Destination slice:", dst)
	fmt.Println("Number of elements copied:", copied)

	// Partial copy (destination smaller than source)
	dst2 := make([]int, 3)
	copied = copy(dst2, src)

	fmt.Println("Partial copy destination:", dst2)
	fmt.Println("Number of elements copied:", copied)

	// Partial copy (source smaller than destination)
	dst3 := make([]int, 10)
	copied = copy(dst3, src)

	fmt.Println("Larger destination:", dst3)
	fmt.Println("Number of elements copied:", copied)

	// SECTION 5: Slice Internals and Shared Memory
	fmt.Println("\n=== Slice Internals and Shared Memory ===\n")

	// Demonstrating that slices share underlying array
	original := []int{1, 2, 3, 4, 5}
	fmt.Println("Original slice:", original)

	// Create a slice that shares memory with original
	shared := original[1:4]
	fmt.Println("Shared slice:", shared)

	// Modifying shared affects original
	shared[0] = 99 // This changes original[1]
	fmt.Println("After modifying shared slice:")
	fmt.Println("Shared slice:", shared)
	fmt.Println("Original slice:", original)

	// Creating a completely independent copy
	independent := make([]int, len(original))
	copy(independent, original)

	// Modifying independent doesn't affect original
	independent[0] = 100
	fmt.Println("\nAfter modifying independent copy:")
	fmt.Println("Independent slice:", independent)
	fmt.Println("Original slice:", original)

	// SECTION 6: Common Slice Operations
	fmt.Println("\n=== Common Slice Operations ===\n")

	// Iterating over a slice
	fmt.Println("Iterating over a slice:")
	for i, color := range colors {
		fmt.Printf("colors[%d] = %s\n", i, color)
	}

	// Sorting a slice
	nums := []int{5, 2, 6, 3, 1, 4}
	fmt.Println("\nUnsorted slice:", nums)

	sort.Ints(nums)
	fmt.Println("Sorted slice:", nums)

	// Sorting strings
	unsortedNames := []string{"Zack", "Bob", "Alice", "Eve", "Charlie"}
	fmt.Println("\nUnsorted names:", unsortedNames)

	sort.Strings(unsortedNames)
	fmt.Println("Sorted names:", unsortedNames)

	// Reversing a slice
	fmt.Println("\nReversing a slice:")
	reverse(nums)
	fmt.Println("Reversed nums:", nums)

	// Filtering a slice
	fmt.Println("\nFiltering a slice:")
	even := filter(nums, func(n int) bool {
		return n%2 == 0
	})
	fmt.Println("Even numbers:", even)

	// SECTION 7: Multi-dimensional Slices
	fmt.Println("\n=== Multi-dimensional Slices ===\n")

	// 2D slice
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	fmt.Println("2D slice (matrix):")
	for _, row := range matrix {
		fmt.Println(row)
	}

	// Creating a dynamic 2D slice
	rows, cols := 3, 4
	dynamicMatrix := make([][]int, rows)
	for i := range dynamicMatrix {
		dynamicMatrix[i] = make([]int, cols)
		for j := range dynamicMatrix[i] {
			dynamicMatrix[i][j] = i*cols + j
		}
	}

	fmt.Println("\nDynamic 2D slice:")
	for _, row := range dynamicMatrix {
		fmt.Println(row)
	}

	// SECTION 8: Practical Examples
	fmt.Println("\n=== Practical Examples ===\n")

	// Example 1: Stack implementation using a slice
	fmt.Println("Stack implementation:")
	stack := []string{}

	// Push
	stack = append(stack, "first")
	stack = append(stack, "second")
	stack = append(stack, "third")
	fmt.Println("Stack after pushes:", stack)

	// Pop
	var item string
	item, stack = stack[len(stack)-1], stack[:len(stack)-1]
	fmt.Println("Popped item:", item)
	fmt.Println("Stack after pop:", stack)

	// Example 2: Queue implementation using a slice
	fmt.Println("\nQueue implementation:")
	queue := []string{}

	// Enqueue
	queue = append(queue, "first")
	queue = append(queue, "second")
	queue = append(queue, "third")
	fmt.Println("Queue after enqueues:", queue)

	// Dequeue
	item, queue = queue[0], queue[1:]
	fmt.Println("Dequeued item:", item)
	fmt.Println("Queue after dequeue:", queue)

	// Example 3: String splitting and joining
	fmt.Println("\nString splitting and joining:")
	sentence := "The quick brown fox jumps over the lazy dog"
	words := strings.Split(sentence, " ")
	fmt.Println("Words slice:", words)

	rejoined := strings.Join(words, "-")
	fmt.Println("Rejoined with hyphens:", rejoined)

	// Example 4: Removing elements from a slice
	fmt.Println("\nRemoving elements from a slice:")
	fruits := []string{"apple", "banana", "cherry", "date", "elderberry"}
	fmt.Println("Original fruits:", fruits)

	// Remove element at index 2
	index := 2
	fruits = append(fruits[:index], fruits[index+1:]...)
	fmt.Println("After removing element at index 2:", fruits)
}

// Helper function to reverse a slice
func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Helper function to filter a slice
func filter(s []int, fn func(int) bool) []int {
	var result []int
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
