# Java Support

🌐 [English](java.md) | [한국어](java.ko.md) | [日本語](java.ja.md) | [हिन्दी](java.hi.md) | [Deutsch](java.de.md)

## Supported Extensions

- `.java`

## Extraction Targets

| Element | Kind | Example |
|---------|------|---------|
| Class | `class` | `public class User { ... }` |
| Interface | `interface` | `public interface Repository<T> { ... }` |
| Method | `method` | `public String getName() { ... }` |
| Constructor | `constructor` | `public User(String name) { ... }` |
| Enum | `enum` | `public enum Status { ... }` |
| Annotation | `annotation` | `public @interface Inject { ... }` |
| Record (Java 14+) | `record` | `public record Point(int x, int y) { ... }` |
| Comment | `doc` | `// Comment` or `/* Block */` |

## Example

### Input

```java
package com.example;

/**
 * User class represents a user in the system.
 */
public class User {
    private String name;

    public User(String name) {
        this.name = name;
    }

    public String getName() {
        return name;
    }

    private void internalMethod() {
        // Private method
    }
}

public interface Repository<T> {
    T findById(String id);
    void save(T entity);
}

public enum Status {
    PENDING, ACTIVE, COMPLETED
}

public record Point(int x, int y) {}
```

### Output (XML)

```xml
<file path="User.java" language="java">
  <signature kind="class" line="6">
    <name>User</name>
    <text>public class User</text>
  </signature>
  <signature kind="constructor" line="9">
    <name>User</name>
    <text>public User(String name)</text>
  </signature>
  <signature kind="method" line="13">
    <name>getName</name>
    <text>public String getName()</text>
  </signature>
  <signature kind="interface" line="22">
    <name>Repository</name>
    <text>public interface Repository&lt;T&gt;</text>
  </signature>
  <signature kind="method" line="23">
    <name>findById</name>
    <text>T findById(String id);</text>
  </signature>
  <signature kind="method" line="24">
    <name>save</name>
    <text>void save(T entity);</text>
  </signature>
  <signature kind="enum" line="27">
    <name>Status</name>
    <text>public enum Status</text>
  </signature>
  <signature kind="record" line="31">
    <name>Point</name>
    <text>public record Point(int x, int y)</text>
  </signature>
</file>
```

## Notes

### Visibility Filtering

- `public`, `protected`, package-private (default): extracted by default
- `private`: filtered out unless `--include-body` is used

### Generics Handling

Generic type parameters are included in the signature:

```java
public class Box<T extends Comparable<T>>  // Fully captured
public <U> U transform(Function<T, U> fn)  // Method type parameters included
```

### Annotations in Output

Method and class annotations are included in the signature text:

```java
@Override
public String toString()  // Signature includes @Override
```

### Record Support (Java 14+)

Records are extracted with their component parameters:

```java
public record User(String name, int age)  // Components preserved
```

### Inner/Nested Classes

All nested classes are extracted as separate signatures:

```java
public class Outer {
    public static class Nested { ... }  // Extracted separately
    public class Inner { ... }          // Also extracted
}
```

### Abstract Methods

Abstract methods in interfaces end with `;` (no body to strip):

```java
interface Foo {
    void bar();  // Captured as-is
}
```

### Body Removal

When `--include-body` flag is not used:

- Methods/Constructors: body removed after opening brace `{`
- Classes/Interfaces/Enums: body removed after opening brace `{`
- Abstract methods: kept as-is (end with `;`)

### Javadoc (Future Support)

- Current version: `//` and `/* */` comments above declarations captured as doc
- Future version: `/** */` Javadoc parsing planned

### Unsupported Elements

- Field declarations
- Static initializer blocks
- Anonymous classes
- Lambda expressions
