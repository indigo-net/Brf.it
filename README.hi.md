# Brf.it

🌐 [English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md) | [हिन्दी](README.hi.md) | [Deutsch](README.de.md)

[![Release](https://img.shields.io/github/v/release/indigo-net/Brf.it)](https://github.com/indigo-net/Brf.it/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/indigo-net/Brf.it)](https://goreportcard.com/report/github.com/indigo-net/Brf.it)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **अपने कोडबेस को AI समझ के लिए पैकेज करें।**
>
> `50 टोकन → 8 टोकन` — समान जानकारी, कम टोकन।

[इंस्टॉलेशन](#इंस्टॉलेशन) · [क्विक स्टार्ट](#क्विक-स्टार्ट) · [समर्थित भाषाएं](#समर्थित-भाषाएं)

---

## यह कैसे काम करता है

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

## क्विक स्टार्ट

### इंस्टॉलेशन

**macOS (Homebrew)**

```bash
brew install indigo-net/tap/brfit
```

**Linux / macOS (स्क्रिप्ट)**

```bash
curl -fsSL https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.sh | sh
```

**Windows (PowerShell)**

```powershell
irm https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.ps1 | iex
```

**सोर्स से बिल्ड**

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o brfit ./cmd/brfit
```

### पहली बार चलाएं

```bash
brfit .                    # वर्तमान डायरेक्टरी का विश्लेषण
brfit . -f md              # Markdown में आउटपुट
brfit . -o briefing.xml    # फाइल में सेव करें
```

---

## फीचर्स

| फीचर | विवरण |
|------|-------|
| Tree-sitter आधारित | सटीक AST पार्सिंग से भाषा संरचना विश्लेषण |
| मल्टीपल फॉर्मेट्स | XML और Markdown आउटपुट सपोर्ट |
| टोकन काउंटिंग | आउटपुट टोकन की स्वचालित गणना |
| Gitignore अवेयर | अनावश्यक फाइल्स को स्वचालित रूप से बाहर करें |
| क्रॉस-प्लेटफॉर्म | Linux, macOS, और Windows सपोर्ट |

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

## CLI रेफरेंस

```bash
brfit [पथ] [विकल्प]
```

### विकल्प

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

### उदाहरण

```bash
# AI सहायकों को भेजें (क्लिपबोर्ड में कॉपी)
brfit . | pbcopy              # macOS
brfit . | xclip               # Linux
brfit . | clip                # Windows

# प्रोजेक्ट का विश्लेषण करें और फाइल में सेव करें
brfit ./my-project -o briefing.xml

# फंक्शन बॉडी शामिल करें (पूरा कोड)
brfit . --include-body

# डायरेक्टरी ट्री आउटपुट स्किप करें
brfit . --no-tree
```

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

[MIT](LICENSE)
