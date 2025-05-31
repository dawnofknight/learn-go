# JSON in Go

This section covers working with JSON (JavaScript Object Notation) in Go, which is a lightweight data interchange format widely used for APIs, configuration files, and data storage. Go provides excellent support for JSON encoding and decoding through the `encoding/json` package.

## Topics Covered

### 1. Basic JSON Marshaling (Go to JSON)

- Converting Go structs to JSON using `json.Marshal`
- Pretty printing JSON with `json.MarshalIndent`
- Understanding the default field naming convention

### 2. JSON Tags and Field Customization

- Using struct tags to customize JSON field names
- Omitting empty fields with `omitempty`
- Ignoring fields with `-`
- Controlling JSON output format

### 3. Basic JSON Unmarshaling (JSON to Go)

- Converting JSON strings to Go structs using `json.Unmarshal`
- Handling errors during unmarshaling

### 4. Working with JSON Arrays

- Marshaling slices to JSON arrays
- Unmarshaling JSON arrays to Go slices

### 5. Working with Nested JSON

- Handling nested objects with embedded structs
- Using pointers for optional nested fields
- Marshaling and unmarshaling complex nested structures

### 6. Working with Unknown JSON Structure

- Using `map[string]interface{}` for dynamic JSON
- Using `json.RawMessage` for partial decoding
- Handling arbitrary JSON data

### 7. Custom JSON Marshaling and Unmarshaling

- Implementing the `json.Marshaler` interface
- Implementing the `json.Unmarshaler` interface
- Creating custom types with special JSON handling

### 8. Working with JSON Files and Streams

- Reading and writing JSON to files
- Using `json.Encoder` and `json.Decoder`
- Processing JSON streams

## Running the Program

To run the program, navigate to this directory and execute:

```bash
go run main.go
```

## Key Concepts

1. **JSON and Go Types Mapping**:

   - Go booleans ↔ JSON booleans
   - Go numbers ↔ JSON numbers
   - Go strings ↔ JSON strings
   - Go arrays and slices ↔ JSON arrays
   - Go structs and maps ↔ JSON objects
   - Go `nil` ↔ JSON `null`

2. **Field Visibility**: Only exported fields (starting with an uppercase letter) in Go structs are marshaled to JSON.

3. **Struct Tags**: Use the `json:"fieldname,options"` tag to customize how struct fields are marshaled and unmarshaled.

4. **Error Handling**: Always check for errors when marshaling and unmarshaling JSON, as invalid JSON or type mismatches can cause failures.

5. **Performance Considerations**:

   - For large JSON data, consider using `json.Encoder` and `json.Decoder` to avoid loading the entire data into memory.
   - For high-performance applications, consider using alternative JSON libraries like `jsoniter` or `easyjson`.

6. **Type Safety**: Go's static typing helps catch many JSON-related errors at compile time, but runtime errors can still occur with dynamic JSON.

7. **Time Handling**: The `time.Time` type is automatically marshaled to RFC 3339 format in JSON.

## Exercises

1. Create a program that reads a JSON configuration file and uses the configuration to set up an application.

2. Implement a simple REST API that accepts and returns JSON data.

3. Create a custom type for handling currency values in JSON with proper formatting and validation.

4. Write a program that can convert between JSON and other formats (like CSV or XML).

5. Implement a JSON diff tool that compares two JSON objects and identifies the differences.
