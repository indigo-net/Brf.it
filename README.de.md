# Brf.it

🌐 [English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md) | [हिन्दी](README.hi.md) | [Deutsch](README.de.md)

[![Release](https://img.shields.io/github/v/release/indigo-net/Brf.it)](https://github.com/indigo-net/Brf.it/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/indigo-net/Brf.it)](https://goreportcard.com/report/github.com/indigo-net/Brf.it)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Code-Briefing-Tool für KI-Assistenten**

Brf.it extrahiert Funktionssignaturen aus Ihrer Codebasis, entfernt Implementierungsdetails und reduziert den Token-Verbrauch drastisch, während die wesentlichen Informationen für KI erhalten bleiben.

---

## Was es macht

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

## Installation

### macOS (Homebrew)

```bash
brew install indigo-net/tap/brfit
```

### Von Release herunterladen

Laden Sie die neueste Binary von [Releases](https://github.com/indigo-net/Brf.it/releases) herunter.

### Aus Quellcode bauen

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o brfit ./cmd/brfit
```

---

## Verwendung

```bash
brfit [Pfad] [Optionen]
```

### Schnelle Beispiele

```bash
# Signaturen aus dem aktuellen Verzeichnis extrahieren
brfit .

# Ausgabe im Markdown-Format
brfit . -f md

# In Datei speichern
brfit . -o output.xml

# Funktionskörper einschließen (vollständiger Code)
brfit . --include-body

# Verzeichnisbaum überspringen
brfit . --no-tree
```

### CLI-Optionen

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

---

## Unterstützte Sprachen

| Sprache | Erweiterungen | Dokumentation |
|---------|---------------|---------------|
| Go | `.go` | [Go-Leitfaden](docs/languages/go.de.md) |
| TypeScript | `.ts`, `.tsx` | [TypeScript-Leitfaden](docs/languages/typescript.de.md) |
| JavaScript | `.js`, `.jsx` | [TypeScript-Leitfaden](docs/languages/typescript.de.md) |
| Python | `.py` | [Python-Leitfaden](docs/languages/python.de.md) |

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

## Lizenz

MIT
