---
title: Elixir
---

# Elixir サポート

[English](../../languages/elixir.md) | [한국어](../../ko/languages/elixir.md) | [日本語](elixir.md) | [हिन्दी](../../hi/languages/elixir.md) | [Deutsch](../../de/languages/elixir.md)

## 対応拡張子

- `.ex`
- `.exs`

## 抽出対象

| 要素 | 種類 | XMLタグ | 例 |
|------|------|---------|-----|
| モジュール | `class` | `<type>` | `defmodule MyModule` |
| プロトコル | `interface` | `<type>` | `defprotocol Printable` |
| プロトコル実装 | `impl` | `<type>` | `defimpl Printable, for: Integer` |
| 公開関数 | `function` | `<function>` | `def hello(name)` |
| 非公開関数 | `function` | `<function>` | `defp validate(x)` |
| 公開マクロ | `macro` | `<variable>` | `defmacro unless(cond, block)` |
| ガード | `function` | `<function>` | `defguard is_positive(x) when is_integer(x) and x > 0` |
| デリゲート | `function` | `<function>` | `defdelegate keys(map), to: Map` |
| 構造体 | `struct` | `<type>` | `defstruct [:name, :email]` |
| 型スペック | `type` | `<type>` | `@spec add(integer(), integer()) :: integer()` |
| 型定義 | `type` | `<type>` | `@type color :: :red \| :green \| :blue` |
| コールバック | `type` | `<type>` | `@callback handle_event(term()) :: {:ok, term()}` |

## 例

### 入力

```elixir
defmodule Calculator do
  @type number :: integer() | float()
  @spec add(number(), number()) :: number()

  # Adds two numbers
  def add(a, b) do
    a + b
  end

  defp validate(x) when is_number(x) do
    :ok
  end

  defstruct [:value, :operation]
end
```

### 出力 (XML)

```xml
<file path="calculator.ex" language="elixir">
  <type>defmodule Calculator</type>
  <type>@type number :: integer() | float()</type>
  <type>@spec add(number(), number()) :: number()</type>
  <doc>Adds two numbers</doc>
  <function>def add(a, b)</function>
  <function>defp validate(x) when is_number(x)</function>
  <type>defstruct [:value, :operation]</type>
</file>
```

## 注意事項

### 可視性

- `def`（公開）と`defp`（非公開）の両方の関数が抽出されます
- 可視性キーワードはシグネチャテキストに保持されます

### モジュール属性

- `@spec`、`@type`、`@typep`、`@opaque`、`@callback`は型宣言として抽出されます
- `@doc`と`@moduledoc`はシグネチャとして抽出されません（ドキュメント用途）

### 本文の除去

`--include-body`フラグを使用しない場合：

- 関数：`do...end`ブロックが除去され、シグネチャ行のみが保持されます
- インライン関数：`, do: expr`が除去されます
- モジュール/プロトコル：本文が除去され、宣言行のみが保持されます
- 型スペックと構造体定義：そのまま保持されます

### インポートの抽出

`--include-imports`を使用すると、以下がキャプチャされます：

- `import Module`
- `alias Module`
- `use Module`
- `require Module`
