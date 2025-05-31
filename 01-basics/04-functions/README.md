# Functions in Go

This directory demonstrates how to define and use functions in Go, covering basic functions, multiple return values, variadic functions, closures, and more.

## Topics Covered

### 1. Basic Function Declaration

Go functions are declared using the `func` keyword:

```go
func functionName(parameter1 Type1, parameter2 Type2) ReturnType {
    // function body
    return value
}
```

Parameters of the same type can be declared together:

```go
func add(a, b int) int {
    return a + b
}
```

### 2. Multiple Return Values

Go functions can return multiple values:

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

You can also use named return values:

```go
func calculateCircle(radius float64) (area, circumference float64) {
    area = math.Pi * radius * radius
    circumference = 2 * math.Pi * radius
    return // naked return - returns the named return values
}
```

### 3. Variadic Functions

Go supports functions with a variable number of arguments:

```go
func sum(numbers ...int) int {
    total := 0
    for _, num := range numbers {
        total += num
    }
    return total
}
```

### 4. Anonymous Functions and Closures

Go supports anonymous functions and closures:

```go
// Anonymous function
func() {
    fmt.Println("This is an anonymous function")
}()

// Closure
func makeMultiplier(factor int) func(int) int {
    return func(x int) int {
        return x * factor
    }
}
```

### 5. Function as a Type

Functions in Go are first-class citizens and can be used as types:

```go
type MathFunc func(int, int) int

func calculate(a, b int, operation MathFunc) int {
    return operation(a, b)
}
```

### 6. Defer Statement

The `defer` statement delays the execution of a function until the surrounding function returns:

```go
func processFile(filename string) {
    fmt.Println("Opening file")
    defer fmt.Println("Closing file") // Executes when processFile returns
    fmt.Println("Processing file")
}
```

Multiple defers are executed in LIFO (last-in, first-out) order.

### 7. Recursion

Go supports recursive functions (functions that call themselves):

```go
func factorial(n uint) uint {
    if n == 0 {
        return 1
    }
    return n * factorial(n-1)
}
```

### 8. Methods

Go doesn't have classes, but you can define methods on types:

```go
type Rectangle struct {
    Width  float64
    Height float64
}

// Method with a value receiver
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Method with a pointer receiver
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}
```

## Running the Program

To run this program, navigate to this directory in your terminal and execute:

```bash
go run main.go
```

## Key Concepts

- Functions in Go are first-class citizens
- Multiple return values make error handling cleaner
- Defer statements simplify resource management
- Methods provide a way to associate behavior with data
- Closures allow functions to maintain state

## Exercises

1. Write a function that takes a slice of integers and returns the sum and average
2. Create a function that returns another function (a closure)
3. Implement a recursive function to solve a problem (e.g., Tower of Hanoi)
4. Define a struct type and add methods to it
5. Write a function that takes a function as an argument (higher-order function)
