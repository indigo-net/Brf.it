# Java 지원

🌐 [English](../../languages/java.md) | [한국어](java.md) | [日本語](../../ja/languages/java.md) | [हिन्दी](../../hi/languages/java.md) | [Deutsch](../../de/languages/java.md)

## 지원 확장자

- `.java`

## 추출 대상

| 요소 | Kind | 예시 |
|------|------|------|
| 클래스 | `class` | `public class User { ... }` |
| 인터페이스 | `interface` | `public interface Repository<T> { ... }` |
| 메서드 | `method` | `public String getName() { ... }` |
| 생성자 | `constructor` | `public User(String name) { ... }` |
| Enum | `enum` | `public enum Status { ... }` |
| 어노테이션 | `annotation` | `public @interface Inject { ... }` |
| Record (Java 14+) | `record` | `public record Point(int x, int y) { ... }` |
| 필드 | `field` | `public static final String API = "..."` |
| 주석 | `doc` | `// Comment` 또는 `/* Block */` |

## 예시

### 입력

```java
package com.example;

/**
 * User class represents a user in the system.
 */
public class User {
    private String name;

    public User(String name) {
        this.name = name;
    }

    public String getName() {
        return name;
    }

    private void internalMethod() {
        // Private method
    }
}

public interface Repository<T> {
    T findById(String id);
    void save(T entity);
}

public enum Status {
    PENDING, ACTIVE, COMPLETED
}

public record Point(int x, int y) {}
```

### 출력 (XML)

```xml
<file path="User.java" language="java">
  <signature kind="class" line="6">
    <name>User</name>
    <text>public class User</text>
  </signature>
  <signature kind="constructor" line="9">
    <name>User</name>
    <text>public User(String name)</text>
  </signature>
  <signature kind="method" line="13">
    <name>getName</name>
    <text>public String getName()</text>
  </signature>
  <signature kind="interface" line="22">
    <name>Repository</name>
    <text>public interface Repository&lt;T&gt;</text>
  </signature>
  <signature kind="method" line="23">
    <name>findById</name>
    <text>T findById(String id);</text>
  </signature>
  <signature kind="method" line="24">
    <name>save</name>
    <text>void save(T entity);</text>
  </signature>
  <signature kind="enum" line="27">
    <name>Status</name>
    <text>public enum Status</text>
  </signature>
  <signature kind="record" line="31">
    <name>Point</name>
    <text>public record Point(int x, int y)</text>
  </signature>
</file>
```

## 특이사항

### 가시성 필터링

- `public`, `protected`, package-private (기본값): 기본적으로 추출됨
- `private`: `--include-body` 사용 시에만 포함됨

### 제네릭 처리

제네릭 타입 매개변수가 시그니처에 포함됨:

```java
public class Box<T extends Comparable<T>>  // 전체 캡처
public <U> U transform(Function<T, U> fn)  // 메서드 타입 매개변수 포함
```

### 어노테이션 출력

메서드와 클래스 어노테이션이 시그니처 텍스트에 포함됨:

```java
@Override
public String toString()  // 시그니처에 @Override 포함
```

### Record 지원 (Java 14+)

Record는 컴포넌트 매개변수와 함께 추출됨:

```java
public record User(String name, int age)  // 컴포넌트 유지
```

### 내부/중첩 클래스

모든 중첩 클래스가 별도 시그니처로 추출됨:

```java
public class Outer {
    public static class Nested { ... }  // 별도로 추출
    public class Inner { ... }          // 역시 추출
}
```

### 추상 메서드

인터페이스의 추상 메서드는 `;`로 끝남 (본문 없음):

```java
interface Foo {
    void bar();  // 그대로 캡처
}
```

### 본문 제거

`--include-body` 플래그 미사용 시:

- 메서드/생성자: 여는 중괄호 `{` 이후 본문 제거
- 클래스/인터페이스/Enum: 여는 중괄호 `{` 이후 본문 제거
- 추상 메서드: 그대로 유지 (`;`로 끝남)

`--include-private`를 사용하여 비공개/unexported 심볼 포함.

### Javadoc (향후 지원)

- 현재 버전: 선언부 위의 `//` 및 `/* */` 주석이 doc으로 캡처됨
- 향후 버전: `/** */` Javadoc 파싱 지원 예정

### 지원하지 않는 요소

- static 초기화 블록
- 익명 클래스
- Lambda 표현식
