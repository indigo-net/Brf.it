# Swift सपोर्ट

[English](../../languages/swift.md) | [한국어](../../ko/languages/swift.md) | [日本語](../../ja/languages/swift.md) | [हिन्दी](swift.md) | [Deutsch](../../de/languages/swift.md)

## समर्थित एक्सटेंशन

- `.swift`

## निष्कर्षण लक्ष्य

| तत्व | Kind | उदाहरण |
|------|------|--------|
| फ़ंक्शन | `function` | `public func greet(name: String) -> String` |
| Async फ़ंक्शन | `function` | `public func fetch() async throws -> Data` |
| स्ट्रक्ट | `struct` | `public struct Point { let x: Double }` |
| क्लास | `class` | `public class Vehicle { var speed: Int }` |
| Enum | `enum` | `public enum Direction { case north }` |
| प्रोटोकॉल | `interface` | `public protocol Drawable { func draw() }` |
| एक्सटेंशन | `type` | `extension Point: Drawable` |
| Type उपनाम | `type` | `public typealias Coordinate = (Double, Double)` |
| प्रॉपर्टी (let/var) | `variable` | `public let PI = 3.14159` |
| इनिशियलाइज़र | `constructor` | `init(value: Int)` |
| डीइनिशियलाइज़र | `destructor` | `deinit` |
| सबस्क्रिप्ट | `method` | `subscript(index: Int) -> Int` |
| ऑपरेटर | `function` | `prefix operator +++` |
| Doc टिप्पणी | `doc` | `/// दस्तावेज़ीकरण` |

## उदाहरण

### इनपुट

```swift
/// एक 2D बिंदु संरचना।
public struct Point {
    public let x: Double
    public let y: Double

    /// मूल बिंदु से दूरी की गणना करता है।
    public func distance() -> Double {
        return (x * x + y * y).squareRoot()
    }
}

/// अधिकतम अनुमत बिंदु।
public let MAX_POINTS = 1000
```

### आउटपुट (XML)

```xml
<file path="point.swift" language="swift">
  <type kind="struct" line="2">
    <name>Point</name>
    <text>public struct Point</text>
    <doc>एक 2D बिंदु संरचना।</doc>
  </type>
  <function kind="function" line="8">
    <name>distance</name>
    <text>public func distance() -> Double</text>
    <doc>मूल बिंदु से दूरी की गणना करता है।</doc>
  </function>
  <variable kind="variable" line="14">
    <name>MAX_POINTS</name>
    <text>public let MAX_POINTS = 1000</text>
    <doc>अधिकतम अनुमत बिंदु।</doc>
  </variable>
</file>
```

## टिप्पणियाँ

### दृश्यता (Visibility)

- सभी घोषणाएं निकाली जाती हैं (`public` और `internal`/`private` दोनों)
- एक्सेस संशोधक (`public`, `internal`, `private`, `fileprivate`, `open`) हस्ताक्षर में संरक्षित हैं

### फ़ंक्शन संशोधक

- `async`, `throws`, `@discardableResult` फ़ंक्शन सभी kind `function` के रूप में वर्गीकृत हैं
- संशोधक हस्ताक्षर पाठ में संरक्षित हैं

### Generics

- जेनेरिक टाइप पैरामीटर (`<T>`, `<T: Equatable>`) पूरी तरह से संरक्षित हैं
- जेनेरिक where क्लॉज हस्ताक्षर में शामिल हैं

### एक्सटेंशन

- `extension` ब्लॉक स्वयं और इसके आंतरिक सदस्य दोनों निकाले जाते हैं
- `extension Type: Protocol` अनुरूपता पैटर्न कैप्चर किए जाते हैं

### बॉडी हटाना

जब `--include-body` फ्लैग का उपयोग नहीं किया जाता है:

- फ़ंक्शन/विधियां: खुलने वाले ब्रेस `{` के बाद बॉडी हटा दी जाती है
- स्ट्रक्ट/क्लास/Enum: खुलने वाले ब्रेस `{` के बाद बॉडी हटा दी जाती है
- एक्सटेंशन: खुलने वाले ब्रेस `{` के बाद बॉडी हटा दी जाती है
- प्रॉपर्टी (let/var): मान अभिव्यक्ति संरक्षित है

`--include-private` का उपयोग करके गैर-निर्यातित/निजी सिंबल शामिल करें।

### दस्तावेज़ टिप्पणियाँ

- केवल `///` दस्तावेज़ टिप्पणियाँ निकाली जाती हैं
- सामान्य `//` टिप्पणियाँ शामिल नहीं हैं
