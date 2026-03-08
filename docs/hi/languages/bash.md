---
layout: default
title: Bash/Shell
parent: भाषा गाइड
nav_order: 13
---

# Bash/Shell समर्थन

[English](../../languages/bash.md) | [한국어](../../ko/languages/bash.md) | [日本語](../../ja/languages/bash.md) | [हिन्दी](bash.md) | [Deutsch](../../de/languages/bash.md)

## समर्थित एक्सटेंशन

- `.sh`
- `.bash`

## निष्कर्षण लक्ष्य

| तत्व | प्रकार | उदाहरण |
|------|-------|--------|
| फंक्शन | `function` | `function greet { ... }` |
| फंक्शन | `function` | `greet() { ... }` |
| वेरिएबल असाइनमेंट | `variable` | `NAME="value"` |
| डिक्लेरेशन | `variable` | `declare VERBOSE` |
| लोकल वेरिएबल | `variable` | `local count=0` |
| रीडओनली वेरिएबल | `variable` | `readonly VERSION="1.0"` |
| टिप्पणी | `doc` | `# विवरण` |
| source स्टेटमेंट | `import` | `source /path/to/lib.sh` |
| dot स्टेटमेंट | `import` | `. ./config.sh` |

## उदाहरण

### इनपुट

```bash
#!/bin/bash

# कॉन्फ़िगरेशन
CONFIG_PATH="/etc/myapp"
VERSION="1.0.0"
declare VERBOSE=false

# एप्लिकेशन डिप्लॉय करें
function deploy {
    local app_name="$1"
    echo "Deploying $app_name"
}

# प्रोजेक्ट बिल्ड करें
build() {
    echo "Building..."
}

source ./utils.sh
. ./config.sh
```

### आउटपुट (XML)

```xml
<file path="deploy.sh" language="bash">
  <variable kind="variable" line="4">
    <name>CONFIG_PATH</name>
    <text>CONFIG_PATH="/etc/myapp"</text>
  </variable>
  <variable kind="variable" line="5">
    <name>VERSION</name>
    <text>VERSION="1.0.0"</text>
  </variable>
  <variable kind="variable" line="6">
    <name>VERBOSE</name>
    <text>declare VERBOSE=false</text>
  </variable>
  <function kind="function" line="9">
    <name>deploy</name>
    <text>function deploy</text>
  </function>
  <function kind="function" line="15">
    <name>build</name>
    <text>build()</text>
  </function>
</file>
```

## टिप्पणियाँ

### दृश्यता

- सभी घोषणाएं निकाली जाती हैं (Bash में कोई एक्सेस मॉडिफायर नहीं हैं)
- फंक्शन के अंदर `local` वेरिएबल भी पार्स समय पर घोषित होने पर निकाले जाते हैं

### फंक्शन सिंटैक्स

Bash दो फंक्शन घोषणा शैलियों का समर्थन करता है:

- `function नाम { ... }` - `function` कीवर्ड के साथ
- `नाम() { ... }` - कोष्ठक के साथ

दोनों `function` प्रकार के रूप में निकाले जाते हैं।

### बॉडी हटाना

जब `--include-body` फ्लैग का उपयोग नहीं किया जाता है:

- फंक्शन: खुलने वाले ब्रेस `{` के बाद बॉडी हटा दी जाती है
- वेरिएबल: केवल पहली पंक्ति रखी जाती है (मल्टी-लाइन असाइनमेंट को संभाला जाता है)

### इंपोर्ट निष्कर्षण

- `source` और `.` कमांड `--include-imports` फ्लैग के साथ निकाले जाते हैं
- क्वोट और बिना क्वोट वाले पथ दोनों का समर्थन करता है

### डॉक्यूमेंटेशन टिप्पणियाँ

- `#` से शुरू होने वाली शेल टिप्पणियां निकाली जाती हैं
- शेबैंग लाइनें (`#!/bin/bash`) टिप्पणियों के रूप में नहीं मानी जाती हैं

### सीमाएँ

- नेस्टेड फंक्शन समर्थित हैं
- फंक्शन बॉडी में हियर-डॉक्यूमेंट्स सही ढंग से संभाले जाते हैं
- जटिल वेरिएबल एक्सपैंशन सिग्नेचर में संरक्षित रहते हैं
