# Python सपोर्ट

🌐 [English](../../languages/python.md) | [한국어](../../ko/languages/python.md) | [日本語](../../ja/languages/python.md) | [हिन्दी](python.md) | [Deutsch](../../de/languages/python.md)

## समर्थित एक्सटेंशन

- `.py`

## निष्कर्षण लक्ष्य

| तत्व | Kind | उदाहरण |
|------|------|--------|
| फंक्शन | `function` | `def greet():` |
| Async फंक्शन | `function` | `async def fetch():` |
| क्लास | `class` | `class User:` |
| मेथड | `method` | `def __init__(self):` |
| क्लास मेथड | `method` | `def method(cls):` |
| मॉड्यूल-लेवल वेरिएबल | `variable` | `API_URL = "..."` |
| कमेंट | `doc` | `# Comment` |

## उदाहरण

### इनपुट

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

### आउटपुट (XML)

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

## नोट्स

### एक्सपोर्ट डिटेक्शन

- Python सभी तत्वों को public मानता है
- `_private` या `__mangled` नेमिंग भी शामिल है

### मेथड vs फंक्शन डिटेक्शन

- अगर पहला पैरामीटर `self` या `cls` है, तो `method` के रूप में वर्गीकृत
- अन्यथा `function` के रूप में वर्गीकृत
- `@staticmethod` डेकोरेटर वाले मेथड `function` के रूप में वर्गीकृत (कोई self नहीं)

### Async हैंडलिंग

- `async def` को `function` kind के रूप में एकीकृत किया गया है
- सिग्नेचर टेक्स्ट में `async` कीवर्ड शामिल है

### बॉडी रिमूवल

`--include-body` फ्लैग का उपयोग न करने पर:

- फंक्शन/मेथड: सिग्नेचर-एंडिंग कोलन (`:`) के बाद बॉडी हटा दी जाती है
- क्लास: केवल क्लास नाम और इनहेरिटेंस इंफो संरक्षित रहती है

`--include-private` का उपयोग करके गैर-निर्यातित/निजी सिंबल शामिल करें।

### टाइप हिंट में कोलन हैंडलिंग

कॉम्प्लेक्स टाइप हिंट (जैसे `Dict[str, int]`) में कोलन फंक्शन-एंडिंग कोलन से अलग होते हैं:

```python
def func(x: Dict[str, List[int]]) -> str:  # केवल आखिरी कोलन फंक्शन समाप्त करता है
```

### Docstring (भविष्य सपोर्ट)

- वर्तमान वर्जन: केवल फंक्शन/क्लास के ऊपर `#` कमेंट्स doc के रूप में कैप्चर होते हैं
- भविष्य वर्जन: ट्रिपल-क्वोटेड docstring (`"""..."""`) सपोर्ट प्लान्ड

### डेकोरेटर्स

डेकोरेटर्स सिग्नेचर में शामिल नहीं होते:

```python
@decorator
def func():  # सिग्नेचर: "def func()"
```

### असमर्थित तत्व

- Lambda एक्सप्रेशन
- नेस्टेड फंक्शन (केवल बाहरी फंक्शन कैप्चर होता है)
