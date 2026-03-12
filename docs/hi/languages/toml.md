---
title: TOML
---

# TOML समर्थन

[English](../../languages/toml.md) | [한국어](../../ko/languages/toml.md) | [日本語](../../ja/languages/toml.md) | [हिन्दी](toml.md) | [Deutsch](../../de/languages/toml.md)

## समर्थित एक्सटेंशन

- `.toml`

## व्याकरण

- [tree-sitter-toml](https://github.com/tree-sitter-grammars/tree-sitter-toml) v0.7.0 by tree-sitter-grammars

## निष्कर्षण लक्ष्य

| तत्व | प्रकार | XML टैग | उदाहरण |
|------|--------|---------|--------|
| टेबल | `namespace` | `<type>` | `[package]` |
| टेबल एरे | `namespace` | `<type>` | `[[bin]]` |
| की-वैल्यू पेयर | `variable` | `<variable>` | `name = "myapp"` |

## उदाहरण

### इनपुट

```toml
# Project configuration
name = "myapp"
version = "1.0.0"

[package]
authors = ["Alice"]
edition = "2024"

[[bin]]
name = "cli"
path = "src/main.rs"
```

### आउटपुट (XML)

```xml
<file path="config.toml" language="toml">
  <variable>name = "myapp"</variable>
  <variable>version = "1.0.0"</variable>
  <type>[package]</type>
  <type>[[bin]]</type>
</file>
```

## नोट्स

### बॉडी हटाना

जब `--include-body` फ़्लैग का उपयोग नहीं किया जाता:

- टेबल सेक्शन (`[table]`) से बॉडी हटा दी जाती है, केवल हेडर दिखाया जाता है
- टेबल एरे (`[[table]]`) से बॉडी हटा दी जाती है, केवल हेडर दिखाया जाता है
- शीर्ष-स्तरीय की-वैल्यू पेयर में मान संरक्षित रहते हैं

`--include-private` का उपयोग करके गैर-निर्यातित/निजी सिंबल शामिल करें।

### टिप्पणियाँ

- एकल पंक्ति टिप्पणियाँ (`# टिप्पणी`) दस्तावेज़ के रूप में निष्कर्षित की जाती हैं

### इम्पोर्ट

- TOML में कोई इम्पोर्ट सिस्टम नहीं है; `--include-imports` का कोई प्रभाव नहीं है

### सीमाएँ

- इनलाइन टेबल और इनलाइन एरे को अलग से निष्कर्षित नहीं किया जाता
- डॉटेड कीज़ (उदा. `physical.color = "orange"`) को एक पेयर के रूप में कैप्चर किया जाता है
