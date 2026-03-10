package languages

import (
	"testing"
	"unsafe"

	tree_sitter_sql "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/sql"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// Helper: extract names captured by @name from SQL query matches.
func extractSQLNames(t *testing.T, code string) []string {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(unsafe.Pointer(tree_sitter_sql.Language()))
	parser.SetLanguage(lang)

	src := []byte(code)
	tree := parser.Parse(src, nil)
	defer tree.Close()

	query := NewSQLQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()
	matches := qc.Matches(q, tree.RootNode(), src)

	captureNames := q.CaptureNames()
	var names []string
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		for _, c := range match.Captures {
			if captureNames[c.Index] == "name" {
				names = append(names, string(src[c.Node.StartByte():c.Node.EndByte()]))
			}
		}
	}
	return names
}

// Helper: extract (name, signature, kind) tuples from SQL query matches.
type sqlCapture struct {
	Name      string
	Signature string
	Kind      string
}

func extractSQLCaptures(t *testing.T, code string) []sqlCapture {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(unsafe.Pointer(tree_sitter_sql.Language()))
	parser.SetLanguage(lang)

	src := []byte(code)
	tree := parser.Parse(src, nil)
	defer tree.Close()

	query := NewSQLQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()
	matches := qc.Matches(q, tree.RootNode(), src)

	captureNames := q.CaptureNames()
	var captures []sqlCapture
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		var sc sqlCapture
		for _, c := range match.Captures {
			text := string(src[c.Node.StartByte():c.Node.EndByte()])
			switch captureNames[c.Index] {
			case "name":
				sc.Name = text
			case "signature":
				sc.Signature = text
			case "kind":
				sc.Kind = c.Node.Kind()
			}
		}
		if sc.Signature != "" || sc.Name != "" {
			captures = append(captures, sc)
		}
	}
	return captures
}

func TestSQLQueryLanguage(t *testing.T) {
	q := NewSQLQuery()
	if q.Language() == nil {
		t.Fatal("Language() returned nil")
	}
}

func TestSQLQueryPattern(t *testing.T) {
	q := NewSQLQuery()
	lang := q.Language()
	_, err := sitter.NewQuery(lang, string(q.Query()))
	if err != nil {
		t.Fatalf("failed to compile SQL query pattern: %v", err)
	}
}

func TestSQLQueryImportPattern(t *testing.T) {
	q := NewSQLQuery()
	if q.ImportQuery() != nil {
		t.Fatal("SQL ImportQuery() should return nil")
	}
}

func TestSQLQueryExtractCreateTable(t *testing.T) {
	code := `
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email TEXT
);

CREATE TABLE IF NOT EXISTS orders (
    order_id INT PRIMARY KEY,
    user_id BIGINT
);
`
	names := extractSQLNames(t, code)
	expected := map[string]bool{"users": false, "orders": false}
	for _, name := range names {
		if _, ok := expected[name]; ok {
			expected[name] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Errorf("expected to find table '%s'", name)
		}
	}

	// Verify kind
	captures := extractSQLCaptures(t, code)
	for _, c := range captures {
		if c.Kind == "create_table" {
			return // found at least one
		}
	}
	t.Error("expected at least one create_table kind")
}

func TestSQLQueryExtractCreateFunction(t *testing.T) {
	code := `
CREATE FUNCTION add_numbers(a INT, b INT) RETURNS INT AS $$
BEGIN
    RETURN a + b;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_user(user_id INT) RETURNS TEXT AS $$
BEGIN
    RETURN 'hello';
END;
$$ LANGUAGE plpgsql;
`
	names := extractSQLNames(t, code)
	expected := map[string]bool{"add_numbers": false, "get_user": false}
	for _, name := range names {
		if _, ok := expected[name]; ok {
			expected[name] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Errorf("expected to find function '%s'", name)
		}
	}
}

func TestSQLQueryExtractCreateView(t *testing.T) {
	code := `
CREATE VIEW active_users AS
SELECT * FROM users WHERE active = true;

CREATE OR REPLACE VIEW user_summary AS
SELECT id, name FROM users;
`
	names := extractSQLNames(t, code)
	expected := map[string]bool{"active_users": false, "user_summary": false}
	for _, name := range names {
		if _, ok := expected[name]; ok {
			expected[name] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Errorf("expected to find view '%s'", name)
		}
	}
}

func TestSQLQueryExtractCreateIndex(t *testing.T) {
	code := `
CREATE INDEX idx_users_name ON users (name);
CREATE UNIQUE INDEX idx_orders_uid ON orders (user_id);
`
	// CREATE INDEX captures whole statement without @name
	// Name extraction happens Go-side (extractSQLDDLName)
	captures := extractSQLCaptures(t, code)
	foundIndex := 0
	for _, c := range captures {
		if c.Kind == "create_index" {
			foundIndex++
		}
	}
	if foundIndex < 2 {
		t.Errorf("expected 2 create_index captures, got %d", foundIndex)
	}
}

func TestSQLQueryExtractCreateType(t *testing.T) {
	code := `
CREATE TYPE mood AS ENUM ('sad', 'ok', 'happy');
`
	names := extractSQLNames(t, code)
	found := false
	for _, name := range names {
		if name == "mood" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find type 'mood'")
	}
}

func TestSQLQueryExtractCreateTrigger(t *testing.T) {
	code := `
CREATE TRIGGER audit_trigger
AFTER INSERT ON users
FOR EACH ROW EXECUTE FUNCTION audit_log();
`
	names := extractSQLNames(t, code)
	// First object_reference should be trigger name
	if len(names) == 0 {
		t.Fatal("expected at least one name capture for trigger")
	}
	if names[0] != "audit_trigger" {
		t.Errorf("expected first captured name to be 'audit_trigger', got '%s'", names[0])
	}
}

func TestSQLQueryExtractCreateSchema(t *testing.T) {
	code := `
CREATE SCHEMA myschema;
`
	names := extractSQLNames(t, code)
	found := false
	for _, name := range names {
		if name == "myschema" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find schema 'myschema'")
	}
}

func TestSQLQueryExtractMaterializedView(t *testing.T) {
	code := `
CREATE MATERIALIZED VIEW user_stats AS
SELECT count(*) as total FROM users;
`
	names := extractSQLNames(t, code)
	found := false
	for _, name := range names {
		if name == "user_stats" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find materialized view 'user_stats'")
	}
}

func TestSQLQueryExtractAlterTable(t *testing.T) {
	code := `
ALTER TABLE users ADD COLUMN email VARCHAR(255);
`
	captures := extractSQLCaptures(t, code)
	found := false
	for _, c := range captures {
		if c.Kind == "alter_table" && c.Name == "users" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find alter_table for 'users'")
	}
}

func TestSQLQueryExtractCreateSequence(t *testing.T) {
	code := `
CREATE SEQUENCE user_id_seq START WITH 1 INCREMENT BY 1;
`
	names := extractSQLNames(t, code)
	found := false
	for _, name := range names {
		if name == "user_id_seq" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find sequence 'user_id_seq'")
	}
}

func TestSQLQueryExtractComments(t *testing.T) {
	code := `
-- This is a single-line comment
CREATE TABLE test (id INT);

/* This is a
   multi-line comment */
CREATE TABLE test2 (id INT);
`
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(unsafe.Pointer(tree_sitter_sql.Language()))
	parser.SetLanguage(lang)

	src := []byte(code)
	tree := parser.Parse(src, nil)
	defer tree.Close()

	query := NewSQLQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()
	matches := qc.Matches(q, tree.RootNode(), src)

	captureNames := q.CaptureNames()
	docCount := 0
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		for _, c := range match.Captures {
			if captureNames[c.Index] == "doc" {
				docCount++
			}
		}
	}
	if docCount < 2 {
		t.Errorf("expected at least 2 comment captures, got %d", docCount)
	}
}

func TestSQLQueryKindMapping(t *testing.T) {
	q := NewSQLQuery()
	km := q.KindMapping()

	expectedKinds := map[string]string{
		"create_table":             "struct",
		"create_function":          "function",
		"create_view":              "type",
		"create_materialized_view": "type",
		"create_index":             "variable",
		"create_trigger":           "function",
		"create_type":              "type",
		"create_schema":            "namespace",
		"create_sequence":          "variable",
		"alter_table":              "type",
	}

	for nodeType, expectedKind := range expectedKinds {
		if got, ok := km[nodeType]; !ok {
			t.Errorf("missing kind mapping for '%s'", nodeType)
		} else if got != expectedKind {
			t.Errorf("kind mapping for '%s': got '%s', want '%s'", nodeType, got, expectedKind)
		}
	}
}

func TestSQLQueryCaptures(t *testing.T) {
	q := NewSQLQuery()
	captures := q.Captures()

	expected := map[string]bool{
		"name":      false,
		"signature": false,
		"doc":       false,
		"kind":      false,
	}
	for _, c := range captures {
		expected[c] = true
	}
	for name, found := range expected {
		if !found {
			t.Errorf("missing capture '%s'", name)
		}
	}
}

func TestSQLQuerySchemaQualifiedName(t *testing.T) {
	code := `
CREATE TABLE myschema.users (
    id INT PRIMARY KEY
);
`
	names := extractSQLNames(t, code)
	found := false
	for _, name := range names {
		if name == "myschema.users" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected schema-qualified name 'myschema.users', got %v", names)
	}
}
