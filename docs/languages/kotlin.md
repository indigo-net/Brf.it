---
layout: default
title: Kotlin
parent: Language Guides
nav_order: 9
---

# Kotlin Support

[English](kotlin.md) | [한국어](../ko/languages/kotlin.md) | [日本語](../ja/languages/kotlin.md) | [हिन्दी](../hi/languages/kotlin.md) | [Deutsch](../de/languages/kotlin.md)

## Supported Extensions

- `.kt`
- `.kts`

## Extraction Targets

| Element | Kind | Example |
|---------|------|---------|
| Function | `function` | `fun add(a: Int, b: Int): Int` |
| Suspend Function | `function` | `suspend fun fetchData(url: String): String` |
| Extension Function | `function` | `fun String.isEmail(): Boolean` |
| Class | `class` | `class User(val name: String)` |
| Data Class | `class` | `data class Point(val x: Double, val y: Double)` |
| Sealed Class | `class` | `sealed class Result<out T>` |
| Enum Class | `enum` | `enum class Color { RED, GREEN, BLUE }` |
| Interface | `interface` | `interface Repository<T>` |
| Object | `class` | `object AppConfig` |
| Companion Object | `class` | `companion object Factory` |
| Property (val/var) | `variable` | `val MAX_SIZE = 100` |
| Type Alias | `type` | `typealias Handler<T> = (T) -> Unit` |
| Enum Entry | `variable` | `RED("#FF0000")` |
| Secondary Constructor | `constructor` | `constructor(name: String)` |
| Doc Comment | `doc` | `/** Documentation */` |

## Example

### Input

```kotlin
/** A user data class for API responses. */
data class User(
    val id: Long,
    val name: String,
    val email: String
) {
    fun isValid(): Boolean = email.contains("@")
}

/** Repository interface for user operations. */
interface UserRepository {
    suspend fun getUser(id: Long): User?
    fun save(user: User): Boolean
}

val DEFAULT_TIMEOUT: Long = 5000L
```

### Output (XML)

```xml
<file path="user.kt" language="kotlin">
  <type kind="class" line="2">
    <name>User</name>
    <text>data class User(
    val id: Long,
    val name: String,
    val email: String
)</text>
    <doc>A user data class for API responses.</doc>
  </type>
  <function kind="function" line="7">
    <name>isValid</name>
    <text>fun isValid(): Boolean = email.contains("@")</text>
  </function>
  <type kind="interface" line="11">
    <name>UserRepository</name>
    <text>interface UserRepository</text>
    <doc>Repository interface for user operations.</doc>
  </type>
  <function kind="function" line="12">
    <name>getUser</name>
    <text>suspend fun getUser(id: Long): User?</text>
  </function>
  <function kind="function" line="13">
    <name>save</name>
    <text>fun save(user: User): Boolean</text>
  </function>
  <variable kind="variable" line="16">
    <name>DEFAULT_TIMEOUT</name>
    <text>val DEFAULT_TIMEOUT: Long = 5000L</text>
  </variable>
</file>
```

## Notes

### Visibility

- All declarations are extracted (Kotlin defaults to `public`)
- Access modifiers (`public`, `internal`, `private`, `protected`) are preserved in signatures

### Function Modifiers

- `suspend`, `inline`, `infix`, `operator`, `tailrec` functions are all classified as kind `function`
- Modifiers are preserved in the signature text
- Single-expression functions (`fun double(x: Int) = x * 2`) are fully preserved

### Generics

- Generic type parameters (`<T>`, `<T : Comparable<T>>`) are fully preserved
- `where` clauses and variance annotations (`in`, `out`) are included in signatures
- `reified` type parameters are preserved

### Classes

- `data class`, `sealed class`, `abstract class`, `open class`, `inner class`, `annotation class`, `value class` are all classified as kind `class`
- `enum class` is classified as kind `enum`
- `interface` and `sealed interface` are classified as kind `interface`

### Objects

- `object` declarations (singletons) are classified as kind `class`
- `companion object` blocks are extracted; unnamed companions get the synthetic name "Companion"

### Body Removal

When `--include-body` flag is not used:

- Functions/Methods: body removed after opening brace `{`
- Single-expression functions: fully preserved (the expression is part of the signature)
- Classes/Interfaces/Enums: body removed after opening brace `{`
- Properties (val/var): value expression is preserved
- Type aliases: fully preserved

### Doc Comments

- Both `/** ... */` (KDoc) and `//` line comments are extracted
- KDoc comments are associated with the following declaration
