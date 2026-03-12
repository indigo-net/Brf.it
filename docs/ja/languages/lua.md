# Lua サポート

[English](../../languages/lua.md) | [한국어](../../ko/languages/lua.md) | [日本語](lua.md) | [हिन्दी](../../hi/languages/lua.md) | [Deutsch](../../de/languages/lua.md)

## サポートされている拡張子

- `.lua`

## 抽出対象

| 要素 | Kind | 例 |
|------|------|-----|
| グローバル関数 | `function` | `function globalFunc(a, b)` |
| ローカル関数 | `local_function` | `local function helper()` |
| モジュール関数 | `module_function` | `function M.create(name)` |
| メソッド | `method` | `function M:init(config)` |
| テーブル代入 | `variable` | `local M = {}` |
| 関数代入 | `variable` | `local foo = function() end` |
| LuaDocコメント | `doc` | `--- 説明` |
| 行コメント | `doc` | `-- 説明` |
| require() | `import` | `local json = require("json")` |

## 例

### 入力

```lua
local M = {}

--- 名前で人に挨拶します。
-- @param name string 人の名前
function M.greet(name)
    print("Hello, " .. name)
end

function M:init(config)
    self.config = config
end

local function helper()
    return 42
end

function globalFunc(a, b)
    return a + b
end

local json = require("json")
```

### 出力 (XML)

```xml
<file path="example.lua" language="lua">
  <variable kind="variable" line="1">
    <name>M</name>
    <text>local M = {}</text>
  </variable>
  <function kind="module_function" line="4">
    <name>greet</name>
    <text>function M.greet(name)</text>
  </function>
  <function kind="method" line="9">
    <name>init</name>
    <text>function M:init(config)</text>
  </function>
  <function kind="local_function" line="13">
    <name>helper</name>
    <text>local function helper()</text>
  </function>
  <function kind="function" line="17">
    <name>globalFunc</name>
    <text>function globalFunc(a, b)</text>
  </function>
</file>
```

## 注意事項

### 可視性

- すべての宣言が抽出されます（Luaにはアクセス修飾子がありません）
- `local`関数/変数もグローバルのものと一緒に含まれます

### 関数の種類

- `function`: グローバル関数宣言 (`function foo()`)
- `local_function`: ローカル関数宣言 (`local function foo()`)
- `module_function`: ドットインデックス関数 (`function M.foo()`)
- `method`: コロンインデックスメソッド (`function M:foo()`)

### 本体の除去

`--include-body`フラグを使用しない場合：

- 関数/メソッド: パラメータリストの閉じ括弧 `)` 以降の本体を除去
- テーブル代入 (`local M = {}`): そのまま保持
- 関数代入 (`local foo = function() end`): そのまま保持

### インポート抽出

- `require()`呼び出しは`--include-imports`フラグで抽出可能
- `--include-private`を使用して非公開/unexportedシンボルを含める
- 形式: `local json = require("json")`（完全な文が保持されます）

### ドキュメントコメント

- LuaDocコメント (`---`) と通常の行コメント (`--`) の両方が抽出されます
- ブロックコメント (`--[[ ... ]]`) もサポートされています
