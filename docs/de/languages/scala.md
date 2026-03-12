# Scala-Unterstützung

[English](../../languages/scala.md) | [한국어](../../ko/languages/scala.md) | [日本語](../../ja/languages/scala.md) | [हिन्दी](../../hi/languages/scala.md) | [Deutsch](scala.md)

## Unterstützte Erweiterungen

- `.scala`
- `.sc`

## Extraktionsziele

| Element | Kind | XML Tag | Beispiel |
|---------|------|---------|----------|
| Methode (mit Körper) | `method` | `<function>` | `def add(a: Int, b: Int): Int` |
| Methode (abstrakt) | `method` | `<function>` | `def greet(name: String): String` |
| Klasse | `class` | `<type>` | `class Person(val name: String)` |
| Trait | `trait` | `<type>` | `trait Greeter` |
| Objekt | `class` | `<type>` | `object MathUtils` |
| val | `variable` | `<variable>` | `val PI: Double = 3.14159` |
| var | `variable` | `<variable>` | `var count: Int = 0` |
| Typ-Alias | `type` | `<type>` | `type StringList = List[String]` |
| Enum (Scala 3) | `enum` | `<type>` | `enum Color` |
| Given (Scala 3) | `variable` | `<variable>` | `given ordering: Ordering[Int]` |

## Beispiel

### Eingabe

```scala
// Benutzerverwaltung
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

### Ausgabe (XML)

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

## Hinweise

### Sichtbarkeit

- Alle Deklarationen werden unabhängig von Sichtbarkeitsmodifikatoren extrahiert
- Sichtbarkeitsmodifikatoren (`private`, `protected`) werden im Signaturtext beibehalten

### Klassenvarianten

- `class`, `abstract class`, `case class`, `sealed class`, `implicit class` werden alle als Kind `class` klassifiziert
- `trait` und `sealed trait` werden als Kind `trait` klassifiziert
- `object` (Singleton und Companion) wird als Kind `class` klassifiziert

### Körperentfernung

Wenn das `--include-body` Flag nicht verwendet wird:

- Methoden: Körper nach `=` wird entfernt, der Rückgabetyp bleibt erhalten
- Klassen/Traits/Objekte: Körper in `{ }` wird entfernt, nur die Deklarationszeile bleibt erhalten
- val/var: Werte werden beibehalten (einschließlich `lazy val` und `implicit val`)
- Typ-Aliase: vollständig beibehalten

Verwenden Sie `--include-private`, um nicht-exportierte/private Symbole einzubeziehen.

### Generics

- Generische Typparameter `[A, B]` werden vollständig in Signaturen beibehalten
- Kontextgrenzen und Ansichtsgrenzen sind enthalten

### Scala 3-Funktionen

- `enum`-Definitionen werden als Kind `enum` klassifiziert
- Benannte `given`-Instanzen werden als Kind `variable` klassifiziert
- `extension`-Methoden werden einzeln als `method` extrahiert (die extension-Deklaration selbst wird nicht erfasst)
