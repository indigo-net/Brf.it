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
	Version string     `json:"version,omitempty"`
	Path    string     `json:"path,omitempty"`
	Tree    string     `json:"tree,omitempty"`
	Files   []jsonFile `json:"files"`
}

// jsonFile represents a single file in the JSON output.
type jsonFile struct {
	Path       string        `json:"path"`
	Language   string        `json:"language"`
	Signatures []jsonSig     `json:"signatures,omitempty"`
	Imports    []string      `json:"imports,omitempty"`
	Error      string        `json:"error,omitempty"`
}

// jsonSig represents a signature in the JSON output.
type jsonSig struct {
	Kind string `json:"kind"`
	Text string `json:"text"`
	Doc  string `json:"doc,omitempty"`
}

// Format implements Formatter interface.
func (f *JSONFormatter) Format(data *PackageData) ([]byte, error) {
	output := jsonOutput{
		Version: data.Version,
		Path:    data.RootPath,
		Tree:    data.Tree,
		Files:   make([]jsonFile, 0, len(data.Files)),
	}

	for _, file := range data.Files {
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
						Kind: normalizeKind(sig.Kind),
						Text: sig.Text,
					}
					if sig.Doc != "" {
						js.Doc = truncateDoc(sig.Doc, data.MaxDocLength)
					}
					jf.Signatures = append(jf.Signatures, js)
				}
			}

			// Add imports if requested
			if data.IncludeImports && len(file.RawImports) > 0 {
				jf.Imports = file.RawImports
			}
		}

		output.Files = append(output.Files, jf)
	}

	return json.Marshal(output)
}

