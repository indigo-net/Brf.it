# C++ Unterstützung

🌐 [English](cpp.md) | [한국어](cpp.ko.md) | [日本語](cpp.ja.md) | [हिन्दी](cpp.hi.md) | [Deutsch](cpp.de.md)

## Unterstützte Erweiterungen

- `.cpp`
- `.hpp`
- `.h`

## Extraktionsziele

| Element | Kind | Beispiel |
|---------|------|----------|
| Klasse | `class` | `class User { ... }` |
| Struct | `struct` | `struct Point { int x, y; }` |
| Methode | `method` | `void User::getName()` |
| Konstruktor | `constructor` | `User(string name)` |
| Destruktor | `destructor` | `~User()` |
| Funktion | `function` | `int add(int a, int b)` |
| Namespace | `namespace` | `namespace utils { }` |
| Template | `template` | `template<typename T> class Box` |
| Enum | `enum` | `enum Color { RED, GREEN }` |
| Typedef | `typedef` | `typedef unsigned int uint` |
| Makro | `macro` | `#define MAX_SIZE 100` |
| Include | (import) | `#include <iostream>` |
| Kommentar | `doc` | `// Comment` |

## Beispiel

### Eingabe

```cpp
#include <iostream>
#include <string>

// User-Klasse zur Verwaltung von Benutzerdaten
class User {
public:
    User(const std::string& name);
    ~User();

    std::string getName() const;
    void setName(const std::string& name);

private:
    std::string name_;
};

namespace utils {
    // Hilfsfunktion
    int calculateHash(const std::string& input);
}

template<typename T>
class Box {
    T value;
public:
    T getValue() const;
};

#define MAX_USERS 100
```

### Ausgabe (XML)

```xml
<file path="example.hpp" language="cpp">
  <signature kind="class" line="5">
    <name>User</name>
    <text>class User</text>
    <doc>User-Klasse zur Verwaltung von Benutzerdaten</doc>
  </signature>
  <signature kind="method" line="11">
    <name>getName</name>
    <text>std::string getName() const;</text>
  </signature>
  <signature kind="method" line="12">
    <name>setName</name>
    <text>void setName(const std::string& name);</text>
  </signature>
  <signature kind="namespace" line="18">
    <name>utils</name>
    <text>namespace utils</text>
  </signature>
  <signature kind="function" line="20">
    <name>calculateHash</name>
    <text>int calculateHash(const std::string& input);</text>
    <doc>Hilfsfunktion</doc>
  </signature>
  <signature kind="template" line="23">
    <name>Box</name>
    <text>template&lt;typename T&gt; class Box</text>
  </signature>
  <signature kind="macro" line="30">
    <name>MAX_USERS</name>
    <text>#define MAX_USERS 100</text>
  </signature>
</file>
```

## Hinweise

### Zugriffskontrolle

- Alle Zugriffsebenen (public, private, protected) werden extrahiert
- Keine Filterung basierend auf Sichtbarkeitsmodifikatoren
- Nützlich für KI, um die vollständige Klassenstruktur zu verstehen

### Body-Entfernung

Wenn `--include-body` Flag nicht verwendet wird:

- Funktionen/Methoden: Body nach öffnender Klammer `{` entfernt
- Klassen/Structs/Namespaces: Body entfernt, nur Deklaration beibehalten
- Templates: Zugrundeliegender Deklarations-Body entfernt
- Enum/Typedef/Macro: vollständiger Text erhalten

### Template-Unterstützung

Grundlegende Template-Unterstützung ist enthalten:

```cpp
template<typename T>
class Box { ... };         // Erfasst

template<typename T>
T getMax(T a, T b) { ... } // Erfasst
```

### Namespace-Unterstützung

Sowohl einfache als auch verschachtelte Namespaces werden erfasst:

```cpp
namespace outer {
    namespace inner {
        void helper();     // Alle drei erfasst
    }
}
```

### Include-Anweisungen

Verwenden Sie `--include-imports`, um `#include`-Direktiven zu extrahieren:

```cpp
#include <iostream>        // System-Include
#include "myheader.h"      // Lokales Include
```

## Nicht unterstützte Elemente (v1)

| Element | Grund |
|---------|-------|
| Operator-Überladung | `operator+`, `operator<<` - Spezialfall, selten |
| Friend-Deklaration | `friend class Bar` - Zugriffskontroll-Ausnahme |
| Using-Deklaration | `using namespace std` - Einfacher Alias |
| Lambda-Ausdruck | `[](int x) { ... }` - Inline-Definition |
| Template-Spezialisierung | `template<> class Box<int>` - Komplexes Parsen |
| Variadic Template | `template<typename... Args>` - Fortgeschrittenes Muster |
| C++20 Concepts | `template<Integral T>` - Begrenzte Compiler-Unterstützung |
| C++20 Modules | `import std;` - Begrenzte Compiler-Unterstützung |
| Globale Variablen | Kann in zukünftigen Versionen hinzugefügt werden |
