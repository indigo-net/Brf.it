---
layout: default
title: C#
parent: Language Guides
nav_order: 10
---

# C# Support

[English](csharp.md) | [한국어](../ko/languages/csharp.md) | [日本語](../ja/languages/csharp.md) | [हिन्दी](../hi/languages/csharp.md) | [Deutsch](../de/languages/csharp.md)

## Supported Extensions

- `.cs`

## Extraction Targets

| Element | Kind | Example |
|---------|------|---------|
| Class | `class` | `public class Calculator` |
| Struct | `struct` | `public struct Point` |
| Interface | `interface` | `public interface IDrawable` |
| Enum | `enum` | `public enum Color` |
| Record | `record` | `public record Person(string Name, int Age)` |
| Record Struct | `struct` | `public record struct Measurement(double Value)` |
| Delegate | `type` | `public delegate void Action<T>(T item)` |
| Method | `method` | `public int Add(int a, int b)` |
| Constructor | `constructor` | `public Calculator()` |
| Destructor | `destructor` | `~Calculator()` |
| Property | `variable` | `public string Name { get; set; }` |
| Field (static/const) | `variable` | `public const int MaxValue = 100` |
| Event | `variable` | `public event EventHandler Changed` |
| Indexer | `method` | `public int this[int index]` |
| Operator | `function` | `public static operator +(...)` |
| Conversion Operator | `function` | `public static implicit operator int(...)` |
| Namespace | `namespace` | `namespace MyApp` |
| Enum Member | `variable` | `Red, Green, Blue` |
| Doc Comment | `doc` | `/// <summary>...</summary>` |

## Example

### Input

```csharp
using System;

namespace MyApp
{
    /// <summary>Calculator class for basic math.</summary>
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

### Output (XML)

```xml
<file path="Calculator.cs" language="csharp">
  <type kind="namespace" line="3">
    <name>MyApp</name>
    <text>namespace MyApp</text>
  </type>
  <type kind="class" line="6">
    <name>Calculator</name>
    <text>public class Calculator</text>
    <doc>Calculator class for basic math.</doc>
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

## Notes

### Visibility

- All declarations are extracted regardless of access modifiers
- Access modifiers (`public`, `private`, `internal`, `protected`) are preserved in signatures

### Fields

- Only `static`, `const`, and `static readonly` fields are extracted
- Instance fields are excluded to reduce noise
- Extracted fields are classified as kind `variable`

### Properties

- Auto-properties (`{ get; set; }`) are fully preserved
- Expression-bodied properties (`=> expr`) have the expression removed in signature mode

### Records

- `record` and `record class` are classified as kind `record`
- `record struct` is classified as kind `struct`
- Parameter lists are preserved in signatures

### Operators

- Operator overloads get synthesized names like `operator+`, `operator==`
- Conversion operators get names like `implicit operator int`, `explicit operator string`
- Indexers get the synthesized name `this`

### Namespaces

- Both block-scoped (`namespace Foo { }`) and file-scoped (`namespace Foo;`) namespaces are extracted
- Nested namespace names (e.g., `A.B.C`) are preserved

### Body Removal

When `--include-body` flag is not used:

- Methods/Constructors/Destructors: body removed after opening brace `{`
- Expression-bodied members: `=>` and expression removed
- Classes/Structs/Interfaces/Enums/Records: body removed after opening brace `{`
- Properties: auto-properties preserved, expression-bodied properties stripped
- Delegates: no body, returned as-is

Use `--include-private` to include non-exported/private symbols.
- Abstract/interface methods ending with `;`: returned as-is

### Doc Comments

- Both `///` XML doc comments and `//` line comments are extracted
- `/* */` block comments are also captured
