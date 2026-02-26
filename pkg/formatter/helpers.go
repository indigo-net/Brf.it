package formatter

// getEmptyComment returns the appropriate empty file comment for a language.
func getEmptyComment(lang string) string {
	switch lang {
	case "python", "ruby":
		return "# (empty)"
	case "html", "xml":
		return "<!-- (empty) -->"
	default:
		return "// (empty)"
	}
}
