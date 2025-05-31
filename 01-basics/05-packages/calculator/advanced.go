package calculator

import (
	"math"
)

// Power calculates a raised to the power of b
func Power(a, b int) int {
	// Check if inputs are within range
	if !checkRange(a) || !checkRange(b) {
		return 0
	}

	return int(math.Pow(float64(a), float64(b)))
}

// SquareRoot calculates the square root of a number
func SquareRoot(a int) float64 {
	// Ensure we don't try to take the square root of a negative number
	if a < 0 {
		return 0
	}
	return math.Sqrt(float64(a))
}

// Factorial calculates the factorial of a number
func Factorial(n int) int {
	if n < 0 {
		return 0
	}
	if n == 0 || n == 1 {
		return 1
	}
	return n * Factorial(n-1)
}

// Absolute returns the absolute value of a number
func Absolute(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// GCD calculates the greatest common divisor of two numbers
func GCD(a, b int) int {
	a, b = Absolute(a), Absolute(b)
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM calculates the least common multiple of two numbers
func LCM(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return Absolute(a*b) / GCD(a, b)
}
