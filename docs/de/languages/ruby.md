# Ruby-Unterstützung

[English](../../languages/ruby.md) | [한국어](../../ko/languages/ruby.md) | [日本語](../../ja/languages/ruby.md) | [हिन्दी](../../hi/languages/ruby.md) | [Deutsch](ruby.md)

## Unterstützte Erweiterungen

- `.rb`

## Extraktionsziele

| Element | Kind | XML Tag | Beispiel |
|---------|------|---------|----------|
| Methode | `method` | `<function>` | `def greet(name)` |
| Klassenmethode | `method` | `<function>` | `def self.create(attrs)` |
| Klasse | `class` | `<type>` | `class User < ActiveRecord::Base` |
| Modul | `namespace` | `<type>` | `module Authentication` |
| Konstante (top-level) | `variable` | `<variable>` | `MAX_RETRIES = 3` |
| YARD-Kommentar | `doc` | | `# Beschreibung` |
| require | | `<imports>` | `require "json"` |
| require_relative | | `<imports>` | `require_relative "helpers"` |

## Beispiel

### Eingabe

```ruby
require "json"
require_relative "helpers"

MAX_RETRIES = 3

# Repräsentiert einen Benutzer im System.
class User
  # Erstellt einen neuen Benutzer aus Attributen.
  def self.create(attrs)
    new(attrs).save
  end

  # Initialisiert den Benutzer.
  def initialize(name)
    @name = name
  end

  # Begrüßt eine andere Person.
  def greet(other)
    "Hello, #{other}! I'm #{@name}."
  end
end

module Authentication
  def authenticate(password)
    password == @secret
  end
end
```

### Ausgabe (XML)

```xml
<file path="example.rb" language="ruby">
  <type>class User</type>
  <function>def self.create(attrs)</function>
  <function>def initialize(name)</function>
  <function>def greet(other)</function>
  <variable>MAX_RETRIES = 3</variable>
  <type>module Authentication</type>
  <function>def authenticate(password)</function>
</file>
```

## Hinweise

### Sichtbarkeit

- Alle Methoden werden unabhängig von der Sichtbarkeit (`public`, `protected`, `private`) extrahiert
- Sowohl Instanzmethoden (`def foo`) als auch Klassenmethoden (`def self.foo`) werden erfasst

### Methodenarten

- Sowohl Instanzmethoden (`def foo`) als auch Klassenmethoden (`def self.foo`) verwenden Kind `method`

### Körperentfernung

Wenn das `--include-body` Flag nicht verwendet wird:

- Methoden: Körper nach der schließenden Klammer `)` der Parameterliste entfernt (oder nach dem Methodennamen, wenn keine Parameter vorhanden)
- Klassen/Module: nur die Deklarationszeile wird beibehalten
- Konstanten: unverändert beibehalten

### Import-Extraktion

- `require`- und `require_relative`-Anweisungen werden mit dem `--include-imports` Flag extrahiert
- Format: `require "json"` / `require_relative "helpers"` (vollständige Anweisung wird beibehalten)

### Dokumentationskommentare

- YARD-Kommentare (`#`) direkt über Methoden/Klassen werden extrahiert
- Mehrzeilige `#`-Kommentare werden unterstützt
- `=begin`...`=end` Blockkommentare werden ebenfalls erkannt
