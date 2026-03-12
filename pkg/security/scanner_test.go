package security

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/indigo-net/Brf.it/pkg/extractor"
	"github.com/indigo-net/Brf.it/pkg/parser"
)

func TestScan_NilResult(t *testing.T) {
	s := NewScanner(&bytes.Buffer{})
	sr := s.Scan(nil)
	if len(sr.Findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(sr.Findings))
	}
}

func TestScan_NoSecrets(t *testing.T) {
	var buf bytes.Buffer
	s := NewScanner(&buf)

	result := &extractor.ExtractResult{
		Files: []extractor.ExtractedFile{
			{
				Path:     "main.go",
				Language: "go",
				Signatures: []parser.Signature{
					{Text: "func main()", Doc: "Entry point."},
				},
				RawImports: []string{`import "fmt"`},
			},
		},
	}

	sr := s.Scan(result)
	if len(sr.Findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(sr.Findings))
	}
	if buf.Len() != 0 {
		t.Errorf("expected no warnings, got: %s", buf.String())
	}
	// Verify content is unchanged
	if sr.RedactedFiles[0].Signatures[0].Text != "func main()" {
		t.Errorf("text should be unchanged")
	}
}

func TestScan_AWSAccessKeyDetected(t *testing.T) {
	var buf bytes.Buffer
	s := NewScanner(&buf)

	result := &extractor.ExtractResult{
		Files: []extractor.ExtractedFile{
			{
				Path:     "config.go",
				Language: "go",
				Signatures: []parser.Signature{
					{Text: `var key = "AKIAIOSFODNN7EXAMPLE"`},
				},
			},
		},
	}

	sr := s.Scan(result)
	if len(sr.Findings) == 0 {
		t.Fatal("expected findings for AWS key")
	}
	if !strings.Contains(sr.RedactedFiles[0].Signatures[0].Text, "[REDACTED]") {
		t.Error("expected [REDACTED] in output")
	}
	if !strings.Contains(buf.String(), "potential secret(s) detected") {
		t.Error("expected warning on stderr")
	}
}

func TestScan_GitHubTokenDetected(t *testing.T) {
	var buf bytes.Buffer
	s := NewScanner(&buf)

	result := &extractor.ExtractResult{
		Files: []extractor.ExtractedFile{
			{
				Path:     "auth.go",
				Language: "go",
				Signatures: []parser.Signature{
					{Text: `var token = "ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh"`},
				},
			},
		},
	}

	sr := s.Scan(result)
	if len(sr.Findings) == 0 {
		t.Fatal("expected findings for GitHub token")
	}
	if !strings.Contains(sr.RedactedFiles[0].Signatures[0].Text, "[REDACTED]") {
		t.Error("expected [REDACTED] in output")
	}
}

func TestScan_GenericAPIKeyDetected(t *testing.T) {
	var buf bytes.Buffer
	s := NewScanner(&buf)

	result := &extractor.ExtractResult{
		Files: []extractor.ExtractedFile{
			{
				Path:     "service.py",
				Language: "python",
				Signatures: []parser.Signature{
					{Text: `API_KEY = "sk_live_abcdef1234567890abcd"`},
				},
			},
		},
	}

	sr := s.Scan(result)
	if len(sr.Findings) == 0 {
		t.Fatal("expected findings for API key")
	}
}

func TestScan_PasswordInDocDetected(t *testing.T) {
	var buf bytes.Buffer
	s := NewScanner(&buf)

	result := &extractor.ExtractResult{
		Files: []extractor.ExtractedFile{
			{
				Path:     "db.go",
				Language: "go",
				Signatures: []parser.Signature{
					{
						Text: "func Connect()",
						Doc:  `password = "supersecretpassword123"`,
					},
				},
			},
		},
	}

	sr := s.Scan(result)
	if len(sr.Findings) == 0 {
		t.Fatal("expected findings for password in doc")
	}
	if !strings.Contains(sr.RedactedFiles[0].Signatures[0].Doc, "[REDACTED]") {
		t.Error("expected [REDACTED] in doc")
	}
	// Text should be unchanged
	if sr.RedactedFiles[0].Signatures[0].Text != "func Connect()" {
		t.Error("text should be unchanged when no secret")
	}
}

func TestScan_PrivateKeyDetected(t *testing.T) {
	var buf bytes.Buffer
	s := NewScanner(&buf)

	result := &extractor.ExtractResult{
		Files: []extractor.ExtractedFile{
			{
				Path:     "cert.go",
				Language: "go",
				Signatures: []parser.Signature{
					{Text: `var cert = "-----BEGIN RSA PRIVATE KEY-----\nMIIE..."`},
				},
			},
		},
	}

	sr := s.Scan(result)
	if len(sr.Findings) == 0 {
		t.Fatal("expected findings for private key")
	}
}

func TestScan_ImportRedacted(t *testing.T) {
	var buf bytes.Buffer
	s := NewScanner(&buf)

	result := &extractor.ExtractResult{
		Files: []extractor.ExtractedFile{
			{
				Path:       "main.go",
				Language:   "go",
				RawImports: []string{`api_key = "my_secret_api_key_value_1234567890"`},
			},
		},
	}

	sr := s.Scan(result)
	if len(sr.Findings) == 0 {
		t.Fatal("expected findings in imports")
	}
	if !strings.Contains(sr.RedactedFiles[0].RawImports[0], "[REDACTED]") {
		t.Error("expected [REDACTED] in import")
	}
}

func TestScan_ErrorFileSkipped(t *testing.T) {
	var buf bytes.Buffer
	s := NewScanner(&buf)

	result := &extractor.ExtractResult{
		Files: []extractor.ExtractedFile{
			{
				Path:  "bad.go",
				Error: fmt.Errorf("read error"),
			},
		},
	}

	sr := s.Scan(result)
	if len(sr.Findings) != 0 {
		t.Errorf("expected 0 findings for error file, got %d", len(sr.Findings))
	}
}

func TestScan_BearerTokenDetected(t *testing.T) {
	var buf bytes.Buffer
	s := NewScanner(&buf)

	result := &extractor.ExtractResult{
		Files: []extractor.ExtractedFile{
			{
				Path:     "http.go",
				Language: "go",
				Signatures: []parser.Signature{
					{Text: `Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U`},
				},
			},
		},
	}

	sr := s.Scan(result)
	if len(sr.Findings) == 0 {
		t.Fatal("expected findings for bearer token")
	}
}

func TestScan_MultipleSecretsInOneFile(t *testing.T) {
	var buf bytes.Buffer
	s := NewScanner(&buf)

	result := &extractor.ExtractResult{
		Files: []extractor.ExtractedFile{
			{
				Path:     "config.go",
				Language: "go",
				Signatures: []parser.Signature{
					{Text: `var awsKey = "AKIAIOSFODNN7EXAMPLE"`},
					{Text: `var token = "ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh"`},
				},
			},
		},
	}

	sr := s.Scan(result)
	if len(sr.Findings) < 2 {
		t.Errorf("expected at least 2 findings, got %d", len(sr.Findings))
	}
}
