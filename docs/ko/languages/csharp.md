# C# 지원

[English](../../languages/csharp.md) | [한국어](csharp.md) | [日本語](../../ja/languages/csharp.md) | [हिन्दी](../../hi/languages/csharp.md) | [Deutsch](../../de/languages/csharp.md)

## 지원 확장자

- `.cs`

## 추출 대상

| 요소 | Kind | 예시 |
|------|------|------|
| 클래스 | `class` | `public class Calculator` |
| 구조체 | `struct` | `public struct Point` |
| 인터페이스 | `interface` | `public interface IDrawable` |
| 열거형 | `enum` | `public enum Color` |
| 레코드 | `record` | `public record Person(string Name, int Age)` |
| 레코드 구조체 | `struct` | `public record struct Measurement(double Value)` |
| 대리자 | `type` | `public delegate void Action<T>(T item)` |
| 메서드 | `method` | `public int Add(int a, int b)` |
| 생성자 | `constructor` | `public Calculator()` |
| 소멸자 | `destructor` | `~Calculator()` |
| 프로퍼티 | `variable` | `public string Name { get; set; }` |
| 필드 (static/const) | `variable` | `public const int MaxValue = 100` |
| 이벤트 | `variable` | `public event EventHandler Changed` |
| 인덱서 | `method` | `public int this[int index]` |
| 연산자 | `function` | `public static operator +(...)` |
| 변환 연산자 | `function` | `public static implicit operator int(...)` |
| 네임스페이스 | `namespace` | `namespace MyApp` |
| 열거형 멤버 | `variable` | `Red, Green, Blue` |
| 문서 주석 | `doc` | `/// <summary>...</summary>` |

## 예시

### 입력

```csharp
using System;

namespace MyApp
{
    /// <summary>기본 수학 연산을 위한 Calculator 클래스.</summary>
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

### 출력 (XML)

```xml
<file path="Calculator.cs" language="csharp">
  <type kind="namespace" line="3">
    <name>MyApp</name>
    <text>namespace MyApp</text>
  </type>
  <type kind="class" line="6">
    <name>Calculator</name>
    <text>public class Calculator</text>
    <doc>기본 수학 연산을 위한 Calculator 클래스.</doc>
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

## 참고 사항

### 가시성

- 접근 제한자에 관계없이 모든 선언이 추출됩니다
- 접근 제한자(`public`, `private`, `internal`, `protected`)는 시그니처에 보존됩니다

### 필드

- `static`, `const`, `static readonly` 필드만 추출됩니다
- 인스턴스 필드는 노이즈 감소를 위해 제외됩니다

### 프로퍼티

- 자동 프로퍼티(`{ get; set; }`)는 전체가 보존됩니다
- 식 본문 프로퍼티(`=> expr`)는 시그니처 모드에서 식이 제거됩니다

### 레코드

- `record`와 `record class`는 kind `record`로 분류됩니다
- `record struct`는 kind `struct`로 분류됩니다

### 연산자

- 연산자 오버로드는 `operator+`, `operator==` 같은 합성 이름을 가집니다
- 변환 연산자는 `implicit operator int`, `explicit operator string` 같은 이름을 가집니다
- 인덱서는 합성 이름 `this`를 가집니다

### 본문 제거

`--include-body` 플래그를 사용하지 않을 때:

- 메서드/생성자/소멸자: `{` 이후 본문 제거
- 식 본문 멤버: `=>` 및 식 제거
- 클래스/구조체/인터페이스/열거형/레코드: `{` 이후 본문 제거
- 프로퍼티: 자동 프로퍼티는 보존, 식 본문 프로퍼티는 제거
- 대리자: 본문 없음, 그대로 반환

### 문서 주석

- `///` XML 문서 주석과 `//` 줄 주석 모두 추출됩니다
- `/* */` 블록 주석도 캡처됩니다
