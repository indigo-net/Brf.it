<p align="center">
  <img src="../../assets/logo.png" alt="Brf.it logo" width="128">
</p>

# Brf.it

🌐 [English](../../README.md) | [한국어](../ko/README.md) | [日本語](../ja/README.md) | [हिन्दी](../hi/README.md) | [Deutsch](README.md)

[![Release](https://img.shields.io/github/v/release/indigo-net/Brf.it)](https://github.com/indigo-net/Brf.it/releases)
[![Downloads](https://img.shields.io/github/downloads/indigo-net/Brf.it/total)](https://github.com/indigo-net/Brf.it/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/indigo-net/Brf.it)](https://goreportcard.com/report/github.com/indigo-net/Brf.it)
[![Languages](https://img.shields.io/badge/languages-19-blue)](https://indigo-net.github.io/Brf.it/docs/languages/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Verpacken Sie Ihre Codebasis für KI-Verständnis.**
>
> `50 Token → 8 Token` — Gleiche Information, weniger Token.

Ein CLI-Tool und MCP-Server, der mit [Tree-sitter](https://tree-sitter.github.io/tree-sitter/) Funktionssignaturen, Typdefinitionen und Dokumentation aus Ihrer Codebasis extrahiert und kompakten XML/Markdown-Kontext für LLMs wie Claude, GPT und Gemini erzeugt. Unterstützt 19 Sprachen, darunter Go, TypeScript, Python, Rust, Java, C/C++ und mehr.

[Installation](#installation) · [Schnellstart](#schnellstart) · [Unterstützte Sprachen](#unterstützte-sprachen)

---

<br/>

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

<br/>

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

<br/>

## In Aktion sehen

**[SAMPLE.md](SAMPLE.md)** | **[SAMPLE.xml](SAMPLE.xml)**

Dieses Projekt, von brfit selbst verpackt. Wird bei jedem Commit automatisch generiert.

---

<br/>

## Funktionen

| Funktion | Beschreibung |
|----------|--------------|
| Tree-sitter-basiert | Präzises AST-Parsing für Sprachstrukturanalyse |
| Mehrere Formate | XML- und Markdown-Ausgabe |
| Token-Zählung | Automatische Berechnung der Ausgabe-Token |
| Gitignore-fähig | Automatischer Ausschluss unnötiger Dateien |
| Plattformübergreifend | Linux-, macOS- und Windows-Unterstützung |
| Sicherheitsprüfung | Erkennt und maskiert sicherheitsrelevante Informationen (AWS-Schlüssel, GitHub-Token, API-Schlüssel usw.) im extrahierten Code mit `[REDACTED]` |
| Aufrufgraph | Tree-sitter-basierte Extraktion von Funktions-/Methodenaufrufbeziehungen (Go, TS, Python, Java, Rust, C) |

---

<br/>

## Unterstützte Sprachen

| Sprache | Erweiterungen | Dokumentation |
|---------|---------------|---------------|
| Go | `.go` | [Go-Leitfaden](languages/go.md) |
| TypeScript | `.ts`, `.tsx` | [TypeScript-Leitfaden](languages/typescript.md) |
| JavaScript | `.js`, `.jsx` | [TypeScript-Leitfaden](languages/typescript.md) |
| Python | `.py` | [Python-Leitfaden](languages/python.md) |
| C | `.c`, `.h` | [C-Leitfaden](languages/c.md) |
| Java | `.java` | [Java-Leitfaden](languages/java.md) |
| Rust | `.rs` | [Rust-Leitfaden](languages/rust.md) |
| Swift | `.swift` | [Swift-Leitfaden](languages/swift.md) |
| Kotlin | `.kt`, `.kts` | [Kotlin-Leitfaden](languages/kotlin.md) |
| C# | `.cs` | [C#-Leitfaden](languages/csharp.md) |
| Lua | `.lua` | [Lua-Leitfaden](languages/lua.md) |
| PHP | `.php` | [PHP-Leitfaden](languages/php.md) |
| Ruby | `.rb` | [Ruby-Leitfaden](languages/ruby.md) |
| Scala | `.scala`, `.sc` | [Scala-Leitfaden](languages/scala.md) |
| Elixir | `.ex`, `.exs` | [Elixir-Leitfaden](languages/elixir.md) |
| SQL | `.sql` | [SQL-Leitfaden](languages/sql.md) |
| YAML | `.yaml`, `.yml` | [YAML-Leitfaden](languages/yaml.md) |
| TOML | `.toml` | [TOML-Leitfaden](languages/toml.md) |

---

<br/>

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
| `--include-imports` | | Import-Anweisungen einschließen | `false` |
| `--include-private` | | Nicht-exportierte/private Symbole einbeziehen | `false` |
| `--ignore` | `-i` | Ignore-Dateipfad (kann mehrfach angegeben werden) | `.gitignore` |
| `--include` | | Glob-Muster zum Einschließen (mehrfach angebbar) | |
| `--exclude` | | Glob-Muster zum Ausschließen (mehrfach angebbar) | |
| `--include-hidden` | | Versteckte Dateien einschließen | `false` |
| `--no-tree` | | Verzeichnisbaum überspringen | `false` |
| `--no-tokens` | | Token-Zählung deaktivieren | `false` |
| `--max-size` | | Maximale Dateigröße (Bytes) | `512000` |
| `--changed` | | Nur git-geänderte Dateien scannen | `false` |
| `--since` | | Nur seit Commit/Tag geänderte Dateien scannen | |
| `--token-tree` | | Token-Anzahl pro Datei als Baum anzeigen | `false` |
| `--security-check` / `--no-security-check` | | Sicherheitsrelevante Informationen erkennen und maskieren (API-Schlüssel, Token usw.) | `true` |
| `--call-graph` | | Funktionsaufrufbeziehungen pro Datei extrahieren | `false` |
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

# Imports einschließen (wörtlich)
brfit . --include-imports
```

---

<br/>

## MCP-Server

`brfit-mcp` ist ein eigenstandiger [Model Context Protocol](https://modelcontextprotocol.io/)-Server, der brfits Code-Analyse als Werkzeuge fur KI-Agenten bereitstellt. Er kommuniziert uber stdio (JSON-RPC).

### Werkzeuge

| Werkzeug | Beschreibung |
|----------|--------------|
| `summarize_project` | Signaturen aus einem Projektverzeichnis extrahieren. Optionen: `path`, `format`, `include_body`, `include_imports`, `call_graph` |
| `summarize_file` | Signaturen aus Dateien extrahieren, die einem Glob-Muster entsprechen. Optionen: `path`, `include`, `format` |

### Verwendung

```bash
brfit-mcp --root /path/to/project
```

### Claude Desktop-Konfiguration

```json
{
  "mcpServers": {
    "brfit": {
      "command": "brfit-mcp",
      "args": ["--root", "/path/to/project"]
    }
  }
}
```

---

<br/>

## Lizenz

MIT-Lizenz — Frei verwendbar in persönlichen und kommerziellen Projekten.

Siehe [LICENSE](LICENSE) für Details.
