---
title: Elixir
---

# Elixir 지원

[English](../../languages/elixir.md) | [한국어](elixir.md) | [日本語](../../ja/languages/elixir.md) | [हिन्दी](../../hi/languages/elixir.md) | [Deutsch](../../de/languages/elixir.md)

## 지원 확장자

- `.ex`
- `.exs`

## 추출 대상

| 요소 | 종류 | XML 태그 | 예시 |
|------|------|----------|------|
| 모듈 | `class` | `<type>` | `defmodule MyModule` |
| 프로토콜 | `interface` | `<type>` | `defprotocol Printable` |
| 프로토콜 구현 | `impl` | `<type>` | `defimpl Printable, for: Integer` |
| 공개 함수 | `function` | `<function>` | `def hello(name)` |
| 비공개 함수 | `function` | `<function>` | `defp validate(x)` |
| 공개 매크로 | `macro` | `<variable>` | `defmacro unless(cond, block)` |
| 가드 | `function` | `<function>` | `defguard is_positive(x) when is_integer(x) and x > 0` |
| 위임 | `function` | `<function>` | `defdelegate keys(map), to: Map` |
| 구조체 | `struct` | `<type>` | `defstruct [:name, :email]` |
| 타입 스펙 | `type` | `<type>` | `@spec add(integer(), integer()) :: integer()` |
| 타입 정의 | `type` | `<type>` | `@type color :: :red \| :green \| :blue` |
| 콜백 | `type` | `<type>` | `@callback handle_event(term()) :: {:ok, term()}` |

## 예시

### 입력

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

### 출력 (XML)

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

## 참고사항

### 가시성

- `def` (공개) 및 `defp` (비공개) 함수 모두 추출됩니다
- 가시성 키워드는 시그니처 텍스트에 보존됩니다

### 모듈 속성

- `@spec`, `@type`, `@typep`, `@opaque`, `@callback`은 타입 선언으로 추출됩니다
- `@doc` 및 `@moduledoc`은 시그니처로 추출되지 않습니다 (문서화 용도)

### 본문 제거

`--include-body` 플래그를 사용하지 않을 때:

- 함수: `do...end` 블록이 제거되고 시그니처 줄만 유지됩니다
- 인라인 함수: `, do: expr`이 제거됩니다
- 모듈/프로토콜: 본문이 제거되고 선언 줄만 유지됩니다
- 타입 스펙과 구조체 정의: 그대로 보존됩니다

### 임포트 추출

`--include-imports` 사용 시 다음이 캡처됩니다:

- `import Module`
- `alias Module`
- `use Module`
- `require Module`
