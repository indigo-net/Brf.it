# Swift-Unterstützung

[English](../../languages/swift.md) | [한국어](../../ko/languages/swift.md) | [日本語](../../ja/languages/swift.md) | [हिन्दी](../../hi/languages/swift.md) | [Deutsch](swift.md)

## Unterstützte Erweiterungen

- `.swift`

## Extraktionsziele

| Element | Kind | Beispiel |
|---------|------|----------|
| Funktion | `function` | `public func greet(name: String) -> String` |
| Async-Funktion | `function` | `public func fetch() async throws -> Data` |
| Struktur | `struct` | `public struct Point { let x: Double }` |
| Klasse | `class` | `public class Vehicle { var speed: Int }` |
| Aufzählung | `enum` | `public enum Direction { case north }` |
| Protokoll | `interface` | `public protocol Drawable { func draw() }` |
| Erweiterung | `type` | `extension Point: Drawable` |
| Typ-Alias | `type` | `public typealias Coordinate = (Double, Double)` |
| Eigenschaft (let/var) | `variable` | `public let PI = 3.14159` |
| Initialisierer | `constructor` | `init(value: Int)` |
| Deinitialisierer | `destructor` | `deinit` |
| Subscript | `method` | `subscript(index: Int) -> Int` |
| Operator | `function` | `prefix operator +++` |
| Dokumentationskommentar | `doc` | `/// Dokumentation` |

## Beispiel

### Eingabe

```swift
/// Eine 2D-Punktstruktur.
public struct Point {
    public let x: Double
    public let y: Double

    /// Berechnet die Entfernung vom Ursprung.
    public func distance() -> Double {
        return (x * x + y * y).squareRoot()
    }
}

/// Maximale Anzahl erlaubter Punkte.
public let MAX_POINTS = 1000
```

### Ausgabe (XML)

```xml
<file path="point.swift" language="swift">
  <type kind="struct" line="2">
    <name>Point</name>
    <text>public struct Point</text>
    <doc>Eine 2D-Punktstruktur.</doc>
  </type>
  <function kind="function" line="8">
    <name>distance</name>
    <text>public func distance() -> Double</text>
    <doc>Berechnet die Entfernung vom Ursprung.</doc>
  </function>
  <variable kind="variable" line="14">
    <name>MAX_POINTS</name>
    <text>public let MAX_POINTS = 1000</text>
    <doc>Maximale Anzahl erlaubter Punkte.</doc>
  </variable>
</file>
```

## Hinweise

### Sichtbarkeit

- Alle Deklarationen werden extrahiert (sowohl `public` als auch `internal`/`private`)
- Zugriffsmodifizierer (`public`, `internal`, `private`, `fileprivate`, `open`) werden in Signaturen beibehalten

### Funktionsmodifizierer

- `async`, `throws`, `@discardableResult` Funktionen werden alle als kind `function` klassifiziert
- Modifizierer werden im Signaturtext beibehalten

### Generics

- Generische Typparameter (`<T>`, `<T: Equatable>`) werden vollständig beibehalten
- Generische Where-Klauseln sind in Signaturen enthalten

### Erweiterungen

- Sowohl der `extension`-Block selbst als auch seine inneren Mitglieder werden extrahiert
- `extension Type: Protocol`-Konformitätsmuster werden erfasst

### Körperentfernung

Wenn das Flag `--include-body` nicht verwendet wird:

- Funktionen/Methoden: Körper nach der öffnenden geschweiften Klammer `{` entfernt
- Strukturen/Klassen/Aufzählungen: Körper nach der öffnenden geschweiften Klammer `{` entfernt
- Erweiterungen: Körper nach der öffnenden geschweiften Klammer `{` entfernt
- Eigenschaften (let/var): Wertausdruck wird beibehalten

Verwenden Sie `--include-private`, um nicht-exportierte/private Symbole einzubeziehen.

### Dokumentationskommentare

- Nur `///` Dokumentationskommentare werden extrahiert
- Normale `//` Kommentare sind nicht enthalten
