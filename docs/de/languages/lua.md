# Lua-Unterstützung

[English](../../languages/lua.md) | [한국어](../../ko/languages/lua.md) | [日本語](../../ja/languages/lua.md) | [हिन्दी](../../hi/languages/lua.md) | [Deutsch](lua.md)

## Unterstützte Erweiterungen

- `.lua`

## Extraktionsziele

| Element | Kind | Beispiel |
|---------|------|----------|
| Globale Funktion | `function` | `function globalFunc(a, b)` |
| Lokale Funktion | `local_function` | `local function helper()` |
| Modulfunktion | `module_function` | `function M.create(name)` |
| Methode | `method` | `function M:init(config)` |
| Tabellenzuweisung | `variable` | `local M = {}` |
| Funktionszuweisung | `variable` | `local foo = function() end` |
| LuaDoc-Kommentar | `doc` | `--- Beschreibung` |
| Zeilenkommentar | `doc` | `-- Beschreibung` |
| require() | `import` | `local json = require("json")` |

## Beispiel

### Eingabe

```lua
local M = {}

--- Begrüßt eine Person mit Namen.
-- @param name string Der Name der Person
function M.greet(name)
    print("Hello, " .. name)
end

function M:init(config)
    self.config = config
end

local function helper()
    return 42
end

function globalFunc(a, b)
    return a + b
end

local json = require("json")
```

### Ausgabe (XML)

```xml
<file path="example.lua" language="lua">
  <variable kind="variable" line="1">
    <name>M</name>
    <text>local M = {}</text>
  </variable>
  <function kind="module_function" line="4">
    <name>greet</name>
    <text>function M.greet(name)</text>
  </function>
  <function kind="method" line="9">
    <name>init</name>
    <text>function M:init(config)</text>
  </function>
  <function kind="local_function" line="13">
    <name>helper</name>
    <text>local function helper()</text>
  </function>
  <function kind="function" line="17">
    <name>globalFunc</name>
    <text>function globalFunc(a, b)</text>
  </function>
</file>
```

## Hinweise

### Sichtbarkeit

- Alle Deklarationen werden extrahiert (Lua hat keine Zugriffsmodifikatoren)
- `local` Funktionen/Variablen werden zusammen mit globalen einbezogen

### Funktionsarten

- `function`: Globale Funktionsdeklarationen (`function foo()`)
- `local_function`: Lokale Funktionsdeklarationen (`local function foo()`)
- `module_function`: Punkt-indizierte Funktionen (`function M.foo()`)
- `method`: Doppelpunkt-indizierte Methoden (`function M:foo()`)

### Körperentfernung

Wenn das `--include-body` Flag nicht verwendet wird:

- Funktionen/Methoden: Körper nach der schließenden Klammer `)` der Parameterliste entfernt
- Tabellenzuweisungen (`local M = {}`): unverändert beibehalten
- Funktionszuweisungen (`local foo = function() end`): unverändert beibehalten

### Import-Extraktion

- `require()`-Aufrufe werden mit dem `--include-imports` Flag extrahiert
- Verwenden Sie `--include-private`, um nicht-exportierte/private Symbole einzubeziehen
- Format: `local json = require("json")` (vollständige Anweisung wird beibehalten)

### Dokumentationskommentare

- Sowohl LuaDoc-Kommentare (`---`) als auch reguläre Zeilenkommentare (`--`) werden extrahiert
- Blockkommentare (`--[[ ... ]]`) werden ebenfalls unterstützt
