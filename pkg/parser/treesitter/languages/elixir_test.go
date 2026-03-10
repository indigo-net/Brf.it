package languages

import (
	"strings"
	"testing"

	tree_sitter_elixir "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/elixir"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// extractElixirNames is a test helper that parses Elixir code and returns
// all captured @name values from the query matches.
func extractElixirNames(t *testing.T, code []byte) map[string]bool {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(tree_sitter_elixir.Language())
	parser.SetLanguage(lang)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewElixirQuery()
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

// extractElixirSignatures is a test helper that parses Elixir code and returns
// all captured @signature texts from the query matches.
func extractElixirSignatures(t *testing.T, code []byte) []string {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(tree_sitter_elixir.Language())
	parser.SetLanguage(lang)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewElixirQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()
	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	var sigs []string

	for {
		match := matches.Next()
		if match == nil {
			break
		}
		for _, c := range match.Captures {
			if captureNames[c.Index] == "signature" {
				sigs = append(sigs, strings.TrimSpace(string(code[c.Node.StartByte():c.Node.EndByte()])))
			}
		}
	}
	return sigs
}

// extractElixirImportNames is a test helper for import query.
func extractElixirImportNames(t *testing.T, code []byte) []string {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(tree_sitter_elixir.Language())
	parser.SetLanguage(lang)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewElixirQuery()
	importPattern := query.ImportQuery()
	q, err := sitter.NewQuery(lang, string(importPattern))
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
				imports = append(imports, strings.TrimSpace(string(code[c.Node.StartByte():c.Node.EndByte()])))
			}
		}
	}
	return imports
}

func TestElixirQueryLanguage(t *testing.T) {
	query := NewElixirQuery()
	lang := query.Language()
	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestElixirQueryPattern(t *testing.T) {
	query := NewElixirQuery()
	pattern := query.Query()
	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}

	lang := sitter.NewLanguage(tree_sitter_elixir.Language())
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid query pattern: %v", err)
	}
	defer q.Close()
}

func TestElixirQueryImportPattern(t *testing.T) {
	query := NewElixirQuery()
	pattern := query.ImportQuery()
	if len(pattern) == 0 {
		t.Fatal("expected non-empty import query pattern")
	}

	lang := sitter.NewLanguage(tree_sitter_elixir.Language())
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid import query pattern: %v", err)
	}
	defer q.Close()
}

func TestElixirQueryExtractFunction(t *testing.T) {
	code := []byte(`
defmodule MyModule do
  def hello(name) do
    IO.puts("Hello, #{name}!")
  end

  defp private_helper(x, y) do
    x + y
  end

  def zero_arity do
    :ok
  end

  def with_guard(x) when is_integer(x) do
    x * 2
  end
end
`)

	names := extractElixirNames(t, code)
	for _, expected := range []string{"hello", "private_helper", "zero_arity", "with_guard"} {
		if !names[expected] {
			t.Errorf("expected to find function '%s', found: %v", expected, names)
		}
	}
}

func TestElixirQueryExtractModule(t *testing.T) {
	code := []byte(`
defmodule MyApp.Accounts do
  defmodule User do
    defstruct [:name, :email, :age]
  end

  def create_user(attrs) do
    %User{name: attrs.name}
  end
end
`)

	names := extractElixirNames(t, code)
	// Nested modules: MyApp.Accounts is an alias node
	for _, expected := range []string{"User", "create_user"} {
		if !names[expected] {
			t.Errorf("expected to find '%s', found: %v", expected, names)
		}
	}
}

func TestElixirQueryExtractProtocol(t *testing.T) {
	code := []byte(`
defprotocol Printable do
  @doc "Converts the data to a printable string"
  def to_string(data)
end

defimpl Printable, for: Integer do
  def to_string(data) do
    Integer.to_string(data)
  end
end
`)

	names := extractElixirNames(t, code)
	if !names["Printable"] {
		t.Errorf("expected to find protocol 'Printable', found: %v", names)
	}
	if !names["to_string"] {
		t.Errorf("expected to find function 'to_string', found: %v", names)
	}
}

func TestElixirQueryExtractMacro(t *testing.T) {
	code := []byte(`
defmodule MyMacros do
  defmacro unless(condition, do: block) do
    quote do
      if !unquote(condition), do: unquote(block)
    end
  end

  defmacrop private_macro(x) do
    quote do: unquote(x) + 1
  end
end
`)

	names := extractElixirNames(t, code)
	if !names["unless"] {
		t.Errorf("expected to find macro 'unless', found: %v", names)
	}
}

func TestElixirQueryExtractGuard(t *testing.T) {
	code := []byte(`
defmodule Guards do
  defguard is_positive(x) when is_integer(x) and x > 0

  defguardp is_even(x) when rem(x, 2) == 0
end
`)

	names := extractElixirNames(t, code)
	if !names["is_positive"] {
		t.Errorf("expected to find guard 'is_positive', found: %v", names)
	}
}

func TestElixirQueryExtractDelegate(t *testing.T) {
	code := []byte(`
defmodule MyModule do
  defdelegate size(map), to: Map
  defdelegate fetch(map, key), to: Map
end
`)

	names := extractElixirNames(t, code)
	if !names["size"] {
		t.Errorf("expected to find delegate 'size', found: %v", names)
	}
	if !names["fetch"] {
		t.Errorf("expected to find delegate 'fetch', found: %v", names)
	}
}

func TestElixirQueryExtractStruct(t *testing.T) {
	code := []byte(`
defmodule User do
  defstruct [:name, :email, :age]
end
`)

	names := extractElixirNames(t, code)
	// defstruct captures "defstruct" as the name
	if !names["defstruct"] {
		t.Errorf("expected to find 'defstruct', found: %v", names)
	}
}

func TestElixirQueryExtractTypeSpec(t *testing.T) {
	code := []byte(`
defmodule Types do
  @type color :: :red | :green | :blue
  @typep internal_state :: map()
  @opaque hidden :: %__MODULE__{}
  @spec add(integer(), integer()) :: integer()
  @callback handle_event(event :: term()) :: {:ok, term()} | {:error, term()}

  def add(a, b), do: a + b
end
`)

	names := extractElixirNames(t, code)
	// Raw query captures the attribute keyword as @name (spec, type, etc.)
	// The real function/type names are extracted in parser.go via refineElixirAttrKind()
	for _, expected := range []string{"type", "typep", "opaque", "spec", "callback"} {
		if !names[expected] {
			t.Errorf("expected to find attribute '%s', found: %v", expected, names)
		}
	}

	// Also verify the function def is captured
	if !names["add"] {
		t.Errorf("expected to find function 'add', found: %v", names)
	}

	// Verify signatures contain the full attribute text
	sigs := extractElixirSignatures(t, code)
	foundSpec := false
	foundType := false
	for _, sig := range sigs {
		if strings.Contains(sig, "@spec add") {
			foundSpec = true
		}
		if strings.Contains(sig, "@type color") {
			foundType = true
		}
	}
	if !foundSpec {
		t.Errorf("expected to find @spec signature in: %v", sigs)
	}
	if !foundType {
		t.Errorf("expected to find @type signature in: %v", sigs)
	}
}

func TestElixirQueryExtractImport(t *testing.T) {
	code := []byte(`
defmodule MyModule do
  import Enum
  import String, only: [trim: 1, split: 2]
  alias MyApp.Accounts.User
  alias MyApp.Repo, as: R
  use GenServer
  require Logger
end
`)

	imports := extractElixirImportNames(t, code)
	if len(imports) == 0 {
		t.Fatal("expected to find imports")
	}

	// Check that import/alias/use/require are captured
	found := map[string]bool{}
	for _, imp := range imports {
		if strings.Contains(imp, "import Enum") {
			found["import_enum"] = true
		}
		if strings.Contains(imp, "import String") {
			found["import_string"] = true
		}
		if strings.Contains(imp, "alias") && strings.Contains(imp, "User") {
			found["alias_user"] = true
		}
		if strings.Contains(imp, "use GenServer") {
			found["use_genserver"] = true
		}
		if strings.Contains(imp, "require Logger") {
			found["require_logger"] = true
		}
	}

	for _, key := range []string{"import_enum", "import_string"} {
		if !found[key] {
			t.Errorf("expected to find %s in imports: %v", key, imports)
		}
	}
}

func TestElixirQueryExtractZeroArityWithGuard(t *testing.T) {
	code := []byte(`
defmodule Example do
  def guarded when node() == :primary do
    :ok
  end
end
`)

	names := extractElixirNames(t, code)
	if !names["guarded"] {
		t.Errorf("expected to find 'guarded', found: %v", names)
	}
}

func TestElixirQueryKindMapping(t *testing.T) {
	query := NewElixirQuery()
	mapping := query.KindMapping()

	expectedKinds := map[string]string{
		"call":           "function",
		"unary_operator": "type",
	}

	for nodeType, expectedKind := range expectedKinds {
		if got, ok := mapping[nodeType]; !ok {
			t.Errorf("expected mapping for %s", nodeType)
		} else if got != expectedKind {
			t.Errorf("expected %s -> %s, got %s", nodeType, expectedKind, got)
		}
	}
}

func TestElixirQueryCaptures(t *testing.T) {
	query := NewElixirQuery()
	captures := query.Captures()

	expected := map[string]bool{
		"name":      false,
		"signature": false,
		"doc":       false,
		"kind":      false,
	}

	for _, c := range captures {
		if _, ok := expected[c]; ok {
			expected[c] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Errorf("missing capture: %s", name)
		}
	}
}
