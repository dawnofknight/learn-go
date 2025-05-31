# Packages and Modules in Go

This directory demonstrates how to organize code with packages and modules in Go.

## Topics Covered

### 1. Package Structure

Go programs are organized into packages. A package is a collection of Go source files in the same directory that are compiled together.

- Every Go file must belong to a package
- The package declaration must be the first line of code in a Go file
- Files in the same directory must have the same package name (with the exception of `*_test.go` files)

### 2. Package Naming Conventions

- Package names should be lowercase, single-word names
- Avoid underscores or mixedCaps
- The package name is the same as the last element of the import path (e.g., `math/rand` has package name `rand`)

### 3. Importing Packages

Go uses the `import` statement to use code from other packages:

```go
import (
    "fmt"           // Standard library package
    "math/rand"     // Nested standard library package
    "example.com/user/project" // External package
)
```

### 4. Exported Names

In Go, a name is exported (visible outside its package) if it begins with a capital letter:

```go
func Hello() string { // Exported function
    return "Hello"
}

func notExported() string { // Not exported
    return "This is not visible outside the package"
}
```

### 5. Go Modules

Go modules are the official dependency management system for Go:

- A module is a collection of related Go packages
- Modules are defined by a `go.mod` file at the root of the project
- The `go.mod` file defines the module path and its dependencies

### 6. Creating a Module

To create a new module:

```bash
go mod init example.com/mymodule
```

This creates a `go.mod` file with the module path and Go version.

### 7. Adding Dependencies

When you import a package that's not in the standard library, Go automatically adds it to your `go.mod` file when you run commands like `go build` or `go test`.

You can also explicitly add dependencies:

```bash
go get github.com/some/dependency
```

## Project Structure

This directory contains a simple project demonstrating package organization:

```
05-packages/
├── go.mod                  # Module definition
├── main.go                 # Main package (entry point)
├── calculator/             # Calculator package
│   ├── calculator.go       # Basic operations
│   └── advanced.go         # Advanced operations
└── geometry/               # Geometry package
    ├── shapes.go           # Shape definitions
    └── utils.go            # Utility functions
```

## Running the Program

To run this program, navigate to this directory in your terminal and execute:

```bash
go run main.go
```

## Key Concepts

- Packages help organize and reuse code
- Exported names (capitalized) are visible outside the package
- Go modules manage dependencies
- The `main` package is special - it defines an executable program
- The `init()` function runs automatically when a package is initialized

## Exercises

1. Create a new module with multiple packages
2. Add a third-party dependency to your module
3. Create a package with both exported and unexported functions
4. Use the `init()` function to initialize a package
5. Create a reusable library package and use it in a main program
