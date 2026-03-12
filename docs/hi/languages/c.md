# C सपोर्ट

🌐 [English](../../languages/c.md) | [한국어](../../ko/languages/c.md) | [日本語](../../ja/languages/c.md) | [हिन्दी](c.md) | [Deutsch](../../de/languages/c.md)

## समर्थित एक्सटेंशन

- `.c`
- `.h`

## एक्सट्रैक्शन टारगेट

| एलिमेंट | Kind | उदाहरण |
|---------|------|--------|
| फंक्शन डेफिनिशन | `function` | `int add(int a, int b) { ... }` |
| फंक्शन डिक्लेरेशन | `function` | `int add(int a, int b);` |
| Struct | `struct` | `struct User { ... };` |
| Enum | `enum` | `enum Color { RED, GREEN, BLUE };` |
| Typedef | `typedef` | `typedef struct { ... } User;` |
| ग्लोबल वेरिएबल | `variable` | `int global_count = 0;` |
| ऑब्जेक्ट मैक्रो | `macro` | `#define MAX_SIZE 100` |
| फंक्शन मैक्रो | `macro` | `#define MIN(a, b) ((a) < (b) ? (a) : (b))` |
| कमेंट | `doc` | `// Comment` |

## उदाहरण

### इनपुट

```c
// User स्ट्रक्चर
typedef struct {
    int id;
    char name[64];
} User;

// नया यूजर बनाएं
User* create_user(const char* name);

// इंटरनल हेल्पर
static void init_user(User* u);

#define MAX_USERS 100
#define INIT_USER(u) memset(u, 0, sizeof(User))
```

### आउटपुट (XML)

```xml
<file path="example.h" language="c">
  <signature kind="typedef" line="2">
    <name>User</name>
    <text>typedef struct { int id; char name[64]; } User;</text>
    <doc>User स्ट्रक्चर</doc>
  </signature>
  <signature kind="function" line="8" exported="true">
    <name>create_user</name>
    <text>User* create_user(const char* name);</text>
    <doc>नया यूजर बनाएं</doc>
  </signature>
  <signature kind="function" line="11" exported="true">
    <name>init_user</name>
    <text>static void init_user(User* u);</text>
    <doc>इंटरनल हेल्पर</doc>
  </signature>
  <signature kind="macro" line="13">
    <name>MAX_USERS</name>
    <text>#define MAX_USERS 100</text>
  </signature>
  <signature kind="macro" line="14">
    <name>INIT_USER</name>
    <text>#define INIT_USER(u) memset(u, 0, sizeof(User))</text>
  </signature>
</file>
```

## नोट्स

### Export डिटेक्शन

- सभी C फंक्शन डिफ़ॉल्ट रूप से exported माने जाते हैं
- `static` फंक्शन भी शामिल हैं (भविष्य में: `exported: false` के रूप में मार्क किया जा सकता है)

### बॉडी रिमूवल

जब `--include-body` फ्लैग का उपयोग नहीं किया जाता:

- फंक्शन: ओपनिंग ब्रेस `{` के बाद की बॉडी हटा दी जाती है
- Struct/Enum/Typedef/Macro: पूरा टेक्स्ट संरक्षित

`--include-private` का उपयोग करके गैर-निर्यातित/निजी सिंबल शामिल करें।

### पॉइंटर रिटर्न टाइप

डायरेक्ट और पॉइंटर रिटर्न टाइप दोनों समर्थित:

```c
int get_value();        // डायरेक्ट रिटर्न टाइप
User* create_user();    // पॉइंटर रिटर्न टाइप
```

### असमर्थित एलिमेंट

- फंक्शन पॉइंटर (वेरिएबल के रूप में)
- नेस्टेड स्ट्रक्चर (केवल टॉप-लेवल एक्सट्रैक्ट होते हैं)
