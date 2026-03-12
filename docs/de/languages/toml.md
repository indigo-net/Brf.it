---
title: TOML
---

# TOML-Unterstützung

[English](../../languages/toml.md) | [한국어](../../ko/languages/toml.md) | [日本語](../../ja/languages/toml.md) | [हिन्दी](../../hi/languages/toml.md) | [Deutsch](toml.md)

## Unterstützte Erweiterungen

- `.toml`

## Grammatik

- [tree-sitter-toml](https://github.com/tree-sitter-grammars/tree-sitter-toml) v0.7.0 von tree-sitter-grammars

## Extraktionsziele

| Element | Art | XML-Tag | Beispiel |
|---------|-----|---------|----------|
| Tabelle | `namespace` | `<type>` | `[package]` |
| Tabellen-Array | `namespace` | `<type>` | `[[bin]]` |
| Schlüssel-Wert-Paar | `variable` | `<variable>` | `name = "myapp"` |

## Beispiel

### Eingabe

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

### Ausgabe (XML)

```xml
<file path="config.toml" language="toml">
  <variable>name = "myapp"</variable>
  <variable>version = "1.0.0"</variable>
  <type>[package]</type>
  <type>[[bin]]</type>
</file>
```

## Hinweise

### Entfernung des Inhalts

Wenn das `--include-body`-Flag nicht verwendet wird:

- Tabellensektionen (`[table]`) werden von ihrem Inhalt befreit, nur der Header wird angezeigt
- Tabellen-Arrays (`[[table]]`) werden von ihrem Inhalt befreit, nur der Header wird angezeigt
- Schlüssel-Wert-Paare der obersten Ebene behalten ihre Werte bei

Verwenden Sie `--include-private`, um nicht-exportierte/private Symbole einzubeziehen.

### Kommentare

- Einzeilige Kommentare (`# Kommentar`) werden als Dokumentation extrahiert

### Importe

- TOML hat kein Importsystem; `--include-imports` hat keine Wirkung

### Einschränkungen

- Inline-Tabellen und Inline-Arrays werden nicht separat extrahiert
- Punktierte Schlüssel (z.B. `physical.color = "orange"`) werden als einzelnes Paar erfasst
