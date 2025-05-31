# Error Handling in Go

This section covers error handling in Go, which is a fundamental aspect of writing robust and reliable Go programs. Go's approach to error handling is explicit and straightforward, using return values rather than exceptions.

## Topics Covered

### 1. Basic Error Handling

- The `error` interface
- Creating errors with `errors.New()` and `fmt.Errorf()`
- Checking and handling errors
- Multiple return values pattern

### 2. Custom Error Types

- Creating custom error types
- Implementing the `Error()` method
- Type assertions with custom errors
- Providing additional context in errors

### 3. Error Wrapping (Go 1.13+)

- Wrapping errors with `fmt.Errorf()` and `%w` verb
- Unwrapping errors with `errors.Unwrap()`
- Checking error chains with `errors.Is()` and `errors.As()`
- Adding context to errors while preserving the original error

### 4. Sentinel Errors

- Predefined error values
- Comparing errors with `errors.Is()`
- Using sentinel errors for specific error conditions

### 5. Panic and Recover

- Understanding panic
- Using recover to handle panics
- When to use panic vs. returning errors
- Converting panics to errors

### 6. Defer Statements

- Using defer for cleanup operations
- Multiple defer statements (LIFO order)
- Defer with anonymous functions
- Practical examples with file handling and transactions

### 7. Error Handling Patterns

- Type switches for different error types
- Retry mechanisms
- Error handling with cleanup
- Propagating errors up the call stack

## Running the Program

To run the program, navigate to this directory and execute:

```bash
go run main.go
```

## Key Concepts

1. **Errors as Values**: In Go, errors are just values that can be returned, checked, and passed around like any other value.

2. **Explicit Error Checking**: Go encourages explicit error checking with if statements rather than using exceptions.

3. **Error Context**: Good error messages provide context about what went wrong and where it happened.

4. **Error Wrapping**: Wrapping errors allows you to add context while preserving the original error information.

5. **Defer for Cleanup**: The defer statement ensures that cleanup code runs even if errors occur.

6. **Panic for Unrecoverable Errors**: Panic should be used only for truly unrecoverable situations, not for normal error handling.

7. **Error Handling Responsibility**: Each function is responsible for handling errors appropriately - either by returning them to the caller or dealing with them directly.

## Best Practices

1. **Always Check Errors**: Never ignore returned errors without a good reason.

2. **Provide Context**: When returning errors, add context about what operation failed.

3. **Don't Overuse Panic**: Reserve panic for truly exceptional situations.

4. **Use Defer for Cleanup**: Ensure resources are properly released with defer statements.

5. **Be Consistent**: Follow consistent error handling patterns throughout your codebase.

6. **Don't Just Log and Continue**: Either handle the error or return it, but don't just log it and continue as if nothing happened.

7. **Use Sentinel Errors Sparingly**: Only define sentinel errors for truly distinct error conditions that callers need to distinguish.

## Exercises

1. Extend the custom error type to include more fields like a timestamp or error code.

2. Implement a function that reads from multiple files and combines the results, with proper error handling.

3. Create a retry mechanism with exponential backoff for a network operation.

4. Implement a function that converts panics to errors using defer and recover.

5. Write a function that performs multiple operations and returns detailed errors about which operations failed.
