# Brf.it

🌐 [English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md) | [हिन्दी](README.hi.md) | [Deutsch](README.de.md)

[![Release](https://img.shields.io/github/v/release/indigo-net/Brf.it)](https://github.com/indigo-net/Brf.it/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/indigo-net/Brf.it)](https://goreportcard.com/report/github.com/indigo-net/Brf.it)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**AI सहायकों के लिए कोड ब्रीफिंग टूल**

Brf.it आपके कोडबेस से फंक्शन सिग्नेचर निकालता है, इम्प्लीमेंटेशन डिटेल्स हटाकर AI को आवश्यक जानकारी प्रदान करता है। टोकन उपयोग को नाटकीय रूप से कम करता है।

---

## यह क्या करता है

AI सहायकों को रॉ कोड देने की बजाय:

<table>
<tr>
<th>पहले (50+ टोकन)</th>
<th>brfit के बाद (8 टोकन)</th>
</tr>
<tr>
<td>

```typescript
export async function fetchUser(
  id: string
): Promise<User> {
  const response = await fetch(
    `${API_URL}/users/${id}`
  );
  if (!response.ok) {
    throw new Error('User not found');
  }
  const data = await response.json();
  return {
    id: data.id,
    name: data.name,
    email: data.email,
    createdAt: new Date(data.created_at)
  };
}
```

</td>
<td>

```xml
<signature>
  export async function fetchUser(
    id: string
  ): Promise<User>
</signature>
```

</td>
</tr>
</table>

---

## इंस्टॉलेशन

### macOS (Homebrew)

```bash
brew install indigo-net/tap/brfit
```

### रिलीज से डाउनलोड

[Releases](https://github.com/indigo-net/Brf.it/releases) से नवीनतम बाइनरी डाउनलोड करें।

### सोर्स से बिल्ड

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o brfit ./cmd/brfit
```

---

## उपयोग

```bash
brfit [पथ] [विकल्प]
```

### त्वरित उदाहरण

```bash
# वर्तमान डायरेक्टरी से सिग्नेचर निकालें
brfit .

# Markdown फॉर्मेट में आउटपुट
brfit . -f md

# फाइल में सेव करें
brfit . -o output.xml

# फंक्शन बॉडी शामिल करें (पूरा कोड)
brfit . --include-body

# डायरेक्टरी ट्री स्किप करें
brfit . --no-tree
```

### CLI विकल्प

| विकल्प | शॉर्ट | विवरण | डिफ़ॉल्ट |
|--------|-------|-------|----------|
| `--format` | `-f` | आउटपुट फॉर्मेट (`xml`, `md`) | `xml` |
| `--output` | `-o` | आउटपुट फाइल पथ | stdout |
| `--include-body` | | फंक्शन बॉडी शामिल करें | `false` |
| `--ignore` | `-i` | ignore फाइल पथ | `.gitignore` |
| `--include-hidden` | | हिडन फाइल्स शामिल करें | `false` |
| `--no-tree` | | डायरेक्टरी ट्री स्किप करें | `false` |
| `--no-tokens` | | टोकन काउंटिंग अक्षम करें | `false` |
| `--max-size` | | अधिकतम फाइल साइज (बाइट्स) | `512000` |
| `--version` | `-v` | वर्जन दिखाएं | |

---

## समर्थित भाषाएं

| भाषा | एक्सटेंशन | डॉक्यूमेंटेशन |
|------|-----------|---------------|
| Go | `.go` | [Go गाइड](docs/languages/go.hi.md) |
| TypeScript | `.ts`, `.tsx` | [TypeScript गाइड](docs/languages/typescript.hi.md) |
| JavaScript | `.js`, `.jsx` | [TypeScript गाइड](docs/languages/typescript.hi.md) |
| Python | `.py` | [Python गाइड](docs/languages/python.hi.md) |
| C | `.c`, `.h` | [C गाइड](docs/languages/c.hi.md) |

---

## आउटपुट उदाहरण

### XML (डिफ़ॉल्ट)

```xml
<?xml version="1.0" encoding="UTF-8"?>
<brfit>
  <metadata>
    <tree>pkg/
└── scanner/
    └── scanner.go</tree>
  </metadata>
  <files>
    <file path="pkg/scanner/scanner.go" language="go">
      <signature>func Scan(root string) (*Result, error)</signature>
      <doc>Scan recursively scans the directory.</doc>
    </file>
  </files>
</brfit>
```

### Markdown

```markdown
# Brf.it Output

## Directory Tree

pkg/
└── scanner/
    └── scanner.go

## Files

### pkg/scanner/scanner.go

\`\`\`go
func Scan(root string) (*Result, error)
\`\`\`

> Scan recursively scans the directory.
```

---

## लाइसेंस

MIT
