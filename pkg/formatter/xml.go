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
	buf.WriteString(`      <tag name="imports" description="Raw import/export statements (verbatim text)" />` + "\n")
	buf.WriteString(`      <tag name="doc" description="Documentation comment" />` + "\n")
	buf.WriteString(`      <tag name="error" description="Parse error message" />` + "\n")
	buf.WriteString("    </schema>\n")

	buf.WriteString("  </metadata>\n")

	// Files section
	buf.WriteString("  <files>\n")
	for _, file := range data.Files {
		buf.WriteString(fmt.Sprintf("    <file path=%q language=%q>\n", file.Path, file.Language))

		// Imports section (within file block)
		hasRenderedImports := false
		if file.Error == nil && data.IncludeImports && len(file.RawImports) > 0 {
			hasRenderedImports = true
			buf.WriteString("      <imports>")
			buf.WriteString(escapeXML(strings.Join(file.RawImports, "\n")))
			buf.WriteString("</imports>\n")
		}

		if file.Error != nil {
			buf.WriteString("      <error>")
			buf.WriteString(escapeXML(file.Error.Error()))
			buf.WriteString("</error>\n")
		} else {
			// 빈 파일 확인
			isEmpty := len(file.Signatures) == 0 && !hasRenderedImports

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
// Optimized to scan the string only once instead of 5 sequential ReplaceAll calls.
func escapeXML(s string) string {
	// Fast path: check if escaping is needed
	var needsEscape bool
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '&' || c == '<' || c == '>' || c == '"' || c == '\'' {
			needsEscape = true
			break
		}
	}
	if !needsEscape {
		return s
	}

	// Escape needed: single-pass replacement
	var buf strings.Builder
	buf.Grow(len(s) + len(s)/10) // ~10% extra capacity for escaped chars

	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case '&':
			buf.WriteString("&amp;")
		case '<':
			buf.WriteString("&lt;")
		case '>':
			buf.WriteString("&gt;")
		case '"':
			buf.WriteString("&quot;")
		case '\'':
			buf.WriteString("&apos;")
		default:
			buf.WriteByte(c)
		}
	}

	return buf.String()
}

// kindToTag maps a signature Kind to the appropriate XML tag name.
func kindToTag(kind string) string {
	switch kind {
	case "function", "method", "constructor", "destructor", "arrow", "local_function", "module_function":
		return "function"
	case "class", "interface", "type", "struct", "enum", "record", "annotation", "typedef", "namespace", "template":
		return "type"
	case "variable", "field", "macro", "export":
		return "variable"
	default:
		return "signature" // fallback for empty or unknown kinds
	}
}
