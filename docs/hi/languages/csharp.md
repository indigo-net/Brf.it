# C# सपोर्ट

[English](../../languages/csharp.md) | [한국어](../../ko/languages/csharp.md) | [日本語](../../ja/languages/csharp.md) | [हिन्दी](csharp.md) | [Deutsch](../../de/languages/csharp.md)

## समर्थित एक्सटेंशन

- `.cs`

## एक्सट्रैक्शन टारगेट

| तत्व | Kind | उदाहरण |
|------|------|--------|
| क्लास | `class` | `public class Calculator` |
| स्ट्रक्ट | `struct` | `public struct Point` |
| इंटरफ़ेस | `interface` | `public interface IDrawable` |
| एनम | `enum` | `public enum Color` |
| रिकॉर्ड | `record` | `public record Person(string Name, int Age)` |
| रिकॉर्ड स्ट्रक्ट | `struct` | `public record struct Measurement(double Value)` |
| डेलीगेट | `type` | `public delegate void Action<T>(T item)` |
| मेथड | `method` | `public int Add(int a, int b)` |
| कंस्ट्रक्टर | `constructor` | `public Calculator()` |
| डिस्ट्रक्टर | `destructor` | `~Calculator()` |
| प्रॉपर्टी | `variable` | `public string Name { get; set; }` |
| फ़ील्ड (static/const) | `variable` | `public const int MaxValue = 100` |
| इवेंट | `variable` | `public event EventHandler Changed` |
| इंडेक्सर | `method` | `public int this[int index]` |
| ऑपरेटर | `function` | `public static operator +(...)` |
| कन्वर्शन ऑपरेटर | `function` | `public static implicit operator int(...)` |
| नेमस्पेस | `namespace` | `namespace MyApp` |
| एनम मेंबर | `variable` | `Red, Green, Blue` |
| डॉक कमेंट | `doc` | `/// <summary>...</summary>` |

## उदाहरण

### इनपुट

```csharp
using System;

namespace MyApp
{
    /// <summary>बेसिक गणित के लिए Calculator क्लास।</summary>
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

### आउटपुट (XML)

```xml
<file path="Calculator.cs" language="csharp">
  <type kind="namespace" line="3">
    <name>MyApp</name>
    <text>namespace MyApp</text>
  </type>
  <type kind="class" line="6">
    <name>Calculator</name>
    <text>public class Calculator</text>
    <doc>बेसिक गणित के लिए Calculator क्लास।</doc>
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

## नोट्स

### विजिबिलिटी

- एक्सेस मॉडिफायर की परवाह किए बिना सभी डिक्लेरेशन एक्सट्रैक्ट किए जाते हैं
- एक्सेस मॉडिफायर (`public`, `private`, `internal`, `protected`) सिग्नेचर में संरक्षित रहते हैं

### फ़ील्ड

- केवल `static`, `const`, `static readonly` फ़ील्ड एक्सट्रैक्ट किए जाते हैं
- इंस्टेंस फ़ील्ड नॉइज़ कम करने के लिए बाहर रखे जाते हैं

### प्रॉपर्टी

- ऑटो-प्रॉपर्टी (`{ get; set; }`) पूरी तरह संरक्षित रहती हैं
- एक्सप्रेशन-बॉडीड प्रॉपर्टी (`=> expr`) सिग्नेचर मोड में एक्सप्रेशन हटा दिया जाता है

### रिकॉर्ड

- `record` और `record class` को kind `record` के रूप में वर्गीकृत किया जाता है
- `record struct` को kind `struct` के रूप में वर्गीकृत किया जाता है

### ऑपरेटर

- ऑपरेटर ओवरलोड को `operator+`, `operator==` जैसे सिंथेसाइज़्ड नाम मिलते हैं
- कन्वर्शन ऑपरेटर को `implicit operator int`, `explicit operator string` जैसे नाम मिलते हैं
- इंडेक्सर को सिंथेसाइज़्ड नाम `this` मिलता है

### बॉडी रिमूवल

`--include-body` फ़्लैग का उपयोग न करने पर:

- मेथड/कंस्ट्रक्टर/डिस्ट्रक्टर: `{` के बाद बॉडी हटाई जाती है
- एक्सप्रेशन-बॉडीड मेंबर: `=>` और एक्सप्रेशन हटाया जाता है
- क्लास/स्ट्रक्ट/इंटरफ़ेस/एनम/रिकॉर्ड: `{` के बाद बॉडी हटाई जाती है
- प्रॉपर्टी: ऑटो-प्रॉपर्टी संरक्षित, एक्सप्रेशन-बॉडीड प्रॉपर्टी हटाई जाती है
- डेलीगेट: कोई बॉडी नहीं, जैसा है वैसा लौटाया जाता है

`--include-private` का उपयोग करके गैर-निर्यातित/निजी सिंबल शामिल करें।

### डॉक कमेंट

- `///` XML डॉक कमेंट और `//` लाइन कमेंट दोनों एक्सट्रैक्ट किए जाते हैं
- `/* */` ब्लॉक कमेंट भी कैप्चर किए जाते हैं
