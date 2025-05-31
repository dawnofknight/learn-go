// Package main demonstrates structs in Go
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// SECTION 1: Basic Struct Declaration

// Person is a basic struct with various field types
type Person struct {
	FirstName string
	LastName  string
	Age       int
	Height    float64 // in meters
	Active    bool
}

// SECTION 2: Nested Structs

// Address represents a physical address
type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
	Country string
}

// Employee extends Person with additional fields including a nested struct
type Employee struct {
	Person        // Embedded struct (anonymous field)
	EmployeeID    string
	Position      string
	Salary        float64
	HomeAddress   Address    // Nested struct
	OfficeAddress Address    // Another nested struct
	HireDate      time.Time
	Manager       *Employee // Pointer to another Employee
}

// SECTION 3: Struct Methods

// FullName returns the person's full name
func (p Person) FullName() string {
	return p.FirstName + " " + p.LastName
}

// Birthday increments the person's age
// Using pointer receiver to modify the struct
func (p *Person) Birthday() {
	p.Age++
}

// PromotionEligible checks if employee is eligible for promotion
func (e Employee) PromotionEligible() bool {
	// Example criteria: at least 2 years of service and age over 25
	yearsOfService := time.Since(e.HireDate).Hours() / (24 * 365)
	return yearsOfService >= 2 && e.Age >= 25
}

// GiveRaise increases the employee's salary by the given percentage
func (e *Employee) GiveRaise(percentage float64) {
	e.Salary = e.Salary * (1 + percentage/100)
}

// SECTION 4: Struct Tags for JSON

// Product demonstrates struct tags for JSON serialization
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"` // omit if empty
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"-"` // ignore this field
}

// SECTION 5: Anonymous Structs

func createAnonymousStruct() {
	// Anonymous struct declaration and initialization
	point := struct {
		X int
		Y int
	}{
		X: 10,
		Y: 20,
	}

	fmt.Printf("Point coordinates: (%d, %d)\n", point.X, point.Y)
}

// SECTION 6: Struct Equality

func compareStructs() {
	person1 := Person{"John", "Doe", 30, 1.75, true}
	person2 := Person{"John", "Doe", 30, 1.75, true}
	person3 := Person{"Jane", "Doe", 28, 1.68, true}

	// Direct comparison with == operator
	fmt.Println("person1 == person2:", person1 == person2) // true
	fmt.Println("person1 == person3:", person1 == person3) // false

	// Note: Structs containing fields that can't be compared with == 
	// (like slices, maps, or functions) can't be compared with ==
}

// SECTION 7: Struct Factory Functions

// NewPerson is a factory function for creating Person instances
func NewPerson(firstName, lastName string, age int) Person {
	return Person{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
		Active:    true, // Default value
	}
}

// NewEmployee is a factory function for creating Employee instances
func NewEmployee(firstName, lastName string, age int, employeeID, position string, salary float64) *Employee {
	return &Employee{
		Person: Person{
			FirstName: firstName,
			LastName:  lastName,
			Age:       age,
			Active:    true,
		},
		EmployeeID: employeeID,
		Position:   position,
		Salary:     salary,
		HireDate:   time.Now(),
	}
}

func main() {
	// SECTION 1: Basic Struct Declaration and Initialization
	fmt.Println("=== Basic Struct Declaration and Initialization ===")

	// Method 1: Struct literal with field names
	person1 := Person{
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
		Height:    1.75,
		Active:    true,
	}
	fmt.Printf("Person 1: %+v\n", person1)

	// Method 2: Struct literal with values in field order (not recommended)
	person2 := Person{"Jane", "Smith", 28, 1.68, true}
	fmt.Printf("Person 2: %+v\n", person2)

	// Method 3: Create empty struct and set fields individually
	var person3 Person
	person3.FirstName = "Bob"
	person3.LastName = "Johnson"
	person3.Age = 35
	person3.Height = 1.82
	person3.Active = false
	fmt.Printf("Person 3: %+v\n", person3)

	// Method 4: Using new() function (returns a pointer)
	person4 := new(Person)
	person4.FirstName = "Alice"
	person4.LastName = "Williams"
	person4.Age = 25
	fmt.Printf("Person 4 (pointer): %+v\n", *person4)

	// Method 5: Using factory function
	person5 := NewPerson("Michael", "Brown", 40)
	fmt.Printf("Person 5 (factory): %+v\n", person5)

	// SECTION 2: Nested Structs
	fmt.Println("\n=== Nested Structs ===")

	// Create an employee with nested structs
	employee := Employee{
		Person: Person{
			FirstName: "Sarah",
			LastName:  "Davis",
			Age:       32,
			Height:    1.70,
			Active:    true,
		},
		EmployeeID: "E12345",
		Position:   "Software Engineer",
		Salary:     85000,
		HomeAddress: Address{
			Street:  "123 Home St",
			City:    "Hometown",
			State:   "CA",
			ZipCode: "12345",
			Country: "USA",
		},
		OfficeAddress: Address{
			Street:  "456 Office Blvd",
			City:    "Worktown",
			State:   "CA",
			ZipCode: "54321",
			Country: "USA",
		},
		HireDate: time.Now().AddDate(-3, -6, 0), // 3 years and 6 months ago
	}

	// Accessing fields in nested structs
	fmt.Printf("Employee: %s %s\n", employee.FirstName, employee.LastName) // Embedded fields accessed directly
	fmt.Printf("Position: %s\n", employee.Position)
	fmt.Printf("Home City: %s\n", employee.HomeAddress.City)
	fmt.Printf("Office City: %s\n", employee.OfficeAddress.City)

	// SECTION 3: Struct Methods
	fmt.Println("\n=== Struct Methods ===")

	// Using value receiver methods
	fmt.Printf("Full Name: %s\n", person1.FullName())

	// Using pointer receiver methods
	fmt.Printf("Age before birthday: %d\n", person1.Age)
	person1.Birthday()
	fmt.Printf("Age after birthday: %d\n", person1.Age)

	// Methods on nested structs
	fmt.Printf("Promotion eligible: %t\n", employee.PromotionEligible())
	fmt.Printf("Salary before raise: $%.2f\n", employee.Salary)
	employee.GiveRaise(10) // 10% raise
	fmt.Printf("Salary after raise: $%.2f\n", employee.Salary)

	// SECTION 4: Struct Tags and JSON
	fmt.Println("\n=== Struct Tags and JSON ===")

	// Create a product
	product := Product{
		ID:          1001,
		Name:        "Laptop",
		Description: "High-performance laptop",
		Price:       1299.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   time.Time{}, // Zero value
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(product, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
	} else {
		fmt.Println("JSON output:")
		fmt.Println(string(jsonData))
	}

	// Create a product with empty description to demonstrate omitempty
	productNoDesc := Product{
		ID:        1002,
		Name:      "Mouse",
		Price:     24.99,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	jsonDataNoDesc, _ := json.MarshalIndent(productNoDesc, "", "  ")
	fmt.Println("\nJSON output with omitempty:")
	fmt.Println(string(jsonDataNoDesc))

	// Unmarshal from JSON
	jsonString := `{"id":1003,"name":"Keyboard","description":"Mechanical keyboard","price":89.99,"created_at":"2023-01-15T12:00:00Z","updated_at":"2023-01-15T12:00:00Z"}`
	var newProduct Product

	err = json.Unmarshal([]byte(jsonString), &newProduct)
	if err != nil {
		fmt.Println("Error unmarshaling from JSON:", err)
	} else {
		fmt.Println("\nUnmarshaled product:")
		fmt.Printf("%+v\n", newProduct)
	}

	// SECTION 5: Anonymous Structs
	fmt.Println("\n=== Anonymous Structs ===")
	createAnonymousStruct()

	// Another example of anonymous struct
	config := struct {
		Host     string
		Port     int
		Debug    bool
		MaxConns int
	}{
		Host:     "localhost",
		Port:     8080,
		Debug:    true,
		MaxConns: 10,
	}

	fmt.Printf("Config: %+v\n", config)

	// SECTION 6: Struct Equality
	fmt.Println("\n=== Struct Equality ===")
	compareStructs()

	// SECTION 7: Struct Factory Functions
	fmt.Println("\n=== Struct Factory Functions ===")

	// Already demonstrated with NewPerson above
	// Now demonstrate NewEmployee
	newEmp := NewEmployee("David", "Wilson", 29, "E67890", "Product Manager", 95000)
	fmt.Printf("New employee: %s %s, Position: %s\n", newEmp.FirstName, newEmp.LastName, newEmp.Position)

	// SECTION 8: Practical Examples
	fmt.Println("\n=== Practical Examples ===")

	// Example 1: Modeling a shape hierarchy
	fmt.Println("\nShape hierarchy example:")

	// Define shapes
	type Shape interface {
		Area() float64
		Perimeter() float64
	}

	type Rectangle struct {
		Width  float64
		Height float64
	}

func (r Rectangle) Area() float64;
		return r.Width * r.Height
	}

	func (r Rectangle) Perimeter() float64 {
		return 2 * (r.Width + r.Height)
	}

	type Circle struct {
		Radius float64
	}

	func (c Circle) Area() float64 {
		return 3.14159 * c.Radius * c.Radius
	}

	func (c Circle) Perimeter() float64 {
		return 2 * 3.14159 * c.Radius
	}

	// Create and use shapes
	rect := Rectangle{Width: 5, Height: 3}
	circ := Circle{Radius: 2}

	fmt.Printf("Rectangle - Area: %.2f, Perimeter: %.2f\n", rect.Area(), rect.Perimeter())
	fmt.Printf("Circle - Area: %.2f, Perimeter: %.2f\n", circ.Area(), circ.Perimeter())

	// Example 2: Modeling a bank account
	fmt.Println("\nBank account example:")

	type Account struct {
		Owner   Person
		Number  string
		Balance float64
	}

	func (a *Account) Deposit(amount float64) float64 {
		if amount <= 0 {
			fmt.Println("Deposit amount must be positive")
			return a.Balance
		}
		a.Balance += amount
		return a.Balance
	}

	func (a *Account) Withdraw(amount float64) float64 {
		if amount <= 0 {
			fmt.Println("Withdrawal amount must be positive")
			return a.Balance
		}
		if amount > a.Balance {
			fmt.Println("Insufficient funds")
			return a.Balance
		}
		a.Balance -= amount
		return a.Balance
	}

	// Create and use an account
	account := Account{
		Owner:   person1,
		Number:  "123456789",
		Balance: 1000,
	}

	fmt.Printf("Initial balance: $%.2f\n", account.Balance)
	fmt.Printf("After deposit: $%.2f\n", account.Deposit(500))
	fmt.Printf("After withdrawal: $%.2f\n", account.Withdraw(200))
	fmt.Printf("Insufficient withdrawal: $%.2f\n", account.Withdraw(2000))
}