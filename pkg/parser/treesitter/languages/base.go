package languages

// BaseQuery provides default implementations for common LanguageQuery methods.
// Embed this struct to get the default Captures() implementation.
type BaseQuery struct{}

// Captures returns the standard capture names used across all language queries.
func (BaseQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}
