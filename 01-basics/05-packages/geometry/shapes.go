// Package geometry provides types and functions for geometric calculations
package geometry

import (
	"fmt"
	"math"
)

// Pi is an exported constant
const Pi = 3.14159265359

// init function runs when the package is initialized
func init() {
	fmt.Println("Geometry package initialized")
}

// Shape is an interface that defines methods for calculating area and perimeter
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle represents a circle with a radius
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with the given radius
func NewCircle(radius float64) Circle {
	return Circle{Radius: radius}
}

// Area calculates the area of the circle
func (c Circle) Area() float64 {
	return Pi * c.Radius * c.Radius
}

// Perimeter calculates the perimeter (circumference) of the circle
func (c Circle) Perimeter() float64 {
	return 2 * Pi * c.Radius
}

// Rectangle represents a rectangle with width and height
type Rectangle struct {
	Width  float64
	Height float64
}

// NewRectangle creates a new Rectangle with the given width and height
func NewRectangle(width, height float64) Rectangle {
	return Rectangle{Width: width, Height: height}
}

// Area calculates the area of the rectangle
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter calculates the perimeter of the rectangle
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Triangle represents a triangle with three sides
type Triangle struct {
	A, B, C float64 // Lengths of the three sides
}

// NewTriangle creates a new Triangle with the given side lengths
func NewTriangle(a, b, c float64) (Triangle, error) {
	// Check if the triangle is valid (sum of any two sides > third side)
	if a+b <= c || a+c <= b || b+c <= a {
		return Triangle{}, fmt.Errorf("invalid triangle: sides %f, %f, %f", a, b, c)
	}
	return Triangle{A: a, B: b, C: c}, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t Triangle) Area() float64 {
	// Semi-perimeter
	s := (t.A + t.B + t.C) / 2
	// Heron's formula
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

// Perimeter calculates the perimeter of the triangle
func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

// calculateArea is an unexported helper function
// It's not accessible outside this package
func calculateArea(width, height float64) float64 {
	return width * height
}
