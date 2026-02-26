# Python Unterstützung

🌐 [English](python.md) | [한국어](python.ko.md) | [日本語](python.ja.md) | [हिन्दी](python.hi.md) | [Deutsch](python.de.md)

## Unterstützte Erweiterungen

- `.py`

## Extraktionsziele

| Element | Kind | Beispiel |
|---------|------|----------|
| Funktion | `function` | `def greet():` |
| Async-Funktion | `function` | `async def fetch():` |
| Klasse | `class` | `class User:` |
| Methode | `method` | `def __init__(self):` |
| Klassenmethode | `method` | `def method(cls):` |
| Kommentar | `doc` | `# Comment` |

## Beispiel

### Eingabe

```python
# User model for the application.
class User:
    """Represents a user in the system."""

    def __init__(self, name: str, email: str):
        """Initialize a new user."""
        self.name = name
        self.email = email

    def get_display_name(self) -> str:
        """Return the display name."""
        return f"{self.name} <{self.email}>"


# Create a new user instance.
def create_user(name: str, email: str) -> User:
    """Factory function to create a user."""
    return User(name, email)


async def fetch_user(user_id: int) -> User:
    """Fetch a user from the database."""
    pass
```

### Ausgabe (XML)

```xml
<file path="user.py" language="python">
  <signature kind="class" line="2">
    <name>User</name>
    <text>class User</text>
    <doc>User model for the application.</doc>
  </signature>
  <signature kind="method" line="5">
    <name>__init__</name>
    <text>def __init__(self, name: str, email: str)</text>
  </signature>
  <signature kind="method" line="10">
    <name>get_display_name</name>
    <text>def get_display_name(self) -> str</text>
  </signature>
  <signature kind="function" line="16">
    <name>create_user</name>
    <text>def create_user(name: str, email: str) -> User</text>
    <doc>Create a new user instance.</doc>
  </signature>
  <signature kind="function" line="21">
    <name>fetch_user</name>
    <text>async def fetch_user(user_id: int) -> User</text>
  </signature>
</file>
```

## Hinweise

### Export-Erkennung

- Python behandelt alle Elemente als public
- `_private` oder `__mangled` Namenskonventionen sind ebenfalls enthalten

### Methode vs Funktion Erkennung

- Wenn der erste Parameter `self` oder `cls` ist, wird als `method` klassifiziert
- Andernfalls als `function` klassifiziert
- Methoden mit `@staticmethod` Dekorator werden als `function` klassifiziert (kein self)

### Async-Behandlung

- `async def` wird als `function` Kind vereinheitlicht
- `async` Schlüsselwort ist im Signaturtext enthalten

### Body-Entfernung

Wenn `--include-body` Flag nicht verwendet wird:

- Funktionen/Methoden: Body nach Signatur-beendendem Doppelpunkt (`:`) entfernt
- Klassen: nur Klassenname und Vererbungsinformationen werden beibehalten

### Doppelpunkt-Behandlung in Type Hints

Doppelpunkte in komplexen Type Hints (z.B. `Dict[str, int]`) werden von Funktions-beendenden Doppelpunkten unterschieden:

```python
def func(x: Dict[str, List[int]]) -> str:  # Nur der letzte Doppelpunkt beendet die Funktion
```

### Docstring (Zukünftige Unterstützung)

- Aktuelle Version: nur `#` Kommentare über Funktionen/Klassen werden als doc erfasst
- Zukünftige Version: Triple-Quoted Docstring (`"""..."""`) Unterstützung geplant

### Dekoratoren

Dekoratoren sind nicht in Signaturen enthalten:

```python
@decorator
def func():  # Signatur: "def func()"
```

### Nicht unterstützte Elemente

- Lambda-Ausdrücke
- Verschachtelte Funktionen (nur äußere Funktion wird erfasst)
- Modul-Level-Variablen
