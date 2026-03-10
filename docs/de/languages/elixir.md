---
title: Elixir
---

# Elixir-Unterstützung

[English](../../languages/elixir.md) | [한국어](../../ko/languages/elixir.md) | [日本語](../../ja/languages/elixir.md) | [हिन्दी](../../hi/languages/elixir.md) | [Deutsch](elixir.md)

## Unterstützte Erweiterungen

- `.ex`
- `.exs`

## Extraktionsziele

| Element | Art | XML-Tag | Beispiel |
|---------|-----|---------|----------|
| Modul | `class` | `<type>` | `defmodule MyModule` |
| Protokoll | `interface` | `<type>` | `defprotocol Printable` |
| Protokollimplementierung | `impl` | `<type>` | `defimpl Printable, for: Integer` |
| Öffentliche Funktion | `function` | `<function>` | `def hello(name)` |
| Private Funktion | `function` | `<function>` | `defp validate(x)` |
| Öffentliches Makro | `macro` | `<variable>` | `defmacro unless(cond, block)` |
| Guard | `function` | `<function>` | `defguard is_positive(x) when is_integer(x) and x > 0` |
| Delegat | `function` | `<function>` | `defdelegate keys(map), to: Map` |
| Struct | `struct` | `<type>` | `defstruct [:name, :email]` |
| Typspezifikation | `type` | `<type>` | `@spec add(integer(), integer()) :: integer()` |
| Typdefinition | `type` | `<type>` | `@type color :: :red \| :green \| :blue` |
| Callback | `type` | `<type>` | `@callback handle_event(term()) :: {:ok, term()}` |

## Beispiel

### Eingabe

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

### Ausgabe (XML)

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

## Hinweise

### Sichtbarkeit

- Sowohl `def` (öffentlich) als auch `defp` (privat) Funktionen werden extrahiert
- Das Sichtbarkeitsschlüsselwort bleibt im Signaturtext erhalten

### Modulattribute

- `@spec`, `@type`, `@typep`, `@opaque` und `@callback` werden als Typdeklarationen extrahiert
- `@doc` und `@moduledoc` werden nicht als Signaturen extrahiert (sie dienen der Dokumentation)

### Entfernung des Funktionskörpers

Wenn das `--include-body`-Flag nicht verwendet wird:

- Funktionen: Der `do...end`-Block wird entfernt, nur die Signaturzeile bleibt erhalten
- Inline-Funktionen: `, do: expr` wird entfernt
- Module/Protokolle: Der Körper wird entfernt, nur die Deklarationszeile bleibt erhalten
- Typspezifikationen und Struct-Definitionen: werden unverändert beibehalten

### Import-Extraktion

Mit `--include-imports` werden folgende erfasst:

- `import Module`
- `alias Module`
- `use Module`
- `require Module`
