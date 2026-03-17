package formatter

import (
	"bytes"
	"strconv"
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

	// Metadata section (only output if any metadata exists)
	hasMetadata := data.Version != "" || data.RootPath != "" || data.Tree != ""
	if hasMetadata {
		buf.WriteString("  <metadata>\n")

		// Version
		if data.Version != "" {
			buf.WriteString("    <version>")
			buf.WriteString(escapeXML(data.Version))
			buf.WriteString("</version>\n")
		}

		// Path
		if data.RootPath != "" {
			buf.WriteString("    <path>")
			buf.WriteString(escapeXML(data.RootPath))
			buf.WriteString("</path>\n")
		}

		// Tree
		if data.Tree != "" {
			buf.WriteString("    <tree>")
			buf.WriteString(escapeXML(data.Tree))
			buf.WriteString("</tree>\n")
		}

		// Schema (optional - can be disabled with --no-schema)
		if !data.NoSchema {
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
			buf.WriteString(`      <tag name="call" description="Function/method call reference within the file" />` + "\n")
			buf.WriteString(`      <tag name="doc" description="Documentation comment" />` + "\n")
			buf.WriteString(`      <tag name="error" description="Parse error message" />` + "\n")
			buf.WriteString("    </schema>\n")
		}

		buf.WriteString("  </metadata>\n")
	}

	// Files section
	buf.WriteString("  <files>\n")
	for _, file := range data.Files {
		// Imports section (within file block)
		hasRenderedImports := false
		if file.Error == nil && data.IncludeImports && len(file.RawImports) > 0 {
			hasRenderedImports = true
		}

		// 빈 파일 확인
		isEmpty := file.Error == nil && len(file.Signatures) == 0 && !hasRenderedImports

		// SkipEmpty가 true이면 빈 파일 전체를 건너뜀
		if data.SkipEmpty && isEmpty {
			continue
		}

		buf.WriteString("    <file path=\"")
		buf.WriteString(escapeXML(file.Path))
		buf.WriteString("\" language=\"")
		buf.WriteString(escapeXML(file.Language))
		buf.WriteString("\">\n")

		// Render imports
		if hasRenderedImports {
			buf.WriteString("      <imports>")
			buf.WriteString(escapeXML(strings.Join(file.RawImports, "\n")))
			buf.WriteString("</imports>\n")
		}

		if file.Error != nil {
			buf.WriteString("      <error>")
			buf.WriteString(escapeXML(file.Error.Error()))
			buf.WriteString("</error>\n")
		} else {
			if isEmpty {
				buf.WriteString("      <!-- empty -->\n")
			} else {
				for _, sig := range file.Signatures {
					tag := kindToTag(sig.Kind)
					buf.WriteString("      <")
					buf.WriteString(tag)
					buf.WriteByte('>')
					buf.WriteString(escapeXML(sig.Text))
					buf.WriteString("</")
					buf.WriteString(tag)
					buf.WriteString(">\n")

					if sig.Doc != "" {
						buf.WriteString("      <doc>")
						buf.WriteString(escapeXML(truncateDoc(sig.Doc, data.MaxDocLength)))
						buf.WriteString("</doc>\n")
					}
				}

				// Call graph section
				if data.IncludeCallGraph && len(file.Calls) > 0 {
					buf.WriteString("      <calls>\n")
					for _, call := range file.Calls {
						buf.WriteString("        <call")
						if call.Caller != "" {
							buf.WriteString(" caller=\"")
							buf.WriteString(escapeXML(call.Caller))
							buf.WriteByte('"')
						}
						buf.WriteString(" callee=\"")
						buf.WriteString(escapeXML(call.Callee))
						buf.WriteString("\" line=\"")
						buf.WriteString(strconv.Itoa(call.Line))
						buf.WriteString("\" />\n")
					}
					buf.WriteString("      </calls>\n")
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
// It uses normalizeKind for the common mapping and falls back to "signature"
// for unknown or empty kinds.
func kindToTag(kind string) string {
	result := normalizeKind(kind)
	switch result {
	case "function", "type", "variable":
		return result
	default:
		return "signature" // fallback for empty or unknown kinds
	}
}
