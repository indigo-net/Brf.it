# Kotlin-Unterstützung

[English](../../languages/kotlin.md) | [한국어](../../ko/languages/kotlin.md) | [日本語](../../ja/languages/kotlin.md) | [हिन्दी](../../hi/languages/kotlin.md) | [Deutsch](kotlin.md)

## Unterstützte Erweiterungen

- `.kt`
- `.kts`

## Extraktionsziele

| Element | Kind | Beispiel |
|---------|------|----------|
| Funktion | `function` | `fun add(a: Int, b: Int): Int` |
| Suspend-Funktion | `function` | `suspend fun fetchData(url: String): String` |
| Erweiterungsfunktion | `function` | `fun String.isEmail(): Boolean` |
| Klasse | `class` | `class User(val name: String)` |
| Datenklasse | `class` | `data class Point(val x: Double, val y: Double)` |
| Sealed-Klasse | `class` | `sealed class Result<out T>` |
| Enum-Klasse | `enum` | `enum class Color { RED, GREEN, BLUE }` |
| Interface | `interface` | `interface Repository<T>` |
| Objekt | `class` | `object AppConfig` |
| Companion-Objekt | `class` | `companion object Factory` |
| Eigenschaft (val/var) | `variable` | `val MAX_SIZE = 100` |
| Typ-Alias | `type` | `typealias Handler<T> = (T) -> Unit` |
| Enum-Eintrag | `variable` | `RED("#FF0000")` |
| Sekundärer Konstruktor | `constructor` | `constructor(name: String)` |
| Dokumentationskommentar | `doc` | `/** Dokumentation */` |

## Beispiel

### Eingabe

```kotlin
/** Eine Benutzer-Datenklasse für API-Antworten. */
data class User(
    val id: Long,
    val name: String,
    val email: String
) {
    fun isValid(): Boolean = email.contains("@")
}

/** Repository-Interface für Benutzeroperationen. */
interface UserRepository {
    suspend fun getUser(id: Long): User?
    fun save(user: User): Boolean
}

val DEFAULT_TIMEOUT: Long = 5000L
```

### Ausgabe (XML)

```xml
<file path="user.kt" language="kotlin">
  <type kind="class" line="2">
    <name>User</name>
    <text>data class User(
    val id: Long,
    val name: String,
    val email: String
)</text>
    <doc>Eine Benutzer-Datenklasse für API-Antworten.</doc>
  </type>
  <function kind="function" line="7">
    <name>isValid</name>
    <text>fun isValid(): Boolean = email.contains("@")</text>
  </function>
  <type kind="interface" line="11">
    <name>UserRepository</name>
    <text>interface UserRepository</text>
    <doc>Repository-Interface für Benutzeroperationen.</doc>
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

## Hinweise

### Sichtbarkeit

- Alle Deklarationen werden extrahiert (Kotlin ist standardmäßig `public`)
- Zugriffsmodifizierer (`public`, `internal`, `private`, `protected`) werden in Signaturen beibehalten

### Funktionsmodifizierer

- `suspend`, `inline`, `infix`, `operator`, `tailrec` Funktionen werden alle als kind `function` klassifiziert
- Modifizierer werden im Signaturtext beibehalten
- Einzeilige Ausdrucksfunktionen (`fun double(x: Int) = x * 2`) werden vollständig beibehalten

### Generics

- Generische Typparameter (`<T>`, `<T : Comparable<T>>`) werden vollständig beibehalten
- `where`-Klauseln und Varianz-Annotationen (`in`, `out`) sind in Signaturen enthalten
- `reified` Typparameter werden beibehalten

### Klassen

- `data class`, `sealed class`, `abstract class`, `open class`, `inner class`, `annotation class`, `value class` werden alle als kind `class` klassifiziert
- `enum class` wird als kind `enum` klassifiziert
- `interface` und `sealed interface` werden als kind `interface` klassifiziert

### Objekte

- `object`-Deklarationen (Singletons) werden als kind `class` klassifiziert
- `companion object`-Blöcke werden extrahiert; unbenannte Companions erhalten den synthetischen Namen „Companion"

### Körperentfernung

Wenn das Flag `--include-body` nicht verwendet wird:

- Funktionen/Methoden: Körper nach der öffnenden geschweiften Klammer `{` entfernt
- Einzeilige Ausdrucksfunktionen: vollständig beibehalten (der Ausdruck ist Teil der Signatur)
- Klassen/Interfaces/Aufzählungen: Körper nach der öffnenden geschweiften Klammer `{` entfernt
- Eigenschaften (val/var): Wertausdruck wird beibehalten
- Typ-Aliase: vollständig beibehalten

### Dokumentationskommentare

- Sowohl `/** ... */` (KDoc) als auch `//` Zeilenkommentare werden extrahiert
- KDoc-Kommentare werden der folgenden Deklaration zugeordnet
