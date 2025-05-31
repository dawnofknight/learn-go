package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// 1. Basic unit tests

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Add(2, 3) = %d; expected %d", result, expected)
	}
}

func TestSubtract(t *testing.T) {
	result := Subtract(5, 3)
	expected := 2

	if result != expected {
		t.Errorf("Subtract(5, 3) = %d; expected %d", result, expected)
	}
}

func TestMultiply(t *testing.T) {
	result := Multiply(4, 3)
	expected := 12

	if result != expected {
		t.Errorf("Multiply(4, 3) = %d; expected %d", result, expected)
	}
}

// 2. Tests with error handling

func TestDivide(t *testing.T) {
	// Test valid division
	result, err := Divide(10, 2)
	expected := 5.0

	if err != nil {
		t.Errorf("Divide(10, 2) returned an error: %v", err)
	}

	if result != expected {
		t.Errorf("Divide(10, 2) = %.2f; expected %.2f", result, expected)
	}

	// Test division by zero
	_, err = Divide(10, 0)
	if err == nil {
		t.Error("Divide(10, 0) did not return an error, but it should")
	}
}

// 3. Table-driven tests

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		name     string
	}{
		{"radar", true, "simple palindrome"},
		{"A man a plan a canal Panama", true, "palindrome with spaces"},
		{"hello", false, "non-palindrome"},
		{"Madam, I'm Adam", true, "palindrome with punctuation"},
		{"", true, "empty string"},
		{"12321", true, "numeric palindrome"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsPalindrome(test.input)
			if result != test.expected {
				t.Errorf("IsPalindrome(%q) = %v; expected %v", test.input, result, test.expected)
			}
		})
	}
}

// 4. Subtests

func TestUser(t *testing.T) {
	t.Run("NewUser", func(t *testing.T) {
		user := NewUser("testuser", "test@example.com", 25)
		if user.Username != "testuser" {
			t.Errorf("Expected username 'testuser', got '%s'", user.Username)
		}
		if user.Email != "test@example.com" {
			t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
		}
		if user.Age != 25 {
			t.Errorf("Expected age 25, got %d", user.Age)
		}
	})

	t.Run("IsAdult", func(t *testing.T) {
		adult := NewUser("adult", "adult@example.com", 21)
		if !adult.IsAdult() {
			t.Error("Expected IsAdult() to return true for age 21")
		}

		minor := NewUser("minor", "minor@example.com", 17)
		if minor.IsAdult() {
			t.Error("Expected IsAdult() to return false for age 17")
		}
	})

	t.Run("UpdateEmail", func(t *testing.T) {
		user := NewUser("user", "old@example.com", 30)

		// Test valid email update
		err := user.UpdateEmail("new@example.com")
		if err != nil {
			t.Errorf("UpdateEmail returned an error: %v", err)
		}
		if user.Email != "new@example.com" {
			t.Errorf("Email not updated, got: %s", user.Email)
		}

		// Test invalid email update
		err = user.UpdateEmail("invalid")
		if err == nil {
			t.Error("UpdateEmail with invalid email did not return an error")
		}
	})
}

// 5. Test setup and teardown

func TestMain(m *testing.M) {
	// Setup code here (if needed)
	fmt.Println("Starting tests...")

	// Run the tests
	code := m.Run()

	// Teardown code here (if needed)
	fmt.Println("Tests completed.")

	// Exit with the test status code
	// In a real test, you would use os.Exit(code) here
	// but we'll just print it for demonstration
	fmt.Printf("Exit code: %d\n", code)
}

// 6. Mocking

// MockDatabase is a mock implementation of the Database interface
type MockDatabase struct {
	users map[string]*User
}

// NewMockDatabase creates a new mock database
func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		users: make(map[string]*User),
	}
}

// GetUser retrieves a user from the mock database
func (m *MockDatabase) GetUser(id string) (*User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// SaveUser saves a user to the mock database
func (m *MockDatabase) SaveUser(user *User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	m.users[user.Username] = user
	return nil
}

func TestUserService(t *testing.T) {
	// Create a mock database
	mockDB := NewMockDatabase()

	// Create a user service with the mock database
	userService := &UserService{DB: mockDB}

	// Test SaveUser
	t.Run("SaveUser", func(t *testing.T) {
		user := NewUser("testuser", "test@example.com", 25)

		// Save the user
		err := userService.SaveUser(user)
		if err != nil {
			t.Errorf("SaveUser returned an error: %v", err)
		}

		// Verify the user was saved
		savedUser, err := mockDB.GetUser("testuser")
		if err != nil {
			t.Errorf("GetUser returned an error: %v", err)
		}

		if !reflect.DeepEqual(user, savedUser) {
			t.Errorf("Saved user does not match original user")
		}
	})

	// Test GetUserByID
	t.Run("GetUserByID", func(t *testing.T) {
		// Add a user to the mock database
		user := NewUser("getuser", "get@example.com", 30)
		mockDB.users["getuser"] = user

		// Get the user
		retrievedUser, err := userService.GetUserByID("getuser")
		if err != nil {
			t.Errorf("GetUserByID returned an error: %v", err)
		}

		if !reflect.DeepEqual(user, retrievedUser) {
			t.Errorf("Retrieved user does not match original user")
		}

		// Test getting a non-existent user
		_, err = userService.GetUserByID("nonexistent")
		if err == nil {
			t.Error("GetUserByID for non-existent user did not return an error")
		}
	})
}

// 7. Benchmarks

func BenchmarkIsPalindrome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPalindrome("A man a plan a canal Panama")
	}
}

func BenchmarkReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Reverse("hello world")
	}
}

// 8. Examples (these also serve as documentation)

func ExampleAdd() {
	result := Add(2, 3)
	fmt.Println(result)
	// Output: 5
}

func ExampleIsPalindrome() {
	fmt.Println(IsPalindrome("radar"))
	fmt.Println(IsPalindrome("hello"))
	// Output:
	// true
	// false
}

func ExampleUser_IsAdult() {
	user := NewUser("example", "example@example.com", 21)
	fmt.Println(user.IsAdult())
	// Output: true
}
