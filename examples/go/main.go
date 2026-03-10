package main

import (
	"fmt"
	"math"
)

// Point represents a 2D coordinate.
type Point struct {
	X, Y float64
}

// Distance calculates the Euclidean distance from the origin.
func (p Point) Distance() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

// Add returns a new point by adding another point.
func (p Point) Add(other Point) Point {
	return Point{X: p.X + other.X, Y: p.Y + other.Y}
}

// Shape defines a geometric shape interface.
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle implements the Shape interface.
type Circle struct {
	Center Point
	Radius float64
}

// Area returns the area of the circle.
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter returns the circumference of the circle.
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// NewCircle creates a circle at the given center with the given radius.
func NewCircle(center Point, radius float64) (*Circle, error) {
	if radius < 0 {
		return nil, fmt.Errorf("radius must be non-negative, got %f", radius)
	}
	return &Circle{Center: center, Radius: radius}, nil
}

func main() {
	p := Point{X: 3, Y: 4}
	fmt.Printf("Distance: %f\n", p.Distance())
}
