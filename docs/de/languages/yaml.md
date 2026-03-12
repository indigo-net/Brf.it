---
title: YAML
---

# YAML-Unterstützung

[English](../../languages/yaml.md) | [한국어](../../ko/languages/yaml.md) | [日本語](../../ja/languages/yaml.md) | [हिन्दी](../../hi/languages/yaml.md) | [Deutsch](yaml.md)

## Unterstützte Erweiterungen

- `.yaml`
- `.yml`

## Grammatik

- [tree-sitter-yaml](https://github.com/tree-sitter-grammars/tree-sitter-yaml) v0.7.2 von tree-sitter-grammars

## Extraktionsziele

| Element | Art | XML-Tag | Beispiel |
|---------|-----|---------|----------|
| Schlüssel-Wert-Paar | `variable` | `<variable>` | `name: value` |

## Beispiel

### Eingabe

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

### Ausgabe (XML)

```xml
<file path="config.yaml" language="yaml">
  <variable>name: myapp</variable>
  <variable>version: 1.0.0</variable>
  <variable>database:</variable>
  <variable>features:</variable>
</file>
```

## Hinweise

### Entfernung des Inhalts

Wenn das `--include-body`-Flag nicht verwendet wird:

- Container-Schlüssel (Mappings mit verschachtelten Werten) werden von ihrem verschachtelten Inhalt befreit, nur der Schlüssel wird angezeigt
- Skalare Schlüssel-Wert-Paare behalten ihre Werte bei

Verwenden Sie `--include-private`, um nicht-exportierte/private Symbole einzubeziehen.

### Kommentare

- Einzeilige Kommentare (`# Kommentar`) werden als Dokumentation extrahiert

### Importe

- YAML hat kein Importsystem; `--include-imports` hat keine Wirkung

### Einschränkungen

- Nur Schlüssel der obersten Ebene werden als Signaturen erfasst, um übermäßiges Rauschen durch tief verschachtelte Werte zu vermeiden
- Anker und Aliase werden nicht speziell behandelt
