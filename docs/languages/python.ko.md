# Python 지원

🌐 [English](python.md) | [한국어](python.ko.md) | [日本語](python.ja.md) | [हिन्दी](python.hi.md) | [Deutsch](python.de.md)

## 지원 확장자

- `.py`

## 추출 대상

| 요소 | Kind | 예시 |
|------|------|------|
| 함수 | `function` | `def greet():` |
| Async 함수 | `function` | `async def fetch():` |
| 클래스 | `class` | `class User:` |
| 메서드 | `method` | `def __init__(self):` |
| 클래스 메서드 | `method` | `def method(cls):` |
| 모듈 레벨 변수 | `variable` | `API_URL = "..."` |
| 주석 | `doc` | `# Comment` |

## 예시

### 입력

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

### 출력 (XML)

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

## 특이사항

### Export 판별

- Python은 모든 요소를 public으로 취급
- `_private` 또는 `__mangled` 네이밍도 포함됨

### 메서드 vs 함수 판별

- 첫 번째 파라미터가 `self` 또는 `cls`이면 `method`로 분류
- 그 외의 경우 `function`으로 분류
- `@staticmethod` 데코레이터가 붙은 메서드는 `function`으로 분류됨 (self 없음)

### Async 처리

- `async def`는 `function` kind로 통일
- 시그니처 텍스트에 `async` 키워드 포함

### 본문 제거

`--include-body` 플래그 미사용 시:

- 함수/메서드: 시그니처 끝 콜론(`:`) 이후 본문 제거
- 클래스: 클래스명과 상속 정보까지만 유지

### 타입 힌트 내 콜론 처리

복잡한 타입 힌트(예: `Dict[str, int]`)의 콜론은 함수 끝 콜론과 구분됨:

```python
def func(x: Dict[str, List[int]]) -> str:  # 마지막 콜론만 함수 끝
```

### Docstring (향후 지원)

- 현재 버전: 함수/클래스 위의 `#` 주석만 doc으로 캡처
- 향후 버전: triple-quoted docstring (`"""..."""`) 지원 예정

### 데코레이터

데코레이터는 시그니처에 포함되지 않음:

```python
@decorator
def func():  # 시그니처: "def func()"
```

### 지원하지 않는 요소

- Lambda 표현식
- 중첩 함수 (외부 함수만 캡처)
