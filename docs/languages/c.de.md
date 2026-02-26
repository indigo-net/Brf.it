# C Unterstützung

🌐 [English](c.md) | [한국어](c.ko.md) | [日本語](c.ja.md) | [हिन्दी](c.hi.md) | [Deutsch](c.de.md)

## Unterstützte Erweiterungen

- `.c`
- `.h`

## Extraktionsziele

| Element | Kind | Beispiel |
|---------|------|----------|
| Funktionsdefinition | `function` | `int add(int a, int b) { ... }` |
| Funktionsdeklaration | `function` | `int add(int a, int b);` |
| Struct | `struct` | `struct User { ... };` |
| Enum | `enum` | `enum Color { RED, GREEN, BLUE };` |
| Typedef | `typedef` | `typedef struct { ... } User;` |
| Objekt-Makro | `macro` | `#define MAX_SIZE 100` |
| Funktions-Makro | `macro` | `#define MIN(a, b) ((a) < (b) ? (a) : (b))` |
| Kommentar | `doc` | `// Comment` |

## Beispiel

### Eingabe

```c
// User-Struktur
typedef struct {
    int id;
    char name[64];
} User;

// Neuen Benutzer erstellen
User* create_user(const char* name);

// Interner Helfer
static void init_user(User* u);

#define MAX_USERS 100
#define INIT_USER(u) memset(u, 0, sizeof(User))
```

### Ausgabe (XML)

```xml
<file path="example.h" language="c">
  <signature kind="typedef" line="2">
    <name>User</name>
    <text>typedef struct { int id; char name[64]; } User;</text>
    <doc>User-Struktur</doc>
  </signature>
  <signature kind="function" line="8" exported="true">
    <name>create_user</name>
    <text>User* create_user(const char* name);</text>
    <doc>Neuen Benutzer erstellen</doc>
  </signature>
  <signature kind="function" line="11" exported="true">
    <name>init_user</name>
    <text>static void init_user(User* u);</text>
    <doc>Interner Helfer</doc>
  </signature>
  <signature kind="macro" line="13">
    <name>MAX_USERS</name>
    <text>#define MAX_USERS 100</text>
  </signature>
  <signature kind="macro" line="14">
    <name>INIT_USER</name>
    <text>#define INIT_USER(u) memset(u, 0, sizeof(User))</text>
  </signature>
</file>
```

## Hinweise

### Export-Erkennung

- Alle C-Funktionen werden standardmäßig als exportiert behandelt
- `static` Funktionen sind ebenfalls enthalten (Zukunft: könnte als `exported: false` markiert werden)

### Body-Entfernung

Wenn `--include-body` Flag nicht verwendet wird:

- Funktionen: Body nach öffnender Klammer `{` entfernt
- Struct/Enum/Typedef/Macro: vollständiger Text erhalten

### Pointer-Rückgabetypen

Sowohl direkte als auch Pointer-Rückgabetypen werden unterstützt:

```c
int get_value();        // Direkter Rückgabetyp
User* create_user();    // Pointer-Rückgabetyp
```

### Nicht unterstützte Elemente

- Globale Variablen
- Funktionszeiger (als Variablen)
- Verschachtelte Strukturen (nur oberste Ebene wird extrahiert)
