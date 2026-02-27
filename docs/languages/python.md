# Python Support

🌐 [English](python.md) | [한국어](python.ko.md) | [日本語](python.ja.md) | [हिन्दी](python.hi.md) | [Deutsch](python.de.md)

## Supported Extensions

- `.py`

## Extraction Targets

| Element | Kind | Example |
|---------|------|---------|
| Function | `function` | `def greet():` |
| Async function | `function` | `async def fetch():` |
| Class | `class` | `class User:` |
| Method | `method` | `def __init__(self):` |
| Class method | `method` | `def method(cls):` |
| Module-level variable | `variable` | `API_URL = "..."` |
| Comment | `doc` | `# Comment` |

## Example

### Input

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

### Output (XML)

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

## Notes

### Export Detection

- Python treats all elements as public
- `_private` or `__mangled` naming conventions are also included

### Method vs Function Detection

- If first parameter is `self` or `cls`, classified as `method`
- Otherwise classified as `function`
- Methods with `@staticmethod` decorator are classified as `function` (no self)

### Async Handling

- `async def` is unified as `function` kind
- `async` keyword is included in signature text

### Body Removal

When `--include-body` flag is not used:

- Functions/Methods: body removed after signature-ending colon (`:`)
- Classes: only class name and inheritance info are preserved

### Colon Handling in Type Hints

Colons in complex type hints (e.g., `Dict[str, int]`) are distinguished from function-ending colons:

```python
def func(x: Dict[str, List[int]]) -> str:  # Only the last colon ends the function
```

### Docstring (Future Support)

- Current version: only `#` comments above functions/classes are captured as doc
- Future version: triple-quoted docstring (`"""..."""`) support planned

### Decorators

Decorators are not included in signatures:

```python
@decorator
def func():  # Signature: "def func()"
```

### Unsupported Elements

- Lambda expressions
- Nested functions (only outer function is captured)
