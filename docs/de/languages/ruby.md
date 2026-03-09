# Ruby-Unterstützung

[English](../../languages/ruby.md) | [한국어](../../ko/languages/ruby.md) | [日本語](../../ja/languages/ruby.md) | [हिन्दी](../../hi/languages/ruby.md) | [Deutsch](ruby.md)

## Unterstützte Erweiterungen

- `.rb`

## Extraktionsziele

| Element | Kind | Beispiel |
|---------|------|----------|
| Methode | `method` | `def greet(name)` |
| Klassenmethode | `class_method` | `def self.create(attrs)` |
| Klasse | `class` | `class User < ActiveRecord::Base` |
| Modul | `module` | `module Authentication` |
| Konstante | `variable` | `MAX_RETRIES = 3` |
| YARD-Kommentar | `doc` | `# Beschreibung` |
| require | `import` | `require "json"` |
| require_relative | `import` | `require_relative "helpers"` |

## Beispiel

### Eingabe

```ruby
require "json"
require_relative "helpers"

# Repräsentiert einen Benutzer im System.
class User
  MAX_RETRIES = 3

  # Erstellt einen neuen Benutzer aus Attributen.
  # @param attrs [Hash] Benutzerattribute
  def self.create(attrs)
    new(attrs).save
  end

  # Initialisiert den Benutzer.
  # @param name [String] der Name des Benutzers
  def initialize(name)
    @name = name
  end

  # Begrüßt eine andere Person.
  # @param other [String] der Name der anderen Person
  # @return [String] eine Grußnachricht
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
  <class kind="class" line="5">
    <name>User</name>
    <text>class User</text>
  </class>
  <variable kind="variable" line="6">
    <name>MAX_RETRIES</name>
    <text>MAX_RETRIES = 3</text>
  </variable>
  <function kind="class_method" line="10">
    <name>create</name>
    <text>def self.create(attrs)</text>
  </function>
  <function kind="method" line="15">
    <name>initialize</name>
    <text>def initialize(name)</text>
  </function>
  <function kind="method" line="21">
    <name>greet</name>
    <text>def greet(other)</text>
  </function>
  <module kind="module" line="27">
    <name>Authentication</name>
    <text>module Authentication</text>
  </module>
  <function kind="method" line="28">
    <name>authenticate</name>
    <text>def authenticate(password)</text>
  </function>
</file>
```

## Hinweise

### Sichtbarkeit

- Alle Methoden werden unabhängig von der Sichtbarkeit (`public`, `protected`, `private`) extrahiert
- Sowohl Instanzmethoden (`def foo`) als auch Klassenmethoden (`def self.foo`) werden erfasst

### Methodenarten

- `method`: Instanzmethoden-Deklarationen (`def foo`)
- `class_method`: Klassenmethoden-Deklarationen (`def self.foo`)

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
