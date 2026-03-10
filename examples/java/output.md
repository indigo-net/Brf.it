# Code Summary: examples/java

*brf.it dev*

---

## Files

### examples/java/ShapeService.java

```java
interface Shape
double area();
double perimeter();
String name();
class Circle implements Shape
public Circle(double radius)
@Override
    public double area()
@Override
    public double perimeter()
@Override
    public String name()
class Rectangle implements Shape
public Rectangle(double width, double height)
@Override
    public double area()
@Override
    public double perimeter()
@Override
    public String name()
public class ShapeService
public void addShape(Shape shape)
public double totalArea()
public Optional<Shape> largestShape()
```

---

