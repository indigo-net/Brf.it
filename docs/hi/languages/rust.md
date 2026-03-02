# Rust समर्थन

[English](../../languages/rust.md) | [한국어](../../ko/languages/rust.md) | [日本語](../../ja/languages/rust.md) | [हिन्दी](rust.md) | [Deutsch](../../de/languages/rust.md)

## समर्थित एक्सटेंशन

- `.rs`

## निष्कर्षण लक्ष्य

| तत्व | Kind | उदाहरण |
|------|------|--------|
| फ़ंक्शन | `function` | `pub fn add(a: i32, b: i32) -> i32` |
| Async फ़ंक्शन | `function` | `pub async fn fetch() -> Result<T>` |
| Unsafe फ़ंक्शन | `function` | `pub unsafe fn dangerous()` |
| Const फ़ंक्शन | `function` | `pub const fn compute() -> i32` |
| स्ट्रक्ट | `struct` | `pub struct Point { x: f64, y: f64 }` |
| Enum | `enum` | `pub enum Color { Red, Green, Blue }` |
| Trait | `trait` | `pub trait Drawable { fn draw(&self); }` |
| Type उपनाम | `type` | `pub type Coordinate = (f64, f64)` |
| Impl ब्लॉक | `impl` | `impl Point { ... }` |
| Constant | `variable` | `pub const MAX_SIZE: usize = 1024` |
| Static | `variable` | `pub static COUNTER: AtomicUsize = ...` |
| Module | `namespace` | `pub mod utils { ... }` |
| Macro | `macro` | `macro_rules! say_hello { ... }` |
| Union | `struct` | `pub union IntOrFloat { i: i32, f: f32 }` |
| Doc टिप्पणी | `doc` | `/// दस्तावेज़ीकरण` |

## उदाहरण

### इनपुट

```rust
/// एक 2D बिंदु संरचना।
pub struct Point {
    pub x: f64,
    pub y: f64,
}

impl Point {
    /// दिए गए निर्देशांकों पर नया Point बनाता है।
    pub fn new(x: f64, y: f64) -> Self {
        Point { x, y }
    }

    /// मूल बिंदु से दूरी की गणना करता है।
    pub fn distance(&self) -> f64 {
        (self.x.powi(2) + self.y.powi(2)).sqrt()
    }
}

/// अधिकतम अनुमत बिंदु।
pub const MAX_POINTS: usize = 1000;
```

### आउटपुट (XML)

```xml
<file path="point.rs" language="rust">
  <type kind="struct" line="2">
    <name>Point</name>
    <text>pub struct Point</text>
    <doc>एक 2D बिंदु संरचना।</doc>
  </type>
  <type kind="impl" line="8">
    <name>Point</name>
    <text>impl Point</text>
  </type>
  <function kind="function" line="10">
    <name>new</name>
    <text>pub fn new(x: f64, y: f64) -> Self</text>
    <doc>दिए गए निर्देशांकों पर नया Point बनाता है।</doc>
  </function>
  <function kind="function" line="16">
    <name>distance</name>
    <text>pub fn distance(&self) -> f64</text>
    <doc>मूल बिंदु से दूरी की गणना करता है।</doc>
  </function>
  <variable kind="variable" line="22">
    <name>MAX_POINTS</name>
    <text>pub const MAX_POINTS: usize = 1000</text>
    <doc>अधिकतम अनुमत बिंदु।</doc>
  </variable>
</file>
```

## टिप्पणियाँ

### दृश्यता (Visibility)

- सभी घोषणाएं निकाली जाती हैं (`pub` और private दोनों)
- दृश्यता संशोधक (`pub`, `pub(crate)`, `pub(super)`) हस्ताक्षर में संरक्षित हैं

### फ़ंक्शन संशोधक

- `async`, `unsafe`, `const`, `extern` फ़ंक्शन सभी kind `function` के रूप में वर्गीकृत हैं
- संशोधक हस्ताक्षर पाठ में संरक्षित हैं

### Generics और Lifetimes

- जेनेरिक टाइप पैरामीटर (`<T>`, `<T: Clone>`) पूरी तरह से संरक्षित हैं
- लाइफटाइम एनोटेशन (`'a`, `'static`) पूरी तरह से संरक्षित हैं
- where क्लॉज हस्ताक्षर में शामिल हैं

### Impl ब्लॉक

- `impl` ब्लॉक स्वयं और इसकी विधियां दोनों निकाली जाती हैं
- `impl Trait for Type` पैटर्न कैप्चर किए जाते हैं

### बॉडी हटाना

जब `--include-body` फ्लैग का उपयोग नहीं किया जाता है:

- फ़ंक्शन/विधियां: खुलने वाले ब्रेस `{` के बाद बॉडी हटा दी जाती है
- स्ट्रक्ट/Enum/Trait: खुलने वाले ब्रेस `{` के बाद बॉडी हटा दी जाती है
- Impl ब्लॉक: खुलने वाले ब्रेस `{` के बाद बॉडी हटा दी जाती है
- Const/Static: मान अभिव्यक्ति संरक्षित है

### दस्तावेज़ टिप्पणियाँ

- केवल `///` और `//!` दस्तावेज़ टिप्पणियाँ निकाली जाती हैं
- सामान्य `//` टिप्पणियाँ शामिल नहीं हैं

### Macros

- `macro_rules!` परिभाषाएं kind `macro` के रूप में निकाली जाती हैं
- प्रक्रियात्मक मैक्रो और मैक्रो कॉल निकाले नहीं जाते हैं
