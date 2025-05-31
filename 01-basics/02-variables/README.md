# Variables, Constants, and Data Types in Go

This directory demonstrates how to work with variables, constants, and different data types in Go.

## Topics Covered

### Variable Declaration

Go provides several ways to declare variables:

1. **Declare then assign** (separate steps):

   ```go
   var age int
   age = 30
   ```

2. **Declare and initialize** in one step:

   ```go
   var name string = "John Doe"
   ```

3. **Type inference** (compiler determines type):

   ```go
   var salary = 75000.50 // float64 inferred
   ```

4. **Short declaration** (most common inside functions):
   ```go
   city := "New York"
   ```

### Data Types

Go has several built-in data types:

- **Numeric types**:

  - Integer: `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`
  - Float: `float32`, `float64`
  - Complex: `complex64`, `complex128`

- **Boolean type**: `bool` (true or false)

- **String type**: `string` (UTF-8 encoded)

### Type Conversion

Go requires explicit type conversion between different types:

```go
var i int = 42
var f float64 = float64(i) // Explicit conversion required
```

### Constants

Constants are values that cannot be changed during program execution:

```go
const Pi = 3.14159
const AppName = "Go Learning App"
```

### Zero Values

Variables declared without initialization get zero values:

- Numeric types: `0`
- Boolean type: `false`
- String type: `""` (empty string)
- Reference types: `nil`

## Running the Program

To run this program, navigate to this directory in your terminal and execute:

```bash
go run main.go
```

## Key Concepts

- Go is a statically typed language
- Variables must be declared before use
- Type inference allows for cleaner code
- Constants cannot be changed after declaration
- Variables have zero values when not initialized

## Exercises

1. Create variables of different types and print their values
2. Try to assign a value of one type to a variable of another type without conversion
3. Create a constant and try to change its value (notice the compiler error)
4. Declare multiple variables in a single statement
5. Use the short declaration operator to create and initialize variables
