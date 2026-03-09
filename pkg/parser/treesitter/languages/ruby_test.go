package languages

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_ruby "github.com/tree-sitter/tree-sitter-ruby/bindings/go"
)

// extractRubyNames is a test helper that parses Ruby code and returns
// all captured @name values from the query matches.
func extractRubyNames(t *testing.T, code []byte) map[string]bool {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(tree_sitter_ruby.Language())
	parser.SetLanguage(lang)
	tree := parser.Parse(code, nil)
	defer tree.Close()
	query := NewRubyQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()
	qc := sitter.NewQueryCursor()
	defer qc.Close()
	matches := qc.Matches(q, tree.RootNode(), code)
	captureNames := q.CaptureNames()
	foundNames := make(map[string]bool)
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		for _, c := range match.Captures {
			if captureNames[c.Index] == "name" {
				foundNames[string(code[c.Node.StartByte():c.Node.EndByte()])] = true
			}
		}
	}
	return foundNames
}

// extractRubyImports is a test helper that parses Ruby code and returns
// all captured @import_path values from the import query matches.
func extractRubyImports(t *testing.T, code []byte) []string {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(tree_sitter_ruby.Language())
	parser.SetLanguage(lang)
	tree := parser.Parse(code, nil)
	defer tree.Close()
	query := NewRubyQuery()
	q, err := sitter.NewQuery(lang, string(query.ImportQuery()))
	if err != nil {
		t.Fatalf("failed to create import query: %v", err)
	}
	defer q.Close()
	qc := sitter.NewQueryCursor()
	defer qc.Close()
	matches := qc.Matches(q, tree.RootNode(), code)
	captureNames := q.CaptureNames()
	var imports []string
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		for _, c := range match.Captures {
			if captureNames[c.Index] == "import_path" {
				imports = append(imports, string(code[c.Node.StartByte():c.Node.EndByte()]))
			}
		}
	}
	return imports
}

func TestRubyQueryLanguage(t *testing.T) {
	query := NewRubyQuery()
	lang := query.Language()
	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestRubyQueryPattern(t *testing.T) {
	query := NewRubyQuery()
	pattern := query.Query()
	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}
	lang := sitter.NewLanguage(tree_sitter_ruby.Language())
	_, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("query pattern should compile: %v", err)
	}
}

func TestRubyQueryImportPattern(t *testing.T) {
	query := NewRubyQuery()
	pattern := query.ImportQuery()
	if len(pattern) == 0 {
		t.Fatal("expected non-empty import query pattern")
	}
	lang := sitter.NewLanguage(tree_sitter_ruby.Language())
	_, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("import query pattern should compile: %v", err)
	}
}

func TestRubyQueryExtractFunction(t *testing.T) {
	code := []byte(`
def greet(name)
  puts "Hello, #{name}!"
end

def calculate(x, y, &block)
  block.call(x, y)
end

def complex_method(a, b = 10, *args, key: "default", **opts)
  # body
end
`)
	names := extractRubyNames(t, code)
	for _, expected := range []string{"greet", "calculate", "complex_method"} {
		if !names[expected] {
			t.Errorf("expected to find function '%s'", expected)
		}
	}
}

func TestRubyQueryExtractTypes(t *testing.T) {
	code := []byte(`
class Animal
  def initialize(name)
    @name = name
  end
end

class Dog < Animal
  def speak
    "Woof!"
  end
end

module Helpers
  def helper_method
    true
  end
end
`)
	names := extractRubyNames(t, code)
	for _, expected := range []string{"Animal", "Dog", "Helpers", "initialize", "speak", "helper_method"} {
		if !names[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestRubyQueryExtractClassMethods(t *testing.T) {
	code := []byte(`
class Factory
  def self.create(attrs)
    new(attrs)
  end

  def self.build
    new
  end
end
`)
	names := extractRubyNames(t, code)
	for _, expected := range []string{"Factory", "create", "build"} {
		if !names[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestRubyQueryExtractModuleMethods(t *testing.T) {
	code := []byte(`
module StringUtils
  def self.capitalize(str)
    str.capitalize
  end

  def pad(str, length)
    str.ljust(length)
  end
end
`)
	names := extractRubyNames(t, code)
	for _, expected := range []string{"StringUtils", "capitalize", "pad"} {
		if !names[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestRubyQueryExtractConstants(t *testing.T) {
	code := []byte(`
MAX_SIZE = 100
DEFAULT_NAME = "World"
`)
	names := extractRubyNames(t, code)
	for _, expected := range []string{"MAX_SIZE", "DEFAULT_NAME"} {
		if !names[expected] {
			t.Errorf("expected to find constant '%s'", expected)
		}
	}
}

func TestRubyQueryExtractAttrAccessors(t *testing.T) {
	// attr_reader/writer/accessor are method calls, not captured as declarations.
	// But the methods they define (initialize, etc.) in the class should be captured.
	code := []byte(`
class User
  attr_reader :name, :email
  attr_writer :password
  attr_accessor :role

  def initialize(name, email)
    @name = name
    @email = email
  end
end
`)
	names := extractRubyNames(t, code)
	for _, expected := range []string{"User", "initialize"} {
		if !names[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestRubyQueryExtractImport(t *testing.T) {
	code := []byte(`
require "json"
require "net/http"
require_relative "helper"
require_relative "./lib/utils"
`)
	imports := extractRubyImports(t, code)
	if len(imports) < 4 {
		t.Errorf("expected at least 4 imports, got %d: %v", len(imports), imports)
	}
}

func TestRubyQueryExtractNestedClasses(t *testing.T) {
	code := []byte(`
module Api
  class Client
    def initialize(url)
      @url = url
    end

    def get(path)
      # request
    end

    class Error < StandardError
      def initialize(message)
        super(message)
      end
    end
  end
end
`)
	names := extractRubyNames(t, code)
	for _, expected := range []string{"Api", "Client", "Error", "get"} {
		if !names[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestRubyQueryKindMapping(t *testing.T) {
	query := NewRubyQuery()
	mapping := query.KindMapping()

	expected := map[string]string{
		"method":           "method",
		"singleton_method": "method",
		"class":            "class",
		"module":           "namespace",
		"assignment":       "variable",
	}

	for nodeType, kind := range expected {
		if mapping[nodeType] != kind {
			t.Errorf("expected KindMapping[%q] = %q, got %q", nodeType, kind, mapping[nodeType])
		}
	}
}

func TestRubyQueryCaptures(t *testing.T) {
	query := NewRubyQuery()
	captures := query.Captures()

	expectedCaptures := []string{"name", "signature", "doc", "kind"}
	if len(captures) != len(expectedCaptures) {
		t.Fatalf("expected %d captures, got %d", len(expectedCaptures), len(captures))
	}

	for i, expected := range expectedCaptures {
		if captures[i] != expected {
			t.Errorf("expected capture[%d] = %q, got %q", i, expected, captures[i])
		}
	}
}
