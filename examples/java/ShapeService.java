package com.example.shapes;

import java.util.List;
import java.util.ArrayList;
import java.util.Optional;

/**
 * Represents a geometric shape with area and perimeter.
 */
interface Shape {
    double area();
    double perimeter();
    String name();
}

/**
 * A circle defined by its radius.
 */
class Circle implements Shape {
    private final double radius;

    public Circle(double radius) {
        if (radius < 0) throw new IllegalArgumentException("Negative radius");
        this.radius = radius;
    }

    @Override
    public double area() {
        return Math.PI * radius * radius;
    }

    @Override
    public double perimeter() {
        return 2 * Math.PI * radius;
    }

    @Override
    public String name() {
        return "Circle(r=" + radius + ")";
    }
}

/**
 * A rectangle defined by width and height.
 */
class Rectangle implements Shape {
    private final double width;
    private final double height;

    public Rectangle(double width, double height) {
        this.width = width;
        this.height = height;
    }

    @Override
    public double area() {
        return width * height;
    }

    @Override
    public double perimeter() {
        return 2 * (width + height);
    }

    @Override
    public String name() {
        return "Rectangle(" + width + "x" + height + ")";
    }
}

/**
 * Service for managing and computing shape properties.
 */
public class ShapeService {
    private final List<Shape> shapes = new ArrayList<>();

    /** Add a shape to the collection. */
    public void addShape(Shape shape) {
        shapes.add(shape);
    }

    /** Get the total area of all shapes. */
    public double totalArea() {
        return shapes.stream().mapToDouble(Shape::area).sum();
    }

    /** Find the shape with the largest area. */
    public Optional<Shape> largestShape() {
        return shapes.stream()
            .max((a, b) -> Double.compare(a.area(), b.area()));
    }
}
