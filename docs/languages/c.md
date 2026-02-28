# C Support

🌐 [English](c.md) | [한국어](../ko/languages/c.md) | [日本語](../ja/languages/c.md) | [हिन्दी](../hi/languages/c.md) | [Deutsch](../de/languages/c.md)

## Supported Extensions

- `.c`
- `.h`

## Extraction Targets

| Element | Kind | Example |
|---------|------|---------|
| Function Definition | `function` | `int add(int a, int b) { ... }` |
| Function Declaration | `function` | `int add(int a, int b);` |
| Struct | `struct` | `struct User { ... };` |
| Enum | `enum` | `enum Color { RED, GREEN, BLUE };` |
| Typedef | `typedef` | `typedef struct { ... } User;` |
| Global Variable | `variable` | `int global_count = 0;` |
| Object-like Macro | `macro` | `#define MAX_SIZE 100` |
| Function-like Macro | `macro` | `#define MIN(a, b) ((a) < (b) ? (a) : (b))` |
| Comment | `doc` | `// Comment` |

## Example

### Input

```c
// User structure
typedef struct {
    int id;
    char name[64];
} User;

// Create a new user
User* create_user(const char* name);

// Internal helper
static void init_user(User* u);

#define MAX_USERS 100
#define INIT_USER(u) memset(u, 0, sizeof(User))
```

### Output (XML)

```xml
<file path="example.h" language="c">
  <signature kind="typedef" line="2">
    <name>User</name>
    <text>typedef struct { int id; char name[64]; } User;</text>
    <doc>User structure</doc>
  </signature>
  <signature kind="function" line="8" exported="true">
    <name>create_user</name>
    <text>User* create_user(const char* name);</text>
    <doc>Create a new user</doc>
  </signature>
  <signature kind="function" line="11" exported="true">
    <name>init_user</name>
    <text>static void init_user(User* u);</text>
    <doc>Internal helper</doc>
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

## Notes

### Export Detection

- All C functions are considered exported by default
- `static` functions are included (future: may be marked as `exported: false`)

### Body Removal

When `--include-body` flag is not used:

- Functions: body removed after opening brace `{`
- Struct/Enum/Typedef/Macro: full text preserved

### Pointer Return Types

Both direct and pointer return types are supported:

```c
int get_value();        // Direct return type
User* create_user();    // Pointer return type
```

### Unsupported Elements

- Function pointers (as variables)
- Nested structs (only top-level extracted)
