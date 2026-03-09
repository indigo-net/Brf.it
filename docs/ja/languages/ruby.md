# Ruby サポート

[English](../../languages/ruby.md) | [한국어](../../ko/languages/ruby.md) | [日本語](ruby.md) | [हिन्दी](../../hi/languages/ruby.md) | [Deutsch](../../de/languages/ruby.md)

## サポートされている拡張子

- `.rb`

## 抽出対象

| 要素 | Kind | XML Tag | 例 |
|------|------|---------|-----|
| メソッド | `method` | `<function>` | `def greet(name)` |
| クラスメソッド | `method` | `<function>` | `def self.create(attrs)` |
| クラス | `class` | `<type>` | `class User < ActiveRecord::Base` |
| モジュール | `namespace` | `<type>` | `module Authentication` |
| 定数 (top-level) | `variable` | `<variable>` | `MAX_RETRIES = 3` |
| YARDコメント | `doc` | | `# 説明` |
| require | | `<imports>` | `require "json"` |
| require_relative | | `<imports>` | `require_relative "helpers"` |

## 例

### 入力

```ruby
require "json"
require_relative "helpers"

MAX_RETRIES = 3

# システム内のユーザーを表します。
class User
  # 属性から新しいユーザーを作成します。
  def self.create(attrs)
    new(attrs).save
  end

  # ユーザーを初期化します。
  def initialize(name)
    @name = name
  end

  # 他の人に挨拶します。
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
  <type>class User</type>
  <function>def self.create(attrs)</function>
  <function>def initialize(name)</function>
  <function>def greet(other)</function>
  <variable>MAX_RETRIES = 3</variable>
  <type>module Authentication</type>
  <function>def authenticate(password)</function>
</file>
```

## 注意事項

### 可視性

- 可視性（`public`、`protected`、`private`）に関係なく、すべてのメソッドが抽出されます
- インスタンスメソッド（`def foo`）とクラスメソッド（`def self.foo`）の両方がキャプチャされます

### メソッドの種類

- インスタンスメソッド（`def foo`）とクラスメソッド（`def self.foo`）の両方がkind `method`を使用します

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
