# Ruby サポート

[English](../../languages/ruby.md) | [한국어](../../ko/languages/ruby.md) | [日本語](ruby.md) | [हिन्दी](../../hi/languages/ruby.md) | [Deutsch](../../de/languages/ruby.md)

## サポートされている拡張子

- `.rb`

## 抽出対象

| 要素 | Kind | 例 |
|------|------|-----|
| メソッド | `method` | `def greet(name)` |
| クラスメソッド | `class_method` | `def self.create(attrs)` |
| クラス | `class` | `class User < ActiveRecord::Base` |
| モジュール | `module` | `module Authentication` |
| 定数 | `variable` | `MAX_RETRIES = 3` |
| YARDコメント | `doc` | `# 説明` |
| require | `import` | `require "json"` |
| require_relative | `import` | `require_relative "helpers"` |

## 例

### 入力

```ruby
require "json"
require_relative "helpers"

# システム内のユーザーを表します。
class User
  MAX_RETRIES = 3

  # 属性から新しいユーザーを作成します。
  # @param attrs [Hash] ユーザー属性
  def self.create(attrs)
    new(attrs).save
  end

  # ユーザーを初期化します。
  # @param name [String] ユーザーの名前
  def initialize(name)
    @name = name
  end

  # 他の人に挨拶します。
  # @param other [String] 相手の名前
  # @return [String] 挨拶メッセージ
  def greet(other)
    "Hello, #{other}! I'm #{@name}."
  end
end

module Authentication
  def authenticate(password)
    password == @secret
  end
end
```

### 出力 (XML)

```xml
<file path="example.rb" language="ruby">
  <class kind="class" line="5">
    <name>User</name>
    <text>class User</text>
  </class>
  <variable kind="variable" line="6">
    <name>MAX_RETRIES</name>
    <text>MAX_RETRIES = 3</text>
  </variable>
  <function kind="class_method" line="10">
    <name>create</name>
    <text>def self.create(attrs)</text>
  </function>
  <function kind="method" line="15">
    <name>initialize</name>
    <text>def initialize(name)</text>
  </function>
  <function kind="method" line="21">
    <name>greet</name>
    <text>def greet(other)</text>
  </function>
  <module kind="module" line="27">
    <name>Authentication</name>
    <text>module Authentication</text>
  </module>
  <function kind="method" line="28">
    <name>authenticate</name>
    <text>def authenticate(password)</text>
  </function>
</file>
```

## 注意事項

### 可視性

- 可視性（`public`、`protected`、`private`）に関係なく、すべてのメソッドが抽出されます
- インスタンスメソッド（`def foo`）とクラスメソッド（`def self.foo`）の両方がキャプチャされます

### メソッドの種類

- `method`: インスタンスメソッド宣言 (`def foo`)
- `class_method`: クラスレベルメソッド宣言 (`def self.foo`)

### 本体の除去

`--include-body`フラグを使用しない場合：

- メソッド: パラメータリストの閉じ括弧 `)` 以降の本体を除去（パラメータがない場合はメソッド名以降）
- クラス/モジュール: 宣言行のみ保持
- 定数: そのまま保持

### インポート抽出

- `require`と`require_relative`文は`--include-imports`フラグで抽出可能
- 形式: `require "json"` / `require_relative "helpers"`（完全な文が保持されます）

### ドキュメントコメント

- YARDスタイルコメント（`#`）がメソッド/クラスの直上にある場合に抽出されます
- 複数行の`#`コメントもサポートされています
- `=begin`...`=end`ブロックコメントも認識されます
