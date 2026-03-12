# Scala 지원

[English](../../languages/scala.md) | [한국어](scala.md) | [日本語](../../ja/languages/scala.md) | [हिन्दी](../../hi/languages/scala.md) | [Deutsch](../../de/languages/scala.md)

## 지원 확장자

- `.scala`
- `.sc`

## 추출 대상

| 요소 | Kind | XML Tag | 예시 |
|------|------|---------|------|
| 메서드 (본문 포함) | `method` | `<function>` | `def add(a: Int, b: Int): Int` |
| 메서드 (추상) | `method` | `<function>` | `def greet(name: String): String` |
| 클래스 | `class` | `<type>` | `class Person(val name: String)` |
| 트레이트 | `trait` | `<type>` | `trait Greeter` |
| 오브젝트 | `class` | `<type>` | `object MathUtils` |
| val | `variable` | `<variable>` | `val PI: Double = 3.14159` |
| var | `variable` | `<variable>` | `var count: Int = 0` |
| 타입 별칭 | `type` | `<type>` | `type StringList = List[String]` |
| Enum (Scala 3) | `enum` | `<type>` | `enum Color` |
| Given (Scala 3) | `variable` | `<variable>` | `given ordering: Ordering[Int]` |

## 예시

### 입력

```scala
// 사용자 관리
trait Greeter {
  def greet(name: String): String
}

class Person(val name: String) extends Greeter {
  def greet(name: String): String = s"Hello, $name"
}

object MathUtils {
  val PI: Double = 3.14159
  def add(a: Int, b: Int): Int = a + b
}

type StringList = List[String]
```

### 출력 (XML)

```xml
<file path="example.scala" language="scala">
  <type>trait Greeter</type>
  <function>def greet(name: String): String</function>
  <type>class Person(val name: String) extends Greeter</type>
  <function>def greet(name: String): String</function>
  <type>object MathUtils</type>
  <variable>val PI: Double = 3.14159</variable>
  <function>def add(a: Int, b: Int): Int</function>
  <type>type StringList = List[String]</type>
</file>
```

## 참고사항

### 가시성 (Visibility)

- 가시성 수정자에 관계없이 모든 선언이 추출됩니다
- 가시성 수정자(`private`, `protected`)는 시그니처 텍스트에 보존됩니다

### 클래스 변형

- `class`, `abstract class`, `case class`, `sealed class`, `implicit class`는 모두 kind `class`로 분류됩니다
- `trait`과 `sealed trait`은 kind `trait`으로 분류됩니다
- `object` (싱글톤 및 컴패니언)는 kind `class`로 분류됩니다

### 본문 제거

`--include-body` 플래그를 사용하지 않을 때:

- 메서드: `=` 이후의 본문이 제거되며, 반환 타입은 보존됩니다
- 클래스/트레이트/오브젝트: `{ }` 안의 본문이 제거되고, 선언 줄만 보존됩니다
- val/var: 값이 보존됩니다 (`lazy val`, `implicit val` 포함)
- 타입 별칭: 전체 보존

`--include-private`를 사용하여 비공개/unexported 심볼 포함.

### 제네릭

- 제네릭 타입 파라미터 `[A, B]`가 시그니처에 완전히 보존됩니다
- 컨텍스트 바운드와 뷰 바운드가 포함됩니다

### Scala 3 기능

- `enum` 정의는 kind `enum`으로 분류됩니다
- 이름이 있는 `given` 인스턴스는 kind `variable`로 분류됩니다
- `extension` 메서드는 개별적으로 `method`로 추출됩니다 (extension 선언 자체는 캡처되지 않음)
