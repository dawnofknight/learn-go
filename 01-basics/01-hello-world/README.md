# Hello World in Go

This directory contains your first Go program - the classic "Hello World" example. This simple program demonstrates the basic structure of a Go application.

## Program Structure

Let's break down the components of a Go program:

### Package Declaration

```go
package main
```

Every Go file starts with a package declaration. The `main` package is special - it tells the Go compiler that this package should compile as an executable program instead of a shared library.

### Import Statements

```go
import "fmt"
```

The import statement brings in external packages. Here, we're importing the `fmt` package, which provides formatting functions for input and output.

### The main Function

```go
func main() {
    // Code goes here
}
```

The `main` function is the entry point of an executable Go program. When you run the program, the code inside this function executes automatically.

## Running the Program

To run this program, navigate to this directory in your terminal and execute:

```bash
go run main.go
```

You should see the output:

```
Hello, World!
Hello, Gopher! Welcome to Go programming.
Go is awesome and fun!
Learning Go:
1. Basics
2. Data Structures
3. Advanced Concepts
```

## Building an Executable

You can also compile the program into an executable file:

```bash
go build
```

This will create an executable file that you can run directly:

```bash
# On Windows
.\01-hello-world.exe

# On macOS/Linux
./01-hello-world
```

## Key Concepts

- **Package Declaration**: Every Go file belongs to a package
- **Import Statement**: How to use external packages
- **main Function**: The entry point for executable programs
- **fmt Package**: Basic functions for formatted I/O

## Exercises

1. Modify the program to print your name instead of "Gopher"
2. Add more `fmt.Println()` statements to print additional messages
3. Try using other functions from the `fmt` package, such as `fmt.Print()` (without a newline)
4. Create a new variable with your age and print it using `fmt.Printf()`
