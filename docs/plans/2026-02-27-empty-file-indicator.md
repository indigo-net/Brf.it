# 빈 파일 표시 기능 구현 계획

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 출력할 내용이 없는 파일에 `// (empty)` 코멘트를 표시하여 파일 존재를 알림

**Architecture:** Formatter 단계에서 시그니처와 imports가 모두 비어있으면 빈 파일 코멘트 출력. 언어별 코멘트 스타일 지원.

**Tech Stack:** Go, pkg/formatter 패키지

---

## Task 1: 언어별 빈 파일 코멘트 헬퍼 함수

**Files:**
- Create: `pkg/formatter/helpers.go`
- Test: `pkg/formatter/helpers_test.go`

**Step 1: 테스트 파일 생성**

```go
// pkg/formatter/helpers_test.go
package formatter

import "testing"

func TestGetEmptyComment(t *testing.T) {
	tests := []struct {
		lang     string
		expected string
	}{
		{"go", "// (empty)"},
		{"typescript", "// (empty)"},
		{"javascript", "// (empty)"},
		{"java", "// (empty)"},
		{"c", "// (empty)"},
		{"python", "# (empty)"},
		{"ruby", "# (empty)"},
		{"html", "<!-- (empty) -->"},
		{"xml", "<!-- (empty) -->"},
		{"unknown", "// (empty)"},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			result := getEmptyComment(tt.lang)
			if result != tt.expected {
				t.Errorf("getEmptyComment(%q) = %q, want %q", tt.lang, result, tt.expected)
			}
		})
	}
}
```

**Step 2: 테스트 실행 (실패 확인)**

Run: `go test ./pkg/formatter/... -run TestGetEmptyComment -v`
Expected: FAIL - `undefined: getEmptyComment`

**Step 3: 헬퍼 함수 구현**

```go
// pkg/formatter/helpers.go
package formatter

// getEmptyComment returns the appropriate empty file comment for a language.
func getEmptyComment(lang string) string {
	switch lang {
	case "python", "ruby":
		return "# (empty)"
	case "html", "xml":
		return "<!-- (empty) -->"
	default:
		return "// (empty)"
	}
}
```

**Step 4: 테스트 실행 (통과 확인)**

Run: `go test ./pkg/formatter/... -run TestGetEmptyComment -v`
Expected: PASS

**Step 5: 커밋**

```bash
git add pkg/formatter/helpers.go pkg/formatter/helpers_test.go
git commit -m "feat: 언어별 빈 파일 코멘트 헬퍼 함수 추가"
```

---

## Task 2: Markdown 포맷터 빈 파일 처리

**Files:**
- Modify: `pkg/formatter/markdown.go:74-94`
- Test: `pkg/formatter/formatter_test.go`

**Step 1: 테스트 추가**

```go
// pkg/formatter/formatter_test.go에 추가
func TestMarkdownFormatterEmptyFile(t *testing.T) {
	formatter := NewMarkdownFormatter()
	data := &PackageData{
		Files: []FileData{
			{
				Path:       "cmd/main.go",
				Language:   "go",
				Signatures: []SignatureData{}, // 빈 시그니처
				Imports:    []ImportData{},    // 빈 imports
			},
		},
		IncludeImports: false,
	}

	result, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	output := string(result)
	if !strings.Contains(output, "// (empty)") {
		t.Errorf("Expected empty comment, got:\n%s", output)
	}
}

func TestMarkdownFormatterEmptyFileWithImports(t *testing.T) {
	formatter := NewMarkdownFormatter()
	data := &PackageData{
		Files: []FileData{
			{
				Path:       "cmd/main.go",
				Language:   "go",
				Signatures: []SignatureData{},
				Imports: []ImportData{
					{Type: "import", Path: "fmt"},
				},
			},
		},
		IncludeImports: true,
	}

	result, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	output := string(result)
	// imports가 있으면 빈 파일이 아님
	if strings.Contains(output, "// (empty)") {
		t.Errorf("Should not show empty comment when imports exist, got:\n%s", output)
	}
}
```

**Step 2: 테스트 실행 (실패 확인)**

Run: `go test ./pkg/formatter/... -run TestMarkdownFormatterEmpty -v`
Expected: FAIL

**Step 3: markdown.go 수정**

`pkg/formatter/markdown.go` 74-94행을 다음으로 교체:

```go
		if file.Error != nil {
			buf.WriteString("> **Error:** ")
			buf.WriteString(escapeMarkdown(file.Error.Error()))
			buf.WriteString("\n\n")
		} else {
			// 빈 파일 확인: 시그니처 없고, (imports 포함 안 함 또는 imports도 없음)
			isEmpty := len(file.Signatures) == 0 && (!data.IncludeImports || len(file.Imports) == 0)

			buf.WriteString(fmt.Sprintf("```%s\n", file.Language))
			if isEmpty {
				buf.WriteString(getEmptyComment(file.Language))
				buf.WriteString("\n")
			} else {
				for _, sig := range file.Signatures {
					buf.WriteString(sig.Text)
					buf.WriteString("\n")
				}
			}
			buf.WriteString("```\n\n")

			// Add docs as quotes (빈 파일이면 건너뜀)
			if !isEmpty {
				for _, sig := range file.Signatures {
					if sig.Doc != "" {
						buf.WriteString("> ")
						buf.WriteString(escapeMarkdown(sig.Doc))
						buf.WriteString("\n\n")
					}
				}
			}
		}
```

**Step 4: 테스트 실행 (통과 확인)**

Run: `go test ./pkg/formatter/... -run TestMarkdownFormatterEmpty -v`
Expected: PASS

**Step 5: 커밋**

```bash
git add pkg/formatter/markdown.go pkg/formatter/formatter_test.go
git commit -m "feat: Markdown 포맷터 빈 파일에 // (empty) 표시"
```

---

## Task 3: XML 포맷터 빈 파일 처리

**Files:**
- Modify: `pkg/formatter/xml.go:70-86`
- Test: `pkg/formatter/formatter_test.go`

**Step 1: 테스트 추가**

```go
// pkg/formatter/formatter_test.go에 추가
func TestXMLFormatterEmptyFile(t *testing.T) {
	formatter := NewXMLFormatter()
	data := &PackageData{
		Files: []FileData{
			{
				Path:       "cmd/main.go",
				Language:   "go",
				Signatures: []SignatureData{},
				Imports:    []ImportData{},
			},
		},
		IncludeImports: false,
	}

	result, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	output := string(result)
	if !strings.Contains(output, "<!-- empty -->") {
		t.Errorf("Expected empty comment, got:\n%s", output)
	}
}
```

**Step 2: 테스트 실행 (실패 확인)**

Run: `go test ./pkg/formatter/... -run TestXMLFormatterEmptyFile -v`
Expected: FAIL

**Step 3: xml.go 수정**

`pkg/formatter/xml.go` 70-86행을 다음으로 교체:

```go
		if file.Error != nil {
			buf.WriteString("      <error>")
			buf.WriteString(escapeXML(file.Error.Error()))
			buf.WriteString("</error>\n")
		} else {
			// 빈 파일 확인
			isEmpty := len(file.Signatures) == 0 && (!data.IncludeImports || len(file.Imports) == 0)

			if isEmpty {
				buf.WriteString("      <!-- empty -->\n")
			} else {
				for _, sig := range file.Signatures {
					buf.WriteString("      <signature>")
					buf.WriteString(escapeXML(sig.Text))
					buf.WriteString("</signature>\n")

					if sig.Doc != "" {
						buf.WriteString("      <doc>")
						buf.WriteString(escapeXML(sig.Doc))
						buf.WriteString("</doc>\n")
					}
				}
			}
		}
```

**Step 4: 테스트 실행 (통과 확인)**

Run: `go test ./pkg/formatter/... -run TestXMLFormatterEmptyFile -v`
Expected: PASS

**Step 5: 커밋**

```bash
git add pkg/formatter/xml.go pkg/formatter/formatter_test.go
git commit -m "feat: XML 포맷터 빈 파일에 <!-- empty --> 표시"
```

---

## Task 4: 통합 테스트 및 검증

**Step 1: 전체 테스트 실행**

Run: `go test ./pkg/formatter/... -v`
Expected: All PASS

**Step 2: 빌드**

Run: `go build -o brfit ./cmd/brfit`

**Step 3: 실제 출력 확인 (Markdown)**

Run: `./brfit . -f md | grep -A3 "main.go"`
Expected:
```
### cmd/brfit/main.go

```go
// (empty)
```

**Step 4: 실제 출력 확인 (XML)**

Run: `./brfit . -f xml | grep -A2 "main.go"`
Expected:
```
    <file path="cmd/brfit/main.go" language="go">
      <!-- empty -->
    </file>
```

**Step 5: 최종 커밋 및 푸시**

```bash
git push origin main
```

---

## 검증 체크리스트

- [ ] `getEmptyComment()` 함수가 언어별 올바른 코멘트 반환
- [ ] Markdown: 빈 파일에 `// (empty)` 출력
- [ ] XML: 빈 파일에 `<!-- empty -->` 출력
- [ ] imports만 있는 파일은 빈 파일로 처리 안 함
- [ ] 에러가 있는 파일은 기존대로 에러 표시
- [ ] 모든 기존 테스트 통과
