# Rust-Unterstützung

[English](../../languages/rust.md) | [한국어](../../ko/languages/rust.md) | [日本語](../../ja/languages/rust.md) | [हिन्दी](../../hi/languages/rust.md) | [Deutsch](rust.md)

## Unterstützte Erweiterungen

- `.rs`

## Extraktionsziele

| Element | Kind | Beispiel |
|---------|------|----------|
| Funktion | `function` | `pub fn add(a: i32, b: i32) -> i32` |
| Async-Funktion | `function` | `pub async fn fetch() -> Result<T>` |
| Unsafe-Funktion | `function` | `pub unsafe fn dangerous()` |
| Const-Funktion | `function` | `pub const fn compute() -> i32` |
| Struktur | `struct` | `pub struct Point { x: f64, y: f64 }` |
| Aufzählung | `enum` | `pub enum Color { Red, Green, Blue }` |
| Trait | `trait` | `pub trait Drawable { fn draw(&self); }` |
| Typ-Alias | `type` | `pub type Coordinate = (f64, f64)` |
| Impl-Block | `impl` | `impl Point { ... }` |
| Konstante | `variable` | `pub const MAX_SIZE: usize = 1024` |
| Statische Variable | `variable` | `pub static COUNTER: AtomicUsize = ...` |
| Modul | `namespace` | `pub mod utils { ... }` |
| Makro | `macro` | `macro_rules! say_hello { ... }` |
| Union | `struct` | `pub union IntOrFloat { i: i32, f: f32 }` |
| Dokumentationskommentar | `doc` | `/// Dokumentation` |

## Beispiel

### Eingabe

```rust
/// Eine 2D-Punktstruktur.
pub struct Point {
    pub x: f64,
    pub y: f64,
}

impl Point {
    /// Erstellt einen neuen Point an den angegebenen Koordinaten.
    pub fn new(x: f64, y: f64) -> Self {
        Point { x, y }
    }

    /// Berechnet die Entfernung vom Ursprung.
    pub fn distance(&self) -> f64 {
        (self.x.powi(2) + self.y.powi(2)).sqrt()
    }
}

/// Maximale Anzahl erlaubter Punkte.
pub const MAX_POINTS: usize = 1000;
```

### Ausgabe (XML)

```xml
<file path="point.rs" language="rust">
  <type kind="struct" line="2">
    <name>Point</name>
    <text>pub struct Point</text>
    <doc>Eine 2D-Punktstruktur.</doc>
  </type>
  <type kind="impl" line="8">
    <name>Point</name>
    <text>impl Point</text>
  </type>
  <function kind="function" line="10">
    <name>new</name>
    <text>pub fn new(x: f64, y: f64) -> Self</text>
    <doc>Erstellt einen neuen Point an den angegebenen Koordinaten.</doc>
  </function>
  <function kind="function" line="16">
    <name>distance</name>
    <text>pub fn distance(&self) -> f64</text>
    <doc>Berechnet die Entfernung vom Ursprung.</doc>
  </function>
  <variable kind="variable" line="22">
    <name>MAX_POINTS</name>
    <text>pub const MAX_POINTS: usize = 1000</text>
    <doc>Maximale Anzahl erlaubter Punkte.</doc>
  </variable>
</file>
```

## Hinweise

### Sichtbarkeit

- Alle Deklarationen werden extrahiert (sowohl `pub` als auch private)
- Sichtbarkeitsmodifizierer (`pub`, `pub(crate)`, `pub(super)`) werden in Signaturen beibehalten

### Funktionsmodifizierer

- `async`, `unsafe`, `const`, `extern` Funktionen werden alle als kind `function` klassifiziert
- Modifizierer werden im Signaturtext beibehalten

### Generics und Lifetimes

- Generische Typparameter (`<T>`, `<T: Clone>`) werden vollständig beibehalten
- Lifetime-Annotationen (`'a`, `'static`) werden vollständig beibehalten
- Where-Klauseln sind in Signaturen enthalten

### Impl-Blöcke

- Sowohl der `impl`-Block selbst als auch seine Methoden werden extrahiert
- `impl Trait for Type`-Muster werden erfasst

### Körperentfernung

Wenn das Flag `--include-body` nicht verwendet wird:

- Funktionen/Methoden: Körper nach der öffnenden geschweiften Klammer `{` entfernt
- Structs/Enums/Traits: Körper nach der öffnenden geschweiften Klammer `{` entfernt
- Impl-Blöcke: Körper nach der öffnenden geschweiften Klammer `{` entfernt
- Const/Static: Wertausruck wird beibehalten

Verwenden Sie `--include-private`, um nicht-exportierte/private Symbole einzubeziehen.

### Dokumentationskommentare

- Nur `///` und `//!` Dokumentationskommentare werden extrahiert
- Normale `//` Kommentare sind nicht enthalten

### Makros

- `macro_rules!` Definitionen werden als kind `macro` extrahiert
- Prozedurale Makros und Makroaufrufe werden nicht extrahiert
