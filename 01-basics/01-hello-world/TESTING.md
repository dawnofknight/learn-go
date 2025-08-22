# Testing in Go

This document explains how to run unit tests for the Hello World program.

## Understanding the Code Structure

The main program has been refactored to make it more testable:

1. `addNumbers(a, b int) int` - A function that adds two integers
2. `getPointerValue(ptr *int) int` - A function that returns the value pointed to by an integer pointer

## Running the Tests

To run the tests, navigate to this directory in your terminal and execute:

```bash
go test
```

For more verbose output, use:

```bash
go test -v
```

To see test coverage:

```bash
go test -cover
```

To generate a coverage report:

```bash
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Cases

The test file `main_test.go` contains the following test cases:

### TestAddNumbers
1. Basic addition (10 + 20 = 30)
2. Addition with zero (5 + 0 = 5)
3. Addition with negative numbers (-10 + 5 = -5)

### TestGetPointerValue
1. Tests that dereferencing a pointer returns the correct value

## Writing More Tests

To add more tests, you can:

1. Add more test cases to existing test functions
2. Create new test functions for new functionality

All test functions must start with `Test` followed by a name that begins with a capital letter.