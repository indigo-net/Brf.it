// Package security provides secret detection and redaction for extracted code.
package security

import (
	"fmt"
	"io"
	"regexp"

	"github.com/indigo-net/Brf.it/pkg/extractor"
	"github.com/indigo-net/Brf.it/pkg/parser"
)

// Pattern represents a secret detection pattern.
type Pattern struct {
	// Name is a human-readable name for the pattern.
	Name string

	// Regex is the compiled regular expression.
	Regex *regexp.Regexp
}

// defaultPatterns returns the built-in secret detection patterns.
func defaultPatterns() []Pattern {
	return []Pattern{
		{Name: "AWS Access Key ID", Regex: regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
		{Name: "AWS Secret Access Key", Regex: regexp.MustCompile(`(?i)aws_secret_access_key\s*[=:]\s*[A-Za-z0-9/+=]{40}`)},
		{Name: "GitHub Token", Regex: regexp.MustCompile(`gh[pousr]_[A-Za-z0-9_]{36,255}`)},
		{Name: "GitHub Personal Access Token (classic)", Regex: regexp.MustCompile(`ghp_[A-Za-z0-9]{36}`)},
		{Name: "Generic API Key", Regex: regexp.MustCompile(`(?i)(?:api[_-]?key|apikey)\s*[=:]\s*["']?[A-Za-z0-9_\-]{20,}["']?`)},
		{Name: "Generic Secret", Regex: regexp.MustCompile(`(?i)(?:secret|password|passwd|pwd)\s*[=:]\s*["']?[^\s"']{8,}["']?`)},
		{Name: "Bearer Token", Regex: regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9\-._~+/]+=*`)},
		{Name: "Private Key", Regex: regexp.MustCompile(`-----BEGIN (?:RSA |EC |DSA )?PRIVATE KEY-----`)},
		{Name: "Slack Token", Regex: regexp.MustCompile(`xox[bporas]-[0-9]{10,}-[A-Za-z0-9-]+`)},
		{Name: "Google API Key", Regex: regexp.MustCompile(`AIza[0-9A-Za-z_-]{35}`)},
		{Name: "Heroku API Key", Regex: regexp.MustCompile(`(?i)heroku.*[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)},
		{Name: "Generic Token Assignment", Regex: regexp.MustCompile(`(?i)(?:token|auth)\s*[=:]\s*["']?[A-Za-z0-9_\-]{20,}["']?`)},
	}
}

// Finding represents a single secret detection finding.
type Finding struct {
	// FilePath is the file where the secret was found.
	FilePath string

	// PatternName is the name of the matched pattern.
	PatternName string

	// MatchedText is the matched secret text.
	MatchedText string
}

// ScanResult contains the results of a security scan.
type ScanResult struct {
	// Findings is the list of detected secrets.
	Findings []Finding

	// RedactedFiles is the extract result with secrets redacted.
	RedactedFiles []extractor.ExtractedFile
}

// Scanner detects and redacts secrets in extracted code.
type Scanner struct {
	patterns []Pattern
	warnings io.Writer
}

// NewScanner creates a new Scanner with default patterns.
// Warnings are written to the given writer (typically os.Stderr).
func NewScanner(warnings io.Writer) *Scanner {
	return &Scanner{
		patterns: defaultPatterns(),
		warnings: warnings,
	}
}

// Scan scans extracted files for secrets, redacts them, and returns findings.
func (s *Scanner) Scan(result *extractor.ExtractResult) *ScanResult {
	if result == nil {
		return &ScanResult{}
	}

	sr := &ScanResult{
		RedactedFiles: make([]extractor.ExtractedFile, len(result.Files)),
	}

	for i, file := range result.Files {
		sr.RedactedFiles[i] = s.scanFile(file, sr)
	}

	// Print warnings for findings
	if len(sr.Findings) > 0 {
		fmt.Fprintf(s.warnings, "[brfit] WARN: %d potential secret(s) detected and redacted:\n", len(sr.Findings))
		for _, f := range sr.Findings {
			fmt.Fprintf(s.warnings, "  - %s: %s\n", f.FilePath, f.PatternName)
		}
	}

	return sr
}

// scanFile scans a single extracted file and returns a redacted copy.
func (s *Scanner) scanFile(file extractor.ExtractedFile, sr *ScanResult) extractor.ExtractedFile {
	if file.Error != nil {
		return file
	}

	redacted := extractor.ExtractedFile{
		Path:     file.Path,
		Language: file.Language,
		Size:     file.Size,
		Error:    file.Error,
	}

	// Scan and redact signatures
	redacted.Signatures = make([]parser.Signature, len(file.Signatures))
	copy(redacted.Signatures, file.Signatures)
	for j := range redacted.Signatures {
		sig := &redacted.Signatures[j]
		sig.Text = s.redactString(file.Path, sig.Text, sr)
		sig.Doc = s.redactString(file.Path, sig.Doc, sr)
	}

	// Scan and redact raw imports
	redacted.RawImports = make([]string, len(file.RawImports))
	for j, imp := range file.RawImports {
		redacted.RawImports[j] = s.redactString(file.Path, imp, sr)
	}

	return redacted
}

// redactString replaces any matched secrets in the string with [REDACTED].
func (s *Scanner) redactString(filePath, text string, sr *ScanResult) string {
	if text == "" {
		return text
	}

	result := text
	for _, p := range s.patterns {
		matches := p.Regex.FindAllString(result, -1)
		if len(matches) > 0 {
			for _, m := range matches {
				sr.Findings = append(sr.Findings, Finding{
					FilePath:    filePath,
					PatternName: p.Name,
					MatchedText: m,
				})
			}
			result = p.Regex.ReplaceAllString(result, "[REDACTED]")
		}
	}

	return result
}
