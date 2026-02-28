# Go सपोर्ट

🌐 [English](../../languages/go.md) | [한국어](../../ko/languages/go.md) | [日本語](../../ja/languages/go.md) | [हिन्दी](go.md) | [Deutsch](../../de/languages/go.md)

## समर्थित एक्सटेंशन

- `.go`

## निष्कर्षण लक्ष्य

| तत्व | Kind | उदाहरण |
|------|------|--------|
| फंक्शन | `function` | `func DoSomething()` |
| मेथड | `method` | `func (s *Server) Start()` |
| टाइप (struct, interface आदि) | `type` | `type User struct {...}` |
| कॉन्स्ट/वार | `variable` | `const MaxSize = 100` |
| कमेंट | `doc` | `// Comment` |

## उदाहरण

### इनपुट

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

### आउटपुट (XML)

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

## नोट्स

### एक्सपोर्ट डिटेक्शन

- Go एक्सपोर्ट नियम लागू: केवल अपरकेस से शुरू होने वाले आइडेंटिफायर निकाले जाते हैं
- लोअरकेस से शुरू होने वाले प्राइवेट फंक्शन/टाइप डिफ़ॉल्ट रूप से बाहर रखे जाते हैं

### मेथड vs फंक्शन

- receiver वाली घोषणाएं `method` के रूप में वर्गीकृत
- बिना receiver वाली घोषणाएं `function` के रूप में वर्गीकृत

### बॉडी रिमूवल

`--include-body` फ्लैग का उपयोग न करने पर:

- फंक्शन/मेथड: ओपनिंग ब्रेस `{` के बाद बॉडी हटा दी जाती है
- टाइप: केवल `struct` या `interface` कीवर्ड संरक्षित रहता है

### असमर्थित तत्व

- एम्बेडेड फंक्शन (फंक्शन के अंदर फंक्शन)
