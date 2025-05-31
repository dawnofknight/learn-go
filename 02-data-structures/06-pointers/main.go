// Package main demonstrates pointers in Go
package main

import (
	"fmt"
)

// SECTION 1: Basic Pointer Declaration and Usage

// Basic function to demonstrate pointer declaration and dereferencing
func basicPointers() {
	// Declare a variable
	x := 42
	fmt.Println("Original value of x:", x)

	// Declare a pointer to x
	var p *int = &x
	fmt.Printf("Pointer p points to memory address: %p\n", p)
	fmt.Printf("Value at memory address %p is: %d\n", p, *p)

	// Modify the value through the pointer (dereferencing)
	*p = 100
	fmt.Println("After modifying through pointer, x is now:", x)

	// Alternative pointer declaration with short syntax
	q := &x
	*q = 200
	fmt.Println("After modifying through another pointer, x is now:", x)
}

// SECTION 2: Zero Value and Nil Pointers

func nilPointers() {
	// Zero value of a pointer is nil
	var p *int
	fmt.Println("Zero value of pointer p:", p)
	fmt.Println("Is p nil?", p == nil)

	// Attempting to dereference a nil pointer causes a panic
	fmt.Println("Dereferencing a nil pointer would cause a panic")
	// Uncomment the next line to see the panic:
	// fmt.Println(*p) // This will panic: runtime error: invalid memory address or nil pointer dereference

	// Safe way to use a potentially nil pointer
	if p != nil {
		fmt.Println("Value pointed to by p:", *p)
	} else {
		fmt.Println("Cannot dereference p because it is nil")
	}
}

// SECTION 3: Pointers to Structs

// Person struct for demonstrating pointers to structs
type Person struct {
	Name string
	Age  int
}

func structPointers() {
	// Create a Person struct
	person := Person{Name: "Alice", Age: 30}
	fmt.Printf("Original person: %+v\n", person)

	// Create a pointer to the Person struct
	personPtr := &person

	// Access and modify fields using the pointer
	// Go allows direct access to fields through pointers without explicit dereferencing
	fmt.Println("Name through pointer:", personPtr.Name) // Syntactic sugar for (*personPtr).Name
	personPtr.Age = 31
	fmt.Printf("After modification: %+v\n", person)

	// Create a new Person directly as a pointer
	bobPtr := &Person{Name: "Bob", Age: 25}
	fmt.Printf("New person created as pointer: %+v\n", *bobPtr)

	// Modify through the pointer
	bobPtr.Age = 26
	fmt.Printf("After modification: %+v\n", *bobPtr)
}

// SECTION 4: Pointers as Function Parameters (Pass by Reference)

// Function that takes a value parameter (pass by value)
func incrementByValue(x int) {
	x++
	fmt.Println("Inside incrementByValue, x is:", x)
	// Changes to x are not visible outside this function
}

// Function that takes a pointer parameter (pass by reference)
func incrementByReference(x *int) {
	*x++
	fmt.Println("Inside incrementByReference, *x is:", *x)
	// Changes to *x are visible outside this function
}

// Function to update a Person struct by value
func updatePersonByValue(p Person) {
	p.Name = "Updated " + p.Name
	p.Age++
	fmt.Printf("Inside updatePersonByValue: %+v\n", p)
	// Changes to p are not visible outside this function
}

// Function to update a Person struct by reference
func updatePersonByReference(p *Person) {
	p.Name = "Updated " + p.Name
	p.Age++
	fmt.Printf("Inside updatePersonByReference: %+v\n", *p)
	// Changes to *p are visible outside this function
}

func functionParameters() {
	// Demonstrating pass by value vs pass by reference with primitive types
	x := 10
	fmt.Println("Original x:", x)

	incrementByValue(x)
	fmt.Println("After incrementByValue, x is still:", x) // x is unchanged

	incrementByReference(&x)
	fmt.Println("After incrementByReference, x is now:", x) // x is changed

	// Demonstrating pass by value vs pass by reference with structs
	person := Person{Name: "Charlie", Age: 35}
	fmt.Printf("Original person: %+v\n", person)

	updatePersonByValue(person)
	fmt.Printf("After updatePersonByValue, person is still: %+v\n", person) // person is unchanged

	updatePersonByReference(&person)
	fmt.Printf("After updatePersonByReference, person is now: %+v\n", person) // person is changed
}

// SECTION 5: Pointers to Arrays vs. Slices

func arraysAndSlices() {
	// Arrays are value types in Go
	arr := [3]int{1, 2, 3}
	fmt.Println("Original array:", arr)

	// Passing array by value (creates a copy)
	arrCopy := arr
	arrCopy[0] = 100
	fmt.Println("After modifying copy, original array is still:", arr)
	fmt.Println("Modified copy:", arrCopy)

	// Passing array by reference (using a pointer)
	arrPtr := &arr
	(*arrPtr)[1] = 200 // Need parentheses because [] has higher precedence than *
	fmt.Println("After modifying through pointer, array is now:", arr)

	// Slices are reference types in Go (they contain a pointer to an underlying array)
	slice := []int{1, 2, 3}
	fmt.Println("\nOriginal slice:", slice)

	// Passing slice (no need for explicit pointer as slice already contains a pointer)
	sliceCopy := slice
	sliceCopy[0] = 100
	fmt.Println("After modifying 'copy', original slice is also modified:", slice)
	fmt.Println("Modified 'copy':", sliceCopy)

	// This is because both slice and sliceCopy point to the same underlying array
	// If you want a true copy, you need to create a new slice and copy the elements
	trueSliceCopy := make([]int, len(slice))
	copy(trueSliceCopy, slice)
	trueSliceCopy[1] = 200
	fmt.Println("After modifying true copy, original slice is unchanged:", slice)
	fmt.Println("True copy:", trueSliceCopy)
}

// SECTION 6: Pointers to Pointers

func pointersToPointers() {
	x := 42
	fmt.Println("Original x:", x)

	// Pointer to x
	p := &x
	fmt.Printf("p points to x. Address: %p, Value: %d\n", p, *p)

	// Pointer to pointer to x
	pp := &p
	fmt.Printf("pp points to p. Address: %p, Value: %p\n", pp, *pp)
	fmt.Printf("Dereferencing pp twice: %d\n", **pp)

	// Modify x through pp
	**pp = 100
	fmt.Println("After modifying through pp, x is now:", x)
}

// SECTION 7: Common Pointer Patterns and Idioms

// NewPerson is a constructor function that returns a pointer to a new Person
func NewPerson(name string, age int) *Person {
	// It's safe to return a pointer to a local variable
	// Go will automatically allocate it on the heap instead of the stack
	return &Person{
		Name: name,
		Age:  age,
	}
}

// Optional parameters pattern using pointers
type Config struct {
	Host     string
	Port     int
	Debug    bool
	MaxConns *int // Optional parameter
	Timeout  *int // Optional parameter
}

func createConfig() {
	// Create a config with only required parameters
	config1 := Config{
		Host:  "localhost",
		Port:  8080,
		Debug: true,
		// MaxConns and Timeout are nil (not specified)
	}

	// Create a config with all parameters
	maxConns := 10
	timeout := 30
	config2 := Config{
		Host:     "localhost",
		Port:     8080,
		Debug:    true,
		MaxConns: &maxConns,
		Timeout:  &timeout,
	}

	// Function to use the config
	printConfig := func(c Config) {
		fmt.Printf("Host: %s, Port: %d, Debug: %t\n", c.Host, c.Port, c.Debug)

		// Safely access optional parameters
		if c.MaxConns != nil {
			fmt.Printf("MaxConns: %d\n", *c.MaxConns)
		} else {
			fmt.Println("MaxConns: not specified (using default)")
		}

		if c.Timeout != nil {
			fmt.Printf("Timeout: %d\n", *c.Timeout)
		} else {
			fmt.Println("Timeout: not specified (using default)")
		}
	}

	fmt.Println("Config 1:")
	printConfig(config1)

	fmt.Println("\nConfig 2:")
	printConfig(config2)
}

// SECTION 8: Practical Examples

// Example 1: Implementing a simple linked list with pointers
type Node struct {
	Value int
	Next  *Node
}

func linkedListExample() {
	// Create a linked list: 1 -> 2 -> 3 -> nil
	node3 := &Node{Value: 3, Next: nil}
	node2 := &Node{Value: 2, Next: node3}
	node1 := &Node{Value: 1, Next: node2}

	// Print the linked list
	fmt.Println("Linked list:")
	current := node1
	for current != nil {
		fmt.Printf("%d -> ", current.Value)
		current = current.Next
	}
	fmt.Println("nil")

	// Insert a new node between 2 and 3
	newNode := &Node{Value: 25, Next: node2.Next}
	node2.Next = newNode

	// Print the modified linked list
	fmt.Println("Modified linked list:")
	current = node1
	for current != nil {
		fmt.Printf("%d -> ", current.Value)
		current = current.Next
	}
	fmt.Println("nil")
}

// Example 2: Implementing a simple binary tree with pointers
type TreeNode struct {
	Value int
	Left  *TreeNode
	Right *TreeNode
}

func binaryTreeExample() {
	// Create a binary tree:
	//       10
	//      /  \
	//     5    15
	//    / \     \
	//   3   7     20

	root := &TreeNode{Value: 10}
	root.Left = &TreeNode{Value: 5}
	root.Right = &TreeNode{Value: 15}
	root.Left.Left = &TreeNode{Value: 3}
	root.Left.Right = &TreeNode{Value: 7}
	root.Right.Right = &TreeNode{Value: 20}

	// Print the tree using in-order traversal
	fmt.Println("Binary tree (in-order traversal):")
	inOrderTraversal(root)
	fmt.Println()
}

// In-order traversal: left -> root -> right
func inOrderTraversal(node *TreeNode) {
	if node == nil {
		return
	}

	inOrderTraversal(node.Left)
	fmt.Printf("%d ", node.Value)
	inOrderTraversal(node.Right)
}

func main() {
	// SECTION 1: Basic Pointer Declaration and Usage
	fmt.Println("=== Basic Pointer Declaration and Usage ===")
	basicPointers()

	// SECTION 2: Zero Value and Nil Pointers
	fmt.Println("\n=== Zero Value and Nil Pointers ===")
	nilPointers()

	// SECTION 3: Pointers to Structs
	fmt.Println("\n=== Pointers to Structs ===")
	structPointers()

	// SECTION 4: Pointers as Function Parameters
	fmt.Println("\n=== Pointers as Function Parameters ===")
	functionParameters()

	// SECTION 5: Pointers to Arrays vs. Slices
	fmt.Println("\n=== Pointers to Arrays vs. Slices ===")
	arraysAndSlices()

	// SECTION 6: Pointers to Pointers
	fmt.Println("\n=== Pointers to Pointers ===")
	pointersToPointers()

	// SECTION 7: Common Pointer Patterns and Idioms
	fmt.Println("\n=== Common Pointer Patterns and Idioms ===")
	// Constructor pattern
	person := NewPerson("Dave", 40)
	fmt.Printf("Person created with constructor: %+v\n", *person)

	// Optional parameters pattern
	fmt.Println("\nOptional parameters pattern:")
	createConfig()

	// SECTION 8: Practical Examples
	fmt.Println("\n=== Practical Examples ===")

	// Example 1: Linked List
	fmt.Println("\nLinked List Example:")
	linkedListExample()

	// Example 2: Binary Tree
	fmt.Println("\nBinary Tree Example:")
	binaryTreeExample()
}
