# Go Support

🌐 [English](go.md) | [한국어](../ko/languages/go.md) | [日本語](../ja/languages/go.md) | [हिन्दी](../hi/languages/go.md) | [Deutsch](../de/languages/go.md)

## Supported Extensions

- `.go`

## Extraction Targets

| Element | Kind | Example |
|---------|------|---------|
| Function | `function` | `func DoSomething()` |
| Method | `method` | `func (s *Server) Start()` |
| Type (struct, interface, etc.) | `type` | `type User struct {...}` |
| Const/Var | `variable` | `const MaxSize = 100` |
| Comment | `doc` | `// Comment` |

## Example

### Input

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

### Output (XML)

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

## Notes

### Export Detection

- Go export rules applied: only identifiers starting with uppercase are extracted
- Private functions/types starting with lowercase are excluded by default

### Method vs Function

- Declarations with receiver are classified as `method`
- Declarations without receiver are classified as `function`

### Body Removal

When `--include-body` flag is not used:

- Functions/Methods: body removed after opening brace `{`
- Types: only `struct` or `interface` keyword is preserved

### Unsupported Elements

- Embedded functions (functions inside functions)
