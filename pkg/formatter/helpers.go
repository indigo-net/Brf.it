package formatter

import "unicode/utf8"

// normalizeKind normalizes a signature kind string to one of the canonical
// categories: "function", "type", or "variable". If the kind does not match
// any known category, it is returned unchanged.
func normalizeKind(kind string) string {
	switch kind {
	case "function", "method", "constructor", "destructor", "arrow", "local_function", "module_function":
		return "function"
	case "class", "interface", "type", "struct", "enum", "record", "annotation", "typedef", "namespace", "template", "trait", "impl":
		return "type"
	case "variable", "field", "macro", "export":
		return "variable"
	default:
		return kind
	}
}

// getEmptyComment returns the appropriate empty file comment for a language.
func getEmptyComment(lang string) string {
	switch lang {
	case "python", "ruby":
		return "# (empty)"
	case "html", "xml":
		return "<!-- (empty) -->"
	case "go", "c", "cpp", "java", "javascript", "typescript":
		return "// (empty)"
	default:
		return "// (empty)"
	}
}

// truncateDoc truncates a documentation string to maxLen characters (Unicode code points).
// If maxLen <= 0 or the doc is shorter than or equal to maxLen, returns doc unchanged.
// If truncation occurs, "..." is appended.
func truncateDoc(doc string, maxLen int) string {
	if maxLen <= 0 {
		return doc
	}

	if utf8.RuneCountInString(doc) <= maxLen {
		return doc
	}

	runes := []rune(doc)
	return string(runes[:maxLen]) + "..."
}
