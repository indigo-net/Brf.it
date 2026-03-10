---
title: Elixir
---

# Elixir समर्थन

[English](../../languages/elixir.md) | [한국어](../../ko/languages/elixir.md) | [日本語](../../ja/languages/elixir.md) | [हिन्दी](elixir.md) | [Deutsch](../../de/languages/elixir.md)

## समर्थित एक्सटेंशन

- `.ex`
- `.exs`

## निष्कर्षण लक्ष्य

| तत्व | प्रकार | XML टैग | उदाहरण |
|------|--------|---------|--------|
| मॉड्यूल | `class` | `<type>` | `defmodule MyModule` |
| प्रोटोकॉल | `interface` | `<type>` | `defprotocol Printable` |
| प्रोटोकॉल कार्यान्वयन | `impl` | `<type>` | `defimpl Printable, for: Integer` |
| सार्वजनिक फ़ंक्शन | `function` | `<function>` | `def hello(name)` |
| निजी फ़ंक्शन | `function` | `<function>` | `defp validate(x)` |
| सार्वजनिक मैक्रो | `macro` | `<variable>` | `defmacro unless(cond, block)` |
| गार्ड | `function` | `<function>` | `defguard is_positive(x) when is_integer(x) and x > 0` |
| डेलिगेट | `function` | `<function>` | `defdelegate keys(map), to: Map` |
| स्ट्रक्ट | `struct` | `<type>` | `defstruct [:name, :email]` |
| टाइप स्पेक | `type` | `<type>` | `@spec add(integer(), integer()) :: integer()` |
| टाइप परिभाषा | `type` | `<type>` | `@type color :: :red \| :green \| :blue` |
| कॉलबैक | `type` | `<type>` | `@callback handle_event(term()) :: {:ok, term()}` |

## उदाहरण

### इनपुट

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

### आउटपुट (XML)

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

## नोट्स

### दृश्यता

- `def` (सार्वजनिक) और `defp` (निजी) दोनों फ़ंक्शन निष्कर्षित किए जाते हैं
- दृश्यता कीवर्ड सिग्नेचर टेक्स्ट में संरक्षित रहता है

### मॉड्यूल एट्रिब्यूट

- `@spec`, `@type`, `@typep`, `@opaque`, और `@callback` टाइप घोषणाओं के रूप में निष्कर्षित किए जाते हैं
- `@doc` और `@moduledoc` सिग्नेचर के रूप में निष्कर्षित नहीं किए जाते (ये दस्तावेज़ीकरण हैं)

### बॉडी हटाना

जब `--include-body` फ़्लैग का उपयोग नहीं किया जाता:

- फ़ंक्शन: `do...end` ब्लॉक हटा दिया जाता है, केवल सिग्नेचर लाइन रखी जाती है
- इनलाइन फ़ंक्शन: `, do: expr` हटा दिया जाता है
- मॉड्यूल/प्रोटोकॉल: बॉडी हटा दी जाती है, केवल घोषणा लाइन रखी जाती है
- टाइप स्पेक और स्ट्रक्ट परिभाषाएँ: जैसी हैं वैसी संरक्षित रहती हैं

### इम्पोर्ट निष्कर्षण

`--include-imports` के साथ, निम्नलिखित कैप्चर किए जाते हैं:

- `import Module`
- `alias Module`
- `use Module`
- `require Module`
