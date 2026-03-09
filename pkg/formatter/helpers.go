package formatter

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
