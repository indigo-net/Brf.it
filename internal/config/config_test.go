package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Mode != "sig" {
		t.Errorf("expected default mode 'sig', got '%s'", cfg.Mode)
	}

	if cfg.Format != "xml" {
		t.Errorf("expected default format 'xml', got '%s'", cfg.Format)
	}

	if cfg.Output != "" {
		t.Errorf("expected default output empty, got '%s'", cfg.Output)
	}

	if cfg.IgnoreFile != ".gitignore" {
		t.Errorf("expected default ignore file '.gitignore', got '%s'", cfg.IgnoreFile)
	}

	if cfg.IncludeHidden {
		t.Error("expected IncludeHidden to be false by default")
	}

	if cfg.NoTree {
		t.Error("expected NoTree to be false by default")
	}

	if cfg.NoTokens {
		t.Error("expected NoTokens to be false by default")
	}

	const expectedMaxSize = 512000 // 500KB
	if cfg.MaxFileSize != expectedMaxSize {
		t.Errorf("expected max file size %d, got %d", expectedMaxSize, cfg.MaxFileSize)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid default config",
			config:    *DefaultConfig(),
			wantError: false,
		},
		{
			name: "valid sig mode with xml format",
			config: Config{
				Mode:        "sig",
				Format:      "xml",
				MaxFileSize: 512000,
			},
			wantError: false,
		},
		{
			name: "valid sig mode with md format",
			config: Config{
				Mode:        "sig",
				Format:      "md",
				MaxFileSize: 512000,
			},
			wantError: false,
		},
		{
			name: "invalid mode",
			config: Config{
				Mode:   "invalid",
				Format: "xml",
			},
			wantError: true,
			errorMsg:  "invalid mode",
		},
		{
			name: "invalid format",
			config: Config{
				Mode:   "sig",
				Format: "json",
			},
			wantError: true,
			errorMsg:  "invalid format",
		},
		{
			name: "negative max file size",
			config: Config{
				Mode:        "sig",
				Format:      "xml",
				MaxFileSize: -1,
			},
			wantError: true,
			errorMsg:  "max file size must be positive",
		},
		{
			name: "zero max file size",
			config: Config{
				Mode:        "sig",
				Format:      "xml",
				MaxFileSize: 0,
			},
			wantError: true,
			errorMsg:  "max file size must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error containing '%s', got nil", tt.errorMsg)
					return
				}
				if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got '%s'", err.Error())
				}
			}
		})
	}
}

func TestConfigSupportedLanguages(t *testing.T) {
	cfg := DefaultConfig()
	langs := cfg.SupportedExtensions()

	expected := map[string]string{
		".go":  "go",
		".ts":  "typescript",
		".tsx": "typescript",
		".js":  "javascript",
		".jsx": "javascript",
	}

	for ext, lang := range expected {
		if gotLang, ok := langs[ext]; !ok {
			t.Errorf("expected extension '%s' to be supported", ext)
		} else if gotLang != lang {
			t.Errorf("expected extension '%s' to map to '%s', got '%s'", ext, lang, gotLang)
		}
	}
}

// containsString checks if s contains substr
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
