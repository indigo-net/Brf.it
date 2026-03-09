package languages

import (
	"testing"
	"unsafe"

	tree_sitter_scala "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/scala"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// extractScalaNames is a test helper that parses Scala code and returns
// all captured @name values from the query matches.
func extractScalaNames(t *testing.T, code []byte) map[string]bool {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(unsafe.Pointer(tree_sitter_scala.Language()))
	parser.SetLanguage(lang)
	tree := parser.Parse(code, nil)
	defer tree.Close()
	query := NewScalaQuery()
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

// extractScalaImports is a test helper that parses Scala code and returns
// all captured import statements.
func extractScalaImports(t *testing.T, code []byte) []string {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(unsafe.Pointer(tree_sitter_scala.Language()))
	parser.SetLanguage(lang)
	tree := parser.Parse(code, nil)
	defer tree.Close()
	query := NewScalaQuery()
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

func TestScalaQueryLanguage(t *testing.T) {
	query := NewScalaQuery()
	lang := query.Language()

	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestScalaQueryPattern(t *testing.T) {
	query := NewScalaQuery()
	pattern := query.Query()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(unsafe.Pointer(tree_sitter_scala.Language()))
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid query pattern: %v", err)
	}
	defer q.Close()
}

func TestScalaQueryImportPattern(t *testing.T) {
	query := NewScalaQuery()
	pattern := query.ImportQuery()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty import query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(unsafe.Pointer(tree_sitter_scala.Language()))
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid import query pattern: %v", err)
	}
	defer q.Close()
}

func TestScalaQueryExtractFunction(t *testing.T) {
	code := []byte(`
def add(a: Int, b: Int): Int = a + b
def greet(name: String): String = s"Hello, $name"
private def helper(): Unit = ()
def noParams: Int = 42
def generic[A](x: A): A = x
def curried(a: Int)(b: Int): Int = a + b
def higherOrder(f: Int => Int): Int = f(42)
`)

	found := extractScalaNames(t, code)

	expected := []string{"add", "greet", "helper", "noParams", "generic", "curried", "higherOrder"}
	for _, name := range expected {
		if !found[name] {
			t.Errorf("expected to find function '%s'", name)
		}
	}
}

func TestScalaQueryExtractTypes(t *testing.T) {
	code := []byte(`
class Person(val name: String, var age: Int)
abstract class Vehicle(val wheels: Int)
case class Point(x: Double, y: Double)
sealed class Node
trait Greeter
sealed trait Shape
object MathUtils
`)

	found := extractScalaNames(t, code)

	expected := []string{"Person", "Vehicle", "Point", "Node", "Greeter", "Shape", "MathUtils"}
	for _, name := range expected {
		if !found[name] {
			t.Errorf("expected to find type '%s'", name)
		}
	}
}

func TestScalaQueryExtractClassBody(t *testing.T) {
	code := []byte(`
class Calculator {
  def add(a: Int, b: Int): Int = a + b
  def subtract(a: Int, b: Int): Int = a - b
  val name: String = "calc"
}
`)

	found := extractScalaNames(t, code)

	expected := []string{"Calculator", "add", "subtract", "name"}
	for _, name := range expected {
		if !found[name] {
			t.Errorf("expected to find '%s'", name)
		}
	}
}

func TestScalaQueryExtractTraitMembers(t *testing.T) {
	code := []byte(`
trait Greeter {
  def greet(name: String): String
  def farewell(name: String): String = s"Goodbye, $name"
  val defaultGreeting: String
}
`)

	found := extractScalaNames(t, code)

	expected := []string{"Greeter", "greet", "farewell", "defaultGreeting"}
	for _, name := range expected {
		if !found[name] {
			t.Errorf("expected to find '%s'", name)
		}
	}
}

func TestScalaQueryExtractObjectMembers(t *testing.T) {
	code := []byte(`
object Config {
  val maxRetries: Int = 3
  var debug: Boolean = false
  def getTimeout: Long = 5000L
  lazy val expensive: String = "computed"
}
`)

	found := extractScalaNames(t, code)

	expected := []string{"Config", "maxRetries", "debug", "getTimeout", "expensive"}
	for _, name := range expected {
		if !found[name] {
			t.Errorf("expected to find '%s'", name)
		}
	}
}

func TestScalaQueryExtractValVar(t *testing.T) {
	code := []byte(`
val immutable: String = "cannot change"
var mutable: Int = 0
lazy val computed: Double = 3.14
implicit val defaultOrd: Ordering[Int] = Ordering.Int
`)

	found := extractScalaNames(t, code)

	expected := []string{"immutable", "mutable", "computed", "defaultOrd"}
	for _, name := range expected {
		if !found[name] {
			t.Errorf("expected to find val/var '%s'", name)
		}
	}
}

func TestScalaQueryExtractTypeAlias(t *testing.T) {
	code := []byte(`
type StringList = List[String]
type Callback[A] = A => Unit
type Pair = (Int, String)
`)

	found := extractScalaNames(t, code)

	expected := []string{"StringList", "Callback", "Pair"}
	for _, name := range expected {
		if !found[name] {
			t.Errorf("expected to find type alias '%s'", name)
		}
	}
}

func TestScalaQueryExtractEnum(t *testing.T) {
	code := []byte(`
enum Color {
  case Red, Green, Blue
}

enum Planet(mass: Double, radius: Double) {
  case Mercury extends Planet(3.303e+23, 2.4397e6)
  case Earth extends Planet(5.976e+24, 6.37814e6)
}
`)

	found := extractScalaNames(t, code)

	expected := []string{"Color", "Planet"}
	for _, name := range expected {
		if !found[name] {
			t.Errorf("expected to find enum '%s'", name)
		}
	}
}

func TestScalaQueryExtractCaseClass(t *testing.T) {
	code := []byte(`
case class Point(x: Double, y: Double) {
  def distance: Double = math.sqrt(x * x + y * y)
}

case class User(name: String, email: String, age: Int)
`)

	found := extractScalaNames(t, code)

	expected := []string{"Point", "distance", "User"}
	for _, name := range expected {
		if !found[name] {
			t.Errorf("expected to find '%s'", name)
		}
	}
}

func TestScalaQueryExtractExtension(t *testing.T) {
	code := []byte(`
extension (s: String)
  def greetWith: String = s"Hello, $s!"

extension [A](list: List[A])
  def secondOption: Option[A] = list.drop(1).headOption
`)

	found := extractScalaNames(t, code)

	// Extension methods capture the inner def names
	expected := []string{"greetWith", "secondOption"}
	for _, name := range expected {
		if !found[name] {
			t.Errorf("expected to find extension method '%s'", name)
		}
	}
}

func TestScalaQueryExtractImport(t *testing.T) {
	code := []byte(`
import scala.collection.mutable
import scala.collection.mutable.{ListBuffer, ArrayBuffer}
import java.util._
import scala.math.sqrt
`)

	imports := extractScalaImports(t, code)

	if len(imports) != 4 {
		t.Fatalf("expected 4 imports, got %d: %v", len(imports), imports)
	}

	expected := []string{
		"import scala.collection.mutable",
		"import scala.collection.mutable.{ListBuffer, ArrayBuffer}",
		"import java.util._",
		"import scala.math.sqrt",
	}
	for i, exp := range expected {
		if imports[i] != exp {
			t.Errorf("import[%d]: expected %q, got %q", i, exp, imports[i])
		}
	}
}

func TestScalaQueryExtractGenerics(t *testing.T) {
	code := []byte(`
def identity[A](x: A): A = x
def map[A, B](list: List[A])(f: A => B): List[B] = list.map(f)
class Container[T](val value: T)
trait Functor[F[_]]
`)

	found := extractScalaNames(t, code)

	expected := []string{"identity", "map", "Container", "Functor"}
	for _, name := range expected {
		if !found[name] {
			t.Errorf("expected to find generic '%s'", name)
		}
	}
}

func TestScalaQueryKindMapping(t *testing.T) {
	query := NewScalaQuery()
	mapping := query.KindMapping()

	expectedMappings := map[string]string{
		"function_definition":  "method",
		"function_declaration": "method",
		"class_definition":     "class",
		"trait_definition":     "trait",
		"object_definition":    "class",
		"val_definition":       "variable",
		"val_declaration":      "variable",
		"var_definition":       "variable",
		"var_declaration":      "variable",
		"type_definition":      "type",
		"enum_definition":      "enum",
		"given_definition":     "variable",
	}

	for key, expectedValue := range expectedMappings {
		if value, ok := mapping[key]; !ok {
			t.Errorf("expected mapping for %q", key)
		} else if value != expectedValue {
			t.Errorf("mapping[%q] = %q, want %q", key, value, expectedValue)
		}
	}
}

func TestScalaQueryCaptures(t *testing.T) {
	query := NewScalaQuery()
	captures := query.Captures()

	expected := []string{"name", "signature", "doc", "kind"}
	if len(captures) != len(expected) {
		t.Fatalf("expected %d captures, got %d", len(expected), len(captures))
	}
	for i, name := range expected {
		if captures[i] != name {
			t.Errorf("captures[%d] = %q, want %q", i, captures[i], name)
		}
	}
}
