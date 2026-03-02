# Brf.it

🌐 [English](../../README.md) | [한국어](../ko/README.md) | [日本語](../ja/README.md) | [हिन्दी](../hi/README.md) | [Deutsch](README.md)

[![Release](https://img.shields.io/github/v/release/indigo-net/Brf.it)](https://github.com/indigo-net/Brf.it/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/indigo-net/Brf.it)](https://goreportcard.com/report/github.com/indigo-net/Brf.it)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Verpacken Sie Ihre Codebasis für KI-Verständnis.**
>
> `50 Token → 8 Token` — Gleiche Information, weniger Token.

[Installation](#installation) · [Schnellstart](#schnellstart) · [Unterstützte Sprachen](#unterstützte-sprachen)

---

## Wie es funktioniert

Anstatt rohen Code an KI-Assistenten zu übergeben:

<table>
<tr>
<th>Vorher (50+ Token)</th>
<th>Nachher mit brfit (8 Token)</th>
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

## Schnellstart

### Installation

**macOS (Homebrew)**

```bash
brew install indigo-net/tap/brfit
```

**Linux / macOS (Skript)**

```bash
curl -fsSL https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.sh | sh
```

**Windows (PowerShell)**

```powershell
irm https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.ps1 | iex
```

**Aus Quellcode bauen**

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o brfit ./cmd/brfit
```

### Erster Start

```bash
brfit .                    # Aktuelles Verzeichnis analysieren
brfit . -f md              # Ausgabe in Markdown
brfit . -o briefing.xml    # In Datei speichern
```

---

## In Aktion sehen

**[SAMPLE.md](SAMPLE.md)** | **[SAMPLE.xml](SAMPLE.xml)**

Dieses Projekt, von brfit selbst verpackt. Wird bei jedem Commit automatisch generiert.

---

## Funktionen

| Funktion | Beschreibung |
|----------|--------------|
| Tree-sitter-basiert | Präzises AST-Parsing für Sprachstrukturanalyse |
| Mehrere Formate | XML- und Markdown-Ausgabe |
| Token-Zählung | Automatische Berechnung der Ausgabe-Token |
| Gitignore-fähig | Automatischer Ausschluss unnötiger Dateien |
| Plattformübergreifend | Linux-, macOS- und Windows-Unterstützung |

---

## Unterstützte Sprachen

| Sprache | Erweiterungen | Dokumentation |
|---------|---------------|---------------|
| Go | `.go` | [Go-Leitfaden](languages/go.md) |
| TypeScript | `.ts`, `.tsx` | [TypeScript-Leitfaden](languages/typescript.md) |
| JavaScript | `.js`, `.jsx` | [TypeScript-Leitfaden](languages/typescript.md) |
| Python | `.py` | [Python-Leitfaden](languages/python.md) |
| C | `.c`, `.h` | [C-Leitfaden](languages/c.md) |
| Java | `.java` | [Java-Leitfaden](languages/java.md) |

---

## CLI-Referenz

```bash
brfit [Pfad] [Optionen]
```

### Optionen

| Option | Kurz | Beschreibung | Standard |
|--------|------|--------------|----------|
| `--format` | `-f` | Ausgabeformat (`xml`, `md`) | `xml` |
| `--output` | `-o` | Ausgabedateipfad | stdout |
| `--include-body` | | Funktionskörper einschließen | `false` |
| `--ignore` | `-i` | Ignore-Dateipfad | `.gitignore` |
| `--include-hidden` | | Versteckte Dateien einschließen | `false` |
| `--no-tree` | | Verzeichnisbaum überspringen | `false` |
| `--no-tokens` | | Token-Zählung deaktivieren | `false` |
| `--max-size` | | Maximale Dateigröße (Bytes) | `512000` |
| `--version` | `-v` | Version anzeigen | |

### Beispiele

```bash
# An KI-Assistenten senden (in Zwischenablage kopieren)
brfit . | pbcopy              # macOS
brfit . | xclip               # Linux
brfit . | clip                # Windows

# Projekt analysieren und in Datei speichern
brfit ./my-project -o briefing.xml

# Funktionskörper einschließen (vollständiger Code)
brfit . --include-body

# Verzeichnisbaum-Ausgabe überspringen
brfit . --no-tree
```

---

## Ausgabebeispiele

### XML (Standard)

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
      <function>func Scan(root string) (*Result, error)</function>
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

## Lizenz

MIT-Lizenz — Frei verwendbar in persönlichen und kommerziellen Projekten.

Siehe [LICENSE](LICENSE) für Details.
