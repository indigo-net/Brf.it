---
layout: default
title: Elixir
parent: Language Guides
nav_order: 17
---

# Elixir Support

[English](elixir.md) | [한국어](../ko/languages/elixir.md) | [日本語](../ja/languages/elixir.md) | [हिन्दी](../hi/languages/elixir.md) | [Deutsch](../de/languages/elixir.md)

## Supported Extensions

- `.ex`
- `.exs`

## Extraction Targets

| Element | Kind | XML Tag | Example |
|---------|------|---------|---------|
| Module | `class` | `<type>` | `defmodule MyModule` |
| Protocol | `interface` | `<type>` | `defprotocol Printable` |
| Protocol Implementation | `impl` | `<type>` | `defimpl Printable, for: Integer` |
| Public Function | `function` | `<function>` | `def hello(name)` |
| Private Function | `function` | `<function>` | `defp validate(x)` |
| Public Macro | `macro` | `<variable>` | `defmacro unless(cond, block)` |
| Guard | `function` | `<function>` | `defguard is_positive(x) when is_integer(x) and x > 0` |
| Delegate | `function` | `<function>` | `defdelegate keys(map), to: Map` |
| Struct | `struct` | `<type>` | `defstruct [:name, :email]` |
| Type Spec | `type` | `<type>` | `@spec add(integer(), integer()) :: integer()` |
| Type Definition | `type` | `<type>` | `@type color :: :red \| :green \| :blue` |
| Callback | `type` | `<type>` | `@callback handle_event(term()) :: {:ok, term()}` |

## Example

### Input

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

### Output (XML)

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

## Notes

### Visibility

- Both `def` (public) and `defp` (private) functions are extracted
- Visibility keyword is preserved in signature text

### Module Attributes

- `@spec`, `@type`, `@typep`, `@opaque`, and `@callback` are extracted as type declarations
- `@doc` and `@moduledoc` are not extracted as signatures (they are documentation)

### Body Removal

When `--include-body` flag is not used:

- Functions: `do...end` block is stripped, keeping only the signature line
- Inline functions: `, do: expr` is stripped
- Modules/protocols: body is stripped, keeping the declaration line
- Type specs and struct definitions: preserved as-is

### Import Extraction

With `--include-imports`, the following are captured:

- `import Module`
- `alias Module`
- `use Module`
- `require Module`
