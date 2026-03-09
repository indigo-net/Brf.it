package treesitter

import (
	"strings"
	"testing"

	"github.com/indigo-net/Brf.it/pkg/parser"
)

// BenchmarkParseGo benchmarks parsing Go code of varying sizes.
func BenchmarkParseGo(b *testing.B) {
	tests := []struct {
		name     string
		funcs    int
		docLines int
	}{
		{"small_10funcs_0doc", 10, 0},
		{"medium_50funcs_5doc", 50, 5},
		{"large_100funcs_10doc", 100, 10},
		{"xlarge_500funcs_20doc", 500, 20},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			// Generate Go code
			var code strings.Builder
			code.WriteString("package main\n\n")

			for i := 0; i < tt.funcs; i++ {
				// Add doc comments
				for j := 0; j < tt.docLines; j++ {
					code.WriteString("// This is documentation line for function.\n")
				}
				code.WriteString("func Func")
				code.WriteString(string(rune('A' + i%26)))
				code.WriteString("(a, b int) string {\n")
				code.WriteString("\treturn \"result\"\n")
				code.WriteString("}\n\n")
			}

			data := []byte(code.String())
			p := NewTreeSitterParser()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, err := p.Parse(data, nil)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkParseTypeScript benchmarks parsing TypeScript code.
func BenchmarkParseTypeScript(b *testing.B) {
	tests := []struct {
		name   string
		funcs  int
	}{
		{"small_10funcs", 10},
		{"medium_50funcs", 50},
		{"large_100funcs", 100},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			var code strings.Builder

			for i := 0; i < tt.funcs; i++ {
				code.WriteString("function func")
				code.WriteString(string(rune('A' + i%26)))
				code.WriteString("(a: number, b: string): void {\n")
				code.WriteString("\tconsole.log(a, b);\n")
				code.WriteString("}\n\n")
			}

			data := []byte(code.String())
			p := NewTreeSitterParser()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, err := p.Parse(data, &parser.Options{Language: "typescript"})
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkParsePython benchmarks parsing Python code.
func BenchmarkParsePython(b *testing.B) {
	tests := []struct {
		name   string
		funcs  int
	}{
		{"small_10funcs", 10},
		{"medium_50funcs", 50},
		{"large_100funcs", 100},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			var code strings.Builder

			for i := 0; i < tt.funcs; i++ {
				code.WriteString("def func")
				code.WriteString(string(rune('a' + i%26)))
				code.WriteString("(self, a, b):\n")
				code.WriteString("    \"\"\"Documentation.\"\"\"\n")
				code.WriteString("    return a + b\n\n")
			}

			data := []byte(code.String())
			p := NewTreeSitterParser()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, err := p.Parse(data, &parser.Options{Language: "python"})
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkParseWithImports benchmarks parsing with import extraction.
func BenchmarkParseWithImports(b *testing.B) {
	var code strings.Builder
	code.WriteString("package main\n\n")

	// Add imports
	for i := 0; i < 20; i++ {
		code.WriteString("import \"fmt\"\n")
	}

	// Add functions
	for i := 0; i < 50; i++ {
		code.WriteString("func F() { fmt.Println() }\n")
	}

	data := []byte(code.String())
	p := NewTreeSitterParser()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := p.Parse(data, &parser.Options{IncludeImports: true})
		if err != nil {
			b.Fatal(err)
		}
	}
}
