package formatter

import (
	"strings"
	"testing"
)

func TestXMLFormatterWithNoSchema(t *testing.T) {
	formatter := NewXMLFormatter()
	data := &PackageData{
		NoSchema: true,
		Files: []FileData{
			{
				Path:     "test.go",
				Language: "go",
			},
		},
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)
	if strings.Contains(outputStr, "<schema>") {
		t.Error("expected no <schema> section with --no-schema flag")
	}
}
