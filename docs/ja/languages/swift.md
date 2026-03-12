# Swift サポート

[English](../../languages/swift.md) | [한국어](../../ko/languages/swift.md) | [日本語](swift.md) | [हिन्दी](../../hi/languages/swift.md) | [Deutsch](../../de/languages/swift.md)

## 対応拡張子

- `.swift`

## 抽出対象

| 要素 | Kind | 例 |
|------|------|------|
| 関数 | `function` | `public func greet(name: String) -> String` |
| Async 関数 | `function` | `public func fetch() async throws -> Data` |
| 構造体 | `struct` | `public struct Point { let x: Double }` |
| クラス | `class` | `public class Vehicle { var speed: Int }` |
| 列挙型 | `enum` | `public enum Direction { case north }` |
| プロトコル | `interface` | `public protocol Drawable { func draw() }` |
| エクステンション | `type` | `extension Point: Drawable` |
| 型エイリアス | `type` | `public typealias Coordinate = (Double, Double)` |
| プロパティ (let/var) | `variable` | `public let PI = 3.14159` |
| イニシャライザ | `constructor` | `init(value: Int)` |
| デイニシャライザ | `destructor` | `deinit` |
| サブスクリプト | `method` | `subscript(index: Int) -> Int` |
| 演算子 | `function` | `prefix operator +++` |
| ドキュメントコメント | `doc` | `/// ドキュメント` |

## 例

### 入力

```swift
/// 2次元の点構造体。
public struct Point {
    public let x: Double
    public let y: Double

    /// 原点からの距離を計算します。
    public func distance() -> Double {
        return (x * x + y * y).squareRoot()
    }
}

/// 許容される最大ポイント数。
public let MAX_POINTS = 1000
```

### 出力 (XML)

```xml
<file path="point.swift" language="swift">
  <type kind="struct" line="2">
    <name>Point</name>
    <text>public struct Point</text>
    <doc>2次元の点構造体。</doc>
  </type>
  <function kind="function" line="8">
    <name>distance</name>
    <text>public func distance() -> Double</text>
    <doc>原点からの距離を計算します。</doc>
  </function>
  <variable kind="variable" line="14">
    <name>MAX_POINTS</name>
    <text>public let MAX_POINTS = 1000</text>
    <doc>許容される最大ポイント数。</doc>
  </variable>
</file>
```

## 注意事項

### 可視性

- すべての宣言が抽出されます（`public` および `internal`/`private` の両方）
- アクセス修飾子（`public`、`internal`、`private`、`fileprivate`、`open`）はシグネチャにそのまま保持されます

### 関数修飾子

- `async`、`throws`、`@discardableResult` 関数はすべて kind `function` として分類されます
- 修飾子はシグネチャテキストに保持されます

### ジェネリクス

- ジェネリック型パラメータ（`<T>`、`<T: Equatable>`）が完全に保持されます
- ジェネリック where 句もシグネチャに含まれます

### エクステンション

- `extension` ブロック自体とその内部メンバーの両方が抽出されます
- `extension Type: Protocol` 準拠パターンもキャプチャされます

### 本体の削除

`--include-body` フラグを使用しない場合:

- 関数/メソッド: 開き括弧 `{` 以降の本体を削除
- 構造体/クラス/列挙型: 開き括弧 `{` 以降の本体を削除
- エクステンション: 開き括弧 `{` 以降の本体を削除
- プロパティ (let/var): 値式は保持

`--include-private`を使用して非公開/unexportedシンボルを含める。

### ドキュメントコメント

- `///` ドキュメントコメントのみが抽出されます
- 通常の `//` コメントは含まれません
