# Scala सपोर्ट

[English](../../languages/scala.md) | [한국어](../../ko/languages/scala.md) | [日本語](../../ja/languages/scala.md) | [हिन्दी](scala.md) | [Deutsch](../../de/languages/scala.md)

## समर्थित एक्सटेंशन

- `.scala`
- `.sc`

## एक्सट्रैक्शन लक्ष्य

| तत्व | Kind | XML Tag | उदाहरण |
|------|------|---------|--------|
| मेथड (बॉडी सहित) | `method` | `<function>` | `def add(a: Int, b: Int): Int` |
| मेथड (एब्स्ट्रैक्ट) | `method` | `<function>` | `def greet(name: String): String` |
| क्लास | `class` | `<type>` | `class Person(val name: String)` |
| ट्रेट | `trait` | `<type>` | `trait Greeter` |
| ऑब्जेक्ट | `class` | `<type>` | `object MathUtils` |
| val | `variable` | `<variable>` | `val PI: Double = 3.14159` |
| var | `variable` | `<variable>` | `var count: Int = 0` |
| टाइप एलियास | `type` | `<type>` | `type StringList = List[String]` |
| Enum (Scala 3) | `enum` | `<type>` | `enum Color` |
| Given (Scala 3) | `variable` | `<variable>` | `given ordering: Ordering[Int]` |
| Extension (Scala 3) | `method` | `<function>` | `extension (s: String)` |

## उदाहरण

### इनपुट

```scala
// उपयोगकर्ता प्रबंधन
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

### आउटपुट (XML)

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

## नोट्स

### दृश्यता (Visibility)

- विजिबिलिटी मॉडिफायर की परवाह किए बिना सभी डिक्लेरेशन एक्सट्रैक्ट किए जाते हैं
- विजिबिलिटी मॉडिफायर (`private`, `protected`) सिग्नेचर टेक्स्ट में संरक्षित रहते हैं

### क्लास वेरिएंट

- `class`, `abstract class`, `case class`, `sealed class`, `implicit class` सभी kind `class` के रूप में वर्गीकृत हैं
- `trait` और `sealed trait` kind `trait` के रूप में वर्गीकृत हैं
- `object` (सिंगलटन और कंपेनियन) kind `class` के रूप में वर्गीकृत है

### बॉडी रिमूवल

`--include-body` फ़्लैग का उपयोग न करने पर:

- मेथड: `=` के बाद बॉडी हटा दी जाती है, रिटर्न टाइप संरक्षित रहता है
- क्लास/ट्रेट/ऑब्जेक्ट: `{ }` में बॉडी हटा दी जाती है, केवल डिक्लेरेशन लाइन संरक्षित रहती है
- val/var: वैल्यू संरक्षित रहती हैं (`lazy val`, `implicit val` सहित)
- टाइप एलियास: पूर्ण रूप से संरक्षित

### जेनेरिक्स

- जेनेरिक टाइप पैरामीटर `[A, B]` सिग्नेचर में पूर्ण रूप से संरक्षित रहते हैं
- कॉन्टेक्स्ट बाउंड और व्यू बाउंड शामिल हैं

### Scala 3 फीचर्स

- `enum` डेफिनिशन kind `enum` के रूप में वर्गीकृत हैं
- नामित `given` इंस्टेंस kind `variable` के रूप में वर्गीकृत हैं
- `extension` मेथड ग्रुप kind `method` के रूप में वर्गीकृत हैं
