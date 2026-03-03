# Swift Support

[English](swift.md) | [한국어](../ko/languages/swift.md) | [日本語](../ja/languages/swift.md) | [हिन्दी](../hi/languages/swift.md) | [Deutsch](../de/languages/swift.md)

## Supported Extensions

- `.swift`

## Extraction Targets

| Element | Kind | Example |
|---------|------|---------|
| Function | `function` | `public func greet(name: String) -> String` |
| Async Function | `function` | `public func fetch() async throws -> Data` |
| Struct | `struct` | `public struct Point { let x: Double }` |
| Class | `class` | `public class Vehicle { var speed: Int }` |
| Enum | `enum` | `public enum Direction { case north }` |
| Protocol | `interface` | `public protocol Drawable { func draw() }` |
| Extension | `type` | `extension Point: Drawable` |
| Type Alias | `type` | `public typealias Coordinate = (Double, Double)` |
| Property (let/var) | `variable` | `public let PI = 3.14159` |
| Initializer | `constructor` | `init(value: Int)` |
| Deinitializer | `destructor` | `deinit` |
| Subscript | `method` | `subscript(index: Int) -> Int` |
| Operator | `function` | `prefix operator +++` |
| Doc Comment | `doc` | `/// Documentation` |

## Example

### Input

```swift
/// A 2D point structure.
public struct Point {
    public let x: Double
    public let y: Double

    /// Calculates the distance from the origin.
    public func distance() -> Double {
        return (x * x + y * y).squareRoot()
    }
}

/// Maximum number of points allowed.
public let MAX_POINTS = 1000
```

### Output (XML)

```xml
<file path="point.swift" language="swift">
  <type kind="struct" line="2">
    <name>Point</name>
    <text>public struct Point</text>
    <doc>A 2D point structure.</doc>
  </type>
  <function kind="function" line="8">
    <name>distance</name>
    <text>public func distance() -> Double</text>
    <doc>Calculates the distance from the origin.</doc>
  </function>
  <variable kind="variable" line="14">
    <name>MAX_POINTS</name>
    <text>public let MAX_POINTS = 1000</text>
    <doc>Maximum number of points allowed.</doc>
  </variable>
</file>
```

## Notes

### Visibility

- All declarations are extracted (both `public` and `internal`/`private`)
- Access modifiers (`public`, `internal`, `private`, `fileprivate`, `open`) are preserved in signatures

### Function Modifiers

- `async`, `throws`, `@discardableResult` functions are all classified as kind `function`
- Modifiers are preserved in the signature text

### Generics

- Generic type parameters (`<T>`, `<T: Equatable>`) are fully preserved
- Generic where clauses are included in signatures

### Extensions

- Both the `extension` block itself and its inner members are extracted
- `extension Type: Protocol` conformance patterns are captured

### Body Removal

When `--include-body` flag is not used:

- Functions/Methods: body removed after opening brace `{`
- Structs/Classes/Enums: body removed after opening brace `{`
- Extensions: body removed after opening brace `{`
- Properties (let/var): value expression is preserved

### Doc Comments

- Only `///` doc comments are extracted
- Regular `//` comments are not included
