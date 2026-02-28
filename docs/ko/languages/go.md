# Go 지원

🌐 [English](../../languages/go.md) | [한국어](go.md) | [日本語](../../ja/languages/go.md) | [हिन्दी](../../hi/languages/go.md) | [Deutsch](../../de/languages/go.md)

## 지원 확장자

- `.go`

## 추출 대상

| 요소 | Kind | 예시 |
|------|------|------|
| 함수 | `function` | `func DoSomething()` |
| 메서드 | `method` | `func (s *Server) Start()` |
| 타입 (struct, interface 등) | `type` | `type User struct {...}` |
| 상수/변수 | `variable` | `const MaxSize = 100` |
| 주석 | `doc` | `// Comment` |

## 예시

### 입력

```go
// Server handles HTTP requests.
type Server struct {
    port int
}

// NewServer creates a new Server instance.
func NewServer(port int) *Server {
    return &Server{port: port}
}

// Start begins listening for requests.
func (s *Server) Start() error {
    return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}
```

### 출력 (XML)

```xml
<file path="server.go" language="go">
  <signature kind="type" line="2">
    <name>Server</name>
    <text>type Server struct</text>
    <doc>Server handles HTTP requests.</doc>
  </signature>
  <signature kind="function" line="7">
    <name>NewServer</name>
    <text>func NewServer(port int) *Server</text>
    <doc>NewServer creates a new Server instance.</doc>
  </signature>
  <signature kind="method" line="12">
    <name>Start</name>
    <text>func (s *Server) Start() error</text>
    <doc>Start begins listening for requests.</doc>
  </signature>
</file>
```

## 특이사항

### Export 판별

- Go의 export 규칙 적용: 대문자로 시작하는 식별자만 추출
- 소문자로 시작하는 private 함수/타입은 기본적으로 제외됨

### 메서드 vs 함수

- receiver가 있는 선언은 `method`로 분류
- receiver가 없는 선언은 `function`으로 분류

### 본문 제거

`--include-body` 플래그 미사용 시:

- 함수/메서드: 중괄호 `{` 이후 본문 제거
- 타입: `struct` 또는 `interface` 키워드까지만 유지

### 지원하지 않는 요소

- 임베디드 함수 (함수 내부 함수)
