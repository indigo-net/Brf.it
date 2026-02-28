# Go Unterstützung

🌐 [English](../../languages/go.md) | [한국어](../../ko/languages/go.md) | [日本語](../../ja/languages/go.md) | [हिन्दी](../../hi/languages/go.md) | [Deutsch](go.md)

## Unterstützte Erweiterungen

- `.go`

## Extraktionsziele

| Element | Kind | Beispiel |
|---------|------|----------|
| Funktion | `function` | `func DoSomething()` |
| Methode | `method` | `func (s *Server) Start()` |
| Typ (struct, interface usw.) | `type` | `type User struct {...}` |
| Const/Var | `variable` | `const MaxSize = 100` |
| Kommentar | `doc` | `// Comment` |

## Beispiel

### Eingabe

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

### Ausgabe (XML)

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

## Hinweise

### Export-Erkennung

- Go-Exportregeln angewendet: nur Bezeichner, die mit Großbuchstaben beginnen, werden extrahiert
- Private Funktionen/Typen, die mit Kleinbuchstaben beginnen, werden standardmäßig ausgeschlossen

### Methode vs Funktion

- Deklarationen mit Receiver werden als `method` klassifiziert
- Deklarationen ohne Receiver werden als `function` klassifiziert

### Body-Entfernung

Wenn `--include-body` Flag nicht verwendet wird:

- Funktionen/Methoden: Body nach öffnender Klammer `{` entfernt
- Typen: nur `struct` oder `interface` Schlüsselwort wird beibehalten

### Nicht unterstützte Elemente

- Eingebettete Funktionen (Funktionen innerhalb von Funktionen)
