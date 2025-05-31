# Maps in Go

This section covers maps, one of Go's built-in data structures that implements a hash table providing fast lookups, adds, and deletes. Maps are used to store key-value pairs where each key is unique.

## Topics Covered

### 1. Map Declaration and Initialization

- Using map literals: `map[KeyType]ValueType{key1: value1, key2: value2}`
- Using the `make` function: `make(map[KeyType]ValueType)`
- Nil maps: `var m map[KeyType]ValueType`

### 2. Map Operations

- Adding elements: `m[key] = value`
- Accessing elements: `value := m[key]`
- Checking if a key exists: `value, exists := m[key]`
- Updating elements: `m[key] = newValue`
- Deleting elements: `delete(m, key)`

### 3. Iterating Over Maps

- Using the `range` keyword (unordered iteration)
- Iterating in a specific order by sorting keys

### 4. Map of Maps

- Creating nested maps
- Accessing nested map values
- Adding nested maps

### 5. Maps as Sets

- Using maps with boolean values as sets
- Using maps with empty struct values for memory-efficient sets

### 6. Map Capacity and Performance

- Creating maps with initial capacity
- Understanding map performance characteristics

### 7. Practical Examples

- Word frequency counter
- Grouping data
- Simple cache implementation

## Running the Program

To run the program, navigate to this directory and execute:

```bash
go run main.go
```

## Key Concepts

1. **Maps are Reference Types**: When you assign a map to a new variable or pass it to a function, you're passing a reference to the same underlying data structure.

2. **Maps are Unordered**: The order of iteration over map elements is not guaranteed and can change between runs.

3. **Zero Value**: The zero value of a map is `nil`. A nil map has no keys and attempting to add keys will result in a runtime panic.

4. **Concurrent Access**: Maps in Go are not safe for concurrent use. If you need to use a map concurrently, you must implement synchronization (e.g., using mutexes).

5. **Key Types**: Map keys must be comparable types (types that support the `==` operator). This includes booleans, numbers, strings, pointers, channels, interfaces, structs, and arrays of comparable types. Slices, maps, and functions are not comparable and cannot be used as map keys.

6. **Value Types**: Map values can be of any type, including other maps or slices.

## Exercises

1. Create a program that counts the frequency of each character in a string.

2. Implement a simple phone book application that allows adding, searching, and deleting contacts.

3. Create a function that takes a slice of strings and returns a map where the keys are the strings and the values are the number of times each string appears in the slice.

4. Implement a simple LRU (Least Recently Used) cache using maps.

5. Create a program that reads a text file and creates a concordance (a map where keys are words and values are lists of line numbers where the word appears).
