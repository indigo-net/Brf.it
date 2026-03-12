# Java सपोर्ट

🌐 [English](../../languages/java.md) | [한국어](../../ko/languages/java.md) | [日本語](../../ja/languages/java.md) | [हिन्दी](java.md) | [Deutsch](../../de/languages/java.md)

## समर्थित एक्सटेंशन

- `.java`

## निष्कर्षण लक्ष्य

| तत्व | Kind | उदाहरण |
|------|------|--------|
| क्लास | `class` | `public class User { ... }` |
| इंटरफेस | `interface` | `public interface Repository<T> { ... }` |
| मेथड | `method` | `public String getName() { ... }` |
| कंस्ट्रक्टर | `constructor` | `public User(String name) { ... }` |
| Enum | `enum` | `public enum Status { ... }` |
| एनोटेशन | `annotation` | `public @interface Inject { ... }` |
| Record (Java 14+) | `record` | `public record Point(int x, int y) { ... }` |
| फील्ड | `field` | `public static final String API = "..."` |
| कमेंट | `doc` | `// Comment` या `/* Block */` |

## उदाहरण

### इनपुट

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

### आउटपुट (XML)

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

## विशेष नोट्स

### विजिबिलिटी फ़िल्टरिंग

- `public`, `protected`, package-private (डिफ़ॉल्ट): डिफ़ॉल्ट रूप से निकाला जाता है
- `private`: केवल `--include-body` उपयोग करने पर शामिल

### जेनेरिक्स हैंडलिंग

जेनेरिक टाइप पैरामीटर सिग्नेचर में शामिल होते हैं:

```java
public class Box<T extends Comparable<T>>  // पूर्ण कैप्चर
public <U> U transform(Function<T, U> fn)  // मेथड टाइप पैरामीटर शामिल
```

### एनोटेशन आउटपुट

मेथड और क्लास एनोटेशन सिग्नेचर टेक्स्ट में शामिल होते हैं:

```java
@Override
public String toString()  // सिग्नेचर में @Override शामिल
```

### Record सपोर्ट (Java 14+)

Records कंपोनेंट पैरामीटर के साथ निकाले जाते हैं:

```java
public record User(String name, int age)  // कंपोनेंट संरक्षित
```

### इनर/नेस्टेड क्लास

सभी नेस्टेड क्लास अलग सिग्नेचर के रूप में निकाले जाते हैं:

```java
public class Outer {
    public static class Nested { ... }  // अलग से निकाला जाता है
    public class Inner { ... }          // यह भी निकाला जाता है
}
```

### एब्स्ट्रैक्ट मेथड्स

इंटरफेस में एब्स्ट्रैक्ट मेथड्स `;` से समाप्त होते हैं (कोई बॉडी नहीं):

```java
interface Foo {
    void bar();  // जैसा है वैसा कैप्चर
}
```

### बॉडी रिमूवल

`--include-body` फ्लैग के बिना:

- मेथड्स/कंस्ट्रक्टर्स: ओपनिंग ब्रेस `{` के बाद बॉडी हटाई जाती है
- क्लास/इंटरफेस/Enum: ओपनिंग ब्रेस `{` के बाद बॉडी हटाई जाती है
- एब्स्ट्रैक्ट मेथड्स: जैसा है वैसा रखा जाता है (`;` से समाप्त)

`--include-private` का उपयोग करके गैर-निर्यातित/निजी सिंबल शामिल करें।

### Javadoc (भविष्य में सपोर्ट)

- वर्तमान संस्करण: डिक्लेरेशन के ऊपर `//` और `/* */` कमेंट्स doc के रूप में कैप्चर
- भविष्य का संस्करण: `/** */` Javadoc पार्सिंग सपोर्ट प्लान

### असमर्थित तत्व

- स्टैटिक इनिशियलाइज़र ब्लॉक
- एनोनिमस क्लास
- Lambda एक्सप्रेशन
