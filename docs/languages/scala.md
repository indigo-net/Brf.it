---
layout: default
title: Scala
parent: Language Guides
nav_order: 16
---

# Scala Support

[English](scala.md) | [한국어](../ko/languages/scala.md) | [日本語](../ja/languages/scala.md) | [हिन्दी](../hi/languages/scala.md) | [Deutsch](../de/languages/scala.md)

## Supported Extensions

- `.scala`
- `.sc`

## Extraction Targets

| Element | Kind | XML Tag | Example |
|---------|------|---------|---------|
| Method (with body) | `method` | `<function>` | `def add(a: Int, b: Int): Int` |
| Method (abstract) | `method` | `<function>` | `def greet(name: String): String` |
| Class | `class` | `<type>` | `class Person(val name: String)` |
| Trait | `trait` | `<type>` | `trait Greeter` |
| Object | `class` | `<type>` | `object MathUtils` |
| val | `variable` | `<variable>` | `val PI: Double = 3.14159` |
| var | `variable` | `<variable>` | `var count: Int = 0` |
| Type Alias | `type` | `<type>` | `type StringList = List[String]` |
| Enum (Scala 3) | `enum` | `<type>` | `enum Color` |
| Given (Scala 3) | `variable` | `<variable>` | `given ordering: Ordering[Int]` |

## Example

### Input

```scala
// User management
trait Greeter {
  def greet(name: String): String
}

class Person(val name: String) extends Greeter {
  def greet(name: String): String = s"Hello, $name"
}

object MathUtils {
  val PI: Double = 3.14159
  def add(a: Int, b: Int): Int = a + b
}

type StringList = List[String]
```

### Output (XML)

```xml
<file path="example.scala" language="scala">
  <type>trait Greeter</type>
  <function>def greet(name: String): String</function>
  <type>class Person(val name: String) extends Greeter</type>
  <function>def greet(name: String): String</function>
  <type>object MathUtils</type>
  <variable>val PI: Double = 3.14159</variable>
  <function>def add(a: Int, b: Int): Int</function>
  <type>type StringList = List[String]</type>
</file>
```

## Notes

### Visibility

- All declarations are extracted regardless of visibility modifiers
- Visibility modifiers (`private`, `protected`) are preserved in signature text

### Class Variants

- `class`, `abstract class`, `case class`, `sealed class`, `implicit class` are all classified as kind `class`
- `trait` and `sealed trait` are classified as kind `trait`
- `object` (singleton and companion) is classified as kind `class`

### Body Removal

When `--include-body` flag is not used:

- Methods: body after `=` is removed, preserving the return type
- Classes/Traits/Objects: body in `{ }` is stripped, keeping only the declaration line
- val/var: values are preserved (including `lazy val` and `implicit val`)
- Type aliases: fully preserved

### Generics

- Generic type parameters `[A, B]` are fully preserved in signatures
- Context bounds and view bounds are included

### Scala 3 Features

- `enum` definitions are classified as kind `enum`
- Named `given` instances are classified as kind `variable`
- `extension` methods are extracted individually as `method` (the extension declaration itself is not captured)
