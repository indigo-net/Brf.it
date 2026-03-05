# Brf.it

🌐 [English](../../README.md) | [한국어](../ko/README.md) | [日本語](../ja/README.md) | [हिन्दी](README.md) | [Deutsch](../de/README.md)

[![Release](https://img.shields.io/github/v/release/indigo-net/Brf.it)](https://github.com/indigo-net/Brf.it/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/indigo-net/Brf.it)](https://goreportcard.com/report/github.com/indigo-net/Brf.it)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **अपने कोडबेस को AI समझ के लिए पैकेज करें।**
>
> `50 टोकन → 8 टोकन` — समान जानकारी, कम टोकन।

[इंस्टॉलेशन](#इंस्टॉलेशन) · [क्विक स्टार्ट](#क्विक-स्टार्ट) · [समर्थित भाषाएं](#समर्थित-भाषाएं)

---

<br/>

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
<function>
  export async function fetchUser(
    id: string
  ): Promise<User>
</function>
```

</td>
</tr>
</table>

---

<br/>

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

<br/>

## इसे एक्शन में देखें

**[SAMPLE.md](SAMPLE.md)** | **[SAMPLE.xml](SAMPLE.xml)**

यह प्रोजेक्ट, brfit द्वारा स्वयं पैकेज किया गया। हर कमिट पर ऑटो-जेनरेट होता है।

---

<br/>

## फीचर्स

| फीचर | विवरण |
|------|-------|
| Tree-sitter आधारित | सटीक AST पार्सिंग से भाषा संरचना विश्लेषण |
| मल्टीपल फॉर्मेट्स | XML और Markdown आउटपुट सपोर्ट |
| टोकन काउंटिंग | आउटपुट टोकन की स्वचालित गणना |
| Gitignore अवेयर | अनावश्यक फाइल्स को स्वचालित रूप से बाहर करें |
| क्रॉस-प्लेटफॉर्म | Linux, macOS, और Windows सपोर्ट |

---

<br/>

## समर्थित भाषाएं

| भाषा | एक्सटेंशन | डॉक्यूमेंटेशन |
|------|-----------|---------------|
| Go | `.go` | [Go गाइड](languages/go.md) |
| TypeScript | `.ts`, `.tsx` | [TypeScript गाइड](languages/typescript.md) |
| JavaScript | `.js`, `.jsx` | [TypeScript गाइड](languages/typescript.md) |
| Python | `.py` | [Python गाइड](languages/python.md) |
| C | `.c`, `.h` | [C गाइड](languages/c.md) |
| Java | `.java` | [Java गाइड](languages/java.md) |
| Rust | `.rs` | [Rust गाइड](languages/rust.md) |
| Swift | `.swift` | [Swift गाइड](languages/swift.md) |
| Kotlin | `.kt`, `.kts` | [Kotlin गाइड](languages/kotlin.md) |
| C# | `.cs` | [C# गाइड](languages/csharp.md) |

---

<br/>

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
| `--include-imports` | | import स्टेटमेंट शामिल करें | `false` |
| `--no-std-imports` | | स्टैंडर्ड लाइब्रेरी import हटाएं | `false` |
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

# import शामिल करें (स्टैंडर्ड लाइब्रेरी हटाएं)
brfit . --include-imports --no-std-imports
```

---

<br/>

## लाइसेंस

MIT लाइसेंस — व्यक्तिगत और व्यावसायिक प्रोजेक्ट में स्वतंत्र रूप से उपयोग करें।

विवरण के लिए [LICENSE](LICENSE) देखें।
