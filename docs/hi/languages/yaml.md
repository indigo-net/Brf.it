---
title: YAML
---

# YAML समर्थन

[English](../../languages/yaml.md) | [한국어](../../ko/languages/yaml.md) | [日本語](../../ja/languages/yaml.md) | [हिन्दी](yaml.md) | [Deutsch](../../de/languages/yaml.md)

## समर्थित एक्सटेंशन

- `.yaml`
- `.yml`

## व्याकरण

- [tree-sitter-yaml](https://github.com/tree-sitter-grammars/tree-sitter-yaml) v0.7.2 by tree-sitter-grammars

## निष्कर्षण लक्ष्य

| तत्व | प्रकार | XML टैग | उदाहरण |
|------|--------|---------|--------|
| की-वैल्यू पेयर | `variable` | `<variable>` | `name: value` |

## उदाहरण

### इनपुट

```yaml
# Application configuration
name: myapp
version: 1.0.0

database:
  host: localhost
  port: 5432

features:
  - logging
  - metrics
```

### आउटपुट (XML)

```xml
<file path="config.yaml" language="yaml">
  <variable>name: myapp</variable>
  <variable>version: 1.0.0</variable>
  <variable>database:</variable>
  <variable>features:</variable>
</file>
```

## नोट्स

### बॉडी हटाना

जब `--include-body` फ़्लैग का उपयोग नहीं किया जाता:

- कंटेनर कीज़ (नेस्टेड मानों वाले मैपिंग) से नेस्टेड सामग्री हटा दी जाती है, केवल की दिखाई जाती है
- स्केलर की-वैल्यू पेयर में मान संरक्षित रहते हैं

`--include-private` का उपयोग करके गैर-निर्यातित/निजी सिंबल शामिल करें।

### टिप्पणियाँ

- एकल पंक्ति टिप्पणियाँ (`# टिप्पणी`) दस्तावेज़ के रूप में निष्कर्षित की जाती हैं

### इम्पोर्ट

- YAML में कोई इम्पोर्ट सिस्टम नहीं है; `--include-imports` का कोई प्रभाव नहीं है

### सीमाएँ

- अत्यधिक शोर से बचने के लिए केवल शीर्ष-स्तरीय कीज़ को सिग्नेचर के रूप में कैप्चर किया जाता है
- एंकर और एलियास को विशेष रूप से संभाला नहीं जाता
