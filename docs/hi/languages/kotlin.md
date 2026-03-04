# Kotlin सपोर्ट

[English](../../languages/kotlin.md) | [한국어](../../ko/languages/kotlin.md) | [日本語](../../ja/languages/kotlin.md) | [हिन्दी](kotlin.md) | [Deutsch](../../de/languages/kotlin.md)

## समर्थित एक्सटेंशन

- `.kt`
- `.kts`

## निष्कर्षण लक्ष्य

| तत्व | Kind | उदाहरण |
|------|------|--------|
| फ़ंक्शन | `function` | `fun add(a: Int, b: Int): Int` |
| Suspend फ़ंक्शन | `function` | `suspend fun fetchData(url: String): String` |
| एक्सटेंशन फ़ंक्शन | `function` | `fun String.isEmail(): Boolean` |
| क्लास | `class` | `class User(val name: String)` |
| डेटा क्लास | `class` | `data class Point(val x: Double, val y: Double)` |
| Sealed क्लास | `class` | `sealed class Result<out T>` |
| Enum क्लास | `enum` | `enum class Color { RED, GREEN, BLUE }` |
| इंटरफ़ेस | `interface` | `interface Repository<T>` |
| ऑब्जेक्ट | `class` | `object AppConfig` |
| कंपेनियन ऑब्जेक्ट | `class` | `companion object Factory` |
| प्रॉपर्टी (val/var) | `variable` | `val MAX_SIZE = 100` |
| Type उपनाम | `type` | `typealias Handler<T> = (T) -> Unit` |
| Enum एंट्री | `variable` | `RED("#FF0000")` |
| सेकेंडरी कंस्ट्रक्टर | `constructor` | `constructor(name: String)` |
| Doc टिप्पणी | `doc` | `/** दस्तावेज़ीकरण */` |

## उदाहरण

### इनपुट

```kotlin
/** API प्रतिक्रिया के लिए उपयोगकर्ता डेटा क्लास। */
data class User(
    val id: Long,
    val name: String,
    val email: String
) {
    fun isValid(): Boolean = email.contains("@")
}

/** उपयोगकर्ता संचालन के लिए रिपोजिटरी इंटरफ़ेस। */
interface UserRepository {
    suspend fun getUser(id: Long): User?
    fun save(user: User): Boolean
}

val DEFAULT_TIMEOUT: Long = 5000L
```

### आउटपुट (XML)

```xml
<file path="user.kt" language="kotlin">
  <type kind="class" line="2">
    <name>User</name>
    <text>data class User(
    val id: Long,
    val name: String,
    val email: String
)</text>
    <doc>API प्रतिक्रिया के लिए उपयोगकर्ता डेटा क्लास।</doc>
  </type>
  <function kind="function" line="7">
    <name>isValid</name>
    <text>fun isValid(): Boolean = email.contains("@")</text>
  </function>
  <type kind="interface" line="11">
    <name>UserRepository</name>
    <text>interface UserRepository</text>
    <doc>उपयोगकर्ता संचालन के लिए रिपोजिटरी इंटरफ़ेस।</doc>
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

## टिप्पणियाँ

### दृश्यता (Visibility)

- सभी घोषणाएं निकाली जाती हैं (Kotlin डिफ़ॉल्ट रूप से `public` है)
- एक्सेस संशोधक (`public`, `internal`, `private`, `protected`) हस्ताक्षर में संरक्षित हैं

### फ़ंक्शन संशोधक

- `suspend`, `inline`, `infix`, `operator`, `tailrec` फ़ंक्शन सभी kind `function` के रूप में वर्गीकृत हैं
- संशोधक हस्ताक्षर पाठ में संरक्षित हैं
- एकल अभिव्यक्ति फ़ंक्शन (`fun double(x: Int) = x * 2`) पूरी तरह से संरक्षित हैं

### Generics

- जेनेरिक टाइप पैरामीटर (`<T>`, `<T : Comparable<T>>`) पूरी तरह से संरक्षित हैं
- `where` क्लॉज और वेरिएंस एनोटेशन (`in`, `out`) हस्ताक्षर में शामिल हैं
- `reified` टाइप पैरामीटर भी संरक्षित हैं

### क्लास

- `data class`, `sealed class`, `abstract class`, `open class`, `inner class`, `annotation class`, `value class` सभी kind `class` के रूप में वर्गीकृत हैं
- `enum class` को kind `enum` के रूप में वर्गीकृत किया गया है
- `interface` और `sealed interface` को kind `interface` के रूप में वर्गीकृत किया गया है

### ऑब्जेक्ट

- `object` घोषणाएं (सिंगलटन) kind `class` के रूप में वर्गीकृत हैं
- `companion object` ब्लॉक निकाले जाते हैं; बिना नाम वाले कंपेनियन को "Companion" सिंथेटिक नाम दिया जाता है

### बॉडी हटाना

जब `--include-body` फ्लैग का उपयोग नहीं किया जाता है:

- फ़ंक्शन/विधियां: खुलने वाले ब्रेस `{` के बाद बॉडी हटा दी जाती है
- एकल अभिव्यक्ति फ़ंक्शन: पूरी तरह से संरक्षित (अभिव्यक्ति स्वयं हस्ताक्षर का हिस्सा है)
- क्लास/इंटरफ़ेस/Enum: खुलने वाले ब्रेस `{` के बाद बॉडी हटा दी जाती है
- प्रॉपर्टी (val/var): मान अभिव्यक्ति संरक्षित है
- Type उपनाम: पूरी तरह से संरक्षित

### दस्तावेज़ टिप्पणियाँ

- `/** ... */` (KDoc) और `//` लाइन टिप्पणियाँ दोनों निकाली जाती हैं
- KDoc टिप्पणियाँ अगली घोषणा से संबद्ध होती हैं
