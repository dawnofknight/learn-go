package geometry

import (
	"math"
)

// Distance calculates the Euclidean distance between two points
func Distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

// IsPointInCircle checks if a point is inside a circle
func IsPointInCircle(px, py, cx, cy, radius float64) bool {
	distance := Distance(px, py, cx, cy)
	return distance <= radius
}

// IsPointInRectangle checks if a point is inside a rectangle
func IsPointInRectangle(px, py, x1, y1, x2, y2 float64) bool {
	return px >= x1 && px <= x2 && py >= y1 && py <= y2
}

// CalculatePolygonArea calculates the area of a polygon given its vertices
func CalculatePolygonArea(vertices [][2]float64) float64 {
	n := len(vertices)
	if n < 3 {
		return 0 // Not a polygon
	}

	area := 0.0
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += vertices[i][0] * vertices[j][1]
		area -= vertices[j][0] * vertices[i][1]
	}

	area = math.Abs(area) / 2.0
	return area
}

// DegreesToRadians converts degrees to radians
func DegreesToRadians(degrees float64) float64 {
	return degrees * Pi / 180.0
}

// RadiansToDegrees converts radians to degrees
func RadiansToDegrees(radians float64) float64 {
	return radians * 180.0 / Pi
}
