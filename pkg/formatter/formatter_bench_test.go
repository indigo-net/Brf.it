package formatter

import (
	"testing"

	"github.com/indigo-net/Brf.it/pkg/parser"
)

// createBenchmarkData creates PackageData for benchmarking.
func createBenchmarkData(numFiles, numSigsPerFile int) *PackageData {
	files := make([]FileData, numFiles)
	for i := 0; i < numFiles; i++ {
		sigs := make([]parser.Signature, numSigsPerFile)
		for j := 0; j < numSigsPerFile; j++ {
			sigs[j] = parser.Signature{
				Name:     "FunctionName",
				Kind:     "function",
				Text:     "func FunctionName(a int, b string) (result error)",
				Doc:      "FunctionName does something useful with the given parameters.",
				Line:     j + 1,
				Language: "go",
				Exported: true,
			}
		}
		files[i] = FileData{
			Path:       "pkg/module/file.go",
			Language:   "go",
			Signatures: sigs,
		}
	}

	return &PackageData{
		Version:         "v0.1.0",
		RootPath:        "/project",
		Tree:            "pkg/\n└── module/\n    └── file.go",
		Files:           files,
		TotalSignatures: numFiles * numSigsPerFile,
	}
}

// BenchmarkXMLFormatter benchmarks XML formatting with varying data sizes.
func BenchmarkXMLFormatter(b *testing.B) {
	tests := []struct {
		name          string
		numFiles      int
		numSigsPerFile int
	}{
		{"small_5files_10sigs", 5, 10},
		{"medium_20files_50sigs", 20, 50},
		{"large_50files_100sigs", 50, 100},
		{"xlarge_100files_200sigs", 100, 200},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			data := createBenchmarkData(tt.numFiles, tt.numSigsPerFile)
			f := NewXMLFormatter()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, err := f.Format(data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkMarkdownFormatter benchmarks Markdown formatting.
func BenchmarkMarkdownFormatter(b *testing.B) {
	tests := []struct {
		name          string
		numFiles      int
		numSigsPerFile int
	}{
		{"small_5files_10sigs", 5, 10},
		{"medium_20files_50sigs", 20, 50},
		{"large_50files_100sigs", 50, 100},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			data := createBenchmarkData(tt.numFiles, tt.numSigsPerFile)
			f := NewMarkdownFormatter()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, err := f.Format(data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkJSONFormatter benchmarks JSON formatting.
func BenchmarkJSONFormatter(b *testing.B) {
	tests := []struct {
		name          string
		numFiles      int
		numSigsPerFile int
	}{
		{"small_5files_10sigs", 5, 10},
		{"medium_20files_50sigs", 20, 50},
		{"large_50files_100sigs", 50, 100},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			data := createBenchmarkData(tt.numFiles, tt.numSigsPerFile)
			f := NewJSONFormatter()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, err := f.Format(data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkXMLFormatterWithImports benchmarks XML formatting with imports.
func BenchmarkXMLFormatterWithImports(b *testing.B) {
	data := createBenchmarkData(10, 20)

	// Add imports to files
	for i := range data.Files {
		data.Files[i].RawImports = []string{
			`import "fmt"`,
			`import "strings"`,
			`import "os"`,
		}
	}
	data.IncludeImports = true

	f := NewXMLFormatter()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := f.Format(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkFormatterComparison compares all formatters.
func BenchmarkFormatterComparison(b *testing.B) {
	data := createBenchmarkData(20, 50)

	formatters := map[string]Formatter{
		"xml":      NewXMLFormatter(),
		"markdown": NewMarkdownFormatter(),
		"json":     NewJSONFormatter(),
	}

	for name, f := range formatters {
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, err := f.Format(data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
