# Rust サポート

[English](../../languages/rust.md) | [한국어](../../ko/languages/rust.md) | [日本語](rust.md) | [हिन्दी](../../hi/languages/rust.md) | [Deutsch](../../de/languages/rust.md)

## 対応拡張子

- `.rs`

## 抽出対象

| 要素 | Kind | 例 |
|------|------|------|
| 関数 | `function` | `pub fn add(a: i32, b: i32) -> i32` |
| Async 関数 | `function` | `pub async fn fetch() -> Result<T>` |
| Unsafe 関数 | `function` | `pub unsafe fn dangerous()` |
| Const 関数 | `function` | `pub const fn compute() -> i32` |
| 構造体 | `struct` | `pub struct Point { x: f64, y: f64 }` |
| 列挙型 | `enum` | `pub enum Color { Red, Green, Blue }` |
| トレイト | `trait` | `pub trait Drawable { fn draw(&self); }` |
| 型エイリアス | `type` | `pub type Coordinate = (f64, f64)` |
| Impl ブロック | `impl` | `impl Point { ... }` |
| 定数 | `variable` | `pub const MAX_SIZE: usize = 1024` |
| スタティック | `variable` | `pub static COUNTER: AtomicUsize = ...` |
| モジュール | `namespace` | `pub mod utils { ... }` |
| マクロ | `macro` | `macro_rules! say_hello { ... }` |
| ユニオン | `struct` | `pub union IntOrFloat { i: i32, f: f32 }` |
| ドキュメントコメント | `doc` | `/// ドキュメント` |

## 例

### 入力

```rust
/// 2次元の点構造体。
pub struct Point {
    pub x: f64,
    pub y: f64,
}

impl Point {
    /// 指定された座標で新しいPointを作成します。
    pub fn new(x: f64, y: f64) -> Self {
        Point { x, y }
    }

    /// 原点からの距離を計算します。
    pub fn distance(&self) -> f64 {
        (self.x.powi(2) + self.y.powi(2)).sqrt()
    }
}

/// 許容される最大ポイント数。
pub const MAX_POINTS: usize = 1000;
```

### 出力 (XML)

```xml
<file path="point.rs" language="rust">
  <type kind="struct" line="2">
    <name>Point</name>
    <text>pub struct Point</text>
    <doc>2次元の点構造体。</doc>
  </type>
  <type kind="impl" line="8">
    <name>Point</name>
    <text>impl Point</text>
  </type>
  <function kind="function" line="10">
    <name>new</name>
    <text>pub fn new(x: f64, y: f64) -> Self</text>
    <doc>指定された座標で新しいPointを作成します。</doc>
  </function>
  <function kind="function" line="16">
    <name>distance</name>
    <text>pub fn distance(&self) -> f64</text>
    <doc>原点からの距離を計算します。</doc>
  </function>
  <variable kind="variable" line="22">
    <name>MAX_POINTS</name>
    <text>pub const MAX_POINTS: usize = 1000</text>
    <doc>許容される最大ポイント数。</doc>
  </variable>
</file>
```

## 注意事項

### 可視性

- すべての宣言が抽出されます（`pub` および private の両方）
- 可視性修飾子（`pub`、`pub(crate)`、`pub(super)`）はシグネチャにそのまま保持されます

### 関数修飾子

- `async`、`unsafe`、`const`、`extern` 関数はすべて kind `function` として分類されます
- 修飾子はシグネチャテキストに保持されます

### ジェネリクスとライフタイム

- ジェネリック型パラメータ（`<T>`、`<T: Clone>`）が完全に保持されます
- ライフタイム注釈（`'a`、`'static`）が完全に保持されます
- where 句もシグネチャに含まれます

### Impl ブロック

- `impl` ブロック自体とそのメソッドの両方が抽出されます
- `impl Trait for Type` パターンもキャプチャされます

### 本体の削除

`--include-body` フラグを使用しない場合:

- 関数/メソッド: 開き括弧 `{` 以降の本体を削除
- 構造体/列挙型/トレイト: 開き括弧 `{` 以降の本体を削除
- Impl ブロック: 開き括弧 `{` 以降の本体を削除
- Const/Static: 値式は保持

### ドキュメントコメント

- `///` および `//!` ドキュメントコメントのみが抽出されます
- 通常の `//` コメントは含まれません

### マクロ

- `macro_rules!` 定義が kind `macro` で抽出されます
- 手続き型マクロとマクロ呼び出しは抽出されません
