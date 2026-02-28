package formatter

import (
	"bytes"
	"fmt"
	"strings"
)

// XMLFormatter implements Formatter for XML output.
type XMLFormatter struct{}

// NewXMLFormatter creates a new XMLFormatter.
func NewXMLFormatter() *XMLFormatter {
	return &XMLFormatter{}
}

// Name returns the formatter name.
func (f *XMLFormatter) Name() string {
	return "xml"
}

// xmlSchemaComment provides schema documentation for the XML output.
const xmlSchemaComment = `<!--
Schema:
| Tag       | Description                              |
|-----------|------------------------------------------|
| file      | Source file (path, language attributes)  |
| signature | Function, type, or variable declaration  |
| imports   | Import statements container              |
| import    | Single import statement                  |
| export    | Single export statement                  |
| doc       | Documentation comment                    |
| error     | Parse error message                      |
-->
`

// Format implements Formatter interface.
func (f *XMLFormatter) Format(data *PackageData) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	buf.WriteByte('\n')

	// Header comment with version and path
	if data.Version != "" && data.RootPath != "" {
		buf.WriteString(fmt.Sprintf("<!-- brf.it %s | Code Summary: %s -->\n",
			data.Version, data.RootPath))
	} else if data.Version != "" {
		buf.WriteString(fmt.Sprintf("<!-- brf.it %s -->\n", data.Version))
	} else if data.RootPath != "" {
		buf.WriteString(fmt.Sprintf("<!-- Code Summary: %s -->\n", data.RootPath))
	}

	// Schema documentation
	buf.WriteString(xmlSchemaComment)

	buf.WriteString("<brfit>\n")

	// Metadata section
	buf.WriteString("  <metadata>\n")

	// Tree
	if data.Tree != "" {
		buf.WriteString("    <tree>")
		buf.WriteString(escapeXML(data.Tree))
		buf.WriteString("</tree>\n")
	}

	buf.WriteString("  </metadata>\n")

	// Files section
	buf.WriteString("  <files>\n")
	for _, file := range data.Files {
		buf.WriteString(fmt.Sprintf("    <file path=%q language=%q>\n", file.Path, file.Language))

		// Imports section (within file block)
		if file.Error == nil && data.IncludeImports && len(file.Imports) > 0 {
			buf.WriteString("      <imports>\n")
			for _, imp := range file.Imports {
				if imp.Type == "import" {
					buf.WriteString("        <import>")
					buf.WriteString(escapeXML(imp.Path))
					buf.WriteString("</import>\n")
				} else if imp.Type == "export" {
					if imp.Name != "" {
						buf.WriteString("        <export>")
						buf.WriteString(escapeXML(imp.Name))
						buf.WriteString("</export>\n")
					} else if imp.Path != "" {
						buf.WriteString("        <export>")
						buf.WriteString(escapeXML(imp.Path))
						buf.WriteString("</export>\n")
					}
				}
			}
			buf.WriteString("      </imports>\n")
		}

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

		buf.WriteString("    </file>\n")
	}
	buf.WriteString("  </files>\n")

	buf.WriteString("</brfit>\n")

	return buf.Bytes(), nil
}

// escapeXML escapes special characters for XML content.
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

