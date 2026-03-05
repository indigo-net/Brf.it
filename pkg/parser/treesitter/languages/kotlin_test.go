package languages

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_kotlin "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/kotlin"
)

// extractKotlinNames is a test helper that parses Kotlin code and returns
// all captured @name values from the query matches.
func extractKotlinNames(t *testing.T, code []byte) map[string]bool {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(tree_sitter_kotlin.Language())
	parser.SetLanguage(lang)
	tree := parser.Parse(code, nil)
	defer tree.Close()
	query := NewKotlinQuery()
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

func TestKotlinQueryLanguage(t *testing.T) {
	query := NewKotlinQuery()
	lang := query.Language()

	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestKotlinQueryPattern(t *testing.T) {
	query := NewKotlinQuery()
	pattern := query.Query()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_kotlin.Language())
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid query pattern: %v", err)
	}
	defer q.Close()
}

func TestKotlinQueryImportPattern(t *testing.T) {
	query := NewKotlinQuery()
	pattern := query.ImportQuery()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty import query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_kotlin.Language())
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid import query pattern: %v", err)
	}
	defer q.Close()
}

func TestKotlinQueryExtractFunction(t *testing.T) {
	code := []byte(`
fun add(a: Int, b: Int): Int {
    return a + b
}

fun double(x: Int) = x * 2

suspend fun fetchData(url: String): String {
    return ""
}

inline fun <reified T> parseJson(json: String): T {
    return json as T
}

infix fun Int.power(exponent: Int): Int {
    return this
}

operator fun Point.plus(other: Point): Point {
    return other
}

tailrec fun factorial(n: Int, acc: Int = 1): Int {
    return if (n <= 1) acc else factorial(n - 1, n * acc)
}

fun String.isEmail(): Boolean {
    return true
}

fun <T, R> transform(input: T, mapper: (T) -> R): R {
    return mapper(input)
}

fun log(message: String, level: String = "INFO", vararg tags: String) {
    println(message)
}

fun buildHtml(init: StringBuilder.() -> Unit): String {
    return StringBuilder().apply(init).toString()
}
`)

	foundNames := extractKotlinNames(t, code)

	expectedNames := []string{"add", "double", "fetchData", "parseJson", "power", "plus", "factorial", "isEmail", "transform", "log", "buildHtml"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find function '%s'", expected)
		}
	}
}

func TestKotlinQueryExtractTypes(t *testing.T) {
	code := []byte(`
class User(val name: String, val age: Int) {
    fun greet(): String = name
}

data class Point(val x: Double, val y: Double)

data class ApiResponse<T>(val status: Int, val data: T) {
    fun isSuccess(): Boolean = status in 200..299
}

sealed class Result<out T> {
    data class Success<T>(val data: T) : Result<T>()
    data class Error(val exception: Throwable) : Result<Nothing>()
    object Loading : Result<Nothing>()
}

abstract class BaseRepository<T> {
    abstract fun fetch(): T
}

open class Vehicle(val brand: String) {
    open fun start() {}
}

class Outer(val value: Int) {
    inner class Inner {
        fun access(): Int = value
    }
}

annotation class ApiEndpoint(val path: String)

value class UserId(val value: Long)

enum class HttpStatus(val code: Int) {
    OK(200),
    NOT_FOUND(404);
    fun isSuccess(): Boolean = code in 200..299
}
`)

	foundNames := extractKotlinNames(t, code)

	expectedNames := []string{
		"User", "Point", "ApiResponse", "Result", "Success", "Error", "Loading",
		"BaseRepository", "Vehicle", "Outer", "Inner", "ApiEndpoint", "UserId", "HttpStatus",
	}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find type '%s'", expected)
		}
	}
}

func TestKotlinQueryExtractInterface(t *testing.T) {
	code := []byte(`
interface Repository<T> {
    fun getAll(): List<T>
    fun getById(id: String): T?
}

interface Logger {
    fun log(message: String)
    fun info(message: String) = log(message)
}

sealed interface State {
    object Initial : State
}

fun interface Predicate<T> {
    fun test(value: T): Boolean
}
`)

	foundNames := extractKotlinNames(t, code)

	expectedNames := []string{"Repository", "Logger", "State", "Predicate"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find interface '%s'", expected)
		}
	}
}

func TestKotlinQueryExtractObject(t *testing.T) {
	code := []byte(`
object DatabaseConnection {
    fun execute(query: String): List<Any> = emptyList()
}

class Config(val apiUrl: String) {
    companion object {
        fun fromEnv(): Config = Config("localhost")
        const val VERSION = "1.0.0"
    }
}

class Database {
    companion object Factory {
        fun create(): Database = Database()
    }
}
`)

	foundNames := extractKotlinNames(t, code)

	expectedNames := []string{"DatabaseConnection", "Config", "Database", "Factory"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find object '%s'", expected)
		}
	}
}

func TestKotlinQueryExtractProperties(t *testing.T) {
	code := []byte(`
val MAX_SIZE = 100
val PI: Double = 3.14159
var counter = 0
const val DEFAULT_TIMEOUT = 5000L
const val API_VERSION = "v1"
lateinit var repository: String
`)

	foundNames := extractKotlinNames(t, code)

	expectedNames := []string{"MAX_SIZE", "PI", "counter", "DEFAULT_TIMEOUT", "API_VERSION", "repository"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find property '%s'", expected)
		}
	}
}

func TestKotlinQueryExtractTypeAlias(t *testing.T) {
	code := []byte(`
typealias UserId = String
typealias Timestamp = Long
typealias Handler<T> = (T) -> Unit
typealias Transform<T, R> = (T) -> R
`)

	foundNames := extractKotlinNames(t, code)

	expectedNames := []string{"UserId", "Timestamp", "Handler", "Transform"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find type alias '%s'", expected)
		}
	}
}

func TestKotlinQueryExtractEnumEntry(t *testing.T) {
	code := []byte(`
enum class Color(val hex: String) {
    RED("#FF0000"),
    GREEN("#00FF00"),
    BLUE("#0000FF");
    fun toRgb(): String = hex
}
`)

	foundNames := extractKotlinNames(t, code)

	expectedNames := []string{"RED", "GREEN", "BLUE"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find enum entry '%s'", expected)
		}
	}
}

func TestKotlinQueryExtractImport(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_kotlin.Language())
	parser.SetLanguage(lang)

	code := []byte(`
import kotlin.collections.List
import com.example.models.User
import com.example.utils.*
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewKotlinQuery()
	q, err := sitter.NewQuery(lang, string(query.ImportQuery()))
	if err != nil {
		t.Fatalf("failed to create import query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	count := 0
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		count++
	}

	// Should find 3 import declarations
	if count < 3 {
		t.Errorf("expected at least 3 import declarations, got %d", count)
	}
}

func TestKotlinQueryExtractGenerics(t *testing.T) {
	code := []byte(`
class Container<T>(val item: T)

fun <T : Comparable<T>> maximum(a: T, b: T): T {
    return if (a > b) a else b
}

interface Producer<out T> {
    fun produce(): T
}

interface Consumer<in T> {
    fun consume(item: T)
}

class Cache<K : Comparable<K>, V : Any> {
    fun put(key: K, value: V) {}
}

inline fun <reified T> isInstance(value: Any): Boolean {
    return value is T
}
`)

	foundNames := extractKotlinNames(t, code)

	expectedNames := []string{"Container", "maximum", "Producer", "Consumer", "Cache", "isInstance"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestKotlinQueryKindMapping(t *testing.T) {
	query := NewKotlinQuery()
	mapping := query.KindMapping()

	expectedMappings := map[string]string{
		"function_declaration":  "function",
		"class_declaration":     "class",
		"object_declaration":    "class",
		"companion_object":      "class",
		"property_declaration":  "variable",
		"type_alias":            "type",
		"enum_entry":            "variable",
		"secondary_constructor": "constructor",
	}

	for nodeType, expected := range expectedMappings {
		if actual, ok := mapping[nodeType]; !ok {
			t.Errorf("expected mapping for '%s'", nodeType)
		} else if actual != expected {
			t.Errorf("expected '%s' -> '%s', got '%s'", nodeType, expected, actual)
		}
	}
}

func TestKotlinQueryCaptures(t *testing.T) {
	query := NewKotlinQuery()
	captures := query.Captures()

	expected := []string{"name", "signature", "doc", "kind"}
	if len(captures) != len(expected) {
		t.Errorf("expected %d captures, got %d", len(expected), len(captures))
	}

	for i, exp := range expected {
		if captures[i] != exp {
			t.Errorf("expected capture[%d] = '%s', got '%s'", i, exp, captures[i])
		}
	}
}
