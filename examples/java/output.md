# Code Summary: examples/java

*brf.it dev*

---

## Files

### examples/java/ShapeService.java

```java
public interface Shape
double area();
double perimeter();
String name();
public class Circle implements Shape
public Circle(double radius)
@Override
    public double area()
@Override
    public double perimeter()
@Override
    public String name()
public class Rectangle implements Shape
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

