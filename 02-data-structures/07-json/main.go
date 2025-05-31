// Package main demonstrates JSON encoding and decoding in Go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// SECTION 1: Basic JSON Marshaling (Go to JSON)

// Simple struct for basic marshaling
type Person struct {
	FirstName string
	LastName  string
	Age       int
	Address   string
	Active    bool
}

func basicMarshaling() {
	// Create a Person struct
	person := Person{
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
		Address:   "123 Main St",
		Active:    true,
	}

	// Marshal the struct to JSON
	jsonData, err := json.Marshal(person)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Print the JSON as a string
	fmt.Println("Basic JSON:")
	fmt.Println(string(jsonData))

	// Marshal with indentation for better readability
	prettyJSON, err := json.MarshalIndent(person, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON with indentation: %v", err)
	}

	// Print the pretty JSON
	fmt.Println("\nPretty JSON:")
	fmt.Println(string(prettyJSON))
}

// SECTION 2: JSON Tags and Field Customization

// Struct with JSON tags for field customization
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"` // Omit if empty
	Price       float64   `json:"price"`
	Tags        []string  `json:"tags,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"-"` // Ignore this field
	InStock     bool      `json:"in_stock"`
	SKU         string    `json:"sku,omitempty"`
}

func jsonTags() {
	// Create a Product with all fields
	product1 := Product{
		ID:          1001,
		Name:        "Laptop",
		Description: "High-performance laptop",
		Price:       1299.99,
		Tags:        []string{"electronics", "computers"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   time.Time{}, // Zero value
		InStock:     true,
		SKU:         "LAP-1001",
	}

	// Create a Product with some empty fields
	product2 := Product{
		ID:        1002,
		Name:      "Mouse",
		Price:     24.99,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		InStock:   false,
		// Description, Tags, and SKU are empty
	}

	// Marshal both products
	jsonProduct1, _ := json.MarshalIndent(product1, "", "  ")
	jsonProduct2, _ := json.MarshalIndent(product2, "", "  ")

	fmt.Println("Product 1 (all fields):")
	fmt.Println(string(jsonProduct1))

	fmt.Println("\nProduct 2 (with omitempty fields):")
	fmt.Println(string(jsonProduct2))
}

// SECTION 3: Basic JSON Unmarshaling (JSON to Go)

func basicUnmarshaling() {
	// JSON string
	jsonStr := `{"FirstName":"Jane","LastName":"Smith","Age":28,"Address":"456 Oak St","Active":true}`

	// Create a Person struct to unmarshal into
	var person Person

	// Unmarshal the JSON into the struct
	err := json.Unmarshal([]byte(jsonStr), &person)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Print the struct
	fmt.Println("Unmarshaled Person:")
	fmt.Printf("%+v\n", person)
}

// SECTION 4: Working with JSON Arrays

func jsonArrays() {
	// Create a slice of Person structs
	people := []Person{
		{FirstName: "John", LastName: "Doe", Age: 30, Address: "123 Main St", Active: true},
		{FirstName: "Jane", LastName: "Smith", Age: 28, Address: "456 Oak St", Active: true},
		{FirstName: "Bob", LastName: "Johnson", Age: 35, Address: "789 Pine Rd", Active: false},
	}

	// Marshal the slice to JSON
	jsonData, err := json.MarshalIndent(people, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON array: %v", err)
	}

	// Print the JSON array
	fmt.Println("JSON Array:")
	fmt.Println(string(jsonData))

	// Unmarshal a JSON array
	jsonArrayStr := `[
		{"FirstName":"Alice","LastName":"Williams","Age":25,"Address":"101 Elm St","Active":true},
		{"FirstName":"Charlie","LastName":"Brown","Age":40,"Address":"202 Maple Ave","Active":false}
	]`

	// Create a slice to unmarshal into
	var newPeople []Person

	// Unmarshal the JSON array
	err = json.Unmarshal([]byte(jsonArrayStr), &newPeople)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON array: %v", err)
	}

	// Print the unmarshaled slice
	fmt.Println("\nUnmarshaled People:")
	for i, p := range newPeople {
		fmt.Printf("%d: %+v\n", i, p)
	}
}

// SECTION 5: Working with Nested JSON

// Nested structs for complex JSON
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

type Contact struct {
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Website string `json:"website,omitempty"`
}

type Employee struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	HomeAddress Address   `json:"home_address"`
	WorkAddress *Address  `json:"work_address,omitempty"` // Pointer for optional nested struct
	Contact     Contact   `json:"contact"`
	HireDate    time.Time `json:"hire_date"`
	Position    string    `json:"position"`
	Salary      float64   `json:"salary"`
	Skills      []string  `json:"skills,omitempty"`
}

func nestedJSON() {
	// Create an Employee with nested structs
	employee := Employee{
		ID:        1001,
		FirstName: "John",
		LastName:  "Doe",
		HomeAddress: Address{
			Street:  "123 Home St",
			City:    "Hometown",
			State:   "CA",
			ZipCode: "12345",
			Country: "USA",
		},
		WorkAddress: &Address{
			Street:  "456 Office Blvd",
			City:    "Worktown",
			State:   "CA",
			ZipCode: "54321",
			Country: "USA",
		},
		Contact: Contact{
			Email:   "john.doe@example.com",
			Phone:   "555-1234",
			Website: "johndoe.com",
		},
		HireDate: time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
		Position: "Software Engineer",
		Salary:   85000,
		Skills:   []string{"Go", "Java", "Python", "SQL"},
	}

	// Marshal the employee to JSON
	jsonData, err := json.MarshalIndent(employee, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling nested JSON: %v", err)
	}

	// Print the nested JSON
	fmt.Println("Nested JSON:")
	fmt.Println(string(jsonData))

	// Create an employee without a work address
	employee2 := employee
	employee2.ID = 1002
	employee2.WorkAddress = nil // This field will be omitted in the JSON
	employee2.Skills = nil      // This field will be omitted in the JSON

	// Marshal the second employee
	jsonData2, _ := json.MarshalIndent(employee2, "", "  ")

	fmt.Println("\nNested JSON with omitted fields:")
	fmt.Println(string(jsonData2))

	// Unmarshal nested JSON
	jsonStr := `{
		"id": 1003,
		"first_name": "Jane",
		"last_name": "Smith",
		"home_address": {
			"street": "789 Home Ave",
			"city": "Othertown",
			"state": "NY",
			"zip_code": "67890",
			"country": "USA"
		},
		"contact": {
			"email": "jane.smith@example.com",
			"phone": "555-5678"
		},
		"hire_date": "2021-03-01T00:00:00Z",
		"position": "Product Manager",
		"salary": 95000
	}`

	var newEmployee Employee
	err = json.Unmarshal([]byte(jsonStr), &newEmployee)
	if err != nil {
		log.Fatalf("Error unmarshaling nested JSON: %v", err)
	}

	fmt.Println("\nUnmarshaled Employee:")
	fmt.Printf("ID: %d\n", newEmployee.ID)
	fmt.Printf("Name: %s %s\n", newEmployee.FirstName, newEmployee.LastName)
	fmt.Printf("Home Address: %s, %s, %s\n",
		newEmployee.HomeAddress.Street,
		newEmployee.HomeAddress.City,
		newEmployee.HomeAddress.State)
	fmt.Printf("Work Address: %v\n", newEmployee.WorkAddress) // nil
	fmt.Printf("Contact: %s, %s\n", newEmployee.Contact.Email, newEmployee.Contact.Phone)
	fmt.Printf("Position: %s\n", newEmployee.Position)
	fmt.Printf("Salary: $%.2f\n", newEmployee.Salary)
}

// SECTION 6: Working with Unknown JSON Structure

func unknownJSON() {
	// JSON with unknown structure
	jsonStr := `{
		"name": "John Doe",
		"age": 30,
		"is_employee": true,
		"address": {
			"street": "123 Main St",
			"city": "Anytown"
		},
		"phone_numbers": ["555-1234", "555-5678"],
		"metadata": {
			"created": "2023-01-15T12:00:00Z",
			"department": "Engineering",
			"projects": ["Alpha", "Beta", "Gamma"]
		}
	}`

	// Method 1: Using map[string]interface{}
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		log.Fatalf("Error unmarshaling to map: %v", err)
	}

	fmt.Println("Using map[string]interface{}:")
	fmt.Printf("Name: %s\n", result["name"])
	fmt.Printf("Age: %.0f\n", result["age"])
	fmt.Printf("Is Employee: %t\n", result["is_employee"])

	// Access nested map
	address := result["address"].(map[string]interface{})
	fmt.Printf("Street: %s\n", address["street"])
	fmt.Printf("City: %s\n", address["city"])

	// Access array
	phoneNumbers := result["phone_numbers"].([]interface{})
	fmt.Println("Phone Numbers:")
	for i, number := range phoneNumbers {
		fmt.Printf("  %d: %s\n", i+1, number)
	}

	// Method 2: Using json.RawMessage for partial decoding
	fmt.Println("\nUsing json.RawMessage:")

	// Define a struct with known fields and RawMessage for unknown parts
	type PartialPerson struct {
		Name       string          `json:"name"`
		Age        int             `json:"age"`
		IsEmployee bool            `json:"is_employee"`
		Address    json.RawMessage `json:"address"`  // Store as raw JSON
		Metadata   json.RawMessage `json:"metadata"` // Store as raw JSON
	}

	var person PartialPerson
	err = json.Unmarshal([]byte(jsonStr), &person)
	if err != nil {
		log.Fatalf("Error unmarshaling with RawMessage: %v", err)
	}

	fmt.Printf("Name: %s\n", person.Name)
	fmt.Printf("Age: %d\n", person.Age)
	fmt.Printf("Is Employee: %t\n", person.IsEmployee)

	// Later, if we need the address details
	type AddressDetails struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	var address2 AddressDetails
	err = json.Unmarshal(person.Address, &address2)
	if err != nil {
		log.Fatalf("Error unmarshaling address: %v", err)
	}

	fmt.Printf("Street: %s\n", address2.Street)
	fmt.Printf("City: %s\n", address2.City)
}

// SECTION 7: Custom JSON Marshaling and Unmarshaling

// Date is a custom type for handling dates in JSON
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// MarshalJSON implements the json.Marshaler interface
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%04d-%02d-%02d\"", d.Year, d.Month, d.Day)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (d *Date) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}

	// Parse the date string
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	d.Year = t.Year()
	d.Month = t.Month()
	d.Day = t.Day()

	return nil
}

// Event with custom date handling
type Event struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Date        Date   `json:"date"` // Custom type
	Location    string `json:"location"`
}

func customMarshaling() {
	// Create an event with a custom date
	event := Event{
		ID:          1001,
		Title:       "Conference",
		Description: "Annual tech conference",
		Date:        Date{2023, 6, 15}, // June 15, 2023
		Location:    "Convention Center",
	}

	// Marshal the event to JSON
	jsonData, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling with custom type: %v", err)
	}

	fmt.Println("Custom Marshaling:")
	fmt.Println(string(jsonData))

	// Unmarshal with custom type
	jsonStr := `{
		"id": 1002,
		"title": "Workshop",
		"description": "Go programming workshop",
		"date": "2023-07-20",
		"location": "Learning Center"
	}`

	var newEvent Event
	err = json.Unmarshal([]byte(jsonStr), &newEvent)
	if err != nil {
		log.Fatalf("Error unmarshaling with custom type: %v", err)
	}

	fmt.Println("\nCustom Unmarshaling:")
	fmt.Printf("ID: %d\n", newEvent.ID)
	fmt.Printf("Title: %s\n", newEvent.Title)
	fmt.Printf("Description: %s\n", newEvent.Description)
	fmt.Printf("Date: %04d-%02d-%02d\n", newEvent.Date.Year, newEvent.Date.Month, newEvent.Date.Day)
	fmt.Printf("Location: %s\n", newEvent.Location)
}

// SECTION 8: Working with JSON Files and Streams

func jsonFilesAndStreams() {
	// Create sample data
	people := []Person{
		{FirstName: "John", LastName: "Doe", Age: 30, Address: "123 Main St", Active: true},
		{FirstName: "Jane", LastName: "Smith", Age: 28, Address: "456 Oak St", Active: true},
		{FirstName: "Bob", LastName: "Johnson", Age: 35, Address: "789 Pine Rd", Active: false},
	}

	// Write JSON to a file
	file, err := os.Create("people.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	// Create an encoder for the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print

	// Encode the data to the file
	err = encoder.Encode(people)
	if err != nil {
		log.Fatalf("Error encoding to file: %v", err)
	}

	fmt.Println("JSON written to people.json")

	// Read JSON from a file
	readFile, err := os.Open("people.json")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer readFile.Close()

	// Create a decoder for the file
	decoder := json.NewDecoder(readFile)

	// Decode the data from the file
	var readPeople []Person
	err = decoder.Decode(&readPeople)
	if err != nil {
		log.Fatalf("Error decoding from file: %v", err)
	}

	fmt.Println("\nJSON read from people.json:")
	for i, p := range readPeople {
		fmt.Printf("%d: %s %s, Age: %d\n", i+1, p.FirstName, p.LastName, p.Age)
	}

	// Working with JSON streams
	fmt.Println("\nWorking with JSON streams:")

	// Create a JSON stream with multiple objects
	jsonStream := strings.NewReader(`
		{"FirstName":"Alice","LastName":"Williams","Age":25}
		{"FirstName":"Charlie","LastName":"Brown","Age":40}
		{"FirstName":"Diana","LastName":"Miller","Age":33}
	`)

	// Create a decoder for the stream
	streamDecoder := json.NewDecoder(jsonStream)

	// Read each object from the stream
	fmt.Println("Reading JSON stream:")
	for {
		var p Person
		if err := streamDecoder.Decode(&p); err == io.EOF {
			break // End of stream
		} else if err != nil {
			log.Fatalf("Error decoding stream: %v", err)
		}

		fmt.Printf("%s %s, Age: %d\n", p.FirstName, p.LastName, p.Age)
	}

	// Creating a JSON stream
	fmt.Println("\nCreating JSON stream:")

	// Create a buffer to hold the stream
	buf := new(bytes.Buffer)

	// Create an encoder for the buffer
	streamEncoder := json.NewEncoder(buf)

	// Encode multiple objects to the stream
	for _, p := range people {
		if err := streamEncoder.Encode(p); err != nil {
			log.Fatalf("Error encoding to stream: %v", err)
		}
	}

	// Print the resulting stream
	fmt.Println(buf.String())
}

func main() {
	// SECTION 1: Basic JSON Marshaling
	fmt.Println("=== Basic JSON Marshaling ===")
	basicMarshaling()

	// SECTION 2: JSON Tags and Field Customization
	fmt.Println("\n=== JSON Tags and Field Customization ===")
	jsonTags()

	// SECTION 3: Basic JSON Unmarshaling
	fmt.Println("\n=== Basic JSON Unmarshaling ===")
	basicUnmarshaling()

	// SECTION 4: Working with JSON Arrays
	fmt.Println("\n=== Working with JSON Arrays ===")
	jsonArrays()

	// SECTION 5: Working with Nested JSON
	fmt.Println("\n=== Working with Nested JSON ===")
	nestedJSON()

	// SECTION 6: Working with Unknown JSON Structure
	fmt.Println("\n=== Working with Unknown JSON Structure ===")
	unknownJSON()

	// SECTION 7: Custom JSON Marshaling and Unmarshaling
	fmt.Println("\n=== Custom JSON Marshaling and Unmarshaling ===")
	customMarshaling()

	// SECTION 8: Working with JSON Files and Streams
	fmt.Println("\n=== Working with JSON Files and Streams ===")
	jsonFilesAndStreams()

	// Clean up the file created during the demo
	os.Remove("people.json")
}
