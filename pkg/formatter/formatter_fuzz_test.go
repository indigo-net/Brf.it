package formatter

import (
	"testing"

	"github.com/indigo-net/Brf.it/pkg/parser"
)

// FuzzXMLFormatter tests that the XML formatter does not panic on arbitrary data.
func FuzzXMLFormatter(f *testing.F) {
	// Seed corpus with various PackageData scenarios
	f.Add([]byte("test"), []byte("func main()"), []byte("doc"))
	f.Add([]byte(""), []byte(""), []byte(""))
	f.Add([]byte("path/with/unicode/한글/路径"), []byte(""), []byte(""))
	f.Add([]byte("path"), []byte("func <script>alert('xss')</script>()"), []byte(""))
	f.Add([]byte("path"), []byte("func test()"), []byte("doc with \"quotes\" and 'apostrophes'"))
	f.Add(make([]byte, 10000), make([]byte, 10000), make([]byte, 10000))

	formatter := NewXMLFormatter()

	f.Fuzz(func(t *testing.T, path, sigText, doc []byte) {
		data := &PackageData{
			Version:  "v0.1.0",
			RootPath: "/test",
			Tree:     "test/\n└── test.go",
			Files: []FileData{
				{
					Path:     string(path),
					Language: "go",
					Signatures: []parser.Signature{
						{
							Name:     "Test",
							Kind:     "function",
							Text:     string(sigText),
							Doc:      string(doc),
							Line:     1,
							Language: "go",
							Exported: true,
						},
					},
				},
			},
			TotalSignatures: 1,
		}

		// The formatter should not panic on any input
		output, err := formatter.Format(data)
		if err != nil {
			// Error is acceptable, panic is not
			return
		}
		// Output should be valid XML bytes
		if output == nil {
			t.Error("expected non-nil output on successful format")
		}
	})
}

// FuzzMarkdownFormatter tests that the Markdown formatter does not panic on arbitrary data.
func FuzzMarkdownFormatter(f *testing.F) {
	f.Add([]byte("test"), []byte("func main()"), []byte("doc"))
	f.Add([]byte(""), []byte(""), []byte(""))
	f.Add([]byte("path/with/special#chars"), []byte(""), []byte(""))
	f.Add([]byte("path"), []byte("func **bold**()"), []byte(""))
	f.Add([]byte("path"), []byte("func test()"), []byte("doc with [link](url)"))
	f.Add(make([]byte, 10000), make([]byte, 10000), make([]byte, 10000))

	formatter := NewMarkdownFormatter()

	f.Fuzz(func(t *testing.T, path, sigText, doc []byte) {
		data := &PackageData{
			Version:  "v0.1.0",
			RootPath: "/test",
			Tree:     "test/\n└── test.go",
			Files: []FileData{
				{
					Path:     string(path),
					Language: "go",
					Signatures: []parser.Signature{
						{
							Name:     "Test",
							Kind:     "function",
							Text:     string(sigText),
							Doc:      string(doc),
							Line:     1,
							Language: "go",
							Exported: true,
						},
					},
				},
			},
			TotalSignatures: 1,
		}

		output, err := formatter.Format(data)
		if err != nil {
			return
		}
		if output == nil {
			t.Error("expected non-nil output on successful format")
		}
	})
}

// FuzzJSONFormatter tests that the JSON formatter does not panic on arbitrary data.
func FuzzJSONFormatter(f *testing.F) {
	f.Add([]byte("test"), []byte("func main()"), []byte("doc"))
	f.Add([]byte(""), []byte(""), []byte(""))
	f.Add([]byte("path/with/unicode/日本語"), []byte(""), []byte(""))
	f.Add([]byte("path"), []byte("func \"quoted\"()"), []byte(""))
	f.Add([]byte("path"), []byte("func test()"), []byte("doc with\nnewlines\tand\ttabs"))
	f.Add(make([]byte, 10000), make([]byte, 10000), make([]byte, 10000))

	formatter := NewJSONFormatter()

	f.Fuzz(func(t *testing.T, path, sigText, doc []byte) {
		data := &PackageData{
			Version:  "v0.1.0",
			RootPath: "/test",
			Tree:     "test/\n└── test.go",
			Files: []FileData{
				{
					Path:     string(path),
					Language: "go",
					Signatures: []parser.Signature{
						{
							Name:     "Test",
							Kind:     "function",
							Text:     string(sigText),
							Doc:      string(doc),
							Line:     1,
							Language: "go",
							Exported: true,
						},
					},
				},
			},
			TotalSignatures: 1,
		}

		output, err := formatter.Format(data)
		if err != nil {
			return
		}
		if output == nil {
			t.Error("expected non-nil output on successful format")
		}
	})
}

// FuzzFormatterWithLargeData tests formatters with large amounts of data.
func FuzzFormatterWithLargeData(f *testing.F) {
	f.Add(100, 100) // numFiles, numSigsPerFile

	formatters := []Formatter{
		NewXMLFormatter(),
		NewMarkdownFormatter(),
		NewJSONFormatter(),
	}

	f.Fuzz(func(t *testing.T, numFiles, numSigsPerFile int) {
		// Limit to reasonable sizes to avoid timeout
		if numFiles > 1000 {
			numFiles = 1000
		}
		if numSigsPerFile > 100 {
			numSigsPerFile = 100
		}
		if numFiles < 0 {
			numFiles = 0
		}
		if numSigsPerFile < 0 {
			numSigsPerFile = 0
		}

		files := make([]FileData, numFiles)
		for i := 0; i < numFiles; i++ {
			sigs := make([]parser.Signature, numSigsPerFile)
			for j := 0; j < numSigsPerFile; j++ {
				sigs[j] = parser.Signature{
					Name:     "Func",
					Kind:     "function",
					Text:     "func Func()",
					Line:     j + 1,
					Language: "go",
					Exported: true,
				}
			}
			files[i] = FileData{
				Path:       "file.go",
				Language:   "go",
				Signatures: sigs,
			}
		}

		data := &PackageData{
			Version:         "v0.1.0",
			RootPath:        "/test",
			Files:           files,
			TotalSignatures: numFiles * numSigsPerFile,
		}

		for _, formatter := range formatters {
			output, err := formatter.Format(data)
			if err != nil {
				continue
			}
			if output == nil && numFiles > 0 {
				t.Errorf("formatter %s: expected non-nil output", formatter.Name())
			}
		}
	})
}

// FuzzFormatterWithImports tests formatters with import data.
func FuzzFormatterWithImports(f *testing.F) {
	f.Add([]byte("import \"fmt\""), []byte("import \"os\""), true)

	formatters := []Formatter{
		NewXMLFormatter(),
		NewMarkdownFormatter(),
		NewJSONFormatter(),
	}

	f.Fuzz(func(t *testing.T, imp1, imp2 []byte, includeImports bool) {
		data := &PackageData{
			Version: "v0.1.0",
			Files: []FileData{
				{
					Path:       "a.go",
					Language:   "go",
					RawImports: []string{string(imp1)},
				},
				{
					Path:       "b.go",
					Language:   "go",
					RawImports: []string{string(imp2)},
				},
			},
			IncludeImports: includeImports,
		}

		for _, formatter := range formatters {
			output, err := formatter.Format(data)
			if err != nil {
				continue
			}
			if output == nil {
				t.Errorf("formatter %s: expected non-nil output", formatter.Name())
			}
		}
	})
}
