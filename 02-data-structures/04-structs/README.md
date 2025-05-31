# Structs in Go

This section covers structs, one of Go's most important composite data types that allow you to group together variables of different types under a single name.

## Topics Covered

### 1. Basic Struct Declaration and Initialization

- Defining struct types
- Different ways to initialize structs:
  - Using struct literals with field names
  - Using struct literals with values in field order (not recommended)
  - Creating empty structs and setting fields individually
  - Using the `new()` function (returns a pointer)
  - Using factory functions

### 2. Nested Structs

- Embedding structs (anonymous fields)
- Nested struct fields
- Accessing fields in nested structs

### 3. Struct Methods

- Value receiver methods
- Pointer receiver methods
- When to use each type of receiver

### 4. Struct Tags and JSON

- Using struct tags for metadata
- JSON serialization and deserialization
- Special tag options like `omitempty` and `-`

### 5. Anonymous Structs

- Creating and using one-off struct types
- Use cases for anonymous structs

### 6. Struct Equality

- Comparing structs with the `==` operator
- Limitations of struct comparison

### 7. Struct Factory Functions

- Creating constructor-like functions
- Encapsulating initialization logic

### 8. Practical Examples

- Modeling a shape hierarchy
- Implementing a bank account

## Running the Program

To run the program, navigate to this directory and execute:

```bash
go run main.go
```

## Key Concepts

1. **Structs as Custom Types**: Structs allow you to create custom data types that group related data together, similar to classes in object-oriented languages (but without inheritance).

2. **Value vs. Reference**: Structs are value types in Go. When you assign a struct to a new variable or pass it to a function, Go creates a copy of the struct. To avoid copying, use pointers to structs.

3. **Method Receivers**: Methods can be defined on structs using either value receivers or pointer receivers. Use pointer receivers when you need to modify the struct or avoid copying large structs.

4. **Embedding**: Go supports composition through embedding, which allows you to include one struct type within another. This provides a mechanism for code reuse without traditional inheritance.

5. **Zero Value**: The zero value of a struct is a struct with all fields set to their zero values.

6. **Struct Tags**: Struct tags provide metadata about struct fields that can be accessed through reflection. They're commonly used for encoding/decoding data formats like JSON, XML, etc.

7. **Exported Fields**: Only fields with names starting with an uppercase letter are exported (visible outside the package).

## Exercises

1. Create a `Book` struct with fields for title, author, publication year, and price. Implement methods to apply discounts and print book information.

2. Implement a `Student` struct with fields for name, ID, and grades (as a slice of integers). Add methods to calculate the average grade and determine if the student is passing.

3. Create a struct hierarchy for different types of vehicles (car, truck, motorcycle) with appropriate fields and methods for each.

4. Implement a simple inventory management system using structs for products, categories, and inventory items.

5. Create a struct for representing a blog post with tags. Implement JSON marshaling and unmarshaling with appropriate struct tags.
