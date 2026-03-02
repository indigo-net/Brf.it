# Rust 지원

[English](../../languages/rust.md) | [한국어](rust.md) | [日本語](../../ja/languages/rust.md) | [हिन्दी](../../hi/languages/rust.md) | [Deutsch](../../de/languages/rust.md)

## 지원 확장자

- `.rs`

## 추출 대상

| 요소 | Kind | 예시 |
|------|------|------|
| 함수 | `function` | `pub fn add(a: i32, b: i32) -> i32` |
| Async 함수 | `function` | `pub async fn fetch() -> Result<T>` |
| Unsafe 함수 | `function` | `pub unsafe fn dangerous()` |
| Const 함수 | `function` | `pub const fn compute() -> i32` |
| 구조체 | `struct` | `pub struct Point { x: f64, y: f64 }` |
| 열거형 | `enum` | `pub enum Color { Red, Green, Blue }` |
| 트레이트 | `trait` | `pub trait Drawable { fn draw(&self); }` |
| 타입 별칭 | `type` | `pub type Coordinate = (f64, f64)` |
| Impl 블록 | `impl` | `impl Point { ... }` |
| 상수 | `variable` | `pub const MAX_SIZE: usize = 1024` |
| 정적 변수 | `variable` | `pub static COUNTER: AtomicUsize = ...` |
| 모듈 | `namespace` | `pub mod utils { ... }` |
| 매크로 | `macro` | `macro_rules! say_hello { ... }` |
| 유니온 | `struct` | `pub union IntOrFloat { i: i32, f: f32 }` |
| 문서 주석 | `doc` | `/// 문서화` |

## 예시

### 입력

```rust
/// 2D 점 구조체.
pub struct Point {
    pub x: f64,
    pub y: f64,
}

impl Point {
    /// 주어진 좌표에 새 Point를 생성합니다.
    pub fn new(x: f64, y: f64) -> Self {
        Point { x, y }
    }

    /// 원점으로부터의 거리를 계산합니다.
    pub fn distance(&self) -> f64 {
        (self.x.powi(2) + self.y.powi(2)).sqrt()
    }
}

/// 허용되는 최대 점의 개수.
pub const MAX_POINTS: usize = 1000;
```

### 출력 (XML)

```xml
<file path="point.rs" language="rust">
  <type kind="struct" line="2">
    <name>Point</name>
    <text>pub struct Point</text>
    <doc>2D 점 구조체.</doc>
  </type>
  <type kind="impl" line="8">
    <name>Point</name>
    <text>impl Point</text>
  </type>
  <function kind="function" line="10">
    <name>new</name>
    <text>pub fn new(x: f64, y: f64) -> Self</text>
    <doc>주어진 좌표에 새 Point를 생성합니다.</doc>
  </function>
  <function kind="function" line="16">
    <name>distance</name>
    <text>pub fn distance(&self) -> f64</text>
    <doc>원점으로부터의 거리를 계산합니다.</doc>
  </function>
  <variable kind="variable" line="22">
    <name>MAX_POINTS</name>
    <text>pub const MAX_POINTS: usize = 1000</text>
    <doc>허용되는 최대 점의 개수.</doc>
  </variable>
</file>
```

## 참고사항

### 가시성 (Visibility)

- 모든 선언이 추출됩니다 (`pub` 및 private 모두)
- 가시성 수식어 (`pub`, `pub(crate)`, `pub(super)`)는 시그니처에 그대로 보존됩니다

### 함수 수식어

- `async`, `unsafe`, `const`, `extern` 함수는 모두 kind `function`으로 분류됩니다
- 수식어는 시그니처 텍스트에 보존됩니다

### 제네릭과 라이프타임

- 제네릭 타입 파라미터 (`<T>`, `<T: Clone>`)가 완전히 보존됩니다
- 라이프타임 어노테이션 (`'a`, `'static`)이 완전히 보존됩니다
- where 절도 시그니처에 포함됩니다

### Impl 블록

- `impl` 블록 자체와 내부 메서드가 모두 추출됩니다
- `impl Trait for Type` 패턴도 캡처됩니다

### 본문 제거

`--include-body` 플래그를 사용하지 않을 때:

- 함수/메서드: 여는 중괄호 `{` 이후 본문 제거
- 구조체/열거형/트레이트: 여는 중괄호 `{` 이후 본문 제거
- Impl 블록: 여는 중괄호 `{` 이후 본문 제거
- Const/Static: 값 표현식은 보존

### 문서 주석

- `///` 및 `//!` 문서 주석만 추출됩니다
- 일반 `//` 주석은 포함되지 않습니다

### 매크로

- `macro_rules!` 정의가 kind `macro`로 추출됩니다
- 프로시저 매크로와 매크로 호출은 추출되지 않습니다
