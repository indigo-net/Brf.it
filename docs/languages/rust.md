# Rust Support

[English](rust.md) | [한국어](../ko/languages/rust.md) | [日本語](../ja/languages/rust.md) | [हिन्दी](../hi/languages/rust.md) | [Deutsch](../de/languages/rust.md)

## Supported Extensions

- `.rs`

## Extraction Targets

| Element | Kind | Example |
|---------|------|---------|
| Function | `function` | `pub fn add(a: i32, b: i32) -> i32` |
| Async Function | `function` | `pub async fn fetch() -> Result<T>` |
| Unsafe Function | `function` | `pub unsafe fn dangerous()` |
| Const Function | `function` | `pub const fn compute() -> i32` |
| Struct | `struct` | `pub struct Point { x: f64, y: f64 }` |
| Enum | `enum` | `pub enum Color { Red, Green, Blue }` |
| Trait | `trait` | `pub trait Drawable { fn draw(&self); }` |
| Type Alias | `type` | `pub type Coordinate = (f64, f64)` |
| Impl Block | `impl` | `impl Point { ... }` |
| Const | `variable` | `pub const MAX_SIZE: usize = 1024` |
| Static | `variable` | `pub static COUNTER: AtomicUsize = ...` |
| Module | `namespace` | `pub mod utils { ... }` |
| Macro | `macro` | `macro_rules! say_hello { ... }` |
| Union | `struct` | `pub union IntOrFloat { i: i32, f: f32 }` |
| Doc Comment | `doc` | `/// Documentation` |

## Example

### Input

```rust
/// A 2D point structure.
pub struct Point {
    pub x: f64,
    pub y: f64,
}

impl Point {
    /// Creates a new Point at the given coordinates.
    pub fn new(x: f64, y: f64) -> Self {
        Point { x, y }
    }

    /// Calculates the distance from the origin.
    pub fn distance(&self) -> f64 {
        (self.x.powi(2) + self.y.powi(2)).sqrt()
    }
}

/// Maximum number of points allowed.
pub const MAX_POINTS: usize = 1000;
```

### Output (XML)

```xml
<file path="point.rs" language="rust">
  <type kind="struct" line="2">
    <name>Point</name>
    <text>pub struct Point</text>
    <doc>A 2D point structure.</doc>
  </type>
  <type kind="impl" line="8">
    <name>Point</name>
    <text>impl Point</text>
  </type>
  <function kind="function" line="10">
    <name>new</name>
    <text>pub fn new(x: f64, y: f64) -> Self</text>
    <doc>Creates a new Point at the given coordinates.</doc>
  </function>
  <function kind="function" line="16">
    <name>distance</name>
    <text>pub fn distance(&self) -> f64</text>
    <doc>Calculates the distance from the origin.</doc>
  </function>
  <variable kind="variable" line="22">
    <name>MAX_POINTS</name>
    <text>pub const MAX_POINTS: usize = 1000</text>
    <doc>Maximum number of points allowed.</doc>
  </variable>
</file>
```

## Notes

### Visibility

- All declarations are extracted (both `pub` and private)
- Visibility modifiers (`pub`, `pub(crate)`, `pub(super)`) are preserved in signatures

### Function Modifiers

- `async`, `unsafe`, `const`, `extern` functions are all classified as kind `function`
- Modifiers are preserved in the signature text

### Generics and Lifetimes

- Generic type parameters (`<T>`, `<T: Clone>`) are fully preserved
- Lifetime annotations (`'a`, `'static`) are fully preserved
- Where clauses are included in signatures

### Impl Blocks

- Both the `impl` block itself and its methods are extracted
- `impl Trait for Type` patterns are captured

### Body Removal

When `--include-body` flag is not used:

- Functions/Methods: body removed after opening brace `{`
- Structs/Enums/Traits: body removed after opening brace `{`
- Impl blocks: body removed after opening brace `{`
- Const/Static: value expression is preserved

### Doc Comments

- Only `///` and `//!` doc comments are extracted
- Regular `//` comments are not included

### Macros

- `macro_rules!` definitions are extracted with kind `macro`
- Procedural macros and macro invocations are not extracted
