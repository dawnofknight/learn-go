# Control Flow in Go

This directory demonstrates how to control the flow of execution in Go programs using conditionals and loops.

## Topics Covered

### 1. If-Else Statements

Go's if statements are similar to those in other languages but don't require parentheses around the condition:

```go
if condition {
    // code to execute if condition is true
} else if anotherCondition {
    // code to execute if anotherCondition is true
} else {
    // code to execute if all conditions are false
}
```

Go also allows a short statement before the condition:

```go
if x := computeValue(); x > 10 {
    // x is available only within the if and else blocks
}
```

### 2. Switch Statements

Go's switch statements are more flexible than in many other languages:

```go
switch expression {
case value1:
    // code
case value2, value3: // Multiple values in one case
    // code
default:
    // code
}
```

Unlike other languages, Go automatically breaks after each case. Use `fallthrough` to execute the next case.

Go also supports switch statements without an expression (acting like if-else chains):

```go
switch {
case condition1:
    // code
case condition2:
    // code
default:
    // code
}
```

Type switches are used to determine the type of an interface value:

```go
switch v := x.(type) {
case nil:
    // x is nil
case int:
    // x is an int
case string:
    // x is a string
default:
    // x is some other type
}
```

### 3. For Loops

Go has only one looping construct: the `for` loop. However, it's flexible and can be used in several ways:

**Basic for loop**:

```go
for i := 0; i < 10; i++ {
    // code
}
```

**For as a while loop**:

```go
for condition {
    // code
}
```

**Infinite loop**:

```go
for {
    // code
    if shouldBreak {
        break
    }
}
```

**For loop with range**:

```go
for index, value := range collection {
    // code
}
```

### 4. Nested Loops and Labels

Go supports nested loops and labeled statements for more complex control flow:

```go
OuterLoop:
for i := 0; i < 5; i++ {
    for j := 0; j < 5; j++ {
        if condition {
            break OuterLoop // Breaks out of the outer loop
        }
    }
}
```

## Running the Program

To run this program, navigate to this directory in your terminal and execute:

```bash
go run main.go
```

## Key Concepts

- Go's syntax is clean and minimalist
- Conditions don't need parentheses
- Switch cases break automatically
- The for loop is versatile and can replace while/do-while loops
- Range is a powerful way to iterate over collections

## Exercises

1. Write a program that uses if-else statements to determine if a number is positive, negative, or zero
2. Create a switch statement that categorizes a person's age group (child, teenager, adult, senior)
3. Write a for loop that prints the first 10 Fibonacci numbers
4. Use nested loops to print a pattern (like a triangle or square)
5. Implement a program that uses break and continue statements in a loop
