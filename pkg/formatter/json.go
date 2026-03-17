package formatter

import (
	"encoding/json"
)

// JSONFormatter implements Formatter for JSON output.
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSONFormatter.
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Name returns the formatter name.
func (f *JSONFormatter) Name() string {
	return "json"
}

// jsonOutput represents the top-level JSON output structure.
type jsonOutput struct {
	Version       string            `json:"version,omitempty"`
	Path          string            `json:"path,omitempty"`
	Tree          string            `json:"tree,omitempty"`
	GlobalImports []jsonImportCount `json:"globalImports,omitempty"`
	Files         []jsonFile        `json:"files"`
}

// jsonImportCount represents a global import with usage count.
type jsonImportCount struct {
	Import string `json:"import"`
	Count  int    `json:"count"`
}

// jsonFile represents a single file in the JSON output.
type jsonFile struct {
	Path       string     `json:"path"`
	Language   string     `json:"language"`
	Signatures []jsonSig  `json:"signatures,omitempty"`
	Imports    []string   `json:"imports,omitempty"`
	Calls      []jsonCall `json:"calls,omitempty"`
	Error      string     `json:"error,omitempty"`
}

// jsonCall represents a function call reference in the JSON output.
type jsonCall struct {
	Caller string `json:"caller,omitempty"`
	Callee string `json:"callee"`
	Line   int    `json:"line"`
}

// jsonSig represents a signature in the JSON output.
type jsonSig struct {
	Kind     string `json:"kind"`
	Text     string `json:"text"`
	Doc      string `json:"doc,omitempty"`
	Line     int    `json:"line,omitempty"`
	Exported bool   `json:"exported,omitempty"`
}

// Format implements Formatter interface.
func (f *JSONFormatter) Format(data *PackageData) ([]byte, error) {
	output := jsonOutput{
		Version: data.Version,
		Path:    data.RootPath,
		Tree:    data.Tree,
		Files:   make([]jsonFile, 0, len(data.Files)),
	}

	// Add global imports when dedupe mode is enabled
	if data.DedupeImports && len(data.GlobalImports) > 0 {
		output.GlobalImports = make([]jsonImportCount, 0, len(data.GlobalImports))
		for _, ic := range data.GlobalImports {
			output.GlobalImports = append(output.GlobalImports, jsonImportCount{
				Import: ic.Import,
				Count:  ic.Count,
			})
		}
	}

	for _, file := range data.Files {
		// SkipEmpty: 빈 파일 건너뜀
		if data.SkipEmpty && file.Error == nil {
			hasImports := data.IncludeImports && len(file.RawImports) > 0 && !data.DedupeImports
			if len(file.Signatures) == 0 && !hasImports {
				continue
			}
		}

		jf := jsonFile{
			Path:     file.Path,
			Language: file.Language,
		}

		if file.Error != nil {
			jf.Error = file.Error.Error()
		} else {
			// Add signatures
			if len(file.Signatures) > 0 {
				jf.Signatures = make([]jsonSig, 0, len(file.Signatures))
				for _, sig := range file.Signatures {
					js := jsonSig{
						Kind:     normalizeKind(sig.Kind),
						Text:     sig.Text,
						Line:     sig.Line,
						Exported: sig.Exported,
					}
					if sig.Doc != "" {
						js.Doc = truncateDoc(sig.Doc, data.MaxDocLength)
					}
					jf.Signatures = append(jf.Signatures, js)
				}
			}

			// Add imports if requested (skip if deduping)
			if data.IncludeImports && len(file.RawImports) > 0 && !data.DedupeImports {
				jf.Imports = file.RawImports
			}

			// Add calls if requested
			if data.IncludeCallGraph && len(file.Calls) > 0 {
				jf.Calls = make([]jsonCall, 0, len(file.Calls))
				for _, call := range file.Calls {
					jf.Calls = append(jf.Calls, jsonCall{
						Caller: call.Caller,
						Callee: call.Callee,
						Line:   call.Line,
					})
				}
			}
		}

		output.Files = append(output.Files, jf)
	}

	return json.Marshal(output)
}
