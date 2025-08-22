package main

import (
	"testing"
)

// TestAddNumbers tests the addition functionality
func TestAddNumbers(t *testing.T) {
	// Test case 1: Basic addition
	a := 10
	b := 20
	expected := 30

	result := addNumbers(a, b)

	if result != expected {
		t.Errorf("addNumbers(%d, %d) = %d; expected %d", a, b, result, expected)
	}

	// Test case 2: Addition with zero
	a = 5
	b = 0
	expected = 5

	result = addNumbers(a, b)

	if result != expected {
		t.Errorf("addNumbers(%d, %d) = %d; expected %d", a, b, result, expected)
	}

	// Test case 3: Addition with negative numbers
	a = -10
	b = 5
	expected = -5

	result = addNumbers(a, b)

	if result != expected {
		t.Errorf("addNumbers(%d, %d) = %d; expected %d", a, b, result, expected)
	}
}

// TestGetPointerValue tests the pointer dereferencing functionality
func TestGetPointerValue(t *testing.T) {
	value := 30
	pointer := &value

	result := getPointerValue(pointer)

	if result != value {
		t.Errorf("getPointerValue(%v) = %d; expected %d", pointer, result, value)
	}
}
