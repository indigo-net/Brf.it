package languages

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_java "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

func TestJavaQueryLanguage(t *testing.T) {
	query := NewJavaQuery()
	lang := query.Language()

	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestJavaQueryPattern(t *testing.T) {
	query := NewJavaQuery()
	pattern := query.Query()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_java.Language())
	_, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid query pattern: %v", err)
	}
}

func TestJavaQueryKindMapping(t *testing.T) {
	query := NewJavaQuery()
	mapping := query.KindMapping()

	expectedMappings := map[string]string{
		"class_declaration":           "class",
		"interface_declaration":       "interface",
		"method_declaration":          "method",
		"constructor_declaration":     "constructor",
		"enum_declaration":            "enum",
		"annotation_type_declaration": "annotation",
		"record_declaration":          "record",
	}

	for nodeType, expected := range expectedMappings {
		if actual, ok := mapping[nodeType]; !ok {
			t.Errorf("expected mapping for '%s'", nodeType)
		} else if actual != expected {
			t.Errorf("expected '%s' for '%s', got '%s'", expected, nodeType, actual)
		}
	}
}

func TestJavaQueryExtractClass(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_java.Language())
	parser.SetLanguage(lang)

	code := []byte(`package com.example;

// User class represents a user.
public class User {
    private String name;

    public User(String name) {
        this.name = name;
    }

    public String getName() {
        return name;
    }
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewJavaQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	var foundClass, foundMethod bool

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		captures := make(map[string]string)
		for _, c := range match.Captures {
			name := captureNames[c.Index]
			captures[name] = string(code[c.Node.StartByte():c.Node.EndByte()])
		}

		switch captures["name"] {
		case "User":
			if _, hasKind := captures["kind"]; hasKind {
				foundClass = true
			}
		case "getName":
			foundMethod = true
		}
	}

	if !foundClass {
		t.Error("expected to find class 'User'")
	}
	if !foundMethod {
		t.Error("expected to find method 'getName'")
	}
}

func TestJavaQueryExtractInterface(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_java.Language())
	parser.SetLanguage(lang)

	code := []byte(`package com.example;

public interface Repository<T> {
    T findById(String id);
    void save(T entity);
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewJavaQuery()
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
			name := captureNames[c.Index]
			if name == "name" {
				foundNames[string(code[c.Node.StartByte():c.Node.EndByte()])] = true
			}
		}
	}

	expected := []string{"Repository", "findById", "save"}
	for _, name := range expected {
		if !foundNames[name] {
			t.Errorf("expected to find '%s'", name)
		}
	}
}

func TestJavaQueryExtractEnum(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_java.Language())
	parser.SetLanguage(lang)

	code := []byte(`package com.example;

public enum Status {
    PENDING,
    ACTIVE,
    COMPLETED
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewJavaQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	var foundEnum bool

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, c := range match.Captures {
			name := captureNames[c.Index]
			if name == "name" {
				text := string(code[c.Node.StartByte():c.Node.EndByte()])
				if text == "Status" {
					foundEnum = true
				}
			}
		}
	}

	if !foundEnum {
		t.Error("expected to find enum 'Status'")
	}
}

func TestJavaQueryExtractAnnotationType(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_java.Language())
	parser.SetLanguage(lang)

	code := []byte(`package com.example;

public @interface Inject {
    String value() default "";
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewJavaQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	var foundAnnotation bool

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, c := range match.Captures {
			name := captureNames[c.Index]
			if name == "name" {
				text := string(code[c.Node.StartByte():c.Node.EndByte()])
				if text == "Inject" {
					foundAnnotation = true
				}
			}
		}
	}

	if !foundAnnotation {
		t.Error("expected to find annotation type 'Inject'")
	}
}

func TestJavaQueryExtractRecord(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_java.Language())
	parser.SetLanguage(lang)

	code := []byte(`package com.example;

public record Point(int x, int y) {
    public double distance() {
        return Math.sqrt(x * x + y * y);
    }
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewJavaQuery()
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
			name := captureNames[c.Index]
			if name == "name" {
				foundNames[string(code[c.Node.StartByte():c.Node.EndByte()])] = true
			}
		}
	}

	if !foundNames["Point"] {
		t.Error("expected to find record 'Point'")
	}
	if !foundNames["distance"] {
		t.Error("expected to find method 'distance'")
	}
}

func TestJavaQueryExtractGenerics(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_java.Language())
	parser.SetLanguage(lang)

	code := []byte(`package com.example;

public class Box<T extends Comparable<T>> {
    private T value;

    public <U> U transform(Function<T, U> fn) {
        return fn.apply(value);
    }
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewJavaQuery()
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
			name := captureNames[c.Index]
			if name == "name" {
				foundNames[string(code[c.Node.StartByte():c.Node.EndByte()])] = true
			}
		}
	}

	if !foundNames["Box"] {
		t.Error("expected to find class 'Box'")
	}
	if !foundNames["transform"] {
		t.Error("expected to find method 'transform'")
	}
}
