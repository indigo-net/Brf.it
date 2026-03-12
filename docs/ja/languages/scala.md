# Scala サポート

[English](../../languages/scala.md) | [한국어](../../ko/languages/scala.md) | [日本語](scala.md) | [हिन्दी](../../hi/languages/scala.md) | [Deutsch](../../de/languages/scala.md)

## サポートされている拡張子

- `.scala`
- `.sc`

## 抽出対象

| 要素 | Kind | XML Tag | 例 |
|------|------|---------|-----|
| メソッド（本体あり） | `method` | `<function>` | `def add(a: Int, b: Int): Int` |
| メソッド（抽象） | `method` | `<function>` | `def greet(name: String): String` |
| クラス | `class` | `<type>` | `class Person(val name: String)` |
| トレイト | `trait` | `<type>` | `trait Greeter` |
| オブジェクト | `class` | `<type>` | `object MathUtils` |
| val | `variable` | `<variable>` | `val PI: Double = 3.14159` |
| var | `variable` | `<variable>` | `var count: Int = 0` |
| 型エイリアス | `type` | `<type>` | `type StringList = List[String]` |
| Enum（Scala 3） | `enum` | `<type>` | `enum Color` |
| Given（Scala 3） | `variable` | `<variable>` | `given ordering: Ordering[Int]` |

## 例

### 入力

```scala
// ユーザー管理
trait Greeter {
  def greet(name: String): String
}

class Person(val name: String) extends Greeter {
  def greet(name: String): String = s"Hello, $name"
}

object MathUtils {
  val PI: Double = 3.14159
  def add(a: Int, b: Int): Int = a + b
}

type StringList = List[String]
```

### 出力 (XML)

```xml
<file path="example.scala" language="scala">
  <type>trait Greeter</type>
  <function>def greet(name: String): String</function>
  <type>class Person(val name: String) extends Greeter</type>
  <function>def greet(name: String): String</function>
  <type>object MathUtils</type>
  <variable>val PI: Double = 3.14159</variable>
  <function>def add(a: Int, b: Int): Int</function>
  <type>type StringList = List[String]</type>
</file>
```

## 注意事項

### 可視性

- 可視性修飾子に関係なく、すべての宣言が抽出されます
- 可視性修飾子（`private`、`protected`）はシグネチャテキストに保持されます

### クラスのバリエーション

- `class`、`abstract class`、`case class`、`sealed class`、`implicit class`はすべてkind `class`に分類されます
- `trait`と`sealed trait`はkind `trait`に分類されます
- `object`（シングルトンおよびコンパニオン）はkind `class`に分類されます

### 本体の除去

`--include-body`フラグを使用しない場合：

- メソッド: `=`以降の本体が除去され、戻り値の型は保持されます
- クラス/トレイト/オブジェクト: `{ }`内の本体が除去され、宣言行のみ保持されます
- val/var: 値が保持されます（`lazy val`、`implicit val`を含む）
- 型エイリアス: 完全に保持

`--include-private`を使用して非公開/unexportedシンボルを含める。

### ジェネリクス

- ジェネリック型パラメータ`[A, B]`はシグネチャに完全に保持されます
- コンテキスト境界とビュー境界が含まれます

### Scala 3機能

- `enum`定義はkind `enum`に分類されます
- 名前付き`given`インスタンスはkind `variable`に分類されます
- `extension`メソッドは個別に`method`として抽出されます（extension宣言自体はキャプチャされません）
