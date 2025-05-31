// Package main demonstrates arrays in Go
package main

import (
	"fmt"
)

func main() {
	// SECTION 1: Array Declaration and Initialization
	fmt.Println("\n=== Array Declaration and Initialization ===\n")

	// Method 1: Declaration with size, initialized with zero values
	var numbers [5]int
	fmt.Println("Empty array:", numbers) // [0 0 0 0 0]

	// Method 2: Declaration with initialization
	names := [3]string{"Alice", "Bob", "Charlie"}
	fmt.Println("Names array:", names) // [Alice Bob Charlie]

	// Method 3: Let the compiler count the elements
	colors := [...]string{"Red", "Green", "Blue", "Yellow"}
	fmt.Println("Colors array:", colors)          // [Red Green Blue Yellow]
	fmt.Println("Number of colors:", len(colors)) // 4

	// Method 4: Sparse array (initialize specific elements)
	var sparse [10]int
	sparse[1] = 10
	sparse[8] = 80
	fmt.Println("Sparse array:", sparse) // [0 10 0 0 0 0 0 0 80 0]

	// SECTION 2: Accessing and Modifying Array Elements
	fmt.Println("\n=== Accessing and Modifying Array Elements ===\n")

	// Accessing elements
	fmt.Println("First name:", names[0])              // Alice
	fmt.Println("Last color:", colors[len(colors)-1]) // Yellow

	// Modifying elements
	numbers[2] = 42
	fmt.Println("Modified numbers array:", numbers) // [0 0 42 0 0]

	// SECTION 3: Array Operations
	fmt.Println("\n=== Array Operations ===\n")

	// Iterating over an array with for loop
	fmt.Println("Iterating with traditional for loop:")
	for i := 0; i < len(names); i++ {
		fmt.Printf("names[%d] = %s\n", i, names[i])
	}

	// Iterating with range
	fmt.Println("\nIterating with range:")
	for index, value := range colors {
		fmt.Printf("colors[%d] = %s\n", index, value)
	}

	// If you don't need the index
	fmt.Println("\nIterating values only:")
	for _, value := range names {
		fmt.Println(value)
	}

	// SECTION 4: Multi-dimensional Arrays
	fmt.Println("\n=== Multi-dimensional Arrays ===\n")

	// 2D array (3x2 matrix)
	matrix := [3][2]int{
		{1, 2},
		{3, 4},
		{5, 6},
	}

	fmt.Println("2D array (matrix):")
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			fmt.Printf("%d ", matrix[i][j])
		}
		fmt.Println()
	}

	// 3D array
	cube := [2][2][2]int{
		{{1, 2}, {3, 4}},
		{{5, 6}, {7, 8}},
	}

	fmt.Println("\n3D array (cube):")
	fmt.Println(cube)

	// SECTION 5: Array Comparisons
	fmt.Println("\n=== Array Comparisons ===\n")

	// Arrays of the same type can be compared
	arr1 := [3]int{1, 2, 3}
	arr2 := [3]int{1, 2, 3}
	arr3 := [3]int{3, 2, 1}

	fmt.Println("arr1 == arr2:", arr1 == arr2) // true
	fmt.Println("arr1 == arr3:", arr1 == arr3) // false

	// Arrays of different sizes are different types
	// This would not compile:
	// fmt.Println(arr1 == [4]int{1, 2, 3, 4})

	// SECTION 6: Arrays as Function Parameters
	fmt.Println("\n=== Arrays as Function Parameters ===\n")

	// Arrays are passed by value (copied)
	originalArray := [5]int{1, 2, 3, 4, 5}
	fmt.Println("Original array before function call:", originalArray)

	modifyArray(originalArray)
	fmt.Println("Original array after function call:", originalArray) // Unchanged

	// To modify the original array, pass a pointer
	modifyArrayByPointer(&originalArray)
	fmt.Println("Original array after pointer function call:", originalArray) // Modified

	// SECTION 7: Array Types and Conversions
	fmt.Println("\n=== Array Types and Conversions ===\n")

	// Arrays of different sizes are different types
	type Vector3 [3]float64
	type Vector4 [4]float64

	vec3 := Vector3{1.0, 2.0, 3.0}
	fmt.Println("Vector3:", vec3)

	vec4 := Vector4{1.0, 2.0, 3.0, 4.0}
	fmt.Println("Vector4:", vec4)

	// You can convert between array types and slices
	sliceFromArray := names[:] // Convert array to slice
	fmt.Println("Slice from array:", sliceFromArray)

	// SECTION 8: Practical Examples
	fmt.Println("\n=== Practical Examples ===\n")

	// Example 1: Using an array as a counter
	letterCounts := [26]int{} // One count for each letter (a-z)
	text := "hello world"

	for _, char := range text {
		if char >= 'a' && char <= 'z' {
			letterCounts[char-'a']++
		}
	}

	fmt.Println("Letter counts in 'hello world':")
	for i, count := range letterCounts {
		if count > 0 {
			fmt.Printf("%c: %d\n", 'a'+i, count)
		}
	}

	// Example 2: Fixed-size buffer
	buffer := [10]byte{}
	message := "Go"

	// Copy message to buffer
	for i := 0; i < len(message) && i < len(buffer); i++ {
		buffer[i] = message[i]
	}

	fmt.Printf("\nBuffer with message: %v\n", buffer)
	fmt.Printf("Buffer as string: %s\n", string(buffer[:len(message)]))
}

// Arrays are passed by value (a copy is made)
func modifyArray(arr [5]int) {
	arr[0] = 100 // This only modifies the copy
	fmt.Println("Inside modifyArray function:", arr)
}

// To modify the original array, use a pointer
func modifyArrayByPointer(arr *[5]int) {
	(*arr)[0] = 100 // This modifies the original array
	fmt.Println("Inside modifyArrayByPointer function:", *arr)
}
