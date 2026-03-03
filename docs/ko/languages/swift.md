# Swift 지원

[English](../../languages/swift.md) | [한국어](swift.md) | [日本語](../../ja/languages/swift.md) | [हिन्दी](../../hi/languages/swift.md) | [Deutsch](../../de/languages/swift.md)

## 지원 확장자

- `.swift`

## 추출 대상

| 요소 | Kind | 예시 |
|------|------|------|
| 함수 | `function` | `public func greet(name: String) -> String` |
| Async 함수 | `function` | `public func fetch() async throws -> Data` |
| 구조체 | `struct` | `public struct Point { let x: Double }` |
| 클래스 | `class` | `public class Vehicle { var speed: Int }` |
| 열거형 | `enum` | `public enum Direction { case north }` |
| 프로토콜 | `interface` | `public protocol Drawable { func draw() }` |
| 확장 | `type` | `extension Point: Drawable` |
| 타입 별칭 | `type` | `public typealias Coordinate = (Double, Double)` |
| 프로퍼티 (let/var) | `variable` | `public let PI = 3.14159` |
| 이니셜라이저 | `constructor` | `init(value: Int)` |
| 디이니셜라이저 | `destructor` | `deinit` |
| 서브스크립트 | `method` | `subscript(index: Int) -> Int` |
| 연산자 | `function` | `prefix operator +++` |
| 문서 주석 | `doc` | `/// 문서화` |

## 예시

### 입력

```swift
/// 2D 점 구조체.
public struct Point {
    public let x: Double
    public let y: Double

    /// 원점으로부터의 거리를 계산합니다.
    public func distance() -> Double {
        return (x * x + y * y).squareRoot()
    }
}

/// 허용되는 최대 점의 개수.
public let MAX_POINTS = 1000
```

### 출력 (XML)

```xml
<file path="point.swift" language="swift">
  <type kind="struct" line="2">
    <name>Point</name>
    <text>public struct Point</text>
    <doc>2D 점 구조체.</doc>
  </type>
  <function kind="function" line="8">
    <name>distance</name>
    <text>public func distance() -> Double</text>
    <doc>원점으로부터의 거리를 계산합니다.</doc>
  </function>
  <variable kind="variable" line="14">
    <name>MAX_POINTS</name>
    <text>public let MAX_POINTS = 1000</text>
    <doc>허용되는 최대 점의 개수.</doc>
  </variable>
</file>
```

## 참고사항

### 가시성 (Visibility)

- 모든 선언이 추출됩니다 (`public` 및 `internal`/`private` 모두)
- 접근 수식어 (`public`, `internal`, `private`, `fileprivate`, `open`)는 시그니처에 그대로 보존됩니다

### 함수 수식어

- `async`, `throws`, `@discardableResult` 함수는 모두 kind `function`으로 분류됩니다
- 수식어는 시그니처 텍스트에 보존됩니다

### 제네릭

- 제네릭 타입 파라미터 (`<T>`, `<T: Equatable>`)가 완전히 보존됩니다
- 제네릭 where 절도 시그니처에 포함됩니다

### 확장 (Extension)

- `extension` 블록 자체와 내부 멤버가 모두 추출됩니다
- `extension Type: Protocol` 준수 패턴도 캡처됩니다

### 본문 제거

`--include-body` 플래그를 사용하지 않을 때:

- 함수/메서드: 여는 중괄호 `{` 이후 본문 제거
- 구조체/클래스/열거형: 여는 중괄호 `{` 이후 본문 제거
- 확장: 여는 중괄호 `{` 이후 본문 제거
- 프로퍼티 (let/var): 값 표현식은 보존

### 문서 주석

- `///` 문서 주석만 추출됩니다
- 일반 `//` 주석은 포함되지 않습니다
