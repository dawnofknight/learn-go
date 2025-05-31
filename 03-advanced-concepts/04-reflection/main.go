package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("=== Go Reflection Examples ===")

	// 1. Basic Type Reflection
	fmt.Println("\n1. Basic Type Reflection:")
	basicTypeReflection()

	// 2. Struct Reflection
	fmt.Println("\n2. Struct Reflection:")
	structReflection()

	// 3. Function Reflection
	fmt.Println("\n3. Function Reflection:")
	functionReflection()

	// 4. Modifying Values with Reflection
	fmt.Println("\n4. Modifying Values with Reflection:")
	modifyingValuesWithReflection()

	// 5. Creating Values with Reflection
	fmt.Println("\n5. Creating Values with Reflection:")
	creatingValuesWithReflection()

	// 6. Practical Example: Simple ORM
	fmt.Println("\n6. Practical Example: Simple ORM:")
	ormExample()

	// 7. Practical Example: Generic Print Function
	fmt.Println("\n7. Practical Example: Generic Print Function:")
	genericPrintExample()

	// 8. Practical Example: Config Parser
	fmt.Println("\n8. Practical Example: Config Parser:")
	configParserExample()

	// 9. Reflection Performance Considerations
	fmt.Println("\n9. Reflection Performance Considerations:")
	reflectionPerformanceExample()

	fmt.Println("\nAll reflection examples completed!")
}

// 1. Basic Type Reflection
func basicTypeReflection() {
	// Reflect on basic types
	var i int = 42
	var f float64 = 3.14
	var s string = "hello"
	var b bool = true

	// Get type and value information
	inspectValue(i)
	inspectValue(f)
	inspectValue(s)
	inspectValue(b)

	// Reflect on a slice
	numbers := []int{1, 2, 3, 4, 5}
	inspectValue(numbers)

	// Reflect on a map
	colors := map[string]string{"red": "#FF0000", "green": "#00FF00", "blue": "#0000FF"}
	inspectValue(colors)
}

// Helper function to inspect a value using reflection
func inspectValue(x interface{}) {
	v := reflect.ValueOf(x)
	t := reflect.TypeOf(x)

	fmt.Printf("Value: %v\n", v)
	fmt.Printf("Type: %v\n", t)
	fmt.Printf("Kind: %v\n", v.Kind())

	// Additional information based on the kind
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Printf("Is integer: true, Value: %d\n", v.Int())
	case reflect.Float32, reflect.Float64:
		fmt.Printf("Is float: true, Value: %f\n", v.Float())
	case reflect.String:
		fmt.Printf("Is string: true, Value: %s, Length: %d\n", v.String(), v.Len())
	case reflect.Bool:
		fmt.Printf("Is bool: true, Value: %t\n", v.Bool())
	case reflect.Slice, reflect.Array:
		fmt.Printf("Is slice/array: true, Length: %d, Capacity: %d\n", v.Len(), v.Cap())
		if v.Len() > 0 {
			fmt.Printf("First element: %v\n", v.Index(0))
		}
	case reflect.Map:
		fmt.Printf("Is map: true, Length: %d\n", v.Len())
		keys := v.MapKeys()
		if len(keys) > 0 {
			fmt.Printf("First key: %v, Value: %v\n", keys[0], v.MapIndex(keys[0]))
		}
	}
	fmt.Println("---")
}

// 2. Struct Reflection
func structReflection() {
	// Define a struct with tags
	type Person struct {
		Name    string `json:"name" validate:"required"`
		Age     int    `json:"age" validate:"min=0,max=130"`
		Address string `json:"address,omitempty"`
		Private string `json:"-"`
	}

	// Create an instance
	p := Person{
		Name:    "John Doe",
		Age:     30,
		Address: "123 Main St",
		Private: "secret",
	}

	// Get type information
	t := reflect.TypeOf(p)
	fmt.Printf("Type: %v\n", t)
	fmt.Printf("Kind: %v\n", t.Kind())
	fmt.Printf("Number of fields: %d\n", t.NumField())

	// Iterate through fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Printf("Field %d: %s (%s)\n", i, field.Name, field.Type)

		// Get field tags
		jsonTag := field.Tag.Get("json")
		validateTag := field.Tag.Get("validate")
		fmt.Printf("  json tag: %s\n", jsonTag)
		if validateTag != "" {
			fmt.Printf("  validate tag: %s\n", validateTag)
		}
	}

	// Get value information
	v := reflect.ValueOf(p)
	fmt.Println("\nField values:")
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		fmt.Printf("  %s: %v\n", fieldType.Name, fieldValue.Interface())
	}

	// Check if a field exists
	_, found := t.FieldByName("Email")
	fmt.Printf("\nHas Email field: %t\n", found)

	_, found = t.FieldByName("Name")
	fmt.Printf("Has Name field: %t\n", found)
}

// 3. Function Reflection
func functionReflection() {
	// Define some functions to reflect on
	add := func(a, b int) int { return a + b }
	greet := func(name string) string { return "Hello, " + name }
	multiReturn := func(x int) (int, error) { return x * 2, nil }

	// Reflect on the add function
	addValue := reflect.ValueOf(add)
	addType := reflect.TypeOf(add)
	fmt.Printf("Function type: %v\n", addType)
	fmt.Printf("Is function: %t\n", addValue.Kind() == reflect.Func)
	fmt.Printf("Number of inputs: %d\n", addType.NumIn())
	fmt.Printf("Number of outputs: %d\n", addType.NumOut())

	// Get parameter and return types
	fmt.Println("\nParameter types:")
	for i := 0; i < addType.NumIn(); i++ {
		fmt.Printf("  Parameter %d: %v\n", i, addType.In(i))
	}

	fmt.Println("Return types:")
	for i := 0; i < addType.NumOut(); i++ {
		fmt.Printf("  Return %d: %v\n", i, addType.Out(i))
	}

	// Call a function using reflection
	fmt.Println("\nCalling function with reflection:")
	args := []reflect.Value{
		reflect.ValueOf(5),
		reflect.ValueOf(7),
	}
	result := addValue.Call(args)
	fmt.Printf("Result of add(5, 7): %v\n", result[0].Interface())

	// Reflect on the greet function
	greetValue := reflect.ValueOf(greet)
	greetArgs := []reflect.Value{reflect.ValueOf("World")}
	greetResult := greetValue.Call(greetArgs)
	fmt.Printf("Result of greet(\"World\"): %v\n", greetResult[0].Interface())

	// Reflect on the multiReturn function
	multiValue := reflect.ValueOf(multiReturn)
	multiArgs := []reflect.Value{reflect.ValueOf(10)}
	multiResults := multiValue.Call(multiArgs)
	fmt.Printf("Results of multiReturn(10): %v, %v\n",
		multiResults[0].Interface(),
		multiResults[1].Interface())
}

// 4. Modifying Values with Reflection
func modifyingValuesWithReflection() {
	// Modify basic types
	var x int = 10
	fmt.Printf("Original x: %d\n", x)

	// Get a reflect.Value for x
	v := reflect.ValueOf(&x) // Note: we need a pointer to modify the value
	indirect := v.Elem()     // Dereference the pointer

	// Check if we can set the value
	if indirect.CanSet() {
		indirect.SetInt(20)
		fmt.Printf("Modified x: %d\n", x)
	} else {
		fmt.Println("Cannot modify x")
	}

	// Modify a struct field
	type User struct {
		Name  string
		Email string
		Age   int
	}

	user := User{"Alice", "alice@example.com", 30}
	fmt.Printf("\nOriginal user: %+v\n", user)

	// Get a reflect.Value for the user struct
	userValue := reflect.ValueOf(&user).Elem()

	// Modify the Name field
	nameField := userValue.FieldByName("Name")
	if nameField.CanSet() {
		nameField.SetString("Bob")
	}

	// Modify the Age field
	ageField := userValue.FieldByName("Age")
	if ageField.CanSet() {
		ageField.SetInt(35)
	}

	fmt.Printf("Modified user: %+v\n", user)

	// Modify a slice
	numbers := []int{1, 2, 3, 4, 5}
	fmt.Printf("\nOriginal slice: %v\n", numbers)

	sliceValue := reflect.ValueOf(numbers)
	if sliceValue.Len() > 0 {
		// Modify the first element
		sliceValue.Index(0).SetInt(10)
		// Modify the last element
		sliceValue.Index(sliceValue.Len() - 1).SetInt(50)
	}

	fmt.Printf("Modified slice: %v\n", numbers)

	// Modify a map
	colors := map[string]string{"red": "#FF0000", "green": "#00FF00"}
	fmt.Printf("\nOriginal map: %v\n", colors)

	mapValue := reflect.ValueOf(colors)
	mapValue.SetMapIndex(reflect.ValueOf("red"), reflect.ValueOf("#FF0001"))
	mapValue.SetMapIndex(reflect.ValueOf("blue"), reflect.ValueOf("#0000FF"))

	fmt.Printf("Modified map: %v\n", colors)
}

// 5. Creating Values with Reflection
func creatingValuesWithReflection() {
	// Create a new int
	intType := reflect.TypeOf(0)
	intValue := reflect.New(intType).Elem() // Create a new int value
	intValue.SetInt(42)                     // Set the value

	// Convert back to a regular int
	i := intValue.Interface().(int)
	fmt.Printf("Created int: %d\n", i)

	// Create a new struct
	type Person struct {
		Name string
		Age  int
	}

	// Get the type of Person
	personType := reflect.TypeOf(Person{})

	// Create a new Person value
	personValue := reflect.New(personType).Elem()

	// Set the fields
	personValue.FieldByName("Name").SetString("John")
	personValue.FieldByName("Age").SetInt(25)

	// Convert back to a regular Person
	person := personValue.Interface().(Person)
	fmt.Printf("\nCreated person: %+v\n", person)

	// Create a slice
	sliceType := reflect.SliceOf(intType)
	sliceValue := reflect.MakeSlice(sliceType, 3, 5)

	// Set the elements
	for i := 0; i < sliceValue.Len(); i++ {
		sliceValue.Index(i).SetInt(int64(i + 1))
	}

	// Convert back to a regular slice
	slice := sliceValue.Interface().([]int)
	fmt.Printf("\nCreated slice: %v, Length: %d, Capacity: %d\n", slice, len(slice), cap(slice))

	// Create a map
	mapType := reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(""))
	mapValue := reflect.MakeMap(mapType)

	// Add entries
	mapValue.SetMapIndex(reflect.ValueOf("key1"), reflect.ValueOf("value1"))
	mapValue.SetMapIndex(reflect.ValueOf("key2"), reflect.ValueOf("value2"))

	// Convert back to a regular map
	m := mapValue.Interface().(map[string]string)
	fmt.Printf("\nCreated map: %v\n", m)
}

// 6. Practical Example: Simple ORM
func ormExample() {
	// Define a struct representing a database table
	type User struct {
		ID        int    `db:"id" primary_key:"true"`
		Username  string `db:"username" max_length:"50"`
		Email     string `db:"email" max_length:"100"`
		Age       int    `db:"age"`
		IsActive  bool   `db:"is_active"`
		CreatedAt string `db:"created_at"`
	}

	// Create a user instance
	user := User{
		ID:        1,
		Username:  "johndoe",
		Email:     "john@example.com",
		Age:       30,
		IsActive:  true,
		CreatedAt: "2023-01-01",
	}

	// Generate SQL for inserting the user
	sql, values := generateInsertSQL(user)
	fmt.Printf("Insert SQL: %s\n", sql)
	fmt.Printf("Values: %v\n", values)

	// Generate SQL for updating the user
	sql, values = generateUpdateSQL(user, "id")
	fmt.Printf("\nUpdate SQL: %s\n", sql)
	fmt.Printf("Values: %v\n", values)

	// Generate SQL for selecting the user
	sql = generateSelectSQL(User{}, "username = ?")
	fmt.Printf("\nSelect SQL: %s\n", sql)
}

// Helper function to generate INSERT SQL
func generateInsertSQL(model interface{}) (string, []interface{}) {
	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)

	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
		modelValue = modelValue.Elem()
	}

	tableName := strings.ToLower(modelType.Name())
	var columns []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		dbTag := field.Tag.Get("db")

		if dbTag != "" && dbTag != "-" {
			columns = append(columns, dbTag)
			placeholders = append(placeholders, "?")
			values = append(values, modelValue.Field(i).Interface())
		}
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	return sql, values
}

// Helper function to generate UPDATE SQL
func generateUpdateSQL(model interface{}, primaryKeyField string) (string, []interface{}) {
	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)

	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
		modelValue = modelValue.Elem()
	}

	tableName := strings.ToLower(modelType.Name())
	var setStatements []string
	var values []interface{}
	var primaryKeyValue interface{}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		dbTag := field.Tag.Get("db")

		if dbTag != "" && dbTag != "-" {
			if dbTag == primaryKeyField {
				primaryKeyValue = modelValue.Field(i).Interface()
				continue
			}
			setStatements = append(setStatements, fmt.Sprintf("%s = ?", dbTag))
			values = append(values, modelValue.Field(i).Interface())
		}
	}

	// Add the primary key value at the end for the WHERE clause
	values = append(values, primaryKeyValue)

	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?",
		tableName,
		strings.Join(setStatements, ", "),
		primaryKeyField)

	return sql, values
}

// Helper function to generate SELECT SQL
func generateSelectSQL(model interface{}, whereClause string) string {
	modelType := reflect.TypeOf(model)

	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	tableName := strings.ToLower(modelType.Name())
	var columns []string

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		dbTag := field.Tag.Get("db")

		if dbTag != "" && dbTag != "-" {
			columns = append(columns, dbTag)
		}
	}

	sql := fmt.Sprintf("SELECT %s FROM %s",
		strings.Join(columns, ", "),
		tableName)

	if whereClause != "" {
		sql += " WHERE " + whereClause
	}

	return sql
}

// 7. Practical Example: Generic Print Function
func genericPrintExample() {
	// Define some values to print
	intValue := 42
	floatValue := 3.14
	stringValue := "hello"
	boolValue := true
	sliceValue := []int{1, 2, 3}
	mapValue := map[string]int{"one": 1, "two": 2}

	// Define a struct
	type Person struct {
		Name string
		Age  int
	}
	structValue := Person{"Alice", 30}

	// Print all values using our generic print function
	fmt.Println("Generic printing:")
	genericPrint(intValue)
	genericPrint(floatValue)
	genericPrint(stringValue)
	genericPrint(boolValue)
	genericPrint(sliceValue)
	genericPrint(mapValue)
	genericPrint(structValue)
}

// Generic print function using reflection
func genericPrint(value interface{}) {
	v := reflect.ValueOf(value)
	t := reflect.TypeOf(value)

	fmt.Printf("Type: %v, ", t)

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Printf("Value: %d\n", v.Int())
	case reflect.Float32, reflect.Float64:
		fmt.Printf("Value: %f\n", v.Float())
	case reflect.String:
		fmt.Printf("Value: \"%s\"\n", v.String())
	case reflect.Bool:
		fmt.Printf("Value: %t\n", v.Bool())
	case reflect.Slice, reflect.Array:
		fmt.Printf("Value: [")
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(v.Index(i))
		}
		fmt.Println("]")
	case reflect.Map:
		fmt.Printf("Value: {")
		keys := v.MapKeys()
		for i, key := range keys {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%v: %v", key, v.MapIndex(key))
		}
		fmt.Println("}")
	case reflect.Struct:
		fmt.Printf("Value: {")
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%s: %v", t.Field(i).Name, v.Field(i))
		}
		fmt.Println("}")
	default:
		fmt.Printf("Value: %v\n", v)
	}
}

// 8. Practical Example: Config Parser
func configParserExample() {
	// Define a struct for configuration
	type Config struct {
		ServerName string `config:"server_name" default:"localhost"`
		Port       int    `config:"port" default:"8080"`
		Debug      bool   `config:"debug" default:"false"`
		LogLevel   string `config:"log_level" default:"info"`
		MaxUsers   int    `config:"max_users" default:"100"`
	}

	// Create a default config
	config := Config{}

	// Parse configuration from a string (simulating a config file)
	configStr := `
	server_name = myserver
	port = 9090
	debug = true
	`

	// Parse the config
	parseConfig(&config, configStr)

	// Print the resulting config
	fmt.Printf("Parsed config: %+v\n", config)
}

// Helper function to parse configuration into a struct
func parseConfig(configPtr interface{}, configStr string) {
	// Get the reflect.Value of the pointer
	v := reflect.ValueOf(configPtr)

	// Check if it's a pointer and not nil
	if v.Kind() != reflect.Ptr || v.IsNil() {
		panic("configPtr must be a non-nil pointer")
	}

	// Get the value that the pointer points to
	v = v.Elem()

	// Check if it's a struct
	if v.Kind() != reflect.Struct {
		panic("configPtr must point to a struct")
	}

	// Get the type of the struct
	t := v.Type()

	// First, set default values
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		defaultValue := field.Tag.Get("default")

		if defaultValue != "" {
			setFieldFromString(v.Field(i), defaultValue)
		}
	}

	// Parse the config string
	lines := strings.Split(configStr, "\n")
	for _, line := range lines {
		// Trim spaces
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split by '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Find the corresponding field
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			configKey := field.Tag.Get("config")

			if configKey == key {
				setFieldFromString(v.Field(i), value)
				break
			}
		}
	}
}

// Helper function to set a field value from a string
func setFieldFromString(field reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			field.SetInt(intValue)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uintValue, err := strconv.ParseUint(value, 10, 64); err == nil {
			field.SetUint(uintValue)
		}
	case reflect.Float32, reflect.Float64:
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			field.SetFloat(floatValue)
		}
	case reflect.Bool:
		if boolValue, err := strconv.ParseBool(value); err == nil {
			field.SetBool(boolValue)
		}
	}
}

// 9. Reflection Performance Considerations
func reflectionPerformanceExample() {
	// Define a struct
	type Person struct {
		Name string
		Age  int
	}

	// Create a person
	person := Person{"John", 30}

	// Serialize to JSON using reflection (simulating what json.Marshal does)
	jsonReflection := serializeWithReflection(person)
	fmt.Printf("JSON using reflection: %s\n", jsonReflection)

	// Serialize to JSON using the standard library
	jsonStandard, _ := json.Marshal(person)
	fmt.Printf("JSON using standard library: %s\n", string(jsonStandard))

	// Note: In a real benchmark, you would measure the time taken by each approach
	fmt.Println("\nNote: Reflection is typically slower than type-specific code.")
	fmt.Println("When performance is critical, consider alternatives to reflection:")
	fmt.Println("1. Code generation (e.g., using 'go generate')")
	fmt.Println("2. Interface-based approaches")
	fmt.Println("3. Type switches for a limited set of types")
}

// Simple JSON serializer using reflection
func serializeWithReflection(v interface{}) string {
	value := reflect.ValueOf(v)
	type_ := reflect.TypeOf(v)

	if value.Kind() != reflect.Struct {
		return fmt.Sprintf("%v", v)
	}

	var result strings.Builder
	result.WriteString("{")

	for i := 0; i < value.NumField(); i++ {
		if i > 0 {
			result.WriteString(",")
		}

		field := type_.Field(i)
		fieldValue := value.Field(i)

		// Write the field name
		result.WriteString(fmt.Sprintf("\"%s\":", field.Name))

		// Write the field value based on its kind
		switch fieldValue.Kind() {
		case reflect.String:
			result.WriteString(fmt.Sprintf("\"%s\"", fieldValue.String()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result.WriteString(fmt.Sprintf("%d", fieldValue.Int()))
		case reflect.Float32, reflect.Float64:
			result.WriteString(fmt.Sprintf("%f", fieldValue.Float()))
		case reflect.Bool:
			result.WriteString(fmt.Sprintf("%t", fieldValue.Bool()))
		default:
			result.WriteString("null")
		}
	}

	result.WriteString("}")
	return result.String()
}
