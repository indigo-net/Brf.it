# C# サポート

[English](../../languages/csharp.md) | [한국어](../../ko/languages/csharp.md) | [日本語](csharp.md) | [हिन्दी](../../hi/languages/csharp.md) | [Deutsch](../../de/languages/csharp.md)

## 対応拡張子

- `.cs`

## 抽出対象

| 要素 | Kind | 例 |
|------|------|------|
| クラス | `class` | `public class Calculator` |
| 構造体 | `struct` | `public struct Point` |
| インターフェース | `interface` | `public interface IDrawable` |
| 列挙型 | `enum` | `public enum Color` |
| レコード | `record` | `public record Person(string Name, int Age)` |
| レコード構造体 | `struct` | `public record struct Measurement(double Value)` |
| デリゲート | `type` | `public delegate void Action<T>(T item)` |
| メソッド | `method` | `public int Add(int a, int b)` |
| コンストラクタ | `constructor` | `public Calculator()` |
| デストラクタ | `destructor` | `~Calculator()` |
| プロパティ | `variable` | `public string Name { get; set; }` |
| フィールド (static/const) | `variable` | `public const int MaxValue = 100` |
| イベント | `variable` | `public event EventHandler Changed` |
| インデクサー | `method` | `public int this[int index]` |
| 演算子 | `function` | `public static operator +(...)` |
| 変換演算子 | `function` | `public static implicit operator int(...)` |
| 名前空間 | `namespace` | `namespace MyApp` |
| 列挙型メンバー | `variable` | `Red, Green, Blue` |
| ドキュメントコメント | `doc` | `/// <summary>...</summary>` |

## 例

### 入力

```csharp
using System;

namespace MyApp
{
    /// <summary>基本的な数学演算のためのCalculatorクラス。</summary>
    public class Calculator
    {
        public const int MaxValue = 100;

        public Calculator() { }

        public int Add(int a, int b)
        {
            return a + b;
        }

        public string Name { get; set; }
    }

    public interface IService
    {
        void Execute();
    }

    public record Person(string Name, int Age);
}
```

### 出力 (XML)

```xml
<file path="Calculator.cs" language="csharp">
  <type kind="namespace" line="3">
    <name>MyApp</name>
    <text>namespace MyApp</text>
  </type>
  <type kind="class" line="6">
    <name>Calculator</name>
    <text>public class Calculator</text>
    <doc>基本的な数学演算のためのCalculatorクラス。</doc>
  </type>
  <variable kind="variable" line="8">
    <name>MaxValue</name>
    <text>public const int MaxValue = 100;</text>
  </variable>
  <function kind="constructor" line="10">
    <name>Calculator</name>
    <text>public Calculator()</text>
  </function>
  <function kind="method" line="12">
    <name>Add</name>
    <text>public int Add(int a, int b)</text>
  </function>
  <variable kind="variable" line="17">
    <name>Name</name>
    <text>public string Name { get; set; }</text>
  </variable>
  <type kind="interface" line="20">
    <name>IService</name>
    <text>public interface IService</text>
  </type>
  <function kind="method" line="22">
    <name>Execute</name>
    <text>void Execute();</text>
  </function>
  <type kind="record" line="25">
    <name>Person</name>
    <text>public record Person(string Name, int Age);</text>
  </type>
</file>
```

## 注意事項

### 可視性

- アクセス修飾子に関係なく、すべての宣言が抽出されます
- アクセス修飾子（`public`、`private`、`internal`、`protected`）はシグネチャに保存されます

### フィールド

- `static`、`const`、`static readonly` フィールドのみ抽出されます
- インスタンスフィールドはノイズ削減のために除外されます

### プロパティ

- 自動プロパティ（`{ get; set; }`）は完全に保存されます
- 式本体プロパティ（`=> expr`）はシグネチャモードで式が削除されます

### レコード

- `record` と `record class` は kind `record` に分類されます
- `record struct` は kind `struct` に分類されます

### 演算子

- 演算子オーバーロードは `operator+`、`operator==` のような合成名を持ちます
- 変換演算子は `implicit operator int`、`explicit operator string` のような名前を持ちます
- インデクサーは合成名 `this` を持ちます

### 本体の削除

`--include-body` フラグを使用しない場合：

- メソッド/コンストラクタ/デストラクタ: `{` 以降の本体を削除
- 式本体メンバー: `=>` と式を削除
- クラス/構造体/インターフェース/列挙型/レコード: `{` 以降の本体を削除
- プロパティ: 自動プロパティは保存、式本体プロパティは削除
- デリゲート: 本体なし、そのまま返却

`--include-private`を使用して非公開/unexportedシンボルを含める。

### ドキュメントコメント

- `///` XMLドキュメントコメントと `//` 行コメントの両方が抽出されます
- `/* */` ブロックコメントもキャプチャされます
