# Go サポート

🌐 [English](go.md) | [한국어](go.ko.md) | [日本語](go.ja.md) | [हिन्दी](go.hi.md) | [Deutsch](go.de.md)

## サポート拡張子

- `.go`

## 抽出対象

| 要素 | Kind | 例 |
|------|------|-----|
| 関数 | `function` | `func DoSomething()` |
| メソッド | `method` | `func (s *Server) Start()` |
| 型（struct、interfaceなど） | `type` | `type User struct {...}` |
| コメント | `doc` | `// Comment` |

## 例

### 入力

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

### 出力（XML）

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

## 注意事項

### エクスポート判定

- Goのエクスポートルール適用：大文字で始まる識別子のみ抽出
- 小文字で始まるプライベート関数/型はデフォルトで除外

### メソッド vs 関数

- receiverがある宣言は`method`に分類
- receiverがない宣言は`function`に分類

### 本体削除

`--include-body`フラグ未使用時：

- 関数/メソッド：中括弧`{`以降の本体を削除
- 型：`struct`または`interface`キーワードまで保持

### サポートされていない要素

- パッケージレベル変数（`var`）
- 定数（`const`）
- 埋め込み関数（関数内の関数）
