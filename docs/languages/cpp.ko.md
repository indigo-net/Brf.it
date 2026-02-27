# C++ 지원

🌐 [English](cpp.md) | [한국어](cpp.ko.md) | [日本語](cpp.ja.md) | [हिन्दी](cpp.hi.md) | [Deutsch](cpp.de.md)

## 지원 확장자

- `.cpp`
- `.hpp`
- `.h`

## 추출 대상

| 요소 | Kind | 예시 |
|------|------|------|
| 클래스 | `class` | `class User { ... }` |
| 구조체 | `struct` | `struct Point { int x, y; }` |
| 메서드 | `method` | `void User::getName()` |
| 생성자 | `constructor` | `User(string name)` |
| 소멸자 | `destructor` | `~User()` |
| 함수 | `function` | `int add(int a, int b)` |
| 네임스페이스 | `namespace` | `namespace utils { }` |
| 템플릿 | `template` | `template<typename T> class Box` |
| 열거형 | `enum` | `enum Color { RED, GREEN }` |
| Typedef | `typedef` | `typedef unsigned int uint` |
| 매크로 | `macro` | `#define MAX_SIZE 100` |
| Include | (import) | `#include <iostream>` |
| 주석 | `doc` | `// Comment` |

## 예시

### 입력

```cpp
#include <iostream>
#include <string>

// 사용자 데이터 관리를 위한 User 클래스
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
    // 헬퍼 함수
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

### 출력 (XML)

```xml
<file path="example.hpp" language="cpp">
  <signature kind="class" line="5">
    <name>User</name>
    <text>class User</text>
    <doc>사용자 데이터 관리를 위한 User 클래스</doc>
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
    <doc>헬퍼 함수</doc>
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

## 특이사항

### 접근 제어

- 모든 접근 수준 (public, private, protected) 추출됨
- 가시성 수정자에 따른 필터링 없음
- AI가 완전한 클래스 구조를 이해하는 데 유용

### 본문 제거

`--include-body` 플래그 미사용 시:

- 함수/메서드: 여는 중괄호 `{` 이후 본문 제거
- 클래스/구조체/네임스페이스: 본문 제거, 선언부만 유지
- 템플릿: 기본 선언 본문 제거
- Enum/Typedef/Macro: 전체 텍스트 유지

### 템플릿 지원

기본 템플릿 지원 포함:

```cpp
template<typename T>
class Box { ... };         // 캡처됨

template<typename T>
T getMax(T a, T b) { ... } // 캡처됨
```

### 네임스페이스 지원

단순 및 중첩 네임스페이스 모두 캡처:

```cpp
namespace outer {
    namespace inner {
        void helper();     // 세 개 모두 캡처
    }
}
```

### Include 문

`--include-imports`를 사용하여 `#include` 지시문 추출:

```cpp
#include <iostream>        // 시스템 include
#include "myheader.h"      // 로컬 include
```

## 지원하지 않는 요소 (v1)

| 요소 | 이유 |
|------|------|
| 연산자 오버로드 | `operator+`, `operator<<` - 특수 케이스, 드문 사용 |
| Friend 선언 | `friend class Bar` - 접근 제어 예외 |
| Using 선언 | `using namespace std` - 단순 별칭 |
| Lambda 표현식 | `[](int x) { ... }` - 인라인 정의 |
| 템플릿 특수화 | `template<> class Box<int>` - 복잡한 파싱 |
| Variadic 템플릿 | `template<typename... Args>` - 고급 패턴 |
| C++20 Concepts | `template<Integral T>` - 제한된 컴파일러 지원 |
| C++20 Modules | `import std;` - 제한된 컴파일러 지원 |
| 전역 변수 | 향후 버전에서 추가될 수 있음 |
