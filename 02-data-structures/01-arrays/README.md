# Arrays in Go

This directory demonstrates how to work with arrays in Go. Arrays are fixed-size sequences of elements of the same type.

## Topics Covered

### 1. Array Declaration and Initialization

There are several ways to declare and initialize arrays in Go:

```go
// Method 1: Declaration with size, initialized with zero values
var numbers [5]int

// Method 2: Declaration with initialization
names := [3]string{"Alice", "Bob", "Charlie"}

// Method 3: Let the compiler count the elements
colors := [...]string{"Red", "Green", "Blue", "Yellow"}

// Method 4: Sparse array (initialize specific elements)
var sparse [10]int
sparse[1] = 10
sparse[8] = 80
```

### 2. Accessing and Modifying Array Elements

Array elements are accessed using zero-based indexing:

```go
firstName := names[0]      // Access the first element
names[1] = "Robert"        // Modify the second element
```

### 3. Array Operations

Common operations with arrays include iteration:

```go
// Using a traditional for loop
for i := 0; i < len(names); i++ {
    fmt.Println(names[i])
}

// Using range
for index, value := range names {
    fmt.Printf("names[%d] = %s\n", index, value)
}
```

### 4. Multi-dimensional Arrays

Go supports multi-dimensional arrays:

```go
// 2D array (3x2 matrix)
matrix := [3][2]int{
    {1, 2},
    {3, 4},
    {5, 6},
}

// 3D array
cube := [2][2][2]int{
    {{1, 2}, {3, 4}},
    {{5, 6}, {7, 8}},
}
```

### 5. Array Comparisons

Arrays of the same type can be compared using the `==` operator:

```go
arr1 := [3]int{1, 2, 3}
arr2 := [3]int{1, 2, 3}
arr3 := [3]int{3, 2, 1}

fmt.Println(arr1 == arr2) // true
fmt.Println(arr1 == arr3) // false
```

### 6. Arrays as Function Parameters

Arrays are passed by value (copied) to functions:

```go
func modifyArray(arr [5]int) {
    arr[0] = 100 // This only modifies the copy
}

// To modify the original array, use a pointer
func modifyArrayByPointer(arr *[5]int) {
    (*arr)[0] = 100 // This modifies the original array
}
```

### 7. Array Types and Conversions

Arrays of different sizes are different types in Go:

```go
type Vector3 [3]float64
type Vector4 [4]float64

// Convert array to slice
sliceFromArray := names[:]
```

## Running the Program

To run this program, navigate to this directory in your terminal and execute:

```bash
go run main.go
```

## Key Concepts

- Arrays in Go have a fixed size that is part of their type
- The size of an array is specified at compile time and cannot change
- Arrays are value types (copied when assigned or passed to functions)
- Go provides slices for more flexible, dynamic collections

## Exercises

1. Create an array of 10 integers and fill it with random numbers
2. Write a function that reverses an array in-place
3. Implement a function that finds the minimum and maximum values in an array
4. Create a 2D array representing a tic-tac-toe board and implement a function to check for a winner
5. Write a program that uses an array to count the frequency of each letter in a string
