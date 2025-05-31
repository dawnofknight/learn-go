# Slices in Go

This directory demonstrates how to work with slices in Go. Slices are flexible, dynamic views into arrays that provide more power and convenience than arrays.

## Topics Covered

### 1. Slice Declaration and Initialization

There are several ways to create slices in Go:

```go
// Method 1: Using slice literal
names := []string{"Alice", "Bob", "Charlie"}

// Method 2: Using make function
numbers := make([]int, 5)      // length=5, capacity=5
numbersWithCap := make([]int, 3, 10) // length=3, capacity=10

// Method 3: Empty slice (nil slice)
var empty []int

// Method 4: Creating a slice from an array
array := [5]int{10, 20, 30, 40, 50}
sliceFromArray := array[1:4] // elements 1, 2, 3 (indices are inclusive:exclusive)
```

### 2. Slice Operations

Slices support various operations including accessing, modifying, and re-slicing:

```go
// Accessing elements
firstName := names[0]

// Modifying elements
numbers[2] = 42

// Slicing a slice
fruit := []string{"apple", "banana", "cherry", "date", "elderberry"}
fmt.Println(fruit[1:3])  // [banana cherry]
fmt.Println(fruit[:2])    // [apple banana]
fmt.Println(fruit[2:])    // [cherry date elderberry]

// Slicing with capacity control
limitedCap := fruit[1:3:4] // elements 1,2 with capacity limited to index 4
```

### 3. Appending to Slices

The `append` function adds elements to a slice and handles capacity growth:

```go
colors := []string{"red", "green"}
colors = append(colors, "blue")  // Add one element
colors = append(colors, "yellow", "purple")  // Add multiple elements

// Append one slice to another
moreColors := []string{"orange", "pink"}
colors = append(colors, moreColors...)  // Note the ... operator
```

### 4. Copying Slices

The `copy` function creates independent copies of slices:

```go
src := []int{1, 2, 3, 4, 5}
dst := make([]int, len(src))
copied := copy(dst, src)  // Returns number of elements copied
```

### 5. Slice Internals and Shared Memory

Slices are references to underlying arrays, which means they can share memory:

```go
original := []int{1, 2, 3, 4, 5}
shared := original[1:4]  // References the same underlying array

shared[0] = 99  // This also changes original[1]

// To create an independent copy:
independent := make([]int, len(original))
copy(independent, original)
```

### 6. Common Slice Operations

Common operations include iteration, sorting, and filtering:

```go
// Iterating
for i, value := range slice {
    // Use i and value
}

// Sorting
sort.Ints(nums)  // Sort integers
sort.Strings(names)  // Sort strings

// Filtering
even := filter(nums, func(n int) bool {
    return n%2 == 0
})
```

### 7. Multi-dimensional Slices

Go supports slices of slices for multi-dimensional data:

```go
// 2D slice
matrix := [][]int{
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9},
}

// Creating a dynamic 2D slice
rows, cols := 3, 4
dynamicMatrix := make([][]int, rows)
for i := range dynamicMatrix {
    dynamicMatrix[i] = make([]int, cols)
}
```

### 8. Practical Examples

The code demonstrates practical uses of slices:

- Stack implementation (LIFO)
- Queue implementation (FIFO)
- String splitting and joining
- Removing elements from a slice

## Running the Program

To run this program, navigate to this directory in your terminal and execute:

```bash
go run main.go
```

## Key Concepts

- Slices are references to arrays, not values
- A slice consists of a pointer to the array, a length, and a capacity
- The `append` function handles dynamic resizing
- Slices can share underlying memory, which can be both powerful and dangerous
- Use `copy` to create independent slices

## Exercises

1. Implement a function that removes duplicates from a slice
2. Create a function that merges two sorted slices into a single sorted slice
3. Implement a circular buffer using a slice
4. Write a function that rotates a slice by n positions
5. Create a function that finds the intersection of two slices
