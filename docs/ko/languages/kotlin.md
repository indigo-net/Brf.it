# Kotlin 지원

[English](../../languages/kotlin.md) | [한국어](kotlin.md) | [日本語](../../ja/languages/kotlin.md) | [हिन्दी](../../hi/languages/kotlin.md) | [Deutsch](../../de/languages/kotlin.md)

## 지원 확장자

- `.kt`
- `.kts`

## 추출 대상

| 요소 | Kind | 예시 |
|------|------|------|
| 함수 | `function` | `fun add(a: Int, b: Int): Int` |
| Suspend 함수 | `function` | `suspend fun fetchData(url: String): String` |
| 확장 함수 | `function` | `fun String.isEmail(): Boolean` |
| 클래스 | `class` | `class User(val name: String)` |
| 데이터 클래스 | `class` | `data class Point(val x: Double, val y: Double)` |
| Sealed 클래스 | `class` | `sealed class Result<out T>` |
| Enum 클래스 | `enum` | `enum class Color { RED, GREEN, BLUE }` |
| 인터페이스 | `interface` | `interface Repository<T>` |
| 오브젝트 | `class` | `object AppConfig` |
| 컴패니언 오브젝트 | `class` | `companion object Factory` |
| 프로퍼티 (val/var) | `variable` | `val MAX_SIZE = 100` |
| 타입 별칭 | `type` | `typealias Handler<T> = (T) -> Unit` |
| Enum 엔트리 | `variable` | `RED("#FF0000")` |
| 보조 생성자 | `constructor` | `constructor(name: String)` |
| 문서 주석 | `doc` | `/** 문서화 */` |

## 예시

### 입력

```kotlin
/** API 응답을 위한 사용자 데이터 클래스. */
data class User(
    val id: Long,
    val name: String,
    val email: String
) {
    fun isValid(): Boolean = email.contains("@")
}

/** 사용자 작업을 위한 리포지토리 인터페이스. */
interface UserRepository {
    suspend fun getUser(id: Long): User?
    fun save(user: User): Boolean
}

val DEFAULT_TIMEOUT: Long = 5000L
```

### 출력 (XML)

```xml
<file path="user.kt" language="kotlin">
  <type kind="class" line="2">
    <name>User</name>
    <text>data class User(
    val id: Long,
    val name: String,
    val email: String
)</text>
    <doc>API 응답을 위한 사용자 데이터 클래스.</doc>
  </type>
  <function kind="function" line="7">
    <name>isValid</name>
    <text>fun isValid(): Boolean = email.contains("@")</text>
  </function>
  <type kind="interface" line="11">
    <name>UserRepository</name>
    <text>interface UserRepository</text>
    <doc>사용자 작업을 위한 리포지토리 인터페이스.</doc>
  </type>
  <function kind="function" line="12">
    <name>getUser</name>
    <text>suspend fun getUser(id: Long): User?</text>
  </function>
  <function kind="function" line="13">
    <name>save</name>
    <text>fun save(user: User): Boolean</text>
  </function>
  <variable kind="variable" line="16">
    <name>DEFAULT_TIMEOUT</name>
    <text>val DEFAULT_TIMEOUT: Long = 5000L</text>
  </variable>
</file>
```

## 참고사항

### 가시성 (Visibility)

- 모든 선언이 추출됩니다 (Kotlin은 기본적으로 `public`)
- 접근 수식어 (`public`, `internal`, `private`, `protected`)는 시그니처에 그대로 보존됩니다

### 함수 수식어

- `suspend`, `inline`, `infix`, `operator`, `tailrec` 함수는 모두 kind `function`으로 분류됩니다
- 수식어는 시그니처 텍스트에 보존됩니다
- 단일 표현식 함수 (`fun double(x: Int) = x * 2`)는 전체가 보존됩니다

### 제네릭

- 제네릭 타입 파라미터 (`<T>`, `<T : Comparable<T>>`)가 완전히 보존됩니다
- `where` 절과 변성 어노테이션 (`in`, `out`)이 시그니처에 포함됩니다
- `reified` 타입 파라미터도 보존됩니다

### 클래스

- `data class`, `sealed class`, `abstract class`, `open class`, `inner class`, `annotation class`, `value class`는 모두 kind `class`로 분류됩니다
- `enum class`는 kind `enum`으로 분류됩니다
- `interface`와 `sealed interface`는 kind `interface`로 분류됩니다

### 오브젝트

- `object` 선언 (싱글턴)은 kind `class`로 분류됩니다
- `companion object` 블록이 추출됩니다; 이름 없는 컴패니언은 "Companion"이라는 합성 이름을 부여받습니다

### 본문 제거

`--include-body` 플래그를 사용하지 않을 때:

- 함수/메서드: 여는 중괄호 `{` 이후 본문 제거
- 단일 표현식 함수: 전체 보존 (표현식 자체가 시그니처의 일부)
- 클래스/인터페이스/열거형: 여는 중괄호 `{` 이후 본문 제거
- 프로퍼티 (val/var): 값 표현식은 보존
- 타입 별칭: 전체 보존

### 문서 주석

- `/** ... */` (KDoc) 및 `//` 줄 주석 모두 추출됩니다
- KDoc 주석은 뒤따르는 선언과 연결됩니다
