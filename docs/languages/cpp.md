# C++ Support

🌐 [English](cpp.md) | [한국어](../ko/languages/cpp.md) | [日本語](../ja/languages/cpp.md) | [हिन्दी](../hi/languages/cpp.md) | [Deutsch](../de/languages/cpp.md)

## Supported Extensions

- `.cpp`
- `.hpp`
- `.h`

## Extraction Targets

| Element | Kind | Example |
|---------|------|---------|
| Class | `class` | `class User { ... }` |
| Struct | `struct` | `struct Point { int x, y; }` |
| Method | `method` | `void User::getName()` |
| Constructor | `constructor` | `User(string name)` |
| Destructor | `destructor` | `~User()` |
| Function | `function` | `int add(int a, int b)` |
| Namespace | `namespace` | `namespace utils { }` |
| Template | `template` | `template<typename T> class Box` |
| Enum | `enum` | `enum Color { RED, GREEN }` |
| Typedef | `typedef` | `typedef unsigned int uint` |
| Macro | `macro` | `#define MAX_SIZE 100` |
| Include | (import) | `#include <iostream>` |
| Comment | `doc` | `// Comment` |

## Example

### Input

```cpp
#include <iostream>
#include <string>

// User class for managing user data
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
    // Helper function
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

### Output (XML)

```xml
<file path="example.hpp" language="cpp">
  <signature kind="class" line="5">
    <name>User</name>
    <text>class User</text>
    <doc>User class for managing user data</doc>
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
    <doc>Helper function</doc>
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

## Notes

### Access Control

- All access levels (public, private, protected) are extracted
- No filtering based on visibility modifiers
- Useful for AI to understand complete class structure

### Body Removal

When `--include-body` flag is not used:

- Functions/Methods: body removed after opening brace `{`
- Classes/Structs/Namespaces: body removed, only declaration kept
- Templates: underlying declaration body removed
- Enum/Typedef/Macro: full text preserved

### Template Support

Basic template support is included:

```cpp
template<typename T>
class Box { ... };         // Captured

template<typename T>
T getMax(T a, T b) { ... } // Captured
```

### Namespace Support

Both simple and nested namespaces are captured:

```cpp
namespace outer {
    namespace inner {
        void helper();     // All three captured
    }
}
```

### Include Statements

Use `--include-imports` to extract `#include` directives:

```cpp
#include <iostream>        // System include
#include "myheader.h"      // Local include
```

## Unsupported Elements (v1)

| Element | Reason |
|---------|--------|
| Operator Overload | `operator+`, `operator<<` - Special case, uncommon |
| Friend Declaration | `friend class Bar` - Access control exception |
| Using Declaration | `using namespace std` - Simple alias |
| Lambda Expression | `[](int x) { ... }` - Inline definition |
| Template Specialization | `template<> class Box<int>` - Complex parsing |
| Variadic Template | `template<typename... Args>` - Advanced pattern |
| C++20 Concepts | `template<Integral T>` - Limited compiler support |
| C++20 Modules | `import std;` - Limited compiler support |
| Global Variables | May be added in future versions |
