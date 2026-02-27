# 헤더 및 스키마 설명 개선 설계

## 요약

출력 헤더를 "Brf.it Output"에서 경로 기반 제목으로 변경하고, XML에 태그 스키마 설명을 추가한다.

## 변경 사항

### 1. Markdown 헤더

**Before:**
```markdown
# Brf.it Output
```

**After:**
```markdown
# Code Summary: /Users/jefflee/projects/Brf.it

*brf.it v0.9.2*
```

### 2. XML 헤더 및 스키마

**Before:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<brfit>
```

**After:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<!-- brf.it v0.9.2 | Code Summary: /Users/jefflee/projects/Brf.it -->
<!--
Schema:
| Tag       | Description                              |
|-----------|------------------------------------------|
| file      | Source file (path, language attributes)  |
| signature | Function, type, or variable declaration  |
| imports   | Import statements container              |
| import    | Single import statement                  |
| doc       | Documentation comment                    |
| error     | Parse error message                      |
-->
<brfit>
```

## 수정 파일

| 파일 | 변경 내용 |
|------|----------|
| `pkg/formatter/types.go` | `PackageData`에 `RootPath`, `Version` 필드 추가 |
| `pkg/formatter/markdown.go` | 헤더를 `# Code Summary: {path}` + 버전으로 변경 |
| `pkg/formatter/xml.go` | 헤더 주석 + 스키마 설명 추가 |
| `internal/context/context.go` | `RootPath`, `Version` 설정 |
| `cmd/brfit/root.go` | 버전 정보를 context로 전달 |

## 구현 세부사항

### types.go 수정

```go
type PackageData struct {
    RootPath string  // 패키징 대상 경로
    Version  string  // brf.it 버전
    // ... 기존 필드
}
```

### markdown.go 수정

```go
func (f *MarkdownFormatter) Format(data *PackageData) ([]byte, error) {
    var buf bytes.Buffer

    // 헤더
    buf.WriteString(fmt.Sprintf("# Code Summary: %s\n\n", data.RootPath))
    buf.WriteString(fmt.Sprintf("*brf.it %s*\n\n", data.Version))
    // ...
}
```

### xml.go 수정

```go
func (f *XMLFormatter) Format(data *PackageData) ([]byte, error) {
    var buf bytes.Buffer

    buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
    buf.WriteByte('\n')
    buf.WriteString(fmt.Sprintf("<!-- brf.it %s | Code Summary: %s -->\n", data.Version, data.RootPath))
    buf.WriteString(xmlSchemaComment)
    buf.WriteString("<brfit>\n")
    // ...
}

const xmlSchemaComment = `<!--
Schema:
| Tag       | Description                              |
|-----------|------------------------------------------|
| file      | Source file (path, language attributes)  |
| signature | Function, type, or variable declaration  |
| imports   | Import statements container              |
| import    | Single import statement                  |
| doc       | Documentation comment                    |
| error     | Parse error message                      |
-->
`
```

## 검증

```bash
# 빌드
go build -o brfit ./cmd/brfit

# Markdown 출력 확인
./brfit . -f md | head -10

# XML 출력 확인
./brfit . -f xml | head -20

# 테스트
go test ./pkg/formatter/...
```
