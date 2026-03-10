# Code Summary: examples/go

*brf.it dev*

---

## Files

### examples/go/main.go

```go
type Point struct {
	X, Y float64
}
func (p Point) Distance() float64
func (p Point) Add(other Point) Point
type Shape interface {
	Area() float64
	Perimeter() float64
}
type Circle struct {
	Center Point
	Radius float64
}
func (c Circle) Area() float64
func (c Circle) Perimeter() float64
func NewCircle(center Point, radius float64) (*Circle, error)
func main()
```

---

