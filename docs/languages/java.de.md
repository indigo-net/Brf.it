# Java Unterstützung

🌐 [English](java.md) | [한국어](java.ko.md) | [日本語](java.ja.md) | [हिन्दी](java.hi.md) | [Deutsch](java.de.md)

## Unterstützte Erweiterungen

- `.java`

## Extraktionsziele

| Element | Kind | Beispiel |
|---------|------|----------|
| Klasse | `class` | `public class User { ... }` |
| Interface | `interface` | `public interface Repository<T> { ... }` |
| Methode | `method` | `public String getName() { ... }` |
| Konstruktor | `constructor` | `public User(String name) { ... }` |
| Enum | `enum` | `public enum Status { ... }` |
| Annotation | `annotation` | `public @interface Inject { ... }` |
| Record (Java 14+) | `record` | `public record Point(int x, int y) { ... }` |
| Feld | `field` | `public static final String API = "..."` |
| Kommentar | `doc` | `// Comment` oder `/* Block */` |

## Beispiel

### Eingabe

```java
package com.example;

/**
 * User class represents a user in the system.
 */
public class User {
    private String name;

    public User(String name) {
        this.name = name;
    }

    public String getName() {
        return name;
    }

    private void internalMethod() {
        // Private method
    }
}

public interface Repository<T> {
    T findById(String id);
    void save(T entity);
}

public enum Status {
    PENDING, ACTIVE, COMPLETED
}

public record Point(int x, int y) {}
```

### Ausgabe (XML)

```xml
<file path="User.java" language="java">
  <signature kind="class" line="6">
    <name>User</name>
    <text>public class User</text>
  </signature>
  <signature kind="constructor" line="9">
    <name>User</name>
    <text>public User(String name)</text>
  </signature>
  <signature kind="method" line="13">
    <name>getName</name>
    <text>public String getName()</text>
  </signature>
  <signature kind="interface" line="22">
    <name>Repository</name>
    <text>public interface Repository&lt;T&gt;</text>
  </signature>
  <signature kind="method" line="23">
    <name>findById</name>
    <text>T findById(String id);</text>
  </signature>
  <signature kind="method" line="24">
    <name>save</name>
    <text>void save(T entity);</text>
  </signature>
  <signature kind="enum" line="27">
    <name>Status</name>
    <text>public enum Status</text>
  </signature>
  <signature kind="record" line="31">
    <name>Point</name>
    <text>public record Point(int x, int y)</text>
  </signature>
</file>
```

## Hinweise

### Sichtbarkeitsfilterung

- `public`, `protected`, package-private (Standard): standardmäßig extrahiert
- `private`: nur bei Verwendung von `--include-body` enthalten

### Generics-Behandlung

Generische Typparameter sind in der Signatur enthalten:

```java
public class Box<T extends Comparable<T>>  // Vollständig erfasst
public <U> U transform(Function<T, U> fn)  // Methoden-Typparameter enthalten
```

### Annotationen in der Ausgabe

Methoden- und Klassenannotationen sind im Signaturtext enthalten:

```java
@Override
public String toString()  // Signatur enthält @Override
```

### Record-Unterstützung (Java 14+)

Records werden mit ihren Komponentenparametern extrahiert:

```java
public record User(String name, int age)  // Komponenten erhalten
```

### Innere/Verschachtelte Klassen

Alle verschachtelten Klassen werden als separate Signaturen extrahiert:

```java
public class Outer {
    public static class Nested { ... }  // Separat extrahiert
    public class Inner { ... }          // Ebenfalls extrahiert
}
```

### Abstrakte Methoden

Abstrakte Methoden in Interfaces enden mit `;` (kein Rumpf):

```java
interface Foo {
    void bar();  // Wie vorhanden erfasst
}
```

### Rumpfentfernung

Wenn das `--include-body`-Flag nicht verwendet wird:

- Methoden/Konstruktoren: Rumpf nach öffnender Klammer `{` entfernt
- Klassen/Interfaces/Enums: Rumpf nach öffnender Klammer `{` entfernt
- Abstrakte Methoden: unverändert beibehalten (enden mit `;`)

### Javadoc (Zukünftige Unterstützung)

- Aktuelle Version: `//` und `/* */` Kommentare über Deklarationen werden als doc erfasst
- Zukünftige Version: `/** */` Javadoc-Parsing geplant

### Nicht unterstützte Elemente

- Statische Initialisierungsblöcke
- Anonyme Klassen
- Lambda-Ausdrücke
