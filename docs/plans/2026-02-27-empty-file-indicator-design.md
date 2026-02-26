# 빈 파일 표시 기능 설계

## 요약

출력할 내용이 없는 파일(시그니처 0개, imports 0개)에 대해 빈 코드블록 대신 `// (empty)` 코멘트를 표시한다.

## 배경

현재 `cmd/brfit/main.go`처럼 export된 심볼이 없는 파일은 빈 코드블록으로 출력됨:
```markdown
### cmd/brfit/main.go

```go
```
```

이는 정보가 없어 보이고, 파일이 실제로 비어있는지 파싱 결과가 없는지 구분이 안 됨.

## 설계

### 판단 기준

"출력할 내용이 없음" = 다음 조건 모두 충족:
- `len(signatures) == 0`
- `!includeImports || len(imports) == 0`
- `error == nil` (에러는 별도 표시)

### 출력 형식

**Markdown:**
```markdown
### cmd/brfit/main.go

```go
// (empty)
```
```

**XML:**
```xml
<file path="cmd/brfit/main.go" language="go">
  <!-- empty -->
</file>
```

### 수정 파일

| 파일 | 변경 내용 |
|------|----------|
| `pkg/formatter/markdown.go` | 빈 파일에 `// (empty)` 출력 |
| `pkg/formatter/xml.go` | 빈 파일에 `<!-- empty -->` 출력 |

## 구현 계획

### markdown.go

```go
// 파일 루프 내에서
if file.Error != nil {
    // 기존 에러 처리
} else if len(file.Signatures) == 0 && (!data.IncludeImports || len(file.Imports) == 0) {
    // 빈 파일: 코멘트만 출력
    buf.WriteString(fmt.Sprintf("```%s\n", file.Language))
    buf.WriteString(getEmptyComment(file.Language))
    buf.WriteString("\n```\n\n")
} else {
    // 기존 시그니처 출력
}
```

### xml.go

```go
// 파일 루프 내에서
if file.Error != nil {
    // 기존 에러 처리
} else if len(file.Signatures) == 0 && (!data.IncludeImports || len(file.Imports) == 0) {
    // 빈 파일: 코멘트만 출력
    buf.WriteString("      <!-- empty -->\n")
} else {
    // 기존 시그니처 출력
}
```

### 언어별 빈 파일 코멘트

```go
func getEmptyComment(lang string) string {
    switch lang {
    case "python":
        return "# (empty)"
    case "html", "xml":
        return "<!-- (empty) -->"
    default:
        return "// (empty)"
    }
}
```

## 검증

```bash
# 빌드
go build -o brfit ./cmd/brfit

# Markdown 확인
./brfit . -f md --include-imports | grep -A3 "main.go"

# XML 확인
./brfit . -f xml --include-imports | grep -A2 "main.go"

# 테스트
go test ./pkg/formatter/...
```
