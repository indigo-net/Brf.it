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

// Format implements Formatter interface.
func (f *XMLFormatter) Format(data *PackageData) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	buf.WriteByte('\n')
	buf.WriteString("<brfit>\n")

	// Metadata section
	buf.WriteString("  <metadata>\n")

	// Version
	if data.Version != "" {
		buf.WriteString(fmt.Sprintf("    <version>%s</version>\n", escapeXML(data.Version)))
	}

	// Path
	if data.RootPath != "" {
		buf.WriteString(fmt.Sprintf("    <path>%s</path>\n", escapeXML(data.RootPath)))
	}

	// Tree
	if data.Tree != "" {
		buf.WriteString("    <tree>")
		buf.WriteString(escapeXML(data.Tree))
		buf.WriteString("</tree>\n")
	}

	// Schema
	buf.WriteString("    <schema>\n")
	buf.WriteString(`      <tag name="metadata" description="Project metadata container" />` + "\n")
	buf.WriteString(`      <tag name="version" description="brf.it version" />` + "\n")
	buf.WriteString(`      <tag name="path" description="Root path of the scanned project" />` + "\n")
	buf.WriteString(`      <tag name="tree" description="Directory tree structure" />` + "\n")
	buf.WriteString(`      <tag name="files" description="Source files container" />` + "\n")
	buf.WriteString(`      <tag name="file" description="Source file (path, language attributes)" />` + "\n")
	buf.WriteString(`      <tag name="function" description="Function, method, or constructor declaration" />` + "\n")
	buf.WriteString(`      <tag name="type" description="Type, class, interface, struct, or enum declaration" />` + "\n")
	buf.WriteString(`      <tag name="variable" description="Variable, constant, or field declaration" />` + "\n")
	buf.WriteString(`      <tag name="signature" description="Fallback for unknown declaration kinds" />` + "\n")
	buf.WriteString(`      <tag name="imports" description="Import statements container" />` + "\n")
	buf.WriteString(`      <tag name="import" description="Single import statement" />` + "\n")
	buf.WriteString(`      <tag name="export" description="Single export statement" />` + "\n")
	buf.WriteString(`      <tag name="doc" description="Documentation comment" />` + "\n")
	buf.WriteString(`      <tag name="error" description="Parse error message" />` + "\n")
	buf.WriteString("    </schema>\n")

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
					tag := kindToTag(sig.Kind)
					buf.WriteString(fmt.Sprintf("      <%s>", tag))
					buf.WriteString(escapeXML(sig.Text))
					buf.WriteString(fmt.Sprintf("</%s>\n", tag))

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

// kindToTag maps a signature Kind to the appropriate XML tag name.
func kindToTag(kind string) string {
	switch kind {
	case "function", "method", "constructor", "destructor", "arrow":
		return "function"
	case "class", "interface", "type", "struct", "enum", "record", "annotation", "typedef", "namespace", "template":
		return "type"
	case "variable", "field", "macro", "export":
		return "variable"
	default:
		return "signature" // fallback for empty or unknown kinds
	}
}
