---
layout: default
title: Bash/Shell
parent: Sprachanleitungen
nav_order: 13
---

# Bash/Shell-Unterstützung

[English](../../languages/bash.md) | [한국어](../../ko/languages/bash.md) | [日本語](../../ja/languages/bash.md) | [हिन्दी](../../hi/languages/bash.md) | [Deutsch](bash.md)

## Unterstützte Erweiterungen

- `.sh`
- `.bash`

## Extraktionsziele

| Element | Art | Beispiel |
|---------|-----|----------|
| Funktion | `function` | `function greet { ... }` |
| Funktion | `function` | `greet() { ... }` |
| Variablenzuweisung | `variable` | `NAME="value"` |
| Deklaration | `variable` | `declare VERBOSE` |
| Lokale Variable | `variable` | `local count=0` |
| Schreibgeschützte Variable | `variable` | `readonly VERSION="1.0"` |
| Kommentar | `doc` | `# Beschreibung` |
| Source-Anweisung | `import` | `source /path/to/lib.sh` |
| Dot-Anweisung | `import` | `. ./config.sh` |

## Beispiel

### Eingabe

```bash
#!/bin/bash

# Konfiguration
CONFIG_PATH="/etc/myapp"
VERSION="1.0.0"
declare VERBOSE=false

# Anwendung bereitstellen
function deploy {
    local app_name="$1"
    echo "Deploying $app_name"
}

# Projekt erstellen
build() {
    echo "Building..."
}

source ./utils.sh
. ./config.sh
```

### Ausgabe (XML)

```xml
<file path="deploy.sh" language="bash">
  <variable kind="variable" line="4">
    <name>CONFIG_PATH</name>
    <text>CONFIG_PATH="/etc/myapp"</text>
  </variable>
  <variable kind="variable" line="5">
    <name>VERSION</name>
    <text>VERSION="1.0.0"</text>
  </variable>
  <variable kind="variable" line="6">
    <name>VERBOSE</name>
    <text>declare VERBOSE=false</text>
  </variable>
  <function kind="function" line="9">
    <name>deploy</name>
    <text>function deploy</text>
  </function>
  <function kind="function" line="15">
    <name>build</name>
    <text>build()</text>
  </function>
</file>
```

## Hinweise

### Sichtbarkeit

- Alle Deklarationen werden extrahiert (Bash hat keine Zugriffsmodifikatoren)
- `local`-Variablen innerhalb von Funktionen werden extrahiert, wenn sie zum Parse-Zeitpunkt deklariert sind

### Funktionssyntax

Bash unterstützt zwei Funktionsdeklarationsstile:

- `function name { ... }` - mit `function`-Schlüsselwort
- `name() { ... }` - mit Klammern

Beide werden als `function`-Art extrahiert.

### Körperentfernung

Wenn das Flag `--include-body` nicht verwendet wird:

- Funktionen: Körper nach der öffnenden geschweiften Klammer `{` entfernt
- Variablen: nur die erste Zeile beibehalten (mehrzeilige Zuweisungen werden behandelt)

### Import-Extraktion

- `source`- und `.`-Befehle werden mit dem Flag `--include-imports` extrahiert
- Unterstützt sowohl zitierte als auch nicht zitierte Pfade

### Dokumentationskommentare

- Shell-Kommentare, die mit `#` beginnen, werden extrahiert
- Shebang-Zeilen (`#!/bin/bash`) werden nicht als Kommentare behandelt

### Einschränkungen

- Verschachtelte Funktionen werden unterstützt
- Here-Dokumente im Funktionskörper werden korrekt behandelt
- Komplexe Variablenexpansionen bleiben in Signaturen erhalten
