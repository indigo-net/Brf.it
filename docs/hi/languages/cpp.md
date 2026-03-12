# C++ सपोर्ट

🌐 [English](../../languages/cpp.md) | [한국어](../../ko/languages/cpp.md) | [日本語](../../ja/languages/cpp.md) | [हिन्दी](cpp.md) | [Deutsch](../../de/languages/cpp.md)

## समर्थित एक्सटेंशन

- `.cpp`
- `.hpp`
- `.h`

## निष्कर्षण लक्ष्य

| तत्व | Kind | उदाहरण |
|------|------|--------|
| क्लास | `class` | `class User { ... }` |
| स्ट्रक्ट | `struct` | `struct Point { int x, y; }` |
| मेथड | `method` | `void User::getName()` |
| कंस्ट्रक्टर | `constructor` | `User(string name)` |
| डिस्ट्रक्टर | `destructor` | `~User()` |
| फंक्शन | `function` | `int add(int a, int b)` |
| नेमस्पेस | `namespace` | `namespace utils { }` |
| टेम्पलेट | `template` | `template<typename T> class Box` |
| एनम | `enum` | `enum Color { RED, GREEN }` |
| Typedef | `typedef` | `typedef unsigned int uint` |
| मैक्रो | `macro` | `#define MAX_SIZE 100` |
| Include | (import) | `#include <iostream>` |
| कमेंट | `doc` | `// Comment` |

## उदाहरण

### इनपुट

```cpp
#include <iostream>
#include <string>

// यूजर डेटा मैनेज करने के लिए User क्लास
class User {
public:
    User(const std::string& name);
    ~User();

    std::string getName() const;
    void setName(const std::string& name);

private:
    std::string name_;
};

namespace utils {
    // हेल्पर फंक्शन
    int calculateHash(const std::string& input);
}

template<typename T>
class Box {
    T value;
public:
    T getValue() const;
};

#define MAX_USERS 100
```

### आउटपुट (XML)

```xml
<file path="example.hpp" language="cpp">
  <signature kind="class" line="5">
    <name>User</name>
    <text>class User</text>
    <doc>यूजर डेटा मैनेज करने के लिए User क्लास</doc>
  </signature>
  <signature kind="method" line="11">
    <name>getName</name>
    <text>std::string getName() const;</text>
  </signature>
  <signature kind="method" line="12">
    <name>setName</name>
    <text>void setName(const std::string& name);</text>
  </signature>
  <signature kind="namespace" line="18">
    <name>utils</name>
    <text>namespace utils</text>
  </signature>
  <signature kind="function" line="20">
    <name>calculateHash</name>
    <text>int calculateHash(const std::string& input);</text>
    <doc>हेल्पर फंक्शन</doc>
  </signature>
  <signature kind="template" line="23">
    <name>Box</name>
    <text>template&lt;typename T&gt; class Box</text>
  </signature>
  <signature kind="macro" line="30">
    <name>MAX_USERS</name>
    <text>#define MAX_USERS 100</text>
  </signature>
</file>
```

## नोट्स

### एक्सेस कंट्रोल

- सभी एक्सेस लेवल (public, private, protected) एक्सट्रैक्ट होते हैं
- विजिबिलिटी मॉडिफायर के आधार पर कोई फ़िल्टरिंग नहीं
- AI के लिए पूर्ण क्लास स्ट्रक्चर समझने में उपयोगी

### बॉडी रिमूवल

`--include-body` फ्लैग का उपयोग न करने पर:

- फंक्शन/मेथड: ओपनिंग ब्रेस `{` के बाद बॉडी हटाई जाती है
- क्लास/स्ट्रक्ट/नेमस्पेस: बॉडी हटाई जाती है, केवल डिक्लेरेशन रखी जाती है
- टेम्पलेट: अंडरलाइंग डिक्लेरेशन बॉडी हटाई जाती है
- Enum/Typedef/Macro: पूरा टेक्स्ट संरक्षित

### टेम्पलेट सपोर्ट

बेसिक टेम्पलेट सपोर्ट शामिल:

```cpp
template<typename T>
class Box { ... };         // कैप्चर होता है

template<typename T>
T getMax(T a, T b) { ... } // कैप्चर होता है
```

### नेमस्पेस सपोर्ट

सिंपल और नेस्टेड दोनों नेमस्पेस कैप्चर होते हैं:

```cpp
namespace outer {
    namespace inner {
        void helper();     // तीनों कैप्चर होते हैं
    }
}
```

### Include स्टेटमेंट

`--include-imports` का उपयोग करके `#include` डायरेक्टिव एक्सट्रैक्ट करें।
`--include-private` का उपयोग करके गैर-निर्यातित/निजी सिंबल शामिल करें।

```cpp
#include <iostream>        // सिस्टम include
#include "myheader.h"      // लोकल include
```

## असमर्थित तत्व (v1)

| तत्व | कारण |
|------|------|
| ऑपरेटर ओवरलोड | `operator+`, `operator<<` - स्पेशल केस, अनकॉमन |
| Friend डिक्लेरेशन | `friend class Bar` - एक्सेस कंट्रोल एक्सेप्शन |
| Using डिक्लेरेशन | `using namespace std` - सिंपल एलियास |
| Lambda एक्सप्रेशन | `[](int x) { ... }` - इनलाइन डेफिनिशन |
| टेम्पलेट स्पेशलाइजेशन | `template<> class Box<int>` - कॉम्प्लेक्स पार्सिंग |
| Variadic टेम्पलेट | `template<typename... Args>` - एडवांस्ड पैटर्न |
| C++20 Concepts | `template<Integral T>` - लिमिटेड कंपाइलर सपोर्ट |
| C++20 Modules | `import std;` - लिमिटेड कंपाइलर सपोर्ट |
| ग्लोबल वेरिएबल | भविष्य के वर्जन में जोड़ा जा सकता है |
