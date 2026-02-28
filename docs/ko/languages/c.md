# C 지원

🌐 [English](../../languages/c.md) | [한국어](c.md) | [日本語](../../ja/languages/c.md) | [हिन्दी](../../hi/languages/c.md) | [Deutsch](../../de/languages/c.md)

## 지원 확장자

- `.c`
- `.h`

## 추출 대상

| 요소 | Kind | 예시 |
|------|------|------|
| 함수 정의 | `function` | `int add(int a, int b) { ... }` |
| 함수 선언 | `function` | `int add(int a, int b);` |
| 구조체 | `struct` | `struct User { ... };` |
| 열거형 | `enum` | `enum Color { RED, GREEN, BLUE };` |
| Typedef | `typedef` | `typedef struct { ... } User;` |
| 전역 변수 | `variable` | `int global_count = 0;` |
| 객체 매크로 | `macro` | `#define MAX_SIZE 100` |
| 함수 매크로 | `macro` | `#define MIN(a, b) ((a) < (b) ? (a) : (b))` |
| 주석 | `doc` | `// Comment` |

## 예시

### 입력

```c
// User 구조체
typedef struct {
    int id;
    char name[64];
} User;

// 새 사용자 생성
User* create_user(const char* name);

// 내부 헬퍼
static void init_user(User* u);

#define MAX_USERS 100
#define INIT_USER(u) memset(u, 0, sizeof(User))
```

### 출력 (XML)

```xml
<file path="example.h" language="c">
  <signature kind="typedef" line="2">
    <name>User</name>
    <text>typedef struct { int id; char name[64]; } User;</text>
    <doc>User 구조체</doc>
  </signature>
  <signature kind="function" line="8" exported="true">
    <name>create_user</name>
    <text>User* create_user(const char* name);</text>
    <doc>새 사용자 생성</doc>
  </signature>
  <signature kind="function" line="11" exported="true">
    <name>init_user</name>
    <text>static void init_user(User* u);</text>
    <doc>내부 헬퍼</doc>
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

## 참고사항

### Export 감지

- 모든 C 함수는 기본적으로 exported로 처리됨
- `static` 함수도 포함됨 (향후: `exported: false`로 표시될 수 있음)

### 본문 제거

`--include-body` 플래그를 사용하지 않을 때:

- 함수: 여는 중괄호 `{` 이후 본문 제거
- Struct/Enum/Typedef/Macro: 전체 텍스트 유지

### 포인터 반환 타입

직접 반환 타입과 포인터 반환 타입 모두 지원:

```c
int get_value();        // 직접 반환 타입
User* create_user();    // 포인터 반환 타입
```

### 지원되지 않는 요소

- 함수 포인터 (변수로서)
- 중첩 구조체 (최상위 레벨만 추출)
