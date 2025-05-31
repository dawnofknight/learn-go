# Interfaces in Go

This section covers interfaces, one of Go's most powerful features for abstraction and polymorphism. Interfaces define behavior by declaring a set of methods, but leave the implementation details to the types that implement them.

## Topics Covered

### 1. Basic Interface Declaration and Implementation

- Defining interface types
- Implementing interfaces implicitly (no `implements` keyword)
- Using interface values
- Interface polymorphism

### 2. Interface Composition

- Embedding interfaces within other interfaces
- Creating larger interfaces from smaller ones
- Implementing multiple interfaces

### 3. Empty Interface

- Using `interface{}` to accept any type
- Common use cases for the empty interface

### 4. Type Assertions and Type Switches

- Extracting the underlying concrete type with type assertions
- Handling multiple possible types with type switches
- Safely checking type with the "comma ok" idiom

### 5. Interfaces with Standard Library

- Implementing standard interfaces like `sort.Interface`
- How the standard library uses interfaces for flexibility

### 6. Implementing the Stringer Interface

- Creating custom string representations with `fmt.Stringer`
- How `fmt` package uses the `String()` method

### 7. Interface Values and Nil

- Understanding nil interface values vs. interfaces containing nil pointers
- Safely working with nil values in interfaces

### 8. Practical Examples

- Geometry shapes with interfaces
- IO interfaces and decorator pattern

## Running the Program

To run the program, navigate to this directory and execute:

```bash
go run main.go
```

## Key Concepts

1. **Implicit Implementation**: In Go, types implement interfaces implicitly. There's no `implements` keyword - if a type defines all the methods of an interface, it automatically implements that interface.

2. **Interface Values**: An interface value consists of two components: a concrete type and a value of that type. This is sometimes called a "type-value pair."

3. **Zero Value**: The zero value of an interface is `nil`, which means it holds neither type nor value.

4. **Interface Composition**: Interfaces can be composed of other interfaces, allowing for modular design.

5. **Empty Interface**: The interface type `interface{}` (or just `any` in Go 1.18+) has no methods, so all types implement it. It's used when you need to work with values of unknown type.

6. **Type Assertions**: Type assertions provide access to an interface value's underlying concrete value: `value, ok := interfaceValue.(ConcreteType)`

7. **Type Switches**: Type switches allow you to check for multiple types in a concise way using the `switch v := x.(type)` syntax.

8. **Interface Design**: Good Go interfaces are small, often with just one or two methods. The standard library follows this pattern with interfaces like `io.Reader`, `io.Writer`, and `fmt.Stringer`.

## Exercises

1. Create a `Logger` interface with a `Log(message string)` method. Implement this interface with different types that log to the console, to a file, and to both.

2. Implement a `Serializer` interface with `Marshal()` and `Unmarshal()` methods. Create implementations for JSON, XML, and a custom format.

3. Design a simple plugin system using interfaces, where each plugin implements a common interface but provides different functionality.

4. Create a `Cache` interface and implement it with different caching strategies (e.g., LRU, FIFO).

5. Implement a middleware pattern using interfaces, where each middleware can process a request and pass it to the next middleware in the chain.
