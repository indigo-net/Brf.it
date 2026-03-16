# Code Summary: /home/runner/work/Brf.it/Brf.it

*brf.it 0.20.0*

## Files

### /home/runner/work/Brf.it/Brf.it/cmd/brfit/main.go

```go
version = "dev"
commit  = "none"
date    = "unknown"
func main()
```

### /home/runner/work/Brf.it/Brf.it/cmd/brfit/root.go

```go
import (
	gocontext "context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"github.com/indigo-net/Brf.it/internal/config"
	"github.com/indigo-net/Brf.it/internal/context"
	"github.com/indigo-net/Brf.it/pkg/scanner"
	"github.com/indigo-net/Brf.it/pkg/tokenizer"
	"github.com/spf13/cobra"
	// Import treesitter parser to register Go/TypeScript parsers
	_ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"
)
Version = "dev"
Commit  = "none"
Date    = "unknown"
func SetBuildInfo(v, c, d string)
cfg *config.Config
rootCmd *cobra.Command
func init()
func Execute()
func NewRootCommand() *cobra.Command
func newRootCommandWithConfig(c *config.Config) *cobra.Command
func addFlags(cmd *cobra.Command, c *config.Config)
func runRoot(cmd *cobra.Command, args []string, c *config.Config) error
changedFiles map[string]bool
func writeOutput(result *context.Result, c *config.Config) error
func runTokenTree(ctx gocontext.Context, scanOpts *scanner.ScanOptions, rootPath string) error
func resolveChangedFiles(rootPath string, changed bool, since string) (map[string]bool, error)
diffArgs []string
func splitNonEmpty(s string) []string
func cloneRemote(ctx gocontext.Context, remote string) (string, func(), error)
func resolveRemoteURL(remote string) string
func writeToFile(path string, content []byte) error
```

### /home/runner/work/Brf.it/Brf.it/cmd/brfit/root_test.go

```go
import (
	"bytes"
	gocontext "context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"github.com/indigo-net/Brf.it/internal/config"
	// Import treesitter parser to register it
	_ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"
)
func TestExecuteHelp(t *testing.T)
func TestExecuteVersion(t *testing.T)
buf bytes.Buffer
func TestNewRootCommand(t *testing.T)
func TestParseFlags(t *testing.T)
func TestRootCommandIntegration(t *testing.T)
buf bytes.Buffer
func TestRootCommandIntegrationMarkdown(t *testing.T)
buf bytes.Buffer
func TestRootCommandPathNotFound(t *testing.T)
func TestResolveChangedFilesPathAnchoring(t *testing.T)
func TestResolveChangedFilesIncludesUntracked(t *testing.T)
func TestResolveChangedFilesEmptyOutput(t *testing.T)
func TestResolveRemoteURL(t *testing.T)
func TestRemoteFlagConflictsWithPath(t *testing.T)
func TestRemoteFlagConflictsWithChanged(t *testing.T)
func TestRemoteFlagConflictsWithSince(t *testing.T)
func TestRemoteFlagInvalidURL(t *testing.T)
func TestCloneRemoteIntegration(t *testing.T)
func TestWriteToFile(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/cmd/brfit-mcp/main.go

```go
import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"github.com/indigo-net/Brf.it/internal/config"
	pkgcontext "github.com/indigo-net/Brf.it/internal/context"
	"github.com/indigo-net/Brf.it/pkg/scanner"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	// Register Tree-sitter parsers.
	_ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"
)
version = "dev"
func main()
type SummarizeProjectInput struct {
	Path          string `json:"path,omitempty" jsonschema:"project directory path (defaults to server root)"`
	Format        string `json:"format,omitempty" jsonschema:"output format: xml, md, or json (default: xml)"`
	IncludeBody   bool   `json:"include_body,omitempty" jsonschema:"include function bodies (default: false)"`
	IncludeImport bool   `json:"include_imports,omitempty" jsonschema:"include import statements (default: false)"`
	CallGraph     bool   `json:"call_graph,omitempty" jsonschema:"include function call graph (default: false)"`
}
type SummarizeProjectOutput struct {
	Content         string `json:"content" jsonschema:"the formatted project summary"`
	TotalFiles      int    `json:"total_files" jsonschema:"number of files processed"`
	TotalSignatures int    `json:"total_signatures" jsonschema:"number of signatures extracted"`
	TokenCount      int    `json:"token_count,omitempty" jsonschema:"estimated token count of output"`
}
func makeSummarizeProject(defaultRoot string) func(context.Context, *mcp.CallToolRequest, SummarizeProjectInput) (*mcp.CallToolResult, SummarizeProjectOutput, error)
type SummarizeFileInput struct {
	Path    string `json:"path" jsonschema:"project directory path"`
	Include string `json:"include" jsonschema:"glob pattern to include (e.g. 'pkg/**/*.go')"`
	Format  string `json:"format,omitempty" jsonschema:"output format: xml, md, or json (default: xml)"`
}
type SummarizeFileOutput struct {
	Content         string `json:"content" jsonschema:"the formatted file summary"`
	TotalFiles      int    `json:"total_files" jsonschema:"number of files processed"`
	TotalSignatures int    `json:"total_signatures" jsonschema:"number of signatures extracted"`
}
func makeSummarizeFile(defaultRoot string) func(context.Context, *mcp.CallToolRequest, SummarizeFileInput) (*mcp.CallToolResult, SummarizeFileOutput, error)
validFormats = map[string]bool{"xml": true, "md": true, "markdown": true, "json": true}
func resolvePath(defaultRoot, inputPath string) (string, error)
func validateFormat(format string) (string, error)
func runPackager(ctx context.Context, cfg *config.Config) (*pkgcontext.Result, error)
```

### /home/runner/work/Brf.it/Brf.it/cmd/brfit-mcp/main_test.go

```go
import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	_ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"
)
func TestSummarizeProject(t *testing.T)
func TestSummarizeFile(t *testing.T)
func TestSummarizeProjectInvalidPath(t *testing.T)
func TestPathTraversal(t *testing.T)
func TestInvalidFormat(t *testing.T)
func TestValidSubdirectoryPath(t *testing.T)
func TestResolvePathAbsoluteRejected(t *testing.T)
func TestResolvePathValidRelative(t *testing.T)
func TestResolvePathSymlinkEscape(t *testing.T)
func TestResolvePathSymlinkWithinRoot(t *testing.T)
func TestResolvePathEmpty(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/examples/go/main.go

```go
import (
	"fmt"
	"math"
)
type Point struct {
	X, Y float64
}
func (p Point) Distance() float64
func (p Point) Add(other Point) Point
type Shape interface {
	Area() float64
	Perimeter() float64
}
type Circle struct {
	Center Point
	Radius float64
}
func (c Circle) Area() float64
func (c Circle) Perimeter() float64
func NewCircle(center Point, radius float64) (*Circle, error)
func main()
```

### /home/runner/work/Brf.it/Brf.it/examples/java/ShapeService.java

```java
import java.util.List;
import java.util.ArrayList;
import java.util.Optional;
interface Shape
double area();
double perimeter();
String name();
class Circle implements Shape
public Circle(double radius)
@Override
    public double area()
@Override
    public double perimeter()
@Override
    public String name()
class Rectangle implements Shape
public Rectangle(double width, double height)
@Override
    public double area()
@Override
    public double perimeter()
@Override
    public String name()
public class ShapeService
public void addShape(Shape shape)
public double totalArea()
public Optional<Shape> largestShape()
```

### /home/runner/work/Brf.it/Brf.it/examples/python/api.py

```python
from dataclasses import dataclass, field
from datetime import datetime
from typing import Optional
from enum import Enum
class TaskStatus(Enum)
class Task
class TaskRepository
def __init__(self)
def create(self, title: str) -> Task
def get(self, task_id: int) -> Optional[Task]
def complete(self, task_id: int) -> bool
def list_by_status(self, status: TaskStatus) -> list[Task]
def format_task(task: Task) -> str
```

### /home/runner/work/Brf.it/Brf.it/examples/rust/lib.rs

```rust
use std::collections::HashMap;
use std::fmt;
pub struct Cache<V>
pub enum CacheError
impl fmt::Display for CacheError
fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result
pub trait Cacheable: Clone + fmt::Debug
fn size(&self) -> usize;
impl<V: Clone> Cache<V>
pub fn new(capacity: usize) -> Self
pub fn get(&self, key: &str) -> Option<&V>
pub fn insert(&mut self, key: String, value: V) -> Result<(), CacheError>
pub fn remove(&mut self, key: &str) -> Result<V, CacheError>
pub fn len(&self) -> usize
pub fn is_empty(&self) -> bool
```

### /home/runner/work/Brf.it/Brf.it/examples/sql/schema.sql

```sql
CREATE TABLE products (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    category VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
)
CREATE TABLE orders (
    id BIGINT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    total DECIMAL(10, 2),
    status VARCHAR(50) DEFAULT 'pending',
    ordered_at TIMESTAMP DEFAULT NOW()
)
CREATE TABLE order_items (
    id BIGINT PRIMARY KEY,
    order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 1
)
CREATE FUNCTION revenue_by_category(cat VARCHAR)
RETURNS DECIMAL LANGUAGE sql
CREATE VIEW top_products
CREATE INDEX idx_products_category ON products (category)
CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'shipped', 'delivered')
```

### /home/runner/work/Brf.it/Brf.it/examples/typescript/app.ts

```typescript
interface AppConfig {
  port: number;
  host: string;
  debug: boolean;
}
interface User {
  id: string;
  email: string;
  name: string;
  createdAt: Date;
}
export function createDefaultConfig(): AppConfig
class Repository<T extends { id: string }>
async findById(id: string): Promise<T | undefined>
async save(item: T): Promise<T>
async delete(id: string): Promise<boolean>
async findAll(): Promise<T[]>
export const formatUser = (user: User): string
export function isUser(value: unknown): value is User
```

### /home/runner/work/Brf.it/Brf.it/install.sh

```shell
set -e
printf "${GREEN}==>${NC} %s\n" "$1"
printf "${YELLOW}==>${NC} %s\n" "$1"
printf "${RED}Error:${NC} %s\n" "$1"
exit 1
uname -s
tr '[:upper:]' '[:lower:]'
echo "$OS"
error "Unsupported OS: $OS. Use Linux or macOS."
uname -m
echo "amd64"
echo "arm64"
error "Unsupported architecture: $ARCH. Use x86_64 or arm64."
command -v curl
curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest"
grep '"tag_name":'
sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget
then
wget -qO- "https://api.github.com/repos/${REPO}/releases/latest"
grep '"tag_name":'
sed -E 's/.*"([^"]+)".*/\1/'
error "Neither curl nor wget found. Please install one of them."
command -v curl
curl -fsSL "$URL" -o "$OUTPUT"
command -v wget
wget -q "$URL" -O "$OUTPUT"
error "Neither curl nor wget found. Please install one of them."
command -v sha256sum
sha256sum "$FILE"
awk '{print $1}'
command -v shasum
shasum -a 256 "$FILE"
awk '{print $1}'
error "Neither sha256sum nor shasum found. Cannot verify checksum."
sha256 "$FILE"
error "Checksum mismatch!\n  Expected: $EXPECTED\n  Actual:   $ACTUAL"
info "Checksum verified"
echo ""
echo "sudo"
return 0
return 1
info "Fetching latest version..."
get_latest_version
error "Failed to get latest version. Please specify version manually."
detect_os
detect_arch
info "Detected: ${OS}/${ARCH}"
info "Installing brfit $VERSION"
mktemp -d
trap 'rm -rf "$TMP_DIR"' EXIT
info "Downloading $ARCHIVE_NAME..."
download "$DOWNLOAD_URL" "$TMP_DIR/$ARCHIVE_NAME"
info "Downloading checksums..."
download "$CHECKSUM_URL" "$TMP_DIR/checksums.txt"
grep "$ARCHIVE_NAME" "$TMP_DIR/checksums.txt"
awk '{print $1}'
error "Checksum not found for $ARCHIVE_NAME"
info "Verifying checksum..."
verify_checksum "$TMP_DIR/$ARCHIVE_NAME" "$EXPECTED_CHECKSUM"
info "Extracting..."
tar -xzf "$TMP_DIR/$ARCHIVE_NAME" -C "$TMP_DIR"
need_sudo
info "Installing to $INSTALL_DIR (requires sudo)..."
info "Installing to $INSTALL_DIR..."
$SUDO mkdir -p "$INSTALL_DIR"
$SUDO mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
$SUDO chmod +x "$INSTALL_DIR/$BINARY_NAME"
info "brfit $VERSION installed successfully!"
echo ""
warn "macOS users: If 'brfit' is blocked, run:"
echo "    xattr -d com.apple.quarantine $INSTALL_DIR/brfit"
check_path
echo ""
warn "$INSTALL_DIR is not in your PATH."
echo ""
echo "Add this line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
echo ""
echo "Then restart your terminal or run:"
echo "  source ~/.bashrc  # or ~/.zshrc"
echo ""
echo "Run 'brfit --help' to get started."
main "$@"
REPO="indigo-net/Brf.it"
INSTALL_DIR="${BRFIT_INSTALL_DIR:-/usr/local/bin}"
BINARY_NAME="brfit"
ARCHIVE_PREFIX="Brf.it"
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'
info()
warn()
error()
detect_os()
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
detect_arch()
ARCH=$(uname -m)
get_latest_version()
download()
URL="$1"
OUTPUT="$2"
sha256()
FILE="$1"
verify_checksum()
FILE="$1"
EXPECTED="$2"
ACTUAL=$(sha256 "$FILE")
need_sudo()
check_path()
main()
VERSION="${1:-}"
VERSION=$(get_latest_version)
VERSION="v$VERSION"
OS=$(detect_os)
ARCH=$(detect_arch)
VERSION_NUM="${VERSION#v}"
ARCHIVE_NAME="${ARCHIVE_PREFIX}_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"
CHECKSUM_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"
TMP_DIR=$(mktemp -d)
EXPECTED_CHECKSUM=$(grep "$ARCHIVE_NAME" "$TMP_DIR/checksums.txt" | awk '{print $1}')
SUDO=$(need_sudo)
```

### /home/runner/work/Brf.it/Brf.it/internal/config/config.go

```go
import (
	"errors"
	"fmt"
	"os"
	pkgcontext "github.com/indigo-net/Brf.it/internal/context"
)
MaxFileSizeUpperBound = 10 * 1024 * 1024
type Config struct {
	// Path is the root directory or file to process.
	Path string

	// Version is the brf.it version string.
	Version string

	// Mode determines what to extract. Currently only "sig" (signature) is supported.
	Mode string

	// Format specifies the output format: "xml" or "md".
	Format string

	// Output is the file path to write output. Empty means stdout.
	Output string

	// IgnoreFiles is the list of ignore file paths (default: [".gitignore"]).
	IgnoreFiles []string

	// IncludePatterns is a list of glob patterns. Only matching files are included.
	IncludePatterns []string

	// ExcludePatterns is a list of glob patterns. Matching files are excluded.
	ExcludePatterns []string

	// IncludeHidden determines whether to include hidden files (dotfiles).
	IncludeHidden bool

	// IncludeBody determines whether to include function/method bodies.
	// When false (default), only signatures are extracted.
	IncludeBody bool

	// IncludeImports determines whether to include import/export statements.
	IncludeImports bool

	// DedupeImports deduplicates imports across files and shows them globally.
	// Requires IncludeImports to be true.
	DedupeImports bool

	// NoTree skips directory tree generation in output.
	NoTree bool

	// NoTokens disables token count calculation.
	NoTokens bool

	// IncludePrivate determines whether to include non-exported/private symbols.
	IncludePrivate bool

	// Changed restricts scanning to files changed in the git working tree.
	Changed bool

	// Since restricts scanning to files changed since the specified commit/tag.
	Since string

	// TokenTree outputs a directory tree with per-file token counts and exits.
	TokenTree bool

	// SecurityCheck enables secret detection and redaction (default: true).
	SecurityCheck bool

	// NoSchema skips the schema section in XML output.
	NoSchema bool

	// CallGraph enables function call graph extraction in output.
	CallGraph bool

	// Remote is a git URL or owner/repo shorthand for remote repository analysis.
	Remote string

	// MaxFileSize is the maximum file size in bytes to process.
	MaxFileSize int64

	// MaxDocLength is the maximum length of documentation comments in characters.
	// 0 means no limit (default).
	MaxDocLength int
}
func DefaultConfig() *Config
func (c *Config) Validate() error
func (c *Config) SupportedExtensions() map[string]string
func (c *Config) ToOptions() *pkgcontext.Options
```

### /home/runner/work/Brf.it/Brf.it/internal/config/config_test.go

```go
import (
	"bytes"
	"os"
	"strings"
	"testing"
)
func TestDefaultConfig(t *testing.T)
expectedMaxSize = 512000
func TestConfigValidate(t *testing.T)
func TestConfigSupportedLanguages(t *testing.T)
func TestValidateMaxFileSizeUpperBound(t *testing.T)
buf bytes.Buffer
func TestToOptionsIncludePrivate(t *testing.T)
func containsString(s, substr string) bool
func containsSubstring(s, substr string) bool
```

### /home/runner/work/Brf.it/Brf.it/internal/context/context.go

```go
import (
	"context"
	"io"
	"os"
	"sort"
	"github.com/indigo-net/Brf.it/pkg/extractor"
	"github.com/indigo-net/Brf.it/pkg/formatter"
	"github.com/indigo-net/Brf.it/pkg/scanner"
	"github.com/indigo-net/Brf.it/pkg/security"
	"github.com/indigo-net/Brf.it/pkg/tokenizer"
)
type Options struct {
	// Path is the target path to scan.
	Path string

	// Version is the brf.it version string.
	Version string

	// Format is the output format ("xml" or "md").
	Format string

	// Output is the output file path (empty = stdout).
	Output string

	// IgnoreFiles is the list of custom ignore file paths.
	IgnoreFiles []string

	// IncludeHidden determines whether to include hidden files.
	IncludeHidden bool

	// IncludeBody determines whether to include function/method bodies.
	IncludeBody bool

	// IncludeImports determines whether to include import/export statements.
	IncludeImports bool

	// DedupeImports deduplicates imports across files and shows them globally.
	// Requires IncludeImports to be true.
	DedupeImports bool

	// IncludeTree determines whether to include directory tree.
	IncludeTree bool

	// IncludePrivate determines whether to include private symbols.
	IncludePrivate bool

	// MaxFileSize is the maximum file size in bytes.
	MaxFileSize int64

	// MaxDocLength is the maximum length of documentation comments.
	// 0 means no limit (default).
	MaxDocLength int

	// NoSchema skips the schema section in XML output.
	NoSchema bool

	// SecurityCheck enables secret detection and redaction.
	SecurityCheck bool

	// IncludeCallGraph enables function call graph extraction.
	IncludeCallGraph bool
}
func DefaultOptions() *Options
type Result struct {
	// Content is the formatted output bytes.
	Content []byte

	// TotalSignatures is the total number of signatures.
	TotalSignatures int

	// TotalFiles is the number of processed files.
	TotalFiles int

	// TotalSize is the total size of processed files.
	TotalSize int64

	// TokenCount is the number of tokens in the output.
	// Returns 0 if token counting is disabled or tokenizer is not set.
	TokenCount int
}
type Packager struct {
	scanner    scanner.Scanner
	extractor  extractor.Extractor
	formatters map[string]formatter.Formatter
	tokenizer  tokenizer.Tokenizer
	warnings   io.Writer
}
func NewPackager(
	s scanner.Scanner,
	e extractor.Extractor,
	f map[string]formatter.Formatter,
) *Packager
func (p *Packager) SetTokenizer(t tokenizer.Tokenizer)
func (p *Packager) Package(ctx context.Context, opts *Options) (*Result, error)
treeStr string
globalImports []formatter.ImportCount
func NewDefaultPackager(scanOpts *scanner.ScanOptions) (*Packager, error)
func normalizeFormat(format string) string
func buildGlobalImports(files []formatter.FileData) []formatter.ImportCount
```

### /home/runner/work/Brf.it/Brf.it/internal/context/context_test.go

```go
import (
	"context"
	"fmt"
	"strings"
	"testing"
	"github.com/indigo-net/Brf.it/pkg/extractor"
	"github.com/indigo-net/Brf.it/pkg/formatter"
	"github.com/indigo-net/Brf.it/pkg/parser"
	"github.com/indigo-net/Brf.it/pkg/scanner"
	"github.com/indigo-net/Brf.it/pkg/tokenizer"
)
type mockScanner struct {
	result *scanner.ScanResult
	err    error
}
func (m *mockScanner) Scan(_ context.Context) (*scanner.ScanResult, error)
type mockExtractor struct {
	result *extractor.ExtractResult
	err    error
}
func (m *mockExtractor) Extract(_ context.Context, _ *scanner.ScanResult, _ *extractor.ExtractOptions) (*extractor.ExtractResult, error)
func TestPackagerPackage(t *testing.T)
func TestPackagerPackageMarkdown(t *testing.T)
func TestPackagerPackageMarkdownFull(t *testing.T)
func TestPackagerUnknownFormat(t *testing.T)
func TestPackagerSetTokenizer(t *testing.T)
func TestPackagerWithTiktokenTokenizer(t *testing.T)
func TestPackagerTokenizerConsistency(t *testing.T)
func TestBuildTree(t *testing.T)
func TestBuildTreeStructure(t *testing.T)
func TestPackagerNoStdImportsPassthrough(t *testing.T)
func TestDefaultOptions(t *testing.T)
func TestNormalizeFormat(t *testing.T)
func TestBuildGlobalImports(t *testing.T)
func TestBuildGlobalImportsSorting(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/internal/context/tree.go

```go
import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)
type treeNode struct {
	children map[string]*treeNode
	tokens   int  // token count for leaf nodes (files)
	isFile   bool // true if this node is a file (leaf)
}
func BuildTree(root string, paths []string) string
buf strings.Builder
func renderNode(buf *strings.Builder, n *treeNode, prefix string, isRoot bool)
newPrefix string
type FileTokenCount struct {
	Path   string
	Tokens int
}
func BuildTokenTree(root string, files []FileTokenCount) string
buf strings.Builder
func calcDirTokens(n *treeNode) int
sum int
func renderTokenNode(buf *strings.Builder, n *treeNode, prefix string, isRoot bool)
newPrefix string
func formatNumber(n int) string
buf strings.Builder
```

### /home/runner/work/Brf.it/Brf.it/internal/context/tree_test.go

```go
import (
	"strings"
	"testing"
)
func TestBuildTokenTree(t *testing.T)
func TestFormatNumber(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/internal/profiling/profiling.go

```go
import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)
type MemoryStats struct {
	// Alloc is bytes of allocated heap objects.
	Alloc uint64

	// TotalAlloc is cumulative bytes allocated for heap objects.
	TotalAlloc uint64

	// Sys is total bytes of memory obtained from the OS.
	Sys uint64

	// NumGC is the number of completed GC cycles.
	NumGC uint32

	// GoroutineCount is the number of goroutines.
	GoroutineCount int

	// HeapObjects is the number of allocated heap objects.
	HeapObjects uint64
}
func GetMemoryStats() MemoryStats
m runtime.MemStats
func FormatBytes(b uint64) string
unit = 1024
func WriteHeapProfile(filename string) error
func StartCPUProfile(filename string) (func(), error)
```

### /home/runner/work/Brf.it/Brf.it/internal/profiling/profiling_test.go

```go
import (
	"os"
	"testing"
)
func TestGetMemoryStats(t *testing.T)
func TestFormatBytes(t *testing.T)
func TestWriteHeapProfile(t *testing.T)
func TestWriteHeapProfileInvalidPath(t *testing.T)
func TestStartCPUProfile(t *testing.T)
func TestStartCPUProfileInvalidPath(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/extractor/example_test.go

```go
import (
	"fmt"
	"github.com/indigo-net/Brf.it/pkg/extractor"
)
func ExampleDefaultExtractOptions()
```

### /home/runner/work/Brf.it/Brf.it/pkg/extractor/extractor.go

```go
import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"github.com/indigo-net/Brf.it/pkg/parser"
	"github.com/indigo-net/Brf.it/pkg/scanner"
)
type ExtractedFile struct {
	// Path is the file path.
	Path string

	// Language is the detected language.
	Language string

	// Signatures is the list of extracted signatures.
	Signatures []parser.Signature

	// RawImports is the list of raw import/export statement text.
	RawImports []string

	// Calls is the list of function call references.
	Calls []parser.FunctionCall

	// Size is the file size in bytes.
	Size int64

	// Error is any error that occurred during extraction.
	Error error
}
type ExtractResult struct {
	// Files is the list of extracted files.
	Files []ExtractedFile

	// TotalSignatures is the total number of signatures extracted.
	TotalSignatures int

	// TotalSize is the total size of processed files.
	TotalSize int64

	// ErrorCount is the number of files that had errors.
	ErrorCount int
}
type ExtractOptions struct {
	// IncludePrivate whether to include non-exported/private signatures.
	IncludePrivate bool

	// IncludeBody whether to include function/method bodies.
	IncludeBody bool

	// IncludeImports whether to include import/export statements.
	IncludeImports bool

	// IncludeCalls whether to include function call references.
	IncludeCalls bool

	// Concurrency is the number of concurrent workers.
	// 0 = auto (runtime.NumCPU()), 1 = sequential.
	Concurrency int

	// MaxFileSize is the maximum file size in bytes for TOCTOU re-check.
	// If positive, file content size is verified after reading.
	MaxFileSize int64
}
type Extractor interface {
	// Extract extracts signatures from the given scan result.
	// The context controls cancellation and timeout for the extraction.
	Extract(ctx context.Context, scanResult *scanner.ScanResult, opts *ExtractOptions) (*ExtractResult, error)
}
type FileExtractor struct {
	registry *parser.Registry
}
func NewFileExtractor(registry *parser.Registry) *FileExtractor
func NewDefaultFileExtractor() *FileExtractor
func DefaultExtractOptions() *ExtractOptions
func (e *FileExtractor) Extract(ctx context.Context, scanResult *scanner.ScanResult, opts *ExtractOptions) (*ExtractResult, error)
wg sync.WaitGroup
cancelErr error
cancelOnce sync.Once
func (e *FileExtractor) extractSequential(ctx context.Context, files []scanner.FileEntry, opts *ExtractOptions) (*ExtractResult, error)
binarySniffSize = 512
func isBinaryContent(content []byte) bool
func (e *FileExtractor) extractFile(ctx context.Context, entry scanner.FileEntry, opts *ExtractOptions) ExtractedFile
err error
```

### /home/runner/work/Brf.it/Brf.it/pkg/extractor/extractor_test.go

```go
import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"github.com/indigo-net/Brf.it/pkg/parser"
	_ "github.com/indigo-net/Brf.it/pkg/parser/treesitter" // Register Tree-sitter parsers
	"github.com/indigo-net/Brf.it/pkg/scanner"
)
func TestFileExtractorImplementsExtractor(t *testing.T)
_ Extractor = (*FileExtractor)(nil)
func TestFileExtractorExtract(t *testing.T)
foundAdd bool
func TestFileExtractorTOCTOUGuard(t *testing.T)
func TestFileExtractorTOCTOUGuardDisabled(t *testing.T)
func TestExtractConcurrencySequential(t *testing.T)
func TestExtractConcurrencyDeterministicOrder(t *testing.T)
entries []scanner.FileEntry
func TestExtractConcurrencyEmptyFiles(t *testing.T)
func TestExtractNilScanResult(t *testing.T)
func TestExtractNegativeConcurrency(t *testing.T)
func TestDefaultExtractOptions(t *testing.T)
func TestExtractConcurrencyWithErrors(t *testing.T)
func TestFileExtractorUnsupportedLanguage(t *testing.T)
func TestExtractCanceledContext(t *testing.T)
func TestExtractDeadlineExceededContext(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/example_test.go

```go
import (
	"fmt"
	"github.com/indigo-net/Brf.it/pkg/formatter"
)
func ExampleNewXMLFormatter()
func ExampleNewMarkdownFormatter()
func ExampleNewJSONFormatter()
func ExampleXMLFormatter_Format()
```

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/formatter.go

```go
import (
	"github.com/indigo-net/Brf.it/pkg/parser"
)
type FileData struct {
	// Path is the file path.
	Path string

	// Language is the detected language.
	Language string

	// Signatures is the list of extracted signatures.
	Signatures []parser.Signature

	// RawImports is the list of raw import/export statement text.
	RawImports []string

	// Calls is the list of function call references.
	Calls []parser.FunctionCall

	// Error is any error that occurred during extraction.
	Error error
}
type PackageData struct {
	// RootPath is the root path being packaged.
	RootPath string

	// Version is the brf.it version string.
	Version string

	// Tree is the directory tree string.
	Tree string

	// Files is the list of file data.
	Files []FileData

	// TotalSignatures is the total number of signatures.
	TotalSignatures int

	// TotalSize is the total size of processed files.
	TotalSize int64

	// IncludeImports indicates whether imports should be rendered.
	IncludeImports bool

	// DedupeImports indicates whether imports should be deduplicated across files.
	// When true, imports are collected globally and shown in a separate section.
	DedupeImports bool

	// GlobalImports holds deduplicated imports with their usage counts.
	// Only populated when DedupeImports is true.
	GlobalImports []ImportCount

	// MaxDocLength is the maximum length of documentation comments.
	// 0 means no limit (default).
	MaxDocLength int

	// NoSchema indicates whether to omit the schema section in output.
	NoSchema bool

	// IncludeCallGraph indicates whether to include function call references.
	IncludeCallGraph bool
}
type ImportCount struct {
	// Import is the raw import statement text.
	Import string
	// Count is the number of files that use this import.
	Count int
}
type Formatter interface {
	// Format formats the package data and returns the output bytes.
	Format(data *PackageData) ([]byte, error)

	// Name returns the formatter name (e.g., "xml", "markdown").
	Name() string
}
```

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/formatter_bench_test.go

```go
import (
	"testing"
	"github.com/indigo-net/Brf.it/pkg/parser"
)
func createBenchmarkData(numFiles, numSigsPerFile int) *PackageData
func BenchmarkXMLFormatter(b *testing.B)
func BenchmarkMarkdownFormatter(b *testing.B)
func BenchmarkJSONFormatter(b *testing.B)
func BenchmarkXMLFormatterWithImports(b *testing.B)
func BenchmarkFormatterComparison(b *testing.B)
```

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/formatter_fuzz_test.go

```go
import (
	"testing"
	"github.com/indigo-net/Brf.it/pkg/parser"
)
func FuzzXMLFormatter(f *testing.F)
func FuzzMarkdownFormatter(f *testing.F)
func FuzzJSONFormatter(f *testing.F)
func FuzzFormatterWithLargeData(f *testing.F)
func FuzzFormatterWithImports(f *testing.F)
```

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/formatter_test.go

```go
import (
	"fmt"
	"strings"
	"testing"
	"github.com/indigo-net/Brf.it/pkg/parser"
)
func TestXMLFormatterImplementsFormatter(t *testing.T)
_ Formatter = (*XMLFormatter)(nil)
func TestMarkdownFormatterImplementsFormatter(t *testing.T)
_ Formatter = (*MarkdownFormatter)(nil)
func TestJSONFormatterImplementsFormatter(t *testing.T)
_ Formatter = (*JSONFormatter)(nil)
func TestXMLFormatterFormat(t *testing.T)
func TestXMLFormatterFormatWithError(t *testing.T)
func TestMarkdownFormatterFormat(t *testing.T)
func TestFormatterNames(t *testing.T)
func TestXMLFormatterEscapeXML(t *testing.T)
func TestMarkdownFormatterEscapeMarkdown(t *testing.T)
func TestXMLFormatterEmptyData(t *testing.T)
func TestMarkdownFormatterEmptyData(t *testing.T)
func TestMarkdownFormatterEmptyFile(t *testing.T)
func TestMarkdownFormatterEmptyFileWithImports(t *testing.T)
func TestXMLFormatterEmptyFile(t *testing.T)
func TestKindToTag(t *testing.T)
func TestXMLFormatterKindTags(t *testing.T)
func TestXMLFormatterVerbatimImports(t *testing.T)
func TestMarkdownFormatterVerbatimImports(t *testing.T)
func TestJSONFormatterFormat(t *testing.T)
func TestJSONFormatterKindNormalization(t *testing.T)
func TestJSONFormatterWithImports(t *testing.T)
func TestJSONFormatterWithError(t *testing.T)
func TestJSONFormatterName(t *testing.T)
func TestXMLFormatterWithNoSchema(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/helpers.go

```go
import "unicode/utf8"
func normalizeKind(kind string) string
func getEmptyComment(lang string) string
func truncateDoc(doc string, maxLen int) string
```

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/helpers_test.go

```go
import (
	"testing"
)
func TestGetEmptyComment(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/json.go

```go
import (
	"encoding/json"
)
type JSONFormatter struct{}
func NewJSONFormatter() *JSONFormatter
func (f *JSONFormatter) Name() string
type jsonOutput struct {
	Version       string            `json:"version,omitempty"`
	Path          string            `json:"path,omitempty"`
	Tree          string            `json:"tree,omitempty"`
	GlobalImports []jsonImportCount `json:"globalImports,omitempty"`
	Files         []jsonFile        `json:"files"`
}
type jsonImportCount struct {
	Import string `json:"import"`
	Count  int    `json:"count"`
}
type jsonFile struct {
	Path       string     `json:"path"`
	Language   string     `json:"language"`
	Signatures []jsonSig  `json:"signatures,omitempty"`
	Imports    []string   `json:"imports,omitempty"`
	Calls      []jsonCall `json:"calls,omitempty"`
	Error      string     `json:"error,omitempty"`
}
type jsonCall struct {
	Caller string `json:"caller,omitempty"`
	Callee string `json:"callee"`
	Line   int    `json:"line"`
}
type jsonSig struct {
	Kind     string `json:"kind"`
	Text     string `json:"text"`
	Doc      string `json:"doc,omitempty"`
	Line     int    `json:"line,omitempty"`
	Exported bool   `json:"exported,omitempty"`
}
func (f *JSONFormatter) Format(data *PackageData) ([]byte, error)
```

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/markdown.go

```go
import (
	"bytes"
	"strconv"
	"strings"
)
type MarkdownFormatter struct{}
func NewMarkdownFormatter() *MarkdownFormatter
func (f *MarkdownFormatter) Name() string
func (f *MarkdownFormatter) Format(data *PackageData) ([]byte, error)
buf bytes.Buffer
func escapeMarkdown(s string) string
```

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/xml.go

```go
import (
	"bytes"
	"strconv"
	"strings"
)
type XMLFormatter struct{}
func NewXMLFormatter() *XMLFormatter
func (f *XMLFormatter) Name() string
func (f *XMLFormatter) Format(data *PackageData) ([]byte, error)
buf bytes.Buffer
func escapeXML(s string) string
needsEscape bool
buf strings.Builder
func kindToTag(kind string) string
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/example_test.go

```go
import (
	"fmt"
	"github.com/indigo-net/Brf.it/pkg/parser"
)
func ExampleDetectLanguage()
func ExampleNewRegistry()
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/parser.go

```go
import (
	"path/filepath"
	"strings"
	"sync"
)
type Signature struct {
	// Name is the identifier name (e.g., "Scan", "FileScanner").
	Name string

	// Kind is the type of signature (e.g., "function", "method", "class", "interface").
	Kind string

	// Text is the full signature text including parameters and return type.
	Text string

	// Doc is the documentation comment (if any).
	Doc string

	// Line is the starting line number (1-indexed).
	Line int

	// EndLine is the ending line number (1-indexed).
	EndLine int

	// Language is the source language (e.g., "go", "typescript").
	Language string

	// Exported indicates whether the signature is exported/public.
	Exported bool
}
type Node struct {
	// Type is the node type (e.g., "function_declaration", "class_definition").
	Type string

	// StartRow is the starting row (0-indexed).
	StartRow int

	// EndRow is the ending row (0-indexed).
	EndRow int

	// StartColumn is the starting column.
	StartColumn int

	// EndColumn is the ending column.
	EndColumn int

	// Text is the source text of the node.
	Text string

	// Children are child nodes.
	Children []Node
}
type FunctionCall struct {
	// Caller is the name of the enclosing function (empty if top-level).
	Caller string

	// Callee is the called function/method name.
	Callee string

	// Line is the line number where the call occurs (1-indexed).
	Line int
}
type ParseResult struct {
	// FilePath is the path to the parsed file.
	FilePath string

	// Language is the detected language.
	Language string

	// Signatures is the list of extracted signatures.
	Signatures []Signature

	// RawImports is the list of raw import/export statement text.
	RawImports []string

	// Calls is the list of function call references.
	Calls []FunctionCall

	// AST is the root node of the parsed AST (optional).
	AST *Node

	// Error is any error that occurred during parsing.
	Error error
}
type Options struct {
	// Language forces a specific language (auto-detected if empty).
	Language string

	// IncludeAST whether to include the full AST in the result.
	IncludeAST bool

	// IncludePrivate whether to include non-exported/private signatures.
	IncludePrivate bool

	// IncludeBody whether to include function/method bodies in the signature text.
	// When false (default), only the signature line is extracted.
	// When true, the full declaration including the body is extracted.
	IncludeBody bool

	// IncludeImports whether to include import/export statements in the result.
	IncludeImports bool

	// IncludeCalls whether to include function call references in the result.
	IncludeCalls bool
}
type Parser interface {
	// Parse parses the given content and returns extracted signatures.
	// Content is passed as []byte to avoid unnecessary string conversion
	// from os.ReadFile output.
	Parse(content []byte, opts *Options) (*ParseResult, error)

	// Languages returns the list of supported languages.
	Languages() []string
}
type Registry struct {
	mu      sync.RWMutex
	parsers map[string]Parser
}
func NewRegistry() *Registry
defaultRegistry = NewRegistry()
func DefaultRegistry() *Registry
func (r *Registry) Register(lang string, parser Parser)
func (r *Registry) Get(lang string) (Parser, bool)
func (r *Registry) Languages() []string
func RegisterParser(lang string, parser Parser)
func GetParser(lang string) (Parser, bool)
languageMapping = map[string]string{
	".go":    "go",
	".ts":    "typescript",
	".tsx":   "typescript",
	".js":    "javascript",
	".jsx":   "javascript",
	".py":    "python",
	".java":  "java",
	".rs":    "rust",
	".rb":    "ruby",
	".php":   "php",
	".c":     "c",
	".cpp":   "cpp",
	".h":     "cpp",
	".hpp":   "cpp",
	".cs":    "csharp",
	".swift": "swift",
	".kt":    "kotlin",
	".kts":   "kotlin",
	".lua":   "lua",
	".sh":    "shell",
	".bash":  "shell",
	".zsh":   "shell",
	".scala": "scala",
	".sc":    "scala",
	".ex":    "elixir",
	".exs":   "elixir",
	".sql":   "sql",
	".yaml":  "yaml",
	".yml":   "yaml",
	".toml":  "toml",
}
func LanguageMapping() map[string]string
func DetectLanguage(path string) string
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/parser_test.go

```go
import (
	"testing"
)
func TestSignatureDefaults(t *testing.T)
func TestParseResultDefaults(t *testing.T)
func TestNodeKind(t *testing.T)
func TestParserInterface(t *testing.T)
_ Parser = (*MockParser)(nil)
type MockParser struct {
	signatures []Signature
	err        error
}
func (m *MockParser) Parse(content []byte, opts *Options) (*ParseResult, error)
func (m *MockParser) Languages() []string
func TestMockParser(t *testing.T)
func TestRegistry(t *testing.T)
func TestDefaultRegistry(t *testing.T)
func TestDetectLanguage(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/csharp/binding.go

```go
import "C"
import "unsafe"
func Language() unsafe.Pointer
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/csharp/scanner.c

```c
#include "tree_sitter/alloc.h"
#include "tree_sitter/array.h"
#include "tree_sitter/parser.h"
#include <wctype.h>
enum TokenType {
    OPT_SEMI,
    INTERPOLATION_REGULAR_START,
    INTERPOLATION_VERBATIM_START,
    INTERPOLATION_RAW_START,
    INTERPOLATION_START_QUOTE,
    INTERPOLATION_END_QUOTE,
    INTERPOLATION_OPEN_BRACE,
    INTERPOLATION_CLOSE_BRACE,
    INTERPOLATION_STRING_CONTENT,
    RAW_STRING_START,
    RAW_STRING_END,
    RAW_STRING_CONTENT,
}
typedef enum {
    REGULAR = 1 << 0,
    VERBATIM = 1 << 1,
    RAW = 1 << 2,
} StringType;
typedef struct {
    uint8_t dollar_count;
    uint8_t open_brace_count;
    uint8_t quote_count;
    StringType string_type;
} Interpolation;
static inline bool is_regular(Interpolation *interpolation)
static inline bool is_verbatim(Interpolation *interpolation)
static inline bool is_raw(Interpolation *interpolation)
typedef struct {
    uint8_t quote_count;
    Array(Interpolation) interpolation_stack;
} Scanner;
static inline void advance(TSLexer *lexer)
static inline void skip(TSLexer *lexer)
void *tree_sitter_c_sharp_external_scanner_create()
void tree_sitter_c_sharp_external_scanner_destroy(void *payload)
unsigned tree_sitter_c_sharp_external_scanner_serialize(void *payload, char *buffer)
void tree_sitter_c_sharp_external_scanner_deserialize(void *payload, const char *buffer, unsigned length)
bool tree_sitter_c_sharp_external_scanner_scan(void *payload, TSLexer *lexer, const bool *valid_symbols)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/csharp/tree_sitter/alloc.h

```cpp
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#define TREE_SITTER_ALLOC_H_
#define ts_malloc  ts_current_malloc
#define ts_calloc  ts_current_calloc
#define ts_realloc ts_current_realloc
#define ts_free    ts_current_free
#define ts_malloc  malloc
#define ts_calloc  calloc
#define ts_realloc realloc
#define ts_free    free
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/csharp/tree_sitter/array.h

```cpp
#include "./alloc.h"
#include <assert.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#define TREE_SITTER_ARRAY_H_
#define Array(T)       \
  struct {             \
    T *contents;       \
    uint32_t size;     \
    uint32_t capacity; \
  }
#define array_init(self) \
  ((self)->size = 0, (self)->capacity = 0, (self)->contents = NULL)
#define array_new() \
  { NULL, 0, 0 }
#define array_get(self, _index) \
  (assert((uint32_t)(_index) < (self)->size), &(self)->contents[_index])
#define array_front(self) array_get(self, 0)
#define array_back(self) array_get(self, (self)->size - 1)
#define array_clear(self) ((self)->size = 0)
#define array_reserve(self, new_capacity)        \
  ((self)->contents = _array__reserve(           \
    (void *)(self)->contents, &(self)->capacity, \
    array_elem_size(self), new_capacity)         \
  )
#define array_delete(self) _array__delete((self), (void *)(self)->contents, sizeof(*self))
#define array_push(self, element)                                 \
  do {                                                            \
    (self)->contents = _array__grow(                              \
      (void *)(self)->contents, (self)->size, &(self)->capacity,  \
      1, array_elem_size(self)                                    \
    );                                                            \
   (self)->contents[(self)->size++] = (element);                  \
  } while(0)
#define array_grow_by(self, count)                                               \
  do {                                                                           \
    if ((count) == 0) break;                                                     \
    (self)->contents = _array__grow(                                             \
      (self)->contents, (self)->size, &(self)->capacity,                         \
      count, array_elem_size(self)                                               \
    );                                                                           \
    memset((self)->contents + (self)->size, 0, (count) * array_elem_size(self)); \
    (self)->size += (count);                                                     \
  } while (0)
#define array_push_all(self, other) \
  array_extend((self), (other)->size, (other)->contents)
#define array_extend(self, count, other_contents)                 \
  (self)->contents = _array__splice(                              \
    (void*)(self)->contents, &(self)->size, &(self)->capacity,    \
    array_elem_size(self), (self)->size, 0, count, other_contents \
  )
#define array_splice(self, _index, old_count, new_count, new_contents) \
  (self)->contents = _array__splice(                                   \
    (void *)(self)->contents, &(self)->size, &(self)->capacity,        \
    array_elem_size(self), _index, old_count, new_count, new_contents  \
  )
#define array_insert(self, _index, element)                     \
  (self)->contents = _array__splice(                            \
    (void *)(self)->contents, &(self)->size, &(self)->capacity, \
    array_elem_size(self), _index, 0, 1, &(element)             \
  )
#define array_erase(self, _index) \
  _array__erase((void *)(self)->contents, &(self)->size, array_elem_size(self), _index)
#define array_pop(self) ((self)->contents[--(self)->size])
#define array_assign(self, other)                                   \
  (self)->contents = _array__assign(                                \
    (void *)(self)->contents, &(self)->size, &(self)->capacity,     \
    (const void *)(other)->contents, (other)->size, array_elem_size(self) \
  )
#define array_swap(self, other)                                     \
  do {                                                              \
    struct Swap swapped_contents = _array__swap(                    \
      (void *)(self)->contents, &(self)->size, &(self)->capacity,   \
      (void *)(other)->contents, &(other)->size, &(other)->capacity \
    );                                                              \
    (self)->contents = swapped_contents.self_contents;              \
    (other)->contents = swapped_contents.other_contents;            \
  } while (0)
#define array_elem_size(self) (sizeof *(self)->contents)
#define array_search_sorted_with(self, compare, needle, _index, _exists) \
  _array__search_sorted(self, 0, compare, , needle, _index, _exists)
#define array_search_sorted_by(self, field, needle, _index, _exists) \
  _array__search_sorted(self, 0, _compare_int, field, needle, _index, _exists)
#define array_insert_sorted_with(self, compare, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_with(self, compare, &(value), &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
#define array_insert_sorted_by(self, field, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_by(self, field, (value) field, &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
static inline void _array__delete(void *self, void *contents, size_t self_size)
static inline void _array__erase(void* self_contents, uint32_t *size,
                                size_t element_size, uint32_t index)
static inline void *_array__reserve(void *contents, uint32_t *capacity,
                                  size_t element_size, uint32_t new_capacity)
static inline void *_array__assign(void* self_contents, uint32_t *self_size, uint32_t *self_capacity,
                                 const void *other_contents, uint32_t other_size, size_t element_size)
struct Swap
struct Swap
struct Swap
static inline void *_array__grow(void *contents, uint32_t size, uint32_t *capacity,
                               uint32_t count, size_t element_size)
static inline void *_array__splice(void *self_contents, uint32_t *size, uint32_t *capacity,
                                 size_t element_size,
                                 uint32_t index, uint32_t old_count,
                                 uint32_t new_count, const void *elements)
#define _array__search_sorted(self, start, compare, suffix, needle, _index, _exists) \
  do { \
    *(_index) = start; \
    *(_exists) = false; \
    uint32_t size = (self)->size - *(_index); \
    if (size == 0) break; \
    int comparison; \
    while (size > 1) { \
      uint32_t half_size = size / 2; \
      uint32_t mid_index = *(_index) + half_size; \
      comparison = compare(&((self)->contents[mid_index] suffix), (needle)); \
      if (comparison <= 0) *(_index) = mid_index; \
      size -= half_size; \
    } \
    comparison = compare(&((self)->contents[*(_index)] suffix), (needle)); \
    if (comparison == 0) *(_exists) = true; \
    else if (comparison < 0) *(_index) += 1; \
  } while (0)
#define _compare_int(a, b) ((int)*(a) - (int)(b))
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/csharp/tree_sitter/parser.h

```cpp
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#define TREE_SITTER_PARSER_H_
#define ts_builtin_sym_error ((TSSymbol)-1)
#define ts_builtin_sym_end 0
#define TREE_SITTER_SERIALIZATION_BUFFER_SIZE 1024
typedef uint16_t TSStateId;
typedef uint16_t TSSymbol;
typedef uint16_t TSFieldId;
struct TSLanguage
struct TSLanguageMetadata
typedef struct {
  TSFieldId field_id;
  uint8_t child_index;
  bool inherited;
} TSFieldMapEntry;
typedef struct {
  uint16_t index;
  uint16_t length;
} TSMapSlice;
typedef struct {
  bool visible;
  bool named;
  bool supertype;
} TSSymbolMetadata;
struct TSLexer
struct TSLexer
typedef enum {
  TSParseActionTypeShift,
  TSParseActionTypeReduce,
  TSParseActionTypeAccept,
  TSParseActionTypeRecover,
} TSParseActionType;
typedef union {
  struct {
    uint8_t type;
    TSStateId state;
    bool extra;
    bool repetition;
  } shift;
  struct {
    uint8_t type;
    uint8_t child_count;
    TSSymbol symbol;
    int16_t dynamic_precedence;
    uint16_t production_id;
  } reduce;
  uint8_t type;
} TSParseAction;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
} TSLexMode;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
  uint16_t reserved_word_set_id;
} TSLexerMode;
typedef union {
  TSParseAction action;
  struct {
    uint8_t count;
    bool reusable;
  } entry;
} TSParseActionEntry;
typedef struct {
  int32_t start;
  int32_t end;
} TSCharacterRange;
struct TSLanguage
static inline bool set_contains(const TSCharacterRange *ranges, uint32_t len, int32_t lookahead)
#define UNUSED __pragma(warning(suppress : 4101))
#define UNUSED __attribute__((unused))
#define START_LEXER()           \
  bool result = false;          \
  bool skip = false;            \
  UNUSED                        \
  bool eof = false;             \
  int32_t lookahead;            \
  goto start;                   \
  next_state:                   \
  lexer->advance(lexer, skip);  \
  start:                        \
  skip = false;                 \
  lookahead = lexer->lookahead;
#define ADVANCE(state_value) \
  {                          \
    state = state_value;     \
    goto next_state;         \
  }
#define ADVANCE_MAP(...)                                              \
  {                                                                   \
    static const uint16_t map[] = { __VA_ARGS__ };                    \
    for (uint32_t i = 0; i < sizeof(map) / sizeof(map[0]); i += 2) {  \
      if (map[i] == lookahead) {                                      \
        state = map[i + 1];                                           \
        goto next_state;                                              \
      }                                                               \
    }                                                                 \
  }
#define SKIP(state_value) \
  {                       \
    skip = true;          \
    state = state_value;  \
    goto next_state;      \
  }
#define ACCEPT_TOKEN(symbol_value)     \
  result = true;                       \
  lexer->result_symbol = symbol_value; \
  lexer->mark_end(lexer);
#define END_STATE() return result;
#define SMALL_STATE(id) ((id) - LARGE_STATE_COUNT)
#define STATE(id) id
#define ACTIONS(id) id
#define SHIFT(state_value)            \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value)          \
    }                                 \
  }}
#define SHIFT_REPEAT(state_value)     \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value),         \
      .repetition = true              \
    }                                 \
  }}
#define SHIFT_EXTRA()                 \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .extra = true                   \
    }                                 \
  }}
#define REDUCE(symbol_name, children, precedence, prod_id) \
  {{                                                       \
    .reduce = {                                            \
      .type = TSParseActionTypeReduce,                     \
      .symbol = symbol_name,                               \
      .child_count = children,                             \
      .dynamic_precedence = precedence,                    \
      .production_id = prod_id                             \
    },                                                     \
  }}
#define RECOVER()                    \
  {{                                 \
    .type = TSParseActionTypeRecover \
  }}
#define ACCEPT_INPUT()              \
  {{                                \
    .type = TSParseActionTypeAccept \
  }}
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/elixir/binding.go

```go
import "C"
import "unsafe"
func Language() unsafe.Pointer
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/elixir/scanner.c

```c
#include "tree_sitter/parser.h"
enum TokenType {
  QUOTED_CONTENT_I_SINGLE,
  QUOTED_CONTENT_I_DOUBLE,
  QUOTED_CONTENT_I_HEREDOC_SINGLE,
  QUOTED_CONTENT_I_HEREDOC_DOUBLE,
  QUOTED_CONTENT_I_PARENTHESIS,
  QUOTED_CONTENT_I_CURLY,
  QUOTED_CONTENT_I_SQUARE,
  QUOTED_CONTENT_I_ANGLE,
  QUOTED_CONTENT_I_BAR,
  QUOTED_CONTENT_I_SLASH,
  QUOTED_CONTENT_SINGLE,
  QUOTED_CONTENT_DOUBLE,
  QUOTED_CONTENT_HEREDOC_SINGLE,
  QUOTED_CONTENT_HEREDOC_DOUBLE,
  QUOTED_CONTENT_PARENTHESIS,
  QUOTED_CONTENT_CURLY,
  QUOTED_CONTENT_SQUARE,
  QUOTED_CONTENT_ANGLE,
  QUOTED_CONTENT_BAR,
  QUOTED_CONTENT_SLASH,

  NEWLINE_BEFORE_DO,
  NEWLINE_BEFORE_BINARY_OPERATOR,
  NEWLINE_BEFORE_COMMENT,

  BEFORE_UNARY_OPERATOR,

  NOT_IN,

  QUOTED_ATOM_START
}
static inline void advance(TSLexer *lexer)
static inline void skip(TSLexer *lexer)
static inline bool is_whitespace(int32_t c)
static inline bool is_inline_whitespace(int32_t c)
static inline bool is_newline(int32_t c)
static inline bool is_digit(int32_t c)
static inline bool check_keyword_end(TSLexer *lexer)
static bool check_operator_end(TSLexer *lexer)
const uint8_t token_terminators_length =
    sizeof(token_terminators) / sizeof(char);
static inline bool is_token_end(int32_t c)
enum TokenType
typedef struct {
  const enum TokenType token_type;
  const bool supports_interpol;
  const int32_t end_delimiter;
  const uint8_t delimiter_length;
} QuotedContentInfo;
const uint8_t quoted_content_infos_length =
    sizeof(quoted_content_infos) / sizeof(QuotedContentInfo);
static inline int8_t find_quoted_token_info(const bool *valid_symbols)
static bool scan_quoted_content(TSLexer *lexer, const QuotedContentInfo *info)
static bool scan_newline(TSLexer *lexer, const bool *valid_symbols)
static bool scan(TSLexer *lexer, const bool *valid_symbols)
void *tree_sitter_elixir_external_scanner_create()
bool tree_sitter_elixir_external_scanner_scan(void *payload, TSLexer *lexer,
                                              const bool *valid_symbols)
unsigned tree_sitter_elixir_external_scanner_serialize(void *payload,
                                                       char *buffer)
void tree_sitter_elixir_external_scanner_deserialize(void *payload,
                                                     const char *buffer,
                                                     unsigned length)
void tree_sitter_elixir_external_scanner_destroy(void *payload)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/elixir/tree_sitter/alloc.h

```cpp
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#define TREE_SITTER_ALLOC_H_
#define ts_malloc  ts_current_malloc
#define ts_calloc  ts_current_calloc
#define ts_realloc ts_current_realloc
#define ts_free    ts_current_free
#define ts_malloc  malloc
#define ts_calloc  calloc
#define ts_realloc realloc
#define ts_free    free
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/elixir/tree_sitter/array.h

```cpp
#include "./alloc.h"
#include <assert.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#define TREE_SITTER_ARRAY_H_
#define Array(T)       \
  struct {             \
    T *contents;       \
    uint32_t size;     \
    uint32_t capacity; \
  }
#define array_init(self) \
  ((self)->size = 0, (self)->capacity = 0, (self)->contents = NULL)
#define array_new() \
  { NULL, 0, 0 }
#define array_get(self, _index) \
  (assert((uint32_t)(_index) < (self)->size), &(self)->contents[_index])
#define array_front(self) array_get(self, 0)
#define array_back(self) array_get(self, (self)->size - 1)
#define array_clear(self) ((self)->size = 0)
#define array_reserve(self, new_capacity) \
  _array__reserve((Array *)(self), array_elem_size(self), new_capacity)
#define array_delete(self) _array__delete((Array *)(self))
#define array_push(self, element)                            \
  (_array__grow((Array *)(self), 1, array_elem_size(self)), \
   (self)->contents[(self)->size++] = (element))
#define array_grow_by(self, count) \
  do { \
    if ((count) == 0) break; \
    _array__grow((Array *)(self), count, array_elem_size(self)); \
    memset((self)->contents + (self)->size, 0, (count) * array_elem_size(self)); \
    (self)->size += (count); \
  } while (0)
#define array_push_all(self, other)                                       \
  array_extend((self), (other)->size, (other)->contents)
#define array_extend(self, count, contents)                    \
  _array__splice(                                               \
    (Array *)(self), array_elem_size(self), (self)->size, \
    0, count,  contents                                        \
  )
#define array_splice(self, _index, old_count, new_count, new_contents)  \
  _array__splice(                                                       \
    (Array *)(self), array_elem_size(self), _index,                \
    old_count, new_count, new_contents                                 \
  )
#define array_insert(self, _index, element) \
  _array__splice((Array *)(self), array_elem_size(self), _index, 0, 1, &(element))
#define array_erase(self, _index) \
  _array__erase((Array *)(self), array_elem_size(self), _index)
#define array_pop(self) ((self)->contents[--(self)->size])
#define array_assign(self, other) \
  _array__assign((Array *)(self), (const Array *)(other), array_elem_size(self))
#define array_swap(self, other) \
  _array__swap((Array *)(self), (Array *)(other))
#define array_elem_size(self) (sizeof *(self)->contents)
#define array_search_sorted_with(self, compare, needle, _index, _exists) \
  _array__search_sorted(self, 0, compare, , needle, _index, _exists)
#define array_search_sorted_by(self, field, needle, _index, _exists) \
  _array__search_sorted(self, 0, _compare_int, field, needle, _index, _exists)
#define array_insert_sorted_with(self, compare, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_with(self, compare, &(value), &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
#define array_insert_sorted_by(self, field, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_by(self, field, (value) field, &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
static inline void _array__delete(Array *self)
static inline void _array__erase(Array *self, size_t element_size,
                                uint32_t index)
static inline void _array__reserve(Array *self, size_t element_size, uint32_t new_capacity)
static inline void _array__assign(Array *self, const Array *other, size_t element_size)
static inline void _array__swap(Array *self, Array *other)
static inline void _array__grow(Array *self, uint32_t count, size_t element_size)
static inline void _array__splice(Array *self, size_t element_size,
                                 uint32_t index, uint32_t old_count,
                                 uint32_t new_count, const void *elements)
#define _array__search_sorted(self, start, compare, suffix, needle, _index, _exists) \
  do { \
    *(_index) = start; \
    *(_exists) = false; \
    uint32_t size = (self)->size - *(_index); \
    if (size == 0) break; \
    int comparison; \
    while (size > 1) { \
      uint32_t half_size = size / 2; \
      uint32_t mid_index = *(_index) + half_size; \
      comparison = compare(&((self)->contents[mid_index] suffix), (needle)); \
      if (comparison <= 0) *(_index) = mid_index; \
      size -= half_size; \
    } \
    comparison = compare(&((self)->contents[*(_index)] suffix), (needle)); \
    if (comparison == 0) *(_exists) = true; \
    else if (comparison < 0) *(_index) += 1; \
  } while (0)
#define _compare_int(a, b) ((int)*(a) - (int)(b))
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/elixir/tree_sitter/parser.h

```cpp
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#define TREE_SITTER_PARSER_H_
#define ts_builtin_sym_error ((TSSymbol)-1)
#define ts_builtin_sym_end 0
#define TREE_SITTER_SERIALIZATION_BUFFER_SIZE 1024
typedef uint16_t TSStateId;
typedef uint16_t TSSymbol;
typedef uint16_t TSFieldId;
struct TSLanguage
typedef struct {
  TSFieldId field_id;
  uint8_t child_index;
  bool inherited;
} TSFieldMapEntry;
typedef struct {
  uint16_t index;
  uint16_t length;
} TSFieldMapSlice;
typedef struct {
  bool visible;
  bool named;
  bool supertype;
} TSSymbolMetadata;
struct TSLexer
struct TSLexer
typedef enum {
  TSParseActionTypeShift,
  TSParseActionTypeReduce,
  TSParseActionTypeAccept,
  TSParseActionTypeRecover,
} TSParseActionType;
typedef union {
  struct {
    uint8_t type;
    TSStateId state;
    bool extra;
    bool repetition;
  } shift;
  struct {
    uint8_t type;
    uint8_t child_count;
    TSSymbol symbol;
    int16_t dynamic_precedence;
    uint16_t production_id;
  } reduce;
  uint8_t type;
} TSParseAction;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
} TSLexMode;
typedef union {
  TSParseAction action;
  struct {
    uint8_t count;
    bool reusable;
  } entry;
} TSParseActionEntry;
typedef struct {
  int32_t start;
  int32_t end;
} TSCharacterRange;
struct TSLanguage
static inline bool set_contains(TSCharacterRange *ranges, uint32_t len, int32_t lookahead)
#define UNUSED __pragma(warning(suppress : 4101))
#define UNUSED __attribute__((unused))
#define START_LEXER()           \
  bool result = false;          \
  bool skip = false;            \
  UNUSED                        \
  bool eof = false;             \
  int32_t lookahead;            \
  goto start;                   \
  next_state:                   \
  lexer->advance(lexer, skip);  \
  start:                        \
  skip = false;                 \
  lookahead = lexer->lookahead;
#define ADVANCE(state_value) \
  {                          \
    state = state_value;     \
    goto next_state;         \
  }
#define ADVANCE_MAP(...)                                              \
  {                                                                   \
    static const uint16_t map[] = { __VA_ARGS__ };                    \
    for (uint32_t i = 0; i < sizeof(map) / sizeof(map[0]); i += 2) {  \
      if (map[i] == lookahead) {                                      \
        state = map[i + 1];                                           \
        goto next_state;                                              \
      }                                                               \
    }                                                                 \
  }
#define SKIP(state_value) \
  {                       \
    skip = true;          \
    state = state_value;  \
    goto next_state;      \
  }
#define ACCEPT_TOKEN(symbol_value)     \
  result = true;                       \
  lexer->result_symbol = symbol_value; \
  lexer->mark_end(lexer);
#define END_STATE() return result;
#define SMALL_STATE(id) ((id) - LARGE_STATE_COUNT)
#define STATE(id) id
#define ACTIONS(id) id
#define SHIFT(state_value)            \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value)          \
    }                                 \
  }}
#define SHIFT_REPEAT(state_value)     \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value),         \
      .repetition = true              \
    }                                 \
  }}
#define SHIFT_EXTRA()                 \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .extra = true                   \
    }                                 \
  }}
#define REDUCE(symbol_name, children, precedence, prod_id) \
  {{                                                       \
    .reduce = {                                            \
      .type = TSParseActionTypeReduce,                     \
      .symbol = symbol_name,                               \
      .child_count = children,                             \
      .dynamic_precedence = precedence,                    \
      .production_id = prod_id                             \
    },                                                     \
  }}
#define RECOVER()                    \
  {{                                 \
    .type = TSParseActionTypeRecover \
  }}
#define ACCEPT_INPUT()              \
  {{                                \
    .type = TSParseActionTypeAccept \
  }}
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/kotlin/binding.go

```go
import "C"
import "unsafe"
func Language() unsafe.Pointer
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/kotlin/scanner.c

```c
#include "tree_sitter/array.h"
#include "tree_sitter/parser.h"
#include <string.h>
#include <wctype.h>
enum TokenType {
  AUTOMATIC_SEMICOLON,
  IMPORT_LIST_DELIMITER,
  SAFE_NAV,
  MULTILINE_COMMENT,
  STRING_START,
  STRING_END,
  STRING_CONTENT,
  PRIMARY_CONSTRUCTOR_KEYWORD,
  IMPORT_DOT,
}
#define DELIMITER_LENGTH 3
typedef char Delimiter;
typedef Array(Delimiter) Stack;
static inline void stack_push(Stack *stack, char chr, bool triple)
static inline Delimiter stack_pop(Stack *stack)
static inline void skip(TSLexer *lexer)
static inline void advance(TSLexer *lexer)
static bool scan_string_start(TSLexer *lexer, Stack *stack)
static bool scan_string_content(TSLexer *lexer, Stack *stack)
static bool scan_multiline_comment(TSLexer *lexer)
static bool scan_whitespace_and_comments(TSLexer *lexer)
static bool is_word_char(int32_t c)
static bool scan_for_word(TSLexer *lexer, const char* word, unsigned len)
static bool check_word(TSLexer *lexer, const char *word, unsigned len)
static bool check_modifier_then_constructor(TSLexer *lexer)
static bool scan_automatic_semicolon(TSLexer *lexer, const bool *valid_symbols)
static bool scan_safe_nav(TSLexer *lexer)
static bool scan_line_sep(TSLexer *lexer)
static bool scan_import_list_delimiter(TSLexer *lexer)
static bool scan_import_dot(TSLexer *lexer)
bool tree_sitter_kotlin_external_scanner_scan(void *payload, TSLexer *lexer, const bool *valid_symbols)
void *tree_sitter_kotlin_external_scanner_create()
void tree_sitter_kotlin_external_scanner_destroy(void *payload)
unsigned tree_sitter_kotlin_external_scanner_serialize(void *payload, char *buffer)
void tree_sitter_kotlin_external_scanner_deserialize(void *payload, const char *buffer, unsigned length)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/kotlin/tree_sitter/alloc.h

```cpp
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#define TREE_SITTER_ALLOC_H_
#define ts_malloc  ts_current_malloc
#define ts_calloc  ts_current_calloc
#define ts_realloc ts_current_realloc
#define ts_free    ts_current_free
#define ts_malloc  malloc
#define ts_calloc  calloc
#define ts_realloc realloc
#define ts_free    free
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/kotlin/tree_sitter/array.h

```cpp
#include "./alloc.h"
#include <assert.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#define TREE_SITTER_ARRAY_H_
#define Array(T)       \
  struct {             \
    T *contents;       \
    uint32_t size;     \
    uint32_t capacity; \
  }
#define array_init(self) \
  ((self)->size = 0, (self)->capacity = 0, (self)->contents = NULL)
#define array_new() \
  { NULL, 0, 0 }
#define array_get(self, _index) \
  (assert((uint32_t)(_index) < (self)->size), &(self)->contents[_index])
#define array_front(self) array_get(self, 0)
#define array_back(self) array_get(self, (self)->size - 1)
#define array_clear(self) ((self)->size = 0)
#define array_reserve(self, new_capacity) \
  _array__reserve((Array *)(self), array_elem_size(self), new_capacity)
#define array_delete(self) _array__delete((Array *)(self))
#define array_push(self, element)                            \
  (_array__grow((Array *)(self), 1, array_elem_size(self)), \
   (self)->contents[(self)->size++] = (element))
#define array_grow_by(self, count) \
  do { \
    if ((count) == 0) break; \
    _array__grow((Array *)(self), count, array_elem_size(self)); \
    memset((self)->contents + (self)->size, 0, (count) * array_elem_size(self)); \
    (self)->size += (count); \
  } while (0)
#define array_push_all(self, other)                                       \
  array_extend((self), (other)->size, (other)->contents)
#define array_extend(self, count, contents)                    \
  _array__splice(                                               \
    (Array *)(self), array_elem_size(self), (self)->size, \
    0, count,  contents                                        \
  )
#define array_splice(self, _index, old_count, new_count, new_contents)  \
  _array__splice(                                                       \
    (Array *)(self), array_elem_size(self), _index,                \
    old_count, new_count, new_contents                                 \
  )
#define array_insert(self, _index, element) \
  _array__splice((Array *)(self), array_elem_size(self), _index, 0, 1, &(element))
#define array_erase(self, _index) \
  _array__erase((Array *)(self), array_elem_size(self), _index)
#define array_pop(self) ((self)->contents[--(self)->size])
#define array_assign(self, other) \
  _array__assign((Array *)(self), (const Array *)(other), array_elem_size(self))
#define array_swap(self, other) \
  _array__swap((Array *)(self), (Array *)(other))
#define array_elem_size(self) (sizeof *(self)->contents)
#define array_search_sorted_with(self, compare, needle, _index, _exists) \
  _array__search_sorted(self, 0, compare, , needle, _index, _exists)
#define array_search_sorted_by(self, field, needle, _index, _exists) \
  _array__search_sorted(self, 0, _compare_int, field, needle, _index, _exists)
#define array_insert_sorted_with(self, compare, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_with(self, compare, &(value), &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
#define array_insert_sorted_by(self, field, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_by(self, field, (value) field, &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
static inline void _array__delete(Array *self)
static inline void _array__erase(Array *self, size_t element_size,
                                uint32_t index)
static inline void _array__reserve(Array *self, size_t element_size, uint32_t new_capacity)
static inline void _array__assign(Array *self, const Array *other, size_t element_size)
static inline void _array__swap(Array *self, Array *other)
static inline void _array__grow(Array *self, uint32_t count, size_t element_size)
static inline void _array__splice(Array *self, size_t element_size,
                                 uint32_t index, uint32_t old_count,
                                 uint32_t new_count, const void *elements)
#define _array__search_sorted(self, start, compare, suffix, needle, _index, _exists) \
  do { \
    *(_index) = start; \
    *(_exists) = false; \
    uint32_t size = (self)->size - *(_index); \
    if (size == 0) break; \
    int comparison; \
    while (size > 1) { \
      uint32_t half_size = size / 2; \
      uint32_t mid_index = *(_index) + half_size; \
      comparison = compare(&((self)->contents[mid_index] suffix), (needle)); \
      if (comparison <= 0) *(_index) = mid_index; \
      size -= half_size; \
    } \
    comparison = compare(&((self)->contents[*(_index)] suffix), (needle)); \
    if (comparison == 0) *(_exists) = true; \
    else if (comparison < 0) *(_index) += 1; \
  } while (0)
#define _compare_int(a, b) ((int)*(a) - (int)(b))
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/kotlin/tree_sitter/parser.h

```cpp
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#define TREE_SITTER_PARSER_H_
#define ts_builtin_sym_error ((TSSymbol)-1)
#define ts_builtin_sym_end 0
#define TREE_SITTER_SERIALIZATION_BUFFER_SIZE 1024
typedef uint16_t TSStateId;
typedef uint16_t TSSymbol;
typedef uint16_t TSFieldId;
struct TSLanguage
typedef struct {
  TSFieldId field_id;
  uint8_t child_index;
  bool inherited;
} TSFieldMapEntry;
typedef struct {
  uint16_t index;
  uint16_t length;
} TSFieldMapSlice;
typedef struct {
  bool visible;
  bool named;
  bool supertype;
} TSSymbolMetadata;
struct TSLexer
struct TSLexer
typedef enum {
  TSParseActionTypeShift,
  TSParseActionTypeReduce,
  TSParseActionTypeAccept,
  TSParseActionTypeRecover,
} TSParseActionType;
typedef union {
  struct {
    uint8_t type;
    TSStateId state;
    bool extra;
    bool repetition;
  } shift;
  struct {
    uint8_t type;
    uint8_t child_count;
    TSSymbol symbol;
    int16_t dynamic_precedence;
    uint16_t production_id;
  } reduce;
  uint8_t type;
} TSParseAction;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
} TSLexMode;
typedef union {
  TSParseAction action;
  struct {
    uint8_t count;
    bool reusable;
  } entry;
} TSParseActionEntry;
typedef struct {
  int32_t start;
  int32_t end;
} TSCharacterRange;
struct TSLanguage
static inline bool set_contains(TSCharacterRange *ranges, uint32_t len, int32_t lookahead)
#define UNUSED __pragma(warning(suppress : 4101))
#define UNUSED __attribute__((unused))
#define START_LEXER()           \
  bool result = false;          \
  bool skip = false;            \
  UNUSED                        \
  bool eof = false;             \
  int32_t lookahead;            \
  goto start;                   \
  next_state:                   \
  lexer->advance(lexer, skip);  \
  start:                        \
  skip = false;                 \
  lookahead = lexer->lookahead;
#define ADVANCE(state_value) \
  {                          \
    state = state_value;     \
    goto next_state;         \
  }
#define ADVANCE_MAP(...)                                              \
  {                                                                   \
    static const uint16_t map[] = { __VA_ARGS__ };                    \
    for (uint32_t i = 0; i < sizeof(map) / sizeof(map[0]); i += 2) {  \
      if (map[i] == lookahead) {                                      \
        state = map[i + 1];                                           \
        goto next_state;                                              \
      }                                                               \
    }                                                                 \
  }
#define SKIP(state_value) \
  {                       \
    skip = true;          \
    state = state_value;  \
    goto next_state;      \
  }
#define ACCEPT_TOKEN(symbol_value)     \
  result = true;                       \
  lexer->result_symbol = symbol_value; \
  lexer->mark_end(lexer);
#define END_STATE() return result;
#define SMALL_STATE(id) ((id) - LARGE_STATE_COUNT)
#define STATE(id) id
#define ACTIONS(id) id
#define SHIFT(state_value)            \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value)          \
    }                                 \
  }}
#define SHIFT_REPEAT(state_value)     \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value),         \
      .repetition = true              \
    }                                 \
  }}
#define SHIFT_EXTRA()                 \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .extra = true                   \
    }                                 \
  }}
#define REDUCE(symbol_name, children, precedence, prod_id) \
  {{                                                       \
    .reduce = {                                            \
      .type = TSParseActionTypeReduce,                     \
      .symbol = symbol_name,                               \
      .child_count = children,                             \
      .dynamic_precedence = precedence,                    \
      .production_id = prod_id                             \
    },                                                     \
  }}
#define RECOVER()                    \
  {{                                 \
    .type = TSParseActionTypeRecover \
  }}
#define ACCEPT_INPUT()              \
  {{                                \
    .type = TSParseActionTypeAccept \
  }}
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/lua/binding.go

```go
import "C"
import "unsafe"
func Language() unsafe.Pointer
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/lua/parser.c

```c
#include "tree_sitter/parser.h"
#define LANGUAGE_VERSION 15
#define STATE_COUNT 282
#define LARGE_STATE_COUNT 2
#define SYMBOL_COUNT 143
#define ALIAS_COUNT 0
#define TOKEN_COUNT 73
#define EXTERNAL_TOKEN_COUNT 6
#define FIELD_COUNT 23
#define MAX_ALIAS_SEQUENCE_LENGTH 7
#define MAX_RESERVED_WORD_SET_SIZE 0
#define PRODUCTION_ID_COUNT 70
#define SUPERTYPE_COUNT 4
enum ts_symbol_identifiers {
  sym_identifier = 1,
  sym_hash_bang_line = 2,
  anon_sym_return = 3,
  anon_sym_SEMI = 4,
  anon_sym_EQ = 5,
  anon_sym_COMMA = 6,
  anon_sym_COLON_COLON = 7,
  sym_break_statement = 8,
  anon_sym_goto = 9,
  anon_sym_do = 10,
  anon_sym_end = 11,
  anon_sym_while = 12,
  anon_sym_repeat = 13,
  anon_sym_until = 14,
  anon_sym_if = 15,
  anon_sym_then = 16,
  anon_sym_elseif = 17,
  anon_sym_else = 18,
  anon_sym_for = 19,
  anon_sym_in = 20,
  anon_sym_function = 21,
  anon_sym_local = 22,
  anon_sym_global = 23,
  anon_sym_DOT = 24,
  anon_sym_COLON = 25,
  anon_sym_STAR = 26,
  anon_sym_LT = 27,
  anon_sym_GT = 28,
  sym_nil = 29,
  sym_false = 30,
  sym_true = 31,
  sym_number = 32,
  anon_sym_DQUOTE = 33,
  anon_sym_SQUOTE = 34,
  aux_sym__doublequote_string_content_token1 = 35,
  aux_sym__singlequote_string_content_token1 = 36,
  sym_escape_sequence = 37,
  sym_vararg_expression = 38,
  anon_sym_LPAREN = 39,
  anon_sym_RPAREN = 40,
  anon_sym_LBRACK = 41,
  anon_sym_RBRACK = 42,
  anon_sym_LBRACE = 43,
  anon_sym_RBRACE = 44,
  anon_sym_or = 45,
  anon_sym_and = 46,
  anon_sym_LT_EQ = 47,
  anon_sym_EQ_EQ = 48,
  anon_sym_TILDE_EQ = 49,
  anon_sym_GT_EQ = 50,
  anon_sym_PIPE = 51,
  anon_sym_TILDE = 52,
  anon_sym_AMP = 53,
  anon_sym_LT_LT = 54,
  anon_sym_GT_GT = 55,
  anon_sym_PLUS = 56,
  anon_sym_DASH = 57,
  anon_sym_SLASH = 58,
  anon_sym_SLASH_SLASH = 59,
  anon_sym_PERCENT = 60,
  anon_sym_DOT_DOT = 61,
  anon_sym_CARET = 62,
  anon_sym_not = 63,
  anon_sym_POUND = 64,
  anon_sym_DASH_DASH = 65,
  aux_sym_comment_token1 = 66,
  sym__block_comment_start = 67,
  sym__block_comment_content = 68,
  sym__block_comment_end = 69,
  sym__block_string_start = 70,
  sym__block_string_content = 71,
  sym__block_string_end = 72,
  sym_chunk = 73,
  sym__block = 74,
  sym_statement = 75,
  sym_return_statement = 76,
  sym_empty_statement = 77,
  sym_assignment_statement = 78,
  sym__variable_assignment_varlist = 79,
  sym__variable_assignment_explist = 80,
  sym_label_statement = 81,
  sym_goto_statement = 82,
  sym_do_statement = 83,
  sym_while_statement = 84,
  sym_repeat_statement = 85,
  sym_if_statement = 86,
  sym_elseif_statement = 87,
  sym_else_statement = 88,
  sym_for_statement = 89,
  sym_for_generic_clause = 90,
  sym_for_numeric_clause = 91,
  sym__name_list = 92,
  sym_declaration = 93,
  sym_function_declaration = 94,
  sym__local_function_declaration = 95,
  sym__global_function_declaration = 96,
  sym__function_name = 97,
  sym__function_name_prefix_expression = 98,
  sym__function_name_dot_index_expression = 99,
  sym__function_name_method_index_expression = 100,
  sym_variable_declaration = 101,
  sym__global_variable_declaration = 102,
  sym__variable_assignment = 103,
  sym__att_name_list = 104,
  sym__global_implicit_variable_declaration = 105,
  sym__attrib = 106,
  sym__expression_list = 107,
  sym_expression = 108,
  sym_string = 109,
  sym__quote_string = 110,
  aux_sym__doublequote_string_content = 111,
  aux_sym__singlequote_string_content = 112,
  sym__block_string = 113,
  sym_function_definition = 114,
  sym__function_body = 115,
  sym_parameters = 116,
  sym__parameter_list = 117,
  sym__vararg_parameter = 118,
  sym__prefix_expression = 119,
  sym_variable = 120,
  sym_bracket_index_expression = 121,
  sym_dot_index_expression = 122,
  sym_function_call = 123,
  sym_method_index_expression = 124,
  sym_arguments = 125,
  sym_parenthesized_expression = 126,
  sym_table_constructor = 127,
  sym__field_list = 128,
  sym__field_sep = 129,
  sym_field = 130,
  sym_binary_expression = 131,
  sym_unary_expression = 132,
  sym_comment = 133,
  sym__contextual_keyword = 134,
  aux_sym_chunk_repeat1 = 135,
  aux_sym__variable_assignment_varlist_repeat1 = 136,
  aux_sym__variable_assignment_explist_repeat1 = 137,
  aux_sym_if_statement_repeat1 = 138,
  aux_sym__name_list_repeat1 = 139,
  aux_sym__att_name_list_repeat1 = 140,
  aux_sym__expression_list_repeat1 = 141,
  aux_sym__field_list_repeat1 = 142,
}
enum ts_field_identifiers {
  field_alternative = 1,
  field_arguments = 2,
  field_attribute = 3,
  field_body = 4,
  field_clause = 5,
  field_condition = 6,
  field_consequence = 7,
  field_content = 8,
  field_end = 9,
  field_field = 10,
  field_global_declaration = 11,
  field_left = 12,
  field_local_declaration = 13,
  field_method = 14,
  field_name = 15,
  field_operand = 16,
  field_operator = 17,
  field_parameters = 18,
  field_right = 19,
  field_start = 20,
  field_step = 21,
  field_table = 22,
  field_value = 23,
}
static bool ts_lex(TSLexer *lexer, TSStateId state)
static bool ts_lex_keywords(TSLexer *lexer, TSStateId state)
enum ts_external_scanner_symbol_identifiers {
  ts_external_token__block_comment_start = 0,
  ts_external_token__block_comment_content = 1,
  ts_external_token__block_comment_end = 2,
  ts_external_token__block_string_start = 3,
  ts_external_token__block_string_content = 4,
  ts_external_token__block_string_end = 5,
}
void *tree_sitter_lua_external_scanner_create(void);
void tree_sitter_lua_external_scanner_destroy(void *);
bool tree_sitter_lua_external_scanner_scan(void *, TSLexer *, const bool *);
unsigned tree_sitter_lua_external_scanner_serialize(void *, char *);
void tree_sitter_lua_external_scanner_deserialize(void *, const char *, unsigned);
#define TS_PUBLIC
#define TS_PUBLIC __declspec(dllexport)
#define TS_PUBLIC __attribute__((visibility("default")))
TS_PUBLIC const TSLanguage *tree_sitter_lua(void)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/lua/scanner.c

```c
#include <stdio.h>
#include "tree_sitter/alloc.h"
#include "tree_sitter/parser.h"
#include <wctype.h>
enum TokenType {
  BLOCK_COMMENT_START,
  BLOCK_COMMENT_CONTENT,
  BLOCK_COMMENT_END,

  BLOCK_STRING_START,
  BLOCK_STRING_CONTENT,
  BLOCK_STRING_END,
}
static inline void consume(TSLexer *lexer)
static inline void skip(TSLexer *lexer)
static inline bool consume_char(char c, TSLexer *lexer)
static inline uint8_t consume_and_count_char(char c, TSLexer *lexer)
static inline void skip_whitespaces(TSLexer *lexer)
typedef struct {
  char ending_char;
  uint8_t level_count;
} Scanner;
static inline void reset_state(Scanner *scanner)
void *tree_sitter_lua_external_scanner_create()
void tree_sitter_lua_external_scanner_destroy(void *payload)
unsigned tree_sitter_lua_external_scanner_serialize(void *payload, char *buffer)
void tree_sitter_lua_external_scanner_deserialize(void *payload, const char *buffer, unsigned length)
static bool scan_block_start(Scanner *scanner, TSLexer *lexer)
static bool scan_block_end(Scanner *scanner, TSLexer *lexer)
static bool scan_block_content(Scanner *scanner, TSLexer *lexer)
static bool scan_comment_start(Scanner *scanner, TSLexer *lexer)
static bool scan_comment_content(Scanner *scanner, TSLexer *lexer)
bool tree_sitter_lua_external_scanner_scan(void *payload, TSLexer *lexer, const bool *valid_symbols)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/lua/tree_sitter/alloc.h

```cpp
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#define TREE_SITTER_ALLOC_H_
#define ts_malloc  ts_current_malloc
#define ts_calloc  ts_current_calloc
#define ts_realloc ts_current_realloc
#define ts_free    ts_current_free
#define ts_malloc  malloc
#define ts_calloc  calloc
#define ts_realloc realloc
#define ts_free    free
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/lua/tree_sitter/array.h

```cpp
#include "./alloc.h"
#include <assert.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#define TREE_SITTER_ARRAY_H_
#define Array(T)       \
  struct {             \
    T *contents;       \
    uint32_t size;     \
    uint32_t capacity; \
  }
#define array_init(self) \
  ((self)->size = 0, (self)->capacity = 0, (self)->contents = NULL)
#define array_new() \
  { NULL, 0, 0 }
#define array_get(self, _index) \
  (assert((uint32_t)(_index) < (self)->size), &(self)->contents[_index])
#define array_front(self) array_get(self, 0)
#define array_back(self) array_get(self, (self)->size - 1)
#define array_clear(self) ((self)->size = 0)
#define array_reserve(self, new_capacity)        \
  ((self)->contents = _array__reserve(           \
    (void *)(self)->contents, &(self)->capacity, \
    array_elem_size(self), new_capacity)         \
  )
#define array_delete(self)                           \
  do {                                               \
    if ((self)->contents) ts_free((self)->contents); \
    (self)->contents = NULL;                         \
    (self)->size = 0;                                \
    (self)->capacity = 0;                            \
  } while (0)
#define array_push(self, element)                                 \
  do {                                                            \
    (self)->contents = _array__grow(                              \
      (void *)(self)->contents, (self)->size, &(self)->capacity,  \
      1, array_elem_size(self)                                    \
    );                                                            \
   (self)->contents[(self)->size++] = (element);                  \
  } while(0)
#define array_grow_by(self, count)                                               \
  do {                                                                           \
    if ((count) == 0) break;                                                     \
    (self)->contents = _array__grow(                                             \
      (self)->contents, (self)->size, &(self)->capacity,                         \
      count, array_elem_size(self)                                               \
    );                                                                           \
    memset((self)->contents + (self)->size, 0, (count) * array_elem_size(self)); \
    (self)->size += (count);                                                     \
  } while (0)
#define array_push_all(self, other) \
  array_extend((self), (other)->size, (other)->contents)
#define array_extend(self, count, other_contents)                 \
  (self)->contents = _array__splice(                              \
    (void*)(self)->contents, &(self)->size, &(self)->capacity,    \
    array_elem_size(self), (self)->size, 0, count, other_contents \
  )
#define array_splice(self, _index, old_count, new_count, new_contents) \
  (self)->contents = _array__splice(                                   \
    (void *)(self)->contents, &(self)->size, &(self)->capacity,        \
    array_elem_size(self), _index, old_count, new_count, new_contents  \
  )
#define array_insert(self, _index, element)                     \
  (self)->contents = _array__splice(                            \
    (void *)(self)->contents, &(self)->size, &(self)->capacity, \
    array_elem_size(self), _index, 0, 1, &(element)             \
  )
#define array_erase(self, _index) \
  _array__erase((void *)(self)->contents, &(self)->size, array_elem_size(self), _index)
#define array_pop(self) ((self)->contents[--(self)->size])
#define array_assign(self, other)                                   \
  (self)->contents = _array__assign(                                \
    (void *)(self)->contents, &(self)->size, &(self)->capacity,     \
    (const void *)(other)->contents, (other)->size, array_elem_size(self) \
  )
#define array_swap(self, other)                                     \
  do {                                                              \
    void *_array_swap_tmp = (void *)(self)->contents;               \
    (self)->contents = (other)->contents;                           \
    (other)->contents = _array_swap_tmp;                            \
    _array__swap(&(self)->size, &(self)->capacity,                  \
                 &(other)->size, &(other)->capacity);               \
  } while (0)
#define array_elem_size(self) (sizeof *(self)->contents)
#define array_search_sorted_with(self, compare, needle, _index, _exists) \
  _array__search_sorted(self, 0, compare, , needle, _index, _exists)
#define array_search_sorted_by(self, field, needle, _index, _exists) \
  _array__search_sorted(self, 0, _compare_int, field, needle, _index, _exists)
#define array_insert_sorted_with(self, compare, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_with(self, compare, &(value), &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
#define array_insert_sorted_by(self, field, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_by(self, field, (value) field, &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
static inline void _array__erase(void* self_contents, uint32_t *size,
                                size_t element_size, uint32_t index)
static inline void *_array__reserve(void *contents, uint32_t *capacity,
                                  size_t element_size, uint32_t new_capacity)
static inline void *_array__assign(void* self_contents, uint32_t *self_size, uint32_t *self_capacity,
                                 const void *other_contents, uint32_t other_size, size_t element_size)
static inline void _array__swap(uint32_t *self_size, uint32_t *self_capacity,
                               uint32_t *other_size, uint32_t *other_capacity)
static inline void *_array__grow(void *contents, uint32_t size, uint32_t *capacity,
                               uint32_t count, size_t element_size)
static inline void *_array__splice(void *self_contents, uint32_t *size, uint32_t *capacity,
                                 size_t element_size,
                                 uint32_t index, uint32_t old_count,
                                 uint32_t new_count, const void *elements)
#define _array__search_sorted(self, start, compare, suffix, needle, _index, _exists) \
  do { \
    *(_index) = start; \
    *(_exists) = false; \
    uint32_t size = (self)->size - *(_index); \
    if (size == 0) break; \
    int comparison; \
    while (size > 1) { \
      uint32_t half_size = size / 2; \
      uint32_t mid_index = *(_index) + half_size; \
      comparison = compare(&((self)->contents[mid_index] suffix), (needle)); \
      if (comparison <= 0) *(_index) = mid_index; \
      size -= half_size; \
    } \
    comparison = compare(&((self)->contents[*(_index)] suffix), (needle)); \
    if (comparison == 0) *(_exists) = true; \
    else if (comparison < 0) *(_index) += 1; \
  } while (0)
#define _compare_int(a, b) ((int)*(a) - (int)(b))
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/lua/tree_sitter/parser.h

```cpp
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#define TREE_SITTER_PARSER_H_
#define ts_builtin_sym_error ((TSSymbol)-1)
#define ts_builtin_sym_end 0
#define TREE_SITTER_SERIALIZATION_BUFFER_SIZE 1024
typedef uint16_t TSStateId;
typedef uint16_t TSSymbol;
typedef uint16_t TSFieldId;
struct TSLanguage
struct TSLanguageMetadata
typedef struct {
  TSFieldId field_id;
  uint8_t child_index;
  bool inherited;
} TSFieldMapEntry;
typedef struct {
  uint16_t index;
  uint16_t length;
} TSMapSlice;
typedef struct {
  bool visible;
  bool named;
  bool supertype;
} TSSymbolMetadata;
struct TSLexer
struct TSLexer
typedef enum {
  TSParseActionTypeShift,
  TSParseActionTypeReduce,
  TSParseActionTypeAccept,
  TSParseActionTypeRecover,
} TSParseActionType;
typedef union {
  struct {
    uint8_t type;
    TSStateId state;
    bool extra;
    bool repetition;
  } shift;
  struct {
    uint8_t type;
    uint8_t child_count;
    TSSymbol symbol;
    int16_t dynamic_precedence;
    uint16_t production_id;
  } reduce;
  uint8_t type;
} TSParseAction;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
} TSLexMode;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
  uint16_t reserved_word_set_id;
} TSLexerMode;
typedef union {
  TSParseAction action;
  struct {
    uint8_t count;
    bool reusable;
  } entry;
} TSParseActionEntry;
typedef struct {
  int32_t start;
  int32_t end;
} TSCharacterRange;
struct TSLanguage
static inline bool set_contains(const TSCharacterRange *ranges, uint32_t len, int32_t lookahead)
#define UNUSED __pragma(warning(suppress : 4101))
#define UNUSED __attribute__((unused))
#define START_LEXER()           \
  bool result = false;          \
  bool skip = false;            \
  UNUSED                        \
  bool eof = false;             \
  int32_t lookahead;            \
  goto start;                   \
  next_state:                   \
  lexer->advance(lexer, skip);  \
  start:                        \
  skip = false;                 \
  lookahead = lexer->lookahead;
#define ADVANCE(state_value) \
  {                          \
    state = state_value;     \
    goto next_state;         \
  }
#define ADVANCE_MAP(...)                                              \
  {                                                                   \
    static const uint16_t map[] = { __VA_ARGS__ };                    \
    for (uint32_t i = 0; i < sizeof(map) / sizeof(map[0]); i += 2) {  \
      if (map[i] == lookahead) {                                      \
        state = map[i + 1];                                           \
        goto next_state;                                              \
      }                                                               \
    }                                                                 \
  }
#define SKIP(state_value) \
  {                       \
    skip = true;          \
    state = state_value;  \
    goto next_state;      \
  }
#define ACCEPT_TOKEN(symbol_value)     \
  result = true;                       \
  lexer->result_symbol = symbol_value; \
  lexer->mark_end(lexer);
#define END_STATE() return result;
#define SMALL_STATE(id) ((id) - LARGE_STATE_COUNT)
#define STATE(id) id
#define ACTIONS(id) id
#define SHIFT(state_value)            \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value)          \
    }                                 \
  }}
#define SHIFT_REPEAT(state_value)     \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value),         \
      .repetition = true              \
    }                                 \
  }}
#define SHIFT_EXTRA()                 \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .extra = true                   \
    }                                 \
  }}
#define REDUCE(symbol_name, children, precedence, prod_id) \
  {{                                                       \
    .reduce = {                                            \
      .type = TSParseActionTypeReduce,                     \
      .symbol = symbol_name,                               \
      .child_count = children,                             \
      .dynamic_precedence = precedence,                    \
      .production_id = prod_id                             \
    },                                                     \
  }}
#define RECOVER()                    \
  {{                                 \
    .type = TSParseActionTypeRecover \
  }}
#define ACCEPT_INPUT()              \
  {{                                \
    .type = TSParseActionTypeAccept \
  }}
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/scala/binding.go

```go
import "C"
import "unsafe"
func Language() unsafe.Pointer
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/scala/scanner.c

```c
#include "tree_sitter/alloc.h"
#include "tree_sitter/array.h"
#include "tree_sitter/parser.h"
#include <wctype.h>
#define LOG(...) fprintf(stderr, __VA_ARGS__)
#define LOG(...)
enum TokenType {
  AUTOMATIC_SEMICOLON,
  INDENT,
  OUTDENT,
  SIMPLE_STRING_START,
  SIMPLE_STRING_MIDDLE,
  SIMPLE_MULTILINE_STRING_START,
  INTERPOLATED_STRING_MIDDLE,
  INTERPOLATED_MULTILINE_STRING_MIDDLE,
  RAW_STRING_START,
  RAW_STRING_MIDDLE,
  RAW_STRING_MULTILINE_MIDDLE,
  SINGLE_LINE_STRING_END,
  MULTILINE_STRING_END,
  ELSE,
  CATCH,
  FINALLY,
  EXTENDS,
  DERIVES,
  WITH,
  ERROR_SENTINEL
}
typedef struct {
  Array(int16_t) indents;
  int16_t last_indentation_size;
  int16_t last_newline_count;
  int16_t last_column;
} Scanner;
void *tree_sitter_scala_external_scanner_create()
void tree_sitter_scala_external_scanner_destroy(void *payload)
unsigned tree_sitter_scala_external_scanner_serialize(void *payload, char *buffer)
void tree_sitter_scala_external_scanner_deserialize(void *payload, const char *buffer,
                                                    unsigned length)
static inline void advance(TSLexer *lexer)
static inline void skip(TSLexer *lexer)
typedef enum {
  STRING_MODE_SIMPLE,
  STRING_MODE_INTERPOLATED,
  STRING_MODE_RAW
} StringMode;
static bool scan_string_content(TSLexer *lexer, bool is_multiline, StringMode string_mode)
static bool detect_comment_start(TSLexer *lexer)
static bool scan_word(TSLexer *lexer, const char* const word)
static inline void debug_indents(Scanner *scanner)
bool tree_sitter_scala_external_scanner_scan(void *payload, TSLexer *lexer,
                                             const bool *valid_symbols)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/scala/tree_sitter/alloc.h

```cpp
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#define TREE_SITTER_ALLOC_H_
#define ts_malloc  ts_current_malloc
#define ts_calloc  ts_current_calloc
#define ts_realloc ts_current_realloc
#define ts_free    ts_current_free
#define ts_malloc  malloc
#define ts_calloc  calloc
#define ts_realloc realloc
#define ts_free    free
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/scala/tree_sitter/array.h

```cpp
#include "./alloc.h"
#include <assert.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#define TREE_SITTER_ARRAY_H_
#define Array(T)       \
  struct {             \
    T *contents;       \
    uint32_t size;     \
    uint32_t capacity; \
  }
#define array_init(self) \
  ((self)->size = 0, (self)->capacity = 0, (self)->contents = NULL)
#define array_new() \
  { NULL, 0, 0 }
#define array_get(self, _index) \
  (assert((uint32_t)(_index) < (self)->size), &(self)->contents[_index])
#define array_front(self) array_get(self, 0)
#define array_back(self) array_get(self, (self)->size - 1)
#define array_clear(self) ((self)->size = 0)
#define array_reserve(self, new_capacity) \
  _array__reserve((Array *)(self), array_elem_size(self), new_capacity)
#define array_delete(self) _array__delete((Array *)(self))
#define array_push(self, element)                            \
  (_array__grow((Array *)(self), 1, array_elem_size(self)), \
   (self)->contents[(self)->size++] = (element))
#define array_grow_by(self, count) \
  do { \
    if ((count) == 0) break; \
    _array__grow((Array *)(self), count, array_elem_size(self)); \
    memset((self)->contents + (self)->size, 0, (count) * array_elem_size(self)); \
    (self)->size += (count); \
  } while (0)
#define array_push_all(self, other)                                       \
  array_extend((self), (other)->size, (other)->contents)
#define array_extend(self, count, contents)                    \
  _array__splice(                                               \
    (Array *)(self), array_elem_size(self), (self)->size, \
    0, count,  contents                                        \
  )
#define array_splice(self, _index, old_count, new_count, new_contents)  \
  _array__splice(                                                       \
    (Array *)(self), array_elem_size(self), _index,                \
    old_count, new_count, new_contents                                 \
  )
#define array_insert(self, _index, element) \
  _array__splice((Array *)(self), array_elem_size(self), _index, 0, 1, &(element))
#define array_erase(self, _index) \
  _array__erase((Array *)(self), array_elem_size(self), _index)
#define array_pop(self) ((self)->contents[--(self)->size])
#define array_assign(self, other) \
  _array__assign((Array *)(self), (const Array *)(other), array_elem_size(self))
#define array_swap(self, other) \
  _array__swap((Array *)(self), (Array *)(other))
#define array_elem_size(self) (sizeof *(self)->contents)
#define array_search_sorted_with(self, compare, needle, _index, _exists) \
  _array__search_sorted(self, 0, compare, , needle, _index, _exists)
#define array_search_sorted_by(self, field, needle, _index, _exists) \
  _array__search_sorted(self, 0, _compare_int, field, needle, _index, _exists)
#define array_insert_sorted_with(self, compare, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_with(self, compare, &(value), &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
#define array_insert_sorted_by(self, field, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_by(self, field, (value) field, &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
static inline void _array__delete(Array *self)
static inline void _array__erase(Array *self, size_t element_size,
                                uint32_t index)
static inline void _array__reserve(Array *self, size_t element_size, uint32_t new_capacity)
static inline void _array__assign(Array *self, const Array *other, size_t element_size)
static inline void _array__swap(Array *self, Array *other)
static inline void _array__grow(Array *self, uint32_t count, size_t element_size)
static inline void _array__splice(Array *self, size_t element_size,
                                 uint32_t index, uint32_t old_count,
                                 uint32_t new_count, const void *elements)
#define _array__search_sorted(self, start, compare, suffix, needle, _index, _exists) \
  do { \
    *(_index) = start; \
    *(_exists) = false; \
    uint32_t size = (self)->size - *(_index); \
    if (size == 0) break; \
    int comparison; \
    while (size > 1) { \
      uint32_t half_size = size / 2; \
      uint32_t mid_index = *(_index) + half_size; \
      comparison = compare(&((self)->contents[mid_index] suffix), (needle)); \
      if (comparison <= 0) *(_index) = mid_index; \
      size -= half_size; \
    } \
    comparison = compare(&((self)->contents[*(_index)] suffix), (needle)); \
    if (comparison == 0) *(_exists) = true; \
    else if (comparison < 0) *(_index) += 1; \
  } while (0)
#define _compare_int(a, b) ((int)*(a) - (int)(b))
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/scala/tree_sitter/parser.h

```cpp
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#define TREE_SITTER_PARSER_H_
#define ts_builtin_sym_error ((TSSymbol)-1)
#define ts_builtin_sym_end 0
#define TREE_SITTER_SERIALIZATION_BUFFER_SIZE 1024
typedef uint16_t TSStateId;
typedef uint16_t TSSymbol;
typedef uint16_t TSFieldId;
struct TSLanguage
typedef struct {
  TSFieldId field_id;
  uint8_t child_index;
  bool inherited;
} TSFieldMapEntry;
typedef struct {
  uint16_t index;
  uint16_t length;
} TSFieldMapSlice;
typedef struct {
  bool visible;
  bool named;
  bool supertype;
} TSSymbolMetadata;
struct TSLexer
struct TSLexer
typedef enum {
  TSParseActionTypeShift,
  TSParseActionTypeReduce,
  TSParseActionTypeAccept,
  TSParseActionTypeRecover,
} TSParseActionType;
typedef union {
  struct {
    uint8_t type;
    TSStateId state;
    bool extra;
    bool repetition;
  } shift;
  struct {
    uint8_t type;
    uint8_t child_count;
    TSSymbol symbol;
    int16_t dynamic_precedence;
    uint16_t production_id;
  } reduce;
  uint8_t type;
} TSParseAction;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
} TSLexMode;
typedef union {
  TSParseAction action;
  struct {
    uint8_t count;
    bool reusable;
  } entry;
} TSParseActionEntry;
typedef struct {
  int32_t start;
  int32_t end;
} TSCharacterRange;
struct TSLanguage
static inline bool set_contains(TSCharacterRange *ranges, uint32_t len, int32_t lookahead)
#define UNUSED __pragma(warning(suppress : 4101))
#define UNUSED __attribute__((unused))
#define START_LEXER()           \
  bool result = false;          \
  bool skip = false;            \
  UNUSED                        \
  bool eof = false;             \
  int32_t lookahead;            \
  goto start;                   \
  next_state:                   \
  lexer->advance(lexer, skip);  \
  start:                        \
  skip = false;                 \
  lookahead = lexer->lookahead;
#define ADVANCE(state_value) \
  {                          \
    state = state_value;     \
    goto next_state;         \
  }
#define ADVANCE_MAP(...)                                              \
  {                                                                   \
    static const uint16_t map[] = { __VA_ARGS__ };                    \
    for (uint32_t i = 0; i < sizeof(map) / sizeof(map[0]); i += 2) {  \
      if (map[i] == lookahead) {                                      \
        state = map[i + 1];                                           \
        goto next_state;                                              \
      }                                                               \
    }                                                                 \
  }
#define SKIP(state_value) \
  {                       \
    skip = true;          \
    state = state_value;  \
    goto next_state;      \
  }
#define ACCEPT_TOKEN(symbol_value)     \
  result = true;                       \
  lexer->result_symbol = symbol_value; \
  lexer->mark_end(lexer);
#define END_STATE() return result;
#define SMALL_STATE(id) ((id) - LARGE_STATE_COUNT)
#define STATE(id) id
#define ACTIONS(id) id
#define SHIFT(state_value)            \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value)          \
    }                                 \
  }}
#define SHIFT_REPEAT(state_value)     \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value),         \
      .repetition = true              \
    }                                 \
  }}
#define SHIFT_EXTRA()                 \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .extra = true                   \
    }                                 \
  }}
#define REDUCE(symbol_name, children, precedence, prod_id) \
  {{                                                       \
    .reduce = {                                            \
      .type = TSParseActionTypeReduce,                     \
      .symbol = symbol_name,                               \
      .child_count = children,                             \
      .dynamic_precedence = precedence,                    \
      .production_id = prod_id                             \
    },                                                     \
  }}
#define RECOVER()                    \
  {{                                 \
    .type = TSParseActionTypeRecover \
  }}
#define ACCEPT_INPUT()              \
  {{                                \
    .type = TSParseActionTypeAccept \
  }}
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/sql/binding.go

```go
import "C"
import "unsafe"
func Language() unsafe.Pointer
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/sql/scanner.c

```c
#include "tree_sitter/parser.h"
#include <stdlib.h>
#include <string.h>
#include <wctype.h>
enum TokenType {
  DOLLAR_QUOTED_STRING_START_TAG,
  DOLLAR_QUOTED_STRING_END_TAG,
  DOLLAR_QUOTED_STRING
}
#define MALLOC_STRING_SIZE 1024
struct LexerState {
  char* start_tag;
}
void *tree_sitter_sql_external_scanner_create()
void tree_sitter_sql_external_scanner_destroy(void *payload)
static char* add_char(char* text, size_t* text_size, char c, int index)
static char* scan_dollar_string_tag(TSLexer *lexer)
bool tree_sitter_sql_external_scanner_scan(void *payload, TSLexer *lexer, const bool *valid_symbols)
unsigned tree_sitter_sql_external_scanner_serialize(void *payload, char *buffer)
void tree_sitter_sql_external_scanner_deserialize(void *payload, const char *buffer, unsigned length)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/sql/tree_sitter/alloc.h

```cpp
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#define TREE_SITTER_ALLOC_H_
#define ts_malloc  ts_current_malloc
#define ts_calloc  ts_current_calloc
#define ts_realloc ts_current_realloc
#define ts_free    ts_current_free
#define ts_malloc  malloc
#define ts_calloc  calloc
#define ts_realloc realloc
#define ts_free    free
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/sql/tree_sitter/array.h

```cpp
#include "./alloc.h"
#include <assert.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#define TREE_SITTER_ARRAY_H_
#define Array(T)       \
  struct {             \
    T *contents;       \
    uint32_t size;     \
    uint32_t capacity; \
  }
#define array_init(self) \
  ((self)->size = 0, (self)->capacity = 0, (self)->contents = NULL)
#define array_new() \
  { NULL, 0, 0 }
#define array_get(self, _index) \
  (assert((uint32_t)(_index) < (self)->size), &(self)->contents[_index])
#define array_front(self) array_get(self, 0)
#define array_back(self) array_get(self, (self)->size - 1)
#define array_clear(self) ((self)->size = 0)
#define array_reserve(self, new_capacity)        \
  ((self)->contents = _array__reserve(           \
    (void *)(self)->contents, &(self)->capacity, \
    array_elem_size(self), new_capacity)         \
  )
#define array_delete(self)                           \
  do {                                               \
    if ((self)->contents) ts_free((self)->contents); \
    (self)->contents = NULL;                         \
    (self)->size = 0;                                \
    (self)->capacity = 0;                            \
  } while (0)
#define array_push(self, element)                                 \
  do {                                                            \
    (self)->contents = _array__grow(                              \
      (void *)(self)->contents, (self)->size, &(self)->capacity,  \
      1, array_elem_size(self)                                    \
    );                                                            \
   (self)->contents[(self)->size++] = (element);                  \
  } while(0)
#define array_grow_by(self, count)                                               \
  do {                                                                           \
    if ((count) == 0) break;                                                     \
    (self)->contents = _array__grow(                                             \
      (self)->contents, (self)->size, &(self)->capacity,                         \
      count, array_elem_size(self)                                               \
    );                                                                           \
    memset((self)->contents + (self)->size, 0, (count) * array_elem_size(self)); \
    (self)->size += (count);                                                     \
  } while (0)
#define array_push_all(self, other) \
  array_extend((self), (other)->size, (other)->contents)
#define array_extend(self, count, other_contents)                 \
  (self)->contents = _array__splice(                              \
    (void*)(self)->contents, &(self)->size, &(self)->capacity,    \
    array_elem_size(self), (self)->size, 0, count, other_contents \
  )
#define array_splice(self, _index, old_count, new_count, new_contents) \
  (self)->contents = _array__splice(                                   \
    (void *)(self)->contents, &(self)->size, &(self)->capacity,        \
    array_elem_size(self), _index, old_count, new_count, new_contents  \
  )
#define array_insert(self, _index, element)                     \
  (self)->contents = _array__splice(                            \
    (void *)(self)->contents, &(self)->size, &(self)->capacity, \
    array_elem_size(self), _index, 0, 1, &(element)             \
  )
#define array_erase(self, _index) \
  _array__erase((void *)(self)->contents, &(self)->size, array_elem_size(self), _index)
#define array_pop(self) ((self)->contents[--(self)->size])
#define array_assign(self, other)                                   \
  (self)->contents = _array__assign(                                \
    (void *)(self)->contents, &(self)->size, &(self)->capacity,     \
    (const void *)(other)->contents, (other)->size, array_elem_size(self) \
  )
#define array_swap(self, other)                                     \
  do {                                                              \
    void *_array_swap_tmp = (void *)(self)->contents;               \
    (self)->contents = (other)->contents;                           \
    (other)->contents = _array_swap_tmp;                            \
    _array__swap(&(self)->size, &(self)->capacity,                  \
                 &(other)->size, &(other)->capacity);               \
  } while (0)
#define array_elem_size(self) (sizeof *(self)->contents)
#define array_search_sorted_with(self, compare, needle, _index, _exists) \
  _array__search_sorted(self, 0, compare, , needle, _index, _exists)
#define array_search_sorted_by(self, field, needle, _index, _exists) \
  _array__search_sorted(self, 0, _compare_int, field, needle, _index, _exists)
#define array_insert_sorted_with(self, compare, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_with(self, compare, &(value), &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
#define array_insert_sorted_by(self, field, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_by(self, field, (value) field, &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
static inline void _array__erase(void* self_contents, uint32_t *size,
                                size_t element_size, uint32_t index)
static inline void *_array__reserve(void *contents, uint32_t *capacity,
                                  size_t element_size, uint32_t new_capacity)
static inline void *_array__assign(void* self_contents, uint32_t *self_size, uint32_t *self_capacity,
                                 const void *other_contents, uint32_t other_size, size_t element_size)
static inline void _array__swap(uint32_t *self_size, uint32_t *self_capacity,
                               uint32_t *other_size, uint32_t *other_capacity)
static inline void *_array__grow(void *contents, uint32_t size, uint32_t *capacity,
                               uint32_t count, size_t element_size)
static inline void *_array__splice(void *self_contents, uint32_t *size, uint32_t *capacity,
                                 size_t element_size,
                                 uint32_t index, uint32_t old_count,
                                 uint32_t new_count, const void *elements)
#define _array__search_sorted(self, start, compare, suffix, needle, _index, _exists) \
  do { \
    *(_index) = start; \
    *(_exists) = false; \
    uint32_t size = (self)->size - *(_index); \
    if (size == 0) break; \
    int comparison; \
    while (size > 1) { \
      uint32_t half_size = size / 2; \
      uint32_t mid_index = *(_index) + half_size; \
      comparison = compare(&((self)->contents[mid_index] suffix), (needle)); \
      if (comparison <= 0) *(_index) = mid_index; \
      size -= half_size; \
    } \
    comparison = compare(&((self)->contents[*(_index)] suffix), (needle)); \
    if (comparison == 0) *(_exists) = true; \
    else if (comparison < 0) *(_index) += 1; \
  } while (0)
#define _compare_int(a, b) ((int)*(a) - (int)(b))
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/sql/tree_sitter/parser.h

```cpp
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#define TREE_SITTER_PARSER_H_
#define ts_builtin_sym_error ((TSSymbol)-1)
#define ts_builtin_sym_end 0
#define TREE_SITTER_SERIALIZATION_BUFFER_SIZE 1024
typedef uint16_t TSStateId;
typedef uint16_t TSSymbol;
typedef uint16_t TSFieldId;
struct TSLanguage
struct TSLanguageMetadata
typedef struct {
  TSFieldId field_id;
  uint8_t child_index;
  bool inherited;
} TSFieldMapEntry;
typedef struct {
  uint16_t index;
  uint16_t length;
} TSMapSlice;
typedef struct {
  bool visible;
  bool named;
  bool supertype;
} TSSymbolMetadata;
struct TSLexer
struct TSLexer
typedef enum {
  TSParseActionTypeShift,
  TSParseActionTypeReduce,
  TSParseActionTypeAccept,
  TSParseActionTypeRecover,
} TSParseActionType;
typedef union {
  struct {
    uint8_t type;
    TSStateId state;
    bool extra;
    bool repetition;
  } shift;
  struct {
    uint8_t type;
    uint8_t child_count;
    TSSymbol symbol;
    int16_t dynamic_precedence;
    uint16_t production_id;
  } reduce;
  uint8_t type;
} TSParseAction;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
} TSLexMode;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
  uint16_t reserved_word_set_id;
} TSLexerMode;
typedef union {
  TSParseAction action;
  struct {
    uint8_t count;
    bool reusable;
  } entry;
} TSParseActionEntry;
typedef struct {
  int32_t start;
  int32_t end;
} TSCharacterRange;
struct TSLanguage
static inline bool set_contains(const TSCharacterRange *ranges, uint32_t len, int32_t lookahead)
#define UNUSED __pragma(warning(suppress : 4101))
#define UNUSED __attribute__((unused))
#define START_LEXER()           \
  bool result = false;          \
  bool skip = false;            \
  UNUSED                        \
  bool eof = false;             \
  int32_t lookahead;            \
  goto start;                   \
  next_state:                   \
  lexer->advance(lexer, skip);  \
  start:                        \
  skip = false;                 \
  lookahead = lexer->lookahead;
#define ADVANCE(state_value) \
  {                          \
    state = state_value;     \
    goto next_state;         \
  }
#define ADVANCE_MAP(...)                                              \
  {                                                                   \
    static const uint16_t map[] = { __VA_ARGS__ };                    \
    for (uint32_t i = 0; i < sizeof(map) / sizeof(map[0]); i += 2) {  \
      if (map[i] == lookahead) {                                      \
        state = map[i + 1];                                           \
        goto next_state;                                              \
      }                                                               \
    }                                                                 \
  }
#define SKIP(state_value) \
  {                       \
    skip = true;          \
    state = state_value;  \
    goto next_state;      \
  }
#define ACCEPT_TOKEN(symbol_value)     \
  result = true;                       \
  lexer->result_symbol = symbol_value; \
  lexer->mark_end(lexer);
#define END_STATE() return result;
#define SMALL_STATE(id) ((id) - LARGE_STATE_COUNT)
#define STATE(id) id
#define ACTIONS(id) id
#define SHIFT(state_value)            \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value)          \
    }                                 \
  }}
#define SHIFT_REPEAT(state_value)     \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value),         \
      .repetition = true              \
    }                                 \
  }}
#define SHIFT_EXTRA()                 \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .extra = true                   \
    }                                 \
  }}
#define REDUCE(symbol_name, children, precedence, prod_id) \
  {{                                                       \
    .reduce = {                                            \
      .type = TSParseActionTypeReduce,                     \
      .symbol = symbol_name,                               \
      .child_count = children,                             \
      .dynamic_precedence = precedence,                    \
      .production_id = prod_id                             \
    },                                                     \
  }}
#define RECOVER()                    \
  {{                                 \
    .type = TSParseActionTypeRecover \
  }}
#define ACCEPT_INPUT()              \
  {{                                \
    .type = TSParseActionTypeAccept \
  }}
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/swift/binding.go

```go
import "C"
import "unsafe"
func Language() unsafe.Pointer
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/swift/scanner.c

```c
#include "tree_sitter/parser.h"
#include <string.h>
#include <wctype.h>
#define TOKEN_COUNT 33
enum TokenType {
    BLOCK_COMMENT,
    RAW_STR_PART,
    RAW_STR_CONTINUING_INDICATOR,
    RAW_STR_END_PART,
    IMPLICIT_SEMI,
    EXPLICIT_SEMI,
    ARROW_OPERATOR,
    DOT_OPERATOR,
    CONJUNCTION_OPERATOR,
    DISJUNCTION_OPERATOR,
    NIL_COALESCING_OPERATOR,
    EQUAL_SIGN,
    EQ_EQ,
    PLUS_THEN_WS,
    MINUS_THEN_WS,
    BANG,
    THROWS_KEYWORD,
    RETHROWS_KEYWORD,
    DEFAULT_KEYWORD,
    WHERE_KEYWORD,
    ELSE_KEYWORD,
    CATCH_KEYWORD,
    AS_KEYWORD,
    AS_QUEST,
    AS_BANG,
    ASYNC_KEYWORD,
    CUSTOM_OPERATOR,
    HASH_SYMBOL,
    DIRECTIVE_IF,
    DIRECTIVE_ELSEIF,
    DIRECTIVE_ELSE,
    DIRECTIVE_ENDIF,
    FAKE_TRY_BANG
}
#define OPERATOR_COUNT 20
enum IllegalTerminatorGroup {
    ALPHANUMERIC,
    OPERATOR_SYMBOLS,
    OPERATOR_OR_DOT,
    NON_WHITESPACE
}
enum IllegalTerminatorGroup
enum TokenType
#define RESERVED_OP_COUNT 31
static bool is_cross_semi_token(enum TokenType op)
#define NON_CONSUMING_CROSS_SEMI_CHAR_COUNT 3
enum ParseDirective {
    CONTINUE_PARSING_NOTHING_FOUND,
    CONTINUE_PARSING_TOKEN_FOUND,
    CONTINUE_PARSING_SLASH_CONSUMED,
    STOP_PARSING_NOTHING_FOUND,
    STOP_PARSING_TOKEN_FOUND,
    STOP_PARSING_END_OF_FILE
}
struct ScannerState {
    uint32_t ongoing_raw_str_hash_count;
}
void *tree_sitter_swift_external_scanner_create()
struct ScannerState
void tree_sitter_swift_external_scanner_destroy(void *payload)
void tree_sitter_swift_external_scanner_reset(void *payload)
struct ScannerState
unsigned tree_sitter_swift_external_scanner_serialize(void *payload, char *buffer)
struct ScannerState
void tree_sitter_swift_external_scanner_deserialize(
    void *payload,
    const char *buffer,
    unsigned length
)
struct ScannerState
static void advance(TSLexer *lexer)
static bool should_treat_as_wspace(int32_t character)
static int32_t encountered_op_count(bool *encountered_operator)
static bool any_reserved_ops(uint8_t *encountered_reserved_ops)
static bool is_legal_custom_operator(
    int32_t char_idx,
    int32_t first_char,
    int32_t cur_char
)
static bool eat_operators(
    TSLexer *lexer,
    const bool *valid_symbols,
    bool mark_end,
    const int32_t prior_char,
    enum TokenType *symbol_result
)
enum TokenType
enum IllegalTerminatorGroup
enum ParseDirective
enum TokenType
enum ParseDirective
enum TokenType
enum ParseDirective
enum ParseDirective
enum TokenType
enum TokenType
#define DIRECTIVE_COUNT 4
enum TokenType
enum TokenType
static bool eat_raw_str_part(
    struct ScannerState *state,
    TSLexer *lexer,
    const bool *valid_symbols,
    enum TokenType *symbol_result
)
struct ScannerState
enum TokenType
bool tree_sitter_swift_external_scanner_scan(
    void *payload,
    TSLexer *lexer,
    const bool *valid_symbols
)
struct ScannerState
enum TokenType
enum ParseDirective
enum TokenType
enum ParseDirective
enum TokenType
enum TokenType
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/swift/tree_sitter/alloc.h

```cpp
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#define TREE_SITTER_ALLOC_H_
#define ts_malloc  ts_current_malloc
#define ts_calloc  ts_current_calloc
#define ts_realloc ts_current_realloc
#define ts_free    ts_current_free
#define ts_malloc  malloc
#define ts_calloc  calloc
#define ts_realloc realloc
#define ts_free    free
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/swift/tree_sitter/array.h

```cpp
#include "./alloc.h"
#include <assert.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#define TREE_SITTER_ARRAY_H_
#define Array(T)       \
  struct {             \
    T *contents;       \
    uint32_t size;     \
    uint32_t capacity; \
  }
#define array_init(self) \
  ((self)->size = 0, (self)->capacity = 0, (self)->contents = NULL)
#define array_new() \
  { NULL, 0, 0 }
#define array_get(self, _index) \
  (assert((uint32_t)(_index) < (self)->size), &(self)->contents[_index])
#define array_front(self) array_get(self, 0)
#define array_back(self) array_get(self, (self)->size - 1)
#define array_clear(self) ((self)->size = 0)
#define array_reserve(self, new_capacity) \
  _array__reserve((Array *)(self), array_elem_size(self), new_capacity)
#define array_delete(self) _array__delete((Array *)(self))
#define array_push(self, element)                            \
  (_array__grow((Array *)(self), 1, array_elem_size(self)), \
   (self)->contents[(self)->size++] = (element))
#define array_grow_by(self, count) \
  do { \
    if ((count) == 0) break; \
    _array__grow((Array *)(self), count, array_elem_size(self)); \
    memset((self)->contents + (self)->size, 0, (count) * array_elem_size(self)); \
    (self)->size += (count); \
  } while (0)
#define array_push_all(self, other)                                       \
  array_extend((self), (other)->size, (other)->contents)
#define array_extend(self, count, contents)                    \
  _array__splice(                                               \
    (Array *)(self), array_elem_size(self), (self)->size, \
    0, count,  contents                                        \
  )
#define array_splice(self, _index, old_count, new_count, new_contents)  \
  _array__splice(                                                       \
    (Array *)(self), array_elem_size(self), _index,                \
    old_count, new_count, new_contents                                 \
  )
#define array_insert(self, _index, element) \
  _array__splice((Array *)(self), array_elem_size(self), _index, 0, 1, &(element))
#define array_erase(self, _index) \
  _array__erase((Array *)(self), array_elem_size(self), _index)
#define array_pop(self) ((self)->contents[--(self)->size])
#define array_assign(self, other) \
  _array__assign((Array *)(self), (const Array *)(other), array_elem_size(self))
#define array_swap(self, other) \
  _array__swap((Array *)(self), (Array *)(other))
#define array_elem_size(self) (sizeof *(self)->contents)
#define array_search_sorted_with(self, compare, needle, _index, _exists) \
  _array__search_sorted(self, 0, compare, , needle, _index, _exists)
#define array_search_sorted_by(self, field, needle, _index, _exists) \
  _array__search_sorted(self, 0, _compare_int, field, needle, _index, _exists)
#define array_insert_sorted_with(self, compare, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_with(self, compare, &(value), &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
#define array_insert_sorted_by(self, field, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_by(self, field, (value) field, &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
static inline void _array__delete(Array *self)
static inline void _array__erase(Array *self, size_t element_size,
                                uint32_t index)
static inline void _array__reserve(Array *self, size_t element_size, uint32_t new_capacity)
static inline void _array__assign(Array *self, const Array *other, size_t element_size)
static inline void _array__swap(Array *self, Array *other)
static inline void _array__grow(Array *self, uint32_t count, size_t element_size)
static inline void _array__splice(Array *self, size_t element_size,
                                 uint32_t index, uint32_t old_count,
                                 uint32_t new_count, const void *elements)
#define _array__search_sorted(self, start, compare, suffix, needle, _index, _exists) \
  do { \
    *(_index) = start; \
    *(_exists) = false; \
    uint32_t size = (self)->size - *(_index); \
    if (size == 0) break; \
    int comparison; \
    while (size > 1) { \
      uint32_t half_size = size / 2; \
      uint32_t mid_index = *(_index) + half_size; \
      comparison = compare(&((self)->contents[mid_index] suffix), (needle)); \
      if (comparison <= 0) *(_index) = mid_index; \
      size -= half_size; \
    } \
    comparison = compare(&((self)->contents[*(_index)] suffix), (needle)); \
    if (comparison == 0) *(_exists) = true; \
    else if (comparison < 0) *(_index) += 1; \
  } while (0)
#define _compare_int(a, b) ((int)*(a) - (int)(b))
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/swift/tree_sitter/parser.h

```cpp
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#define TREE_SITTER_PARSER_H_
#define ts_builtin_sym_error ((TSSymbol)-1)
#define ts_builtin_sym_end 0
#define TREE_SITTER_SERIALIZATION_BUFFER_SIZE 1024
typedef uint16_t TSStateId;
typedef uint16_t TSSymbol;
typedef uint16_t TSFieldId;
struct TSLanguage
typedef struct {
  TSFieldId field_id;
  uint8_t child_index;
  bool inherited;
} TSFieldMapEntry;
typedef struct {
  uint16_t index;
  uint16_t length;
} TSFieldMapSlice;
typedef struct {
  bool visible;
  bool named;
  bool supertype;
} TSSymbolMetadata;
struct TSLexer
struct TSLexer
typedef enum {
  TSParseActionTypeShift,
  TSParseActionTypeReduce,
  TSParseActionTypeAccept,
  TSParseActionTypeRecover,
} TSParseActionType;
typedef union {
  struct {
    uint8_t type;
    TSStateId state;
    bool extra;
    bool repetition;
  } shift;
  struct {
    uint8_t type;
    uint8_t child_count;
    TSSymbol symbol;
    int16_t dynamic_precedence;
    uint16_t production_id;
  } reduce;
  uint8_t type;
} TSParseAction;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
} TSLexMode;
typedef union {
  TSParseAction action;
  struct {
    uint8_t count;
    bool reusable;
  } entry;
} TSParseActionEntry;
typedef struct {
  int32_t start;
  int32_t end;
} TSCharacterRange;
struct TSLanguage
static inline bool set_contains(TSCharacterRange *ranges, uint32_t len, int32_t lookahead)
#define UNUSED __pragma(warning(suppress : 4101))
#define UNUSED __attribute__((unused))
#define START_LEXER()           \
  bool result = false;          \
  bool skip = false;            \
  UNUSED                        \
  bool eof = false;             \
  int32_t lookahead;            \
  goto start;                   \
  next_state:                   \
  lexer->advance(lexer, skip);  \
  start:                        \
  skip = false;                 \
  lookahead = lexer->lookahead;
#define ADVANCE(state_value) \
  {                          \
    state = state_value;     \
    goto next_state;         \
  }
#define ADVANCE_MAP(...)                                              \
  {                                                                   \
    static const uint16_t map[] = { __VA_ARGS__ };                    \
    for (uint32_t i = 0; i < sizeof(map) / sizeof(map[0]); i += 2) {  \
      if (map[i] == lookahead) {                                      \
        state = map[i + 1];                                           \
        goto next_state;                                              \
      }                                                               \
    }                                                                 \
  }
#define SKIP(state_value) \
  {                       \
    skip = true;          \
    state = state_value;  \
    goto next_state;      \
  }
#define ACCEPT_TOKEN(symbol_value)     \
  result = true;                       \
  lexer->result_symbol = symbol_value; \
  lexer->mark_end(lexer);
#define END_STATE() return result;
#define SMALL_STATE(id) ((id) - LARGE_STATE_COUNT)
#define STATE(id) id
#define ACTIONS(id) id
#define SHIFT(state_value)            \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value)          \
    }                                 \
  }}
#define SHIFT_REPEAT(state_value)     \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value),         \
      .repetition = true              \
    }                                 \
  }}
#define SHIFT_EXTRA()                 \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .extra = true                   \
    }                                 \
  }}
#define REDUCE(symbol_name, children, precedence, prod_id) \
  {{                                                       \
    .reduce = {                                            \
      .type = TSParseActionTypeReduce,                     \
      .symbol = symbol_name,                               \
      .child_count = children,                             \
      .dynamic_precedence = precedence,                    \
      .production_id = prod_id                             \
    },                                                     \
  }}
#define RECOVER()                    \
  {{                                 \
    .type = TSParseActionTypeRecover \
  }}
#define ACCEPT_INPUT()              \
  {{                                \
    .type = TSParseActionTypeAccept \
  }}
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/toml/binding.go

```go
import "C"
import "unsafe"
func Language() unsafe.Pointer
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/toml/parser.c

```c
#include "tree_sitter/parser.h"
#define LANGUAGE_VERSION 14
#define STATE_COUNT 152
#define LARGE_STATE_COUNT 2
#define SYMBOL_COUNT 66
#define ALIAS_COUNT 0
#define TOKEN_COUNT 40
#define EXTERNAL_TOKEN_COUNT 5
#define FIELD_COUNT 0
#define MAX_ALIAS_SEQUENCE_LENGTH 8
#define PRODUCTION_ID_COUNT 2
enum ts_symbol_identifiers {
  aux_sym_document_token1 = 1,
  sym_comment = 2,
  anon_sym_LBRACK = 3,
  anon_sym_RBRACK = 4,
  anon_sym_LBRACK_LBRACK = 5,
  anon_sym_RBRACK_RBRACK = 6,
  anon_sym_EQ = 7,
  anon_sym_DOT = 8,
  sym_bare_key = 9,
  anon_sym_DQUOTE = 10,
  aux_sym__basic_string_token1 = 11,
  anon_sym_DQUOTE2 = 12,
  anon_sym_DQUOTE_DQUOTE_DQUOTE = 13,
  aux_sym__multiline_basic_string_token1 = 14,
  sym_escape_sequence = 15,
  sym__escape_line_ending = 16,
  anon_sym_SQUOTE = 17,
  aux_sym__literal_string_token1 = 18,
  anon_sym_SQUOTE2 = 19,
  anon_sym_SQUOTE_SQUOTE_SQUOTE = 20,
  aux_sym_integer_token1 = 21,
  aux_sym_integer_token2 = 22,
  aux_sym_integer_token3 = 23,
  aux_sym_integer_token4 = 24,
  aux_sym_float_token1 = 25,
  aux_sym_float_token2 = 26,
  sym_boolean = 27,
  sym_offset_date_time = 28,
  sym_local_date_time = 29,
  sym_local_date = 30,
  sym_local_time = 31,
  anon_sym_COMMA = 32,
  anon_sym_LBRACE = 33,
  anon_sym_RBRACE = 34,
  sym__line_ending_or_eof = 35,
  sym__multiline_basic_string_content = 36,
  sym__multiline_basic_string_end = 37,
  sym__multiline_literal_string_content = 38,
  sym__multiline_literal_string_end = 39,
  sym_document = 40,
  sym_table = 41,
  sym_table_array_element = 42,
  sym_pair = 43,
  sym__inline_pair = 44,
  sym__key = 45,
  sym_dotted_key = 46,
  sym_quoted_key = 47,
  sym__inline_value = 48,
  sym_string = 49,
  sym__basic_string = 50,
  sym__multiline_basic_string = 51,
  sym__literal_string = 52,
  sym__multiline_literal_string = 53,
  sym_integer = 54,
  sym_float = 55,
  sym_array = 56,
  sym_inline_table = 57,
  aux_sym_document_repeat1 = 58,
  aux_sym_document_repeat2 = 59,
  aux_sym__basic_string_repeat1 = 60,
  aux_sym__multiline_basic_string_repeat1 = 61,
  aux_sym__multiline_literal_string_repeat1 = 62,
  aux_sym_array_repeat1 = 63,
  aux_sym_array_repeat2 = 64,
  aux_sym_inline_table_repeat1 = 65,
}
static bool ts_lex(TSLexer *lexer, TSStateId state)
enum ts_external_scanner_symbol_identifiers {
  ts_external_token__line_ending_or_eof = 0,
  ts_external_token__multiline_basic_string_content = 1,
  ts_external_token__multiline_basic_string_end = 2,
  ts_external_token__multiline_literal_string_content = 3,
  ts_external_token__multiline_literal_string_end = 4,
}
void *tree_sitter_toml_external_scanner_create(void);
void tree_sitter_toml_external_scanner_destroy(void *);
bool tree_sitter_toml_external_scanner_scan(void *, TSLexer *, const bool *);
unsigned tree_sitter_toml_external_scanner_serialize(void *, char *);
void tree_sitter_toml_external_scanner_deserialize(void *, const char *, unsigned);
#define TS_PUBLIC
#define TS_PUBLIC __declspec(dllexport)
#define TS_PUBLIC __attribute__((visibility("default")))
TS_PUBLIC const TSLanguage *tree_sitter_toml(void)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/toml/scanner.c

```c
#include "tree_sitter/parser.h"
typedef enum {
    LINE_ENDING_OR_EOF,
    MULTILINE_BASIC_STRING_CONTENT,
    MULTILINE_BASIC_STRING_END,
    MULTILINE_LITERAL_STRING_CONTENT,
    MULTILINE_LITERAL_STRING_END,
} TokenType;
void *tree_sitter_toml_external_scanner_create()
void tree_sitter_toml_external_scanner_destroy(void *payload)
unsigned tree_sitter_toml_external_scanner_serialize(void *payload, char *buffer)
void tree_sitter_toml_external_scanner_deserialize(void *payload, const char *buffer, unsigned length)
bool tree_sitter_toml_external_scanner_scan_multiline_string_end(TSLexer *lexer, const bool *valid_symbols,
                                                                 int32_t delimiter, TokenType content_symbol,
                                                                 TokenType end_symbol)
bool tree_sitter_toml_external_scanner_scan(void *payload, TSLexer *lexer, const bool *valid_symbols)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/toml/tree_sitter/alloc.h

```cpp
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#define TREE_SITTER_ALLOC_H_
#define ts_malloc  ts_current_malloc
#define ts_calloc  ts_current_calloc
#define ts_realloc ts_current_realloc
#define ts_free    ts_current_free
#define ts_malloc  malloc
#define ts_calloc  calloc
#define ts_realloc realloc
#define ts_free    free
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/toml/tree_sitter/array.h

```cpp
#include "./alloc.h"
#include <assert.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#define TREE_SITTER_ARRAY_H_
#define Array(T)       \
  struct {             \
    T *contents;       \
    uint32_t size;     \
    uint32_t capacity; \
  }
#define array_init(self) \
  ((self)->size = 0, (self)->capacity = 0, (self)->contents = NULL)
#define array_new() \
  { NULL, 0, 0 }
#define array_get(self, _index) \
  (assert((uint32_t)(_index) < (self)->size), &(self)->contents[_index])
#define array_front(self) array_get(self, 0)
#define array_back(self) array_get(self, (self)->size - 1)
#define array_clear(self) ((self)->size = 0)
#define array_reserve(self, new_capacity) \
  _array__reserve((Array *)(self), array_elem_size(self), new_capacity)
#define array_delete(self) _array__delete((Array *)(self))
#define array_push(self, element)                            \
  (_array__grow((Array *)(self), 1, array_elem_size(self)), \
   (self)->contents[(self)->size++] = (element))
#define array_grow_by(self, count) \
  do { \
    if ((count) == 0) break; \
    _array__grow((Array *)(self), count, array_elem_size(self)); \
    memset((self)->contents + (self)->size, 0, (count) * array_elem_size(self)); \
    (self)->size += (count); \
  } while (0)
#define array_push_all(self, other)                                       \
  array_extend((self), (other)->size, (other)->contents)
#define array_extend(self, count, contents)                    \
  _array__splice(                                               \
    (Array *)(self), array_elem_size(self), (self)->size, \
    0, count,  contents                                        \
  )
#define array_splice(self, _index, old_count, new_count, new_contents)  \
  _array__splice(                                                       \
    (Array *)(self), array_elem_size(self), _index,                \
    old_count, new_count, new_contents                                 \
  )
#define array_insert(self, _index, element) \
  _array__splice((Array *)(self), array_elem_size(self), _index, 0, 1, &(element))
#define array_erase(self, _index) \
  _array__erase((Array *)(self), array_elem_size(self), _index)
#define array_pop(self) ((self)->contents[--(self)->size])
#define array_assign(self, other) \
  _array__assign((Array *)(self), (const Array *)(other), array_elem_size(self))
#define array_swap(self, other) \
  _array__swap((Array *)(self), (Array *)(other))
#define array_elem_size(self) (sizeof *(self)->contents)
#define array_search_sorted_with(self, compare, needle, _index, _exists) \
  _array__search_sorted(self, 0, compare, , needle, _index, _exists)
#define array_search_sorted_by(self, field, needle, _index, _exists) \
  _array__search_sorted(self, 0, _compare_int, field, needle, _index, _exists)
#define array_insert_sorted_with(self, compare, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_with(self, compare, &(value), &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
#define array_insert_sorted_by(self, field, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_by(self, field, (value) field, &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
static inline void _array__delete(Array *self)
static inline void _array__erase(Array *self, size_t element_size,
                                uint32_t index)
static inline void _array__reserve(Array *self, size_t element_size, uint32_t new_capacity)
static inline void _array__assign(Array *self, const Array *other, size_t element_size)
static inline void _array__swap(Array *self, Array *other)
static inline void _array__grow(Array *self, uint32_t count, size_t element_size)
static inline void _array__splice(Array *self, size_t element_size,
                                 uint32_t index, uint32_t old_count,
                                 uint32_t new_count, const void *elements)
#define _array__search_sorted(self, start, compare, suffix, needle, _index, _exists) \
  do { \
    *(_index) = start; \
    *(_exists) = false; \
    uint32_t size = (self)->size - *(_index); \
    if (size == 0) break; \
    int comparison; \
    while (size > 1) { \
      uint32_t half_size = size / 2; \
      uint32_t mid_index = *(_index) + half_size; \
      comparison = compare(&((self)->contents[mid_index] suffix), (needle)); \
      if (comparison <= 0) *(_index) = mid_index; \
      size -= half_size; \
    } \
    comparison = compare(&((self)->contents[*(_index)] suffix), (needle)); \
    if (comparison == 0) *(_exists) = true; \
    else if (comparison < 0) *(_index) += 1; \
  } while (0)
#define _compare_int(a, b) ((int)*(a) - (int)(b))
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/toml/tree_sitter/parser.h

```cpp
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#define TREE_SITTER_PARSER_H_
#define ts_builtin_sym_error ((TSSymbol)-1)
#define ts_builtin_sym_end 0
#define TREE_SITTER_SERIALIZATION_BUFFER_SIZE 1024
typedef uint16_t TSStateId;
typedef uint16_t TSSymbol;
typedef uint16_t TSFieldId;
struct TSLanguage
typedef struct {
  TSFieldId field_id;
  uint8_t child_index;
  bool inherited;
} TSFieldMapEntry;
typedef struct {
  uint16_t index;
  uint16_t length;
} TSFieldMapSlice;
typedef struct {
  bool visible;
  bool named;
  bool supertype;
} TSSymbolMetadata;
struct TSLexer
struct TSLexer
typedef enum {
  TSParseActionTypeShift,
  TSParseActionTypeReduce,
  TSParseActionTypeAccept,
  TSParseActionTypeRecover,
} TSParseActionType;
typedef union {
  struct {
    uint8_t type;
    TSStateId state;
    bool extra;
    bool repetition;
  } shift;
  struct {
    uint8_t type;
    uint8_t child_count;
    TSSymbol symbol;
    int16_t dynamic_precedence;
    uint16_t production_id;
  } reduce;
  uint8_t type;
} TSParseAction;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
} TSLexMode;
typedef union {
  TSParseAction action;
  struct {
    uint8_t count;
    bool reusable;
  } entry;
} TSParseActionEntry;
typedef struct {
  int32_t start;
  int32_t end;
} TSCharacterRange;
struct TSLanguage
static inline bool set_contains(TSCharacterRange *ranges, uint32_t len, int32_t lookahead)
#define UNUSED __pragma(warning(suppress : 4101))
#define UNUSED __attribute__((unused))
#define START_LEXER()           \
  bool result = false;          \
  bool skip = false;            \
  UNUSED                        \
  bool eof = false;             \
  int32_t lookahead;            \
  goto start;                   \
  next_state:                   \
  lexer->advance(lexer, skip);  \
  start:                        \
  skip = false;                 \
  lookahead = lexer->lookahead;
#define ADVANCE(state_value) \
  {                          \
    state = state_value;     \
    goto next_state;         \
  }
#define ADVANCE_MAP(...)                                              \
  {                                                                   \
    static const uint16_t map[] = { __VA_ARGS__ };                    \
    for (uint32_t i = 0; i < sizeof(map) / sizeof(map[0]); i += 2) {  \
      if (map[i] == lookahead) {                                      \
        state = map[i + 1];                                           \
        goto next_state;                                              \
      }                                                               \
    }                                                                 \
  }
#define SKIP(state_value) \
  {                       \
    skip = true;          \
    state = state_value;  \
    goto next_state;      \
  }
#define ACCEPT_TOKEN(symbol_value)     \
  result = true;                       \
  lexer->result_symbol = symbol_value; \
  lexer->mark_end(lexer);
#define END_STATE() return result;
#define SMALL_STATE(id) ((id) - LARGE_STATE_COUNT)
#define STATE(id) id
#define ACTIONS(id) id
#define SHIFT(state_value)            \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value)          \
    }                                 \
  }}
#define SHIFT_REPEAT(state_value)     \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value),         \
      .repetition = true              \
    }                                 \
  }}
#define SHIFT_EXTRA()                 \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .extra = true                   \
    }                                 \
  }}
#define REDUCE(symbol_name, children, precedence, prod_id) \
  {{                                                       \
    .reduce = {                                            \
      .type = TSParseActionTypeReduce,                     \
      .symbol = symbol_name,                               \
      .child_count = children,                             \
      .dynamic_precedence = precedence,                    \
      .production_id = prod_id                             \
    },                                                     \
  }}
#define RECOVER()                    \
  {{                                 \
    .type = TSParseActionTypeRecover \
  }}
#define ACCEPT_INPUT()              \
  {{                                \
    .type = TSParseActionTypeAccept \
  }}
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/yaml/binding.go

```go
import "C"
import "unsafe"
func Language() unsafe.Pointer
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/yaml/scanner.c

```c
#include "tree_sitter/array.h"
#include "tree_sitter/parser.h"
#include _file(YAML_SCHEMA)
#define _str(x) #x
#define _file(x) _str(schema.x.c)
#define YAML_SCHEMA core
typedef enum {
    END_OF_FILE,

    S_DIR_YML_BGN,  R_DIR_YML_VER,
    S_DIR_TAG_BGN,  R_DIR_TAG_HDL,  R_DIR_TAG_PFX,
    S_DIR_RSV_BGN,  R_DIR_RSV_PRM,
    S_DRS_END,
    S_DOC_END,
    R_BLK_SEQ_BGN,  BR_BLK_SEQ_BGN, B_BLK_SEQ_BGN,
    R_BLK_KEY_BGN,  BR_BLK_KEY_BGN, B_BLK_KEY_BGN,
    R_BLK_VAL_BGN,  BR_BLK_VAL_BGN, B_BLK_VAL_BGN,
    R_BLK_IMP_BGN,
    R_BLK_LIT_BGN,  BR_BLK_LIT_BGN,
    R_BLK_FLD_BGN,  BR_BLK_FLD_BGN,
    BR_BLK_STR_CTN,
    R_FLW_SEQ_BGN,  BR_FLW_SEQ_BGN, B_FLW_SEQ_BGN,
    R_FLW_SEQ_END,  BR_FLW_SEQ_END, B_FLW_SEQ_END,
    R_FLW_MAP_BGN,  BR_FLW_MAP_BGN, B_FLW_MAP_BGN,
    R_FLW_MAP_END,  BR_FLW_MAP_END, B_FLW_MAP_END,
    R_FLW_SEP_BGN,  BR_FLW_SEP_BGN,
    R_FLW_KEY_BGN,  BR_FLW_KEY_BGN,
    R_FLW_JSV_BGN,  BR_FLW_JSV_BGN,
    R_FLW_NJV_BGN,  BR_FLW_NJV_BGN,
    R_DQT_STR_BGN,  BR_DQT_STR_BGN, B_DQT_STR_BGN,
    R_DQT_STR_CTN,  BR_DQT_STR_CTN,
    R_DQT_ESC_NWL,  BR_DQT_ESC_NWL,
    R_DQT_ESC_SEQ,  BR_DQT_ESC_SEQ,
    R_DQT_STR_END,  BR_DQT_STR_END,
    R_SQT_STR_BGN,  BR_SQT_STR_BGN, B_SQT_STR_BGN,
    R_SQT_STR_CTN,  BR_SQT_STR_CTN,
    R_SQT_ESC_SQT,  BR_SQT_ESC_SQT,
    R_SQT_STR_END,  BR_SQT_STR_END,

    R_SGL_PLN_NUL_BLK, BR_SGL_PLN_NUL_BLK, B_SGL_PLN_NUL_BLK, R_SGL_PLN_NUL_FLW, BR_SGL_PLN_NUL_FLW,
    R_SGL_PLN_BOL_BLK, BR_SGL_PLN_BOL_BLK, B_SGL_PLN_BOL_BLK, R_SGL_PLN_BOL_FLW, BR_SGL_PLN_BOL_FLW,
    R_SGL_PLN_INT_BLK, BR_SGL_PLN_INT_BLK, B_SGL_PLN_INT_BLK, R_SGL_PLN_INT_FLW, BR_SGL_PLN_INT_FLW,
    R_SGL_PLN_FLT_BLK, BR_SGL_PLN_FLT_BLK, B_SGL_PLN_FLT_BLK, R_SGL_PLN_FLT_FLW, BR_SGL_PLN_FLT_FLW,
    R_SGL_PLN_TMS_BLK, BR_SGL_PLN_TMS_BLK, B_SGL_PLN_TMS_BLK, R_SGL_PLN_TMS_FLW, BR_SGL_PLN_TMS_FLW,
    R_SGL_PLN_STR_BLK, BR_SGL_PLN_STR_BLK, B_SGL_PLN_STR_BLK, R_SGL_PLN_STR_FLW, BR_SGL_PLN_STR_FLW,

    R_MTL_PLN_STR_BLK,  BR_MTL_PLN_STR_BLK,
    R_MTL_PLN_STR_FLW,  BR_MTL_PLN_STR_FLW,

    R_TAG,     BR_TAG,     B_TAG,
    R_ACR_BGN, BR_ACR_BGN, B_ACR_BGN, R_ACR_CTN,
    R_ALS_BGN, BR_ALS_BGN, B_ALS_BGN, R_ALS_CTN,

    BL,
    COMMENT,

    ERR_REC,
} TokenType;
#define SCN_SUCC 1
#define SCN_STOP 0
#define SCN_FAIL (-1)
#define IND_ROT 'r'
#define IND_MAP 'm'
#define IND_SEQ 'q'
#define IND_STR 's'
#define RET_SYM(RESULT_SYMBOL)                                                                                         \
    {                                                                                                                  \
        flush(scanner);                                                                                                \
        lexer->result_symbol = RESULT_SYMBOL;                                                                          \
        return true;                                                                                                   \
    }
#define POP_IND()                                                                                                      \
    {                                                                                                                  \
        /* incorrect status caused by error recovering */
#define PUSH_IND(TYP, LEN) push_ind(scanner, TYP, LEN)
#define PUSH_BGN_IND(TYP)                                                                                              \
    {                                                                                                                  \
        if (has_tab_ind)                                                                                               \
            return false;                                                                                              \
        push_ind(scanner, TYP, bgn_col);                                                                               \
    }
#define MAY_PUSH_IMP_IND(TYP)                                                                                          \
    {                                                                                                                  \
        if (cur_ind != scanner->blk_imp_col) {                                                                         \
            if (scanner->blk_imp_tab)                                                                                  \
                return false;                                                                                          \
            push_ind(scanner, IND_MAP, scanner->blk_imp_col);                                                          \
        }                                                                                                              \
    }
#define MAY_PUSH_SPC_SEQ_IND()                                                                                         \
    {                                                                                                                  \
        if (cur_ind_typ == IND_MAP) {                                                                                  \
            push_ind(scanner, IND_SEQ, bgn_col);                                                                       \
        }                                                                                                              \
    }
#define MAY_UPD_IMP_COL()                                                                                              \
    {                                                                                                                  \
        if (scanner->blk_imp_row != bgn_row) {                                                                         \
            scanner->blk_imp_row = bgn_row;                                                                            \
            scanner->blk_imp_col = bgn_col;                                                                            \
            scanner->blk_imp_tab = has_tab_ind;                                                                        \
        }                                                                                                              \
    }
#define SGL_PLN_SYM(POS, CTX)                                                                                          \
    (scanner->rlt_sch == RS_NULL        ? POS##_SGL_PLN_NUL_##CTX                                                      \
     : scanner->rlt_sch == RS_BOOL      ? POS##_SGL_PLN_BOL_##CTX                                                      \
     : scanner->rlt_sch == RS_INT       ? POS##_SGL_PLN_INT_##CTX                                                      \
     : scanner->rlt_sch == RS_FLOAT     ? POS##_SGL_PLN_FLT_##CTX                                                      \
     : scanner->rlt_sch == RS_TIMESTAMP ? POS##_SGL_PLN_TMS_##CTX                                                      \
                                        : POS##_SGL_PLN_STR_##CTX)
#define SGL_PLN_SYM(POS, CTX)                                                                                          \
    (scanner->rlt_sch == RS_NULL        ? POS##_SGL_PLN_NUL_##CTX                                                      \
     : scanner->rlt_sch == RS_BOOL      ? POS##_SGL_PLN_BOL_##CTX                                                      \
     : scanner->rlt_sch == RS_INT       ? POS##_SGL_PLN_INT_##CTX                                                      \
     : scanner->rlt_sch == RS_FLOAT     ? POS##_SGL_PLN_FLT_##CTX                                                      \
                                        : POS##_SGL_PLN_STR_##CTX)
typedef struct {
    int16_t row;
    int16_t col;
    int16_t blk_imp_row;
    int16_t blk_imp_col;
    int16_t blk_imp_tab;
    Array(int16_t) ind_typ_stk;
    Array(int16_t) ind_len_stk;

    // temp
    int16_t end_row;
    int16_t end_col;
    int16_t cur_row;
    int16_t cur_col;
    int32_t cur_chr;
    int8_t sch_stt;
    ResultSchema rlt_sch;
} Scanner;
static unsigned serialize(Scanner *scanner, char *buffer)
static void deserialize(Scanner *scanner, const char *buffer, unsigned length)
static inline void adv(Scanner *scanner, TSLexer *lexer)
static inline void adv_nwl(Scanner *scanner, TSLexer *lexer)
static inline void skp(Scanner *scanner, TSLexer *lexer)
static inline void skp_nwl(Scanner *scanner, TSLexer *lexer)
static inline void mrk_end(Scanner *scanner, TSLexer *lexer)
static inline void init(Scanner *scanner)
static inline void flush(Scanner *scanner)
static inline void pop_ind(Scanner *scanner)
static inline void push_ind(Scanner *scanner, int16_t typ, int16_t len)
static inline bool is_wsp(int32_t c)
static inline bool is_nwl(int32_t c)
static inline bool is_wht(int32_t c)
static inline bool is_ns_dec_digit(int32_t c)
static inline bool is_ns_hex_digit(int32_t c)
static inline bool is_ns_word_char(int32_t c)
static inline bool is_nb_json(int32_t c)
static inline bool is_nb_double_char(int32_t c)
static inline bool is_nb_single_char(int32_t c)
static inline bool is_ns_char(int32_t c)
static inline bool is_c_indicator(int32_t c)
static inline bool is_c_flow_indicator(int32_t c)
static inline bool is_plain_safe_in_block(int32_t c)
static inline bool is_plain_safe_in_flow(int32_t c)
static inline bool is_ns_uri_char(int32_t c)
static inline bool is_ns_tag_char(int32_t c)
static inline bool is_ns_anchor_char(int32_t c)
static char scn_uri_esc(Scanner *scanner, TSLexer *lexer)
static char scn_ns_uri_char(Scanner *scanner, TSLexer *lexer)
static char scn_ns_tag_char(Scanner *scanner, TSLexer *lexer)
static bool scn_dir_bgn(Scanner *scanner, TSLexer *lexer)
static bool scn_dir_yml_ver(Scanner *scanner, TSLexer *lexer, TSSymbol result_symbol)
static bool scn_tag_hdl_tal(Scanner *scanner, TSLexer *lexer)
static bool scn_dir_tag_hdl(Scanner *scanner, TSLexer *lexer, TSSymbol result_symbol)
static bool scn_dir_tag_pfx(Scanner *scanner, TSLexer *lexer, TSSymbol result_symbol)
static bool scn_dir_rsv_prm(Scanner *scanner, TSLexer *lexer, TSSymbol result_symbol)
static bool scn_tag(Scanner *scanner, TSLexer *lexer, TSSymbol result_symbol)
static bool scn_acr_bgn(Scanner *scanner, TSLexer *lexer, TSSymbol result_symbol)
static bool scn_acr_ctn(Scanner *scanner, TSLexer *lexer, TSSymbol result_symbol)
static bool scn_als_bgn(Scanner *scanner, TSLexer *lexer, TSSymbol result_symbol)
static bool scn_als_ctn(Scanner *scanner, TSLexer *lexer, TSSymbol result_symbol)
static bool scn_dqt_esc_seq(Scanner *scanner, TSLexer *lexer, TSSymbol result_symbol)
static bool scn_drs_doc_end(Scanner *scanner, TSLexer *lexer)
static bool scn_dqt_str_cnt(Scanner *scanner, TSLexer *lexer, TSSymbol result_symbol)
static bool scn_sqt_str_cnt(Scanner *scanner, TSLexer *lexer, TSSymbol result_symbol)
static bool scn_blk_str_bgn(Scanner *scanner, TSLexer *lexer, TSSymbol result_symbol)
static bool scn_blk_str_cnt(Scanner *scanner, TSLexer *lexer, TSSymbol result_symbol)
static char scn_pln_cnt(Scanner *scanner, TSLexer *lexer, bool (*is_plain_safe)(int32_t))
static bool scan(Scanner *scanner, TSLexer *lexer, const bool *valid_symbols)
void *tree_sitter_yaml_external_scanner_create()
void tree_sitter_yaml_external_scanner_destroy(void *payload)
unsigned tree_sitter_yaml_external_scanner_serialize(void *payload, char *buffer)
void tree_sitter_yaml_external_scanner_deserialize(void *payload, const char *buffer, unsigned length)
bool tree_sitter_yaml_external_scanner_scan(void *payload, TSLexer *lexer, const bool *valid_symbols)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/yaml/schema.core.c

```c
#include <stdint.h>
#include <stdlib.h>
#define SCH_STT_FRZ -1
#define HAS_TIMESTAMP 0
typedef enum {
  RS_STR,
  RS_INT,
  RS_NULL,
  RS_BOOL,
  RS_FLOAT,
} ResultSchema;
static int8_t adv_sch_stt(int8_t sch_stt, int32_t cur_chr, ResultSchema *rlt_sch)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/yaml/schema.json.c

```c
#include <stdint.h>
#include <stdlib.h>
#define SCH_STT_FRZ -1
#define HAS_TIMESTAMP 0
typedef enum {
  RS_STR,
  RS_INT,
  RS_BOOL,
  RS_NULL,
  RS_FLOAT,
} ResultSchema;
static int8_t adv_sch_stt(int8_t sch_stt, int32_t cur_chr, ResultSchema *rlt_sch)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/yaml/schema.legacy.c

```c
#include <stdint.h>
#include <stdlib.h>
#define SCH_STT_FRZ -1
#define HAS_TIMESTAMP 1
typedef enum {
  RS_STR,
  RS_FLOAT,
  RS_INT,
  RS_BOOL,
  RS_NULL,
  RS_TIMESTAMP,
} ResultSchema;
static int8_t adv_sch_stt(int8_t sch_stt, int32_t cur_chr, ResultSchema *rlt_sch)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/yaml/tree_sitter/alloc.h

```cpp
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#define TREE_SITTER_ALLOC_H_
#define ts_malloc  ts_current_malloc
#define ts_calloc  ts_current_calloc
#define ts_realloc ts_current_realloc
#define ts_free    ts_current_free
#define ts_malloc  malloc
#define ts_calloc  calloc
#define ts_realloc realloc
#define ts_free    free
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/yaml/tree_sitter/array.h

```cpp
#include "./alloc.h"
#include <assert.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#define TREE_SITTER_ARRAY_H_
#define Array(T)       \
  struct {             \
    T *contents;       \
    uint32_t size;     \
    uint32_t capacity; \
  }
#define array_init(self) \
  ((self)->size = 0, (self)->capacity = 0, (self)->contents = NULL)
#define array_new() \
  { NULL, 0, 0 }
#define array_get(self, _index) \
  (assert((uint32_t)(_index) < (self)->size), &(self)->contents[_index])
#define array_front(self) array_get(self, 0)
#define array_back(self) array_get(self, (self)->size - 1)
#define array_clear(self) ((self)->size = 0)
#define array_reserve(self, new_capacity) \
  _array__reserve((Array *)(self), array_elem_size(self), new_capacity)
#define array_delete(self) _array__delete((Array *)(self))
#define array_push(self, element)                            \
  (_array__grow((Array *)(self), 1, array_elem_size(self)), \
   (self)->contents[(self)->size++] = (element))
#define array_grow_by(self, count) \
  do { \
    if ((count) == 0) break; \
    _array__grow((Array *)(self), count, array_elem_size(self)); \
    memset((self)->contents + (self)->size, 0, (count) * array_elem_size(self)); \
    (self)->size += (count); \
  } while (0)
#define array_push_all(self, other)                                       \
  array_extend((self), (other)->size, (other)->contents)
#define array_extend(self, count, contents)                    \
  _array__splice(                                               \
    (Array *)(self), array_elem_size(self), (self)->size, \
    0, count,  contents                                        \
  )
#define array_splice(self, _index, old_count, new_count, new_contents)  \
  _array__splice(                                                       \
    (Array *)(self), array_elem_size(self), _index,                \
    old_count, new_count, new_contents                                 \
  )
#define array_insert(self, _index, element) \
  _array__splice((Array *)(self), array_elem_size(self), _index, 0, 1, &(element))
#define array_erase(self, _index) \
  _array__erase((Array *)(self), array_elem_size(self), _index)
#define array_pop(self) ((self)->contents[--(self)->size])
#define array_assign(self, other) \
  _array__assign((Array *)(self), (const Array *)(other), array_elem_size(self))
#define array_swap(self, other) \
  _array__swap((Array *)(self), (Array *)(other))
#define array_elem_size(self) (sizeof *(self)->contents)
#define array_search_sorted_with(self, compare, needle, _index, _exists) \
  _array__search_sorted(self, 0, compare, , needle, _index, _exists)
#define array_search_sorted_by(self, field, needle, _index, _exists) \
  _array__search_sorted(self, 0, _compare_int, field, needle, _index, _exists)
#define array_insert_sorted_with(self, compare, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_with(self, compare, &(value), &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
#define array_insert_sorted_by(self, field, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_by(self, field, (value) field, &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
static inline void _array__delete(Array *self)
static inline void _array__erase(Array *self, size_t element_size,
                                uint32_t index)
static inline void _array__reserve(Array *self, size_t element_size, uint32_t new_capacity)
static inline void _array__assign(Array *self, const Array *other, size_t element_size)
static inline void _array__swap(Array *self, Array *other)
static inline void _array__grow(Array *self, uint32_t count, size_t element_size)
static inline void _array__splice(Array *self, size_t element_size,
                                 uint32_t index, uint32_t old_count,
                                 uint32_t new_count, const void *elements)
#define _array__search_sorted(self, start, compare, suffix, needle, _index, _exists) \
  do { \
    *(_index) = start; \
    *(_exists) = false; \
    uint32_t size = (self)->size - *(_index); \
    if (size == 0) break; \
    int comparison; \
    while (size > 1) { \
      uint32_t half_size = size / 2; \
      uint32_t mid_index = *(_index) + half_size; \
      comparison = compare(&((self)->contents[mid_index] suffix), (needle)); \
      if (comparison <= 0) *(_index) = mid_index; \
      size -= half_size; \
    } \
    comparison = compare(&((self)->contents[*(_index)] suffix), (needle)); \
    if (comparison == 0) *(_exists) = true; \
    else if (comparison < 0) *(_index) += 1; \
  } while (0)
#define _compare_int(a, b) ((int)*(a) - (int)(b))
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/yaml/tree_sitter/parser.h

```cpp
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#define TREE_SITTER_PARSER_H_
#define ts_builtin_sym_error ((TSSymbol)-1)
#define ts_builtin_sym_end 0
#define TREE_SITTER_SERIALIZATION_BUFFER_SIZE 1024
typedef uint16_t TSStateId;
typedef uint16_t TSSymbol;
typedef uint16_t TSFieldId;
struct TSLanguage
struct TSLanguageMetadata
typedef struct {
  TSFieldId field_id;
  uint8_t child_index;
  bool inherited;
} TSFieldMapEntry;
typedef struct {
  uint16_t index;
  uint16_t length;
} TSMapSlice;
typedef struct {
  bool visible;
  bool named;
  bool supertype;
} TSSymbolMetadata;
struct TSLexer
struct TSLexer
typedef enum {
  TSParseActionTypeShift,
  TSParseActionTypeReduce,
  TSParseActionTypeAccept,
  TSParseActionTypeRecover,
} TSParseActionType;
typedef union {
  struct {
    uint8_t type;
    TSStateId state;
    bool extra;
    bool repetition;
  } shift;
  struct {
    uint8_t type;
    uint8_t child_count;
    TSSymbol symbol;
    int16_t dynamic_precedence;
    uint16_t production_id;
  } reduce;
  uint8_t type;
} TSParseAction;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
} TSLexMode;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
  uint16_t reserved_word_set_id;
} TSLexerMode;
typedef union {
  TSParseAction action;
  struct {
    uint8_t count;
    bool reusable;
  } entry;
} TSParseActionEntry;
typedef struct {
  int32_t start;
  int32_t end;
} TSCharacterRange;
struct TSLanguage
static inline bool set_contains(const TSCharacterRange *ranges, uint32_t len, int32_t lookahead)
#define UNUSED __pragma(warning(suppress : 4101))
#define UNUSED __attribute__((unused))
#define START_LEXER()           \
  bool result = false;          \
  bool skip = false;            \
  UNUSED                        \
  bool eof = false;             \
  int32_t lookahead;            \
  goto start;                   \
  next_state:                   \
  lexer->advance(lexer, skip);  \
  start:                        \
  skip = false;                 \
  lookahead = lexer->lookahead;
#define ADVANCE(state_value) \
  {                          \
    state = state_value;     \
    goto next_state;         \
  }
#define ADVANCE_MAP(...)                                              \
  {                                                                   \
    static const uint16_t map[] = { __VA_ARGS__ };                    \
    for (uint32_t i = 0; i < sizeof(map) / sizeof(map[0]); i += 2) {  \
      if (map[i] == lookahead) {                                      \
        state = map[i + 1];                                           \
        goto next_state;                                              \
      }                                                               \
    }                                                                 \
  }
#define SKIP(state_value) \
  {                       \
    skip = true;          \
    state = state_value;  \
    goto next_state;      \
  }
#define ACCEPT_TOKEN(symbol_value)     \
  result = true;                       \
  lexer->result_symbol = symbol_value; \
  lexer->mark_end(lexer);
#define END_STATE() return result;
#define SMALL_STATE(id) ((id) - LARGE_STATE_COUNT)
#define STATE(id) id
#define ACTIONS(id) id
#define SHIFT(state_value)            \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value)          \
    }                                 \
  }}
#define SHIFT_REPEAT(state_value)     \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value),         \
      .repetition = true              \
    }                                 \
  }}
#define SHIFT_EXTRA()                 \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .extra = true                   \
    }                                 \
  }}
#define REDUCE(symbol_name, children, precedence, prod_id) \
  {{                                                       \
    .reduce = {                                            \
      .type = TSParseActionTypeReduce,                     \
      .symbol = symbol_name,                               \
      .child_count = children,                             \
      .dynamic_precedence = precedence,                    \
      .production_id = prod_id                             \
    },                                                     \
  }}
#define RECOVER()                    \
  {{                                 \
    .type = TSParseActionTypeRecover \
  }}
#define ACCEPT_INPUT()              \
  {{                                \
    .type = TSParseActionTypeAccept \
  }}
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/base.go

```go
import "strings"
func hasVisibilityPrefix(sigText, modifier string) bool
type BaseQuery struct{}
func (BaseQuery) Captures() []string
func (BaseQuery) ImportQuery() []byte
func (BaseQuery) CallQuery() []byte
func (BaseQuery) IsExported(name, _ string) bool
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/c.go

```go
import (
	"strings"
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_c "github.com/tree-sitter/tree-sitter-c/bindings/go"
)
type CQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewCQuery() *CQuery
func (q *CQuery) Language() *sitter.Language
func (q *CQuery) Query() []byte
cKindMapping = map[string]string{
	"function_definition":  "function",
	"declaration":          "function", // function prototypes
	"struct_specifier":     "struct",
	"enum_specifier":       "enum",
	"type_definition":      "typedef",
	"preproc_function_def": "macro",
	"preproc_def":          "macro",
	"global_variable":      "variable", // mapped from declaration patterns
}
func (q *CQuery) KindMapping() map[string]string
func (q *CQuery) ImportQuery() []byte
func (q *CQuery) CallQuery() []byte
func (q *CQuery) IsExported(name, sigText string) bool
cCallQueryPattern = `
; Direct function calls (e.g., foo())
(call_expression
  function: (identifier) @callee
) @call_node
`
cImportQueryPattern = `
; #include directives (capture full statement)
(preproc_include) @import_path
`
cQueryPattern = `
; Function definitions - direct declarator
(function_definition
  declarator: (function_declarator
    declarator: (identifier) @name
  )
) @signature @kind

; Function definitions - pointer return type
(function_definition
  declarator: (pointer_declarator
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Function declarations (prototypes) - direct declarator
(declaration
  declarator: (function_declarator
    declarator: (identifier) @name
  )
) @signature @kind

; Function declarations (prototypes) - pointer return type
(declaration
  declarator: (pointer_declarator
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Struct specifiers
(struct_specifier
  name: (type_identifier) @name
) @signature @kind

; Enum specifiers
(enum_specifier
  name: (type_identifier) @name
) @signature @kind

; Typedef
(type_definition
  declarator: (type_identifier) @name
) @signature @kind

; Function-like macros
(preproc_function_def
  name: (identifier) @name
) @signature @kind

; Object-like macros
(preproc_def
  name: (identifier) @name
) @signature @kind

; Global variable declarations (with initializer)
(translation_unit
  (declaration
    declarator: (init_declarator
      declarator: (identifier) @name
    )
  ) @signature @kind
)

; Global variable declarations (simple identifier, e.g., extern)
(translation_unit
  (declaration
    declarator: (identifier) @name
  ) @signature @kind
)

; Global pointer variable declarations
(translation_unit
  (declaration
    declarator: (pointer_declarator
      declarator: (identifier) @name
    )
  ) @signature @kind
)

; Global pointer variable declarations (with initializer)
(translation_unit
  (declaration
    declarator: (init_declarator
      declarator: (pointer_declarator
        declarator: (identifier) @name
      )
    )
  ) @signature @kind
)

; Comments
(comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/c_test.go

```go
import (
	"testing"
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_c "github.com/tree-sitter/tree-sitter-c/bindings/go"
)
func TestCQueryLanguage(t *testing.T)
func TestCQueryPattern(t *testing.T)
func TestCQueryExtractFunction(t *testing.T)
funcCaptures map[string]string
func TestCQueryExtractStruct(t *testing.T)
func TestCQueryExtractMacro(t *testing.T)
func TestCQueryExtractEnum(t *testing.T)
func TestCQueryExtractTypedef(t *testing.T)
func TestCQueryExtractGlobalVariables(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/call_query_test.go

```go
import (
	"testing"
	"github.com/indigo-net/Brf.it/pkg/parser"
	// Blank import to register tree-sitter parsers
	_ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"
)
func TestGoCallExtraction(t *testing.T)
func TestTypeScriptCallExtraction(t *testing.T)
func TestPythonCallExtraction(t *testing.T)
func TestCallExtractionDisabledByDefault(t *testing.T)
func TestCallExtractionTopLevel(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/cpp.go

```go
import (
	"strings"
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_cpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"
)
type CppQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewCppQuery() *CppQuery
func (q *CppQuery) Language() *sitter.Language
func (q *CppQuery) Query() []byte
cppKindMapping = map[string]string{
	"function_definition":  "function",
	"declaration":          "function",
	"struct_specifier":     "struct",
	"enum_specifier":       "enum",
	"type_definition":      "typedef",
	"preproc_function_def": "macro",
	"preproc_def":          "macro",
	"class_specifier":      "class",
	"field_declaration":    "method",
	"template_declaration": "template",
	"namespace_definition": "namespace",
}
func (q *CppQuery) KindMapping() map[string]string
func (q *CppQuery) ImportQuery() []byte
func (q *CppQuery) IsExported(name, sigText string) bool
cppImportQueryPattern = `
; #include directives (capture full statement)
(preproc_include) @import_path
`
cppQueryPattern = `
; Function definitions - direct declarator
(function_definition
  declarator: (function_declarator
    declarator: (identifier) @name
  )
) @signature @kind

; Function definitions - pointer return type
(function_definition
  declarator: (pointer_declarator
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Function definitions - reference return type
(function_definition
  declarator: (reference_declarator
    (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Function declarations (prototypes) - direct declarator
(declaration
  declarator: (function_declarator
    declarator: (identifier) @name
  )
) @signature @kind

; Function declarations (prototypes) - pointer return type
(declaration
  declarator: (pointer_declarator
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Class definitions
(class_specifier
  name: (type_identifier) @name
) @signature @kind

; Struct specifiers
(struct_specifier
  name: (type_identifier) @name
) @signature @kind

; Enum specifiers
(enum_specifier
  name: (type_identifier) @name
) @signature @kind

; Typedef
(type_definition
  declarator: (type_identifier) @name
) @signature @kind

; Function-like macros
(preproc_function_def
  name: (identifier) @name
) @signature @kind

; Object-like macros
(preproc_def
  name: (identifier) @name
) @signature @kind

; Method declarations in class (regular methods)
(field_declaration
  declarator: (function_declarator
    declarator: (field_identifier) @name
  )
) @signature @kind

; Method declarations with pointer return type
(field_declaration
  declarator: (pointer_declarator
    declarator: (function_declarator
      declarator: (field_identifier) @name
    )
  )
) @signature @kind

; Method declarations with reference return type
(field_declaration
  declarator: (reference_declarator
    (function_declarator
      declarator: (field_identifier) @name
    )
  )
) @signature @kind

; Constructor declarations (in class body)
(function_definition
  declarator: (function_declarator
    declarator: (qualified_identifier
      name: (identifier) @name
    )
  )
) @signature @kind

; Destructor definitions (outside class)
(function_definition
  declarator: (function_declarator
    declarator: (destructor_name
      (identifier) @name
    )
  )
) @signature @kind

; Destructor declarations in class (captured via declaration node)
(declaration
  declarator: (function_declarator
    declarator: (destructor_name
      (identifier) @name
    )
  )
) @signature @kind

; Namespace definitions
(namespace_definition
  name: (namespace_identifier) @name
) @signature @kind

; Template function definitions
(template_declaration
  (function_definition
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Template function definitions - pointer return type
(template_declaration
  (function_definition
    declarator: (pointer_declarator
      declarator: (function_declarator
        declarator: (identifier) @name
      )
    )
  )
) @signature @kind

; Template class definitions
(template_declaration
  (class_specifier
    name: (type_identifier) @name
  )
) @signature @kind

; Template struct definitions
(template_declaration
  (struct_specifier
    name: (type_identifier) @name
  )
) @signature @kind

; Template declarations (standalone)
(template_declaration
  (declaration
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Comments
(comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/cpp_test.go

```go
import (
	"testing"
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_cpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"
)
func TestCppQueryLanguage(t *testing.T)
func TestCppQueryPattern(t *testing.T)
func TestCppQueryExtractFunction(t *testing.T)
funcCaptures map[string]string
func TestCppQueryExtractClass(t *testing.T)
func TestCppQueryExtractMethod(t *testing.T)
func TestCppQueryExtractConstructorDestructor(t *testing.T)
func TestCppQueryExtractNamespace(t *testing.T)
func TestCppQueryExtractTemplate(t *testing.T)
func TestCppQueryExtractStruct(t *testing.T)
func TestCppQueryExtractEnum(t *testing.T)
func TestCppQueryExtractMacro(t *testing.T)
func TestCppQueryExtractTypedef(t *testing.T)
func TestCppQueryExtractIncludes(t *testing.T)
imports []string
func TestCppQueryNestedNamespaces(t *testing.T)
func TestCppQueryMultipleInheritance(t *testing.T)
func TestCppQueryEmptyFile(t *testing.T)
func TestCppQueryOnlyComments(t *testing.T)
nameCount int
docCount int
func TestCppQueryKindMapping(t *testing.T)
func TestCppQueryCaptures(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/csharp.go

```go
import (
	tree_sitter_c_sharp "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/csharp"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
type CSharpQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewCSharpQuery() *CSharpQuery
func (q *CSharpQuery) Language() *sitter.Language
func (q *CSharpQuery) Query() []byte
csharpKindMapping = map[string]string{
	"class_declaration":                 "class",
	"struct_declaration":                "struct",
	"interface_declaration":             "interface",
	"enum_declaration":                  "enum",
	"record_declaration":                "record",
	"delegate_declaration":              "type",
	"method_declaration":                "method",
	"constructor_declaration":           "constructor",
	"destructor_declaration":            "destructor",
	"property_declaration":              "variable",
	"field_declaration":                 "field",
	"event_declaration":                 "variable",
	"event_field_declaration":           "variable",
	"indexer_declaration":               "method",
	"operator_declaration":              "function",
	"conversion_operator_declaration":   "function",
	"namespace_declaration":             "namespace",
	"file_scoped_namespace_declaration": "namespace",
	"enum_member_declaration":           "variable",
}
func (q *CSharpQuery) KindMapping() map[string]string
func (q *CSharpQuery) ImportQuery() []byte
func (q *CSharpQuery) IsExported(name, sigText string) bool
csharpImportQueryPattern = `
; using directives (capture full declaration)
(using_directive) @import_path
`
csharpQueryPattern = `
; Class declarations
(class_declaration
  name: (identifier) @name
) @signature @kind

; Struct declarations
(struct_declaration
  name: (identifier) @name
) @signature @kind

; Interface declarations
(interface_declaration
  name: (identifier) @name
) @signature @kind

; Enum declarations
(enum_declaration
  name: (identifier) @name
) @signature @kind

; Record declarations (record, record class, record struct)
(record_declaration
  name: (identifier) @name
) @signature @kind

; Delegate declarations
(delegate_declaration
  name: (identifier) @name
) @signature @kind

; Method declarations
(method_declaration
  name: (identifier) @name
) @signature @kind

; Constructor declarations
(constructor_declaration
  name: (identifier) @name
) @signature @kind

; Destructor declarations
(destructor_declaration
  name: (identifier) @name
) @signature @kind

; Property declarations
(property_declaration
  name: (identifier) @name
) @signature @kind

; Field declarations (static/const filtered in parser.go)
(field_declaration
  (variable_declaration
    (variable_declarator
      name: (identifier) @name
    )
  )
) @signature @kind

; Event field declarations (e.g., public event EventHandler Changed;)
(event_field_declaration
  (variable_declaration
    (variable_declarator
      name: (identifier) @name
    )
  )
) @signature @kind

; Event declarations with accessor body
(event_declaration
  name: (identifier) @name
) @signature @kind

; Indexer declarations (no name capture — synthesized in parser.go)
(indexer_declaration) @signature @kind

; Operator declarations (no name capture — synthesized in parser.go)
(operator_declaration) @signature @kind

; Conversion operator declarations (no name capture — synthesized in parser.go)
(conversion_operator_declaration) @signature @kind

; Namespace declarations
(namespace_declaration
  name: (_) @name
) @signature @kind

; File-scoped namespace declarations (C# 10+)
(file_scoped_namespace_declaration
  name: (_) @name
) @signature @kind

; Enum member declarations
(enum_member_declaration
  name: (identifier) @name
) @signature @kind

; Comments (XML doc comments and regular)
(comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/csharp_test.go

```go
import (
	"testing"
	tree_sitter_c_sharp "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/csharp"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
func extractCSharpNames(t *testing.T, code []byte) map[string]bool
func TestCSharpQueryLanguage(t *testing.T)
func TestCSharpQueryPattern(t *testing.T)
func TestCSharpQueryImportPattern(t *testing.T)
func TestCSharpQueryExtractFunction(t *testing.T)
func TestCSharpQueryExtractTypes(t *testing.T)
func TestCSharpQueryExtractConstructorDestructor(t *testing.T)
func TestCSharpQueryExtractProperties(t *testing.T)
func TestCSharpQueryExtractFields(t *testing.T)
func TestCSharpQueryExtractEvents(t *testing.T)
func TestCSharpQueryExtractOperators(t *testing.T)
foundOperator, foundConversion, foundIndexer bool
func TestCSharpQueryExtractImport(t *testing.T)
func TestCSharpQueryExtractGenerics(t *testing.T)
func TestCSharpQueryKindMapping(t *testing.T)
func TestCSharpQueryCaptures(t *testing.T)
func TestCSharpQueryExtractNamespace(t *testing.T)
func TestCSharpQueryExtractRecords(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/elixir.go

```go
import (
	"strings"
	tree_sitter_elixir "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/elixir"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
type ElixirQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewElixirQuery() *ElixirQuery
func (q *ElixirQuery) Language() *sitter.Language
func (q *ElixirQuery) Query() []byte
elixirKindMapping = map[string]string{
	"call":           "function",
	"unary_operator": "type",
}
func (q *ElixirQuery) KindMapping() map[string]string
func (q *ElixirQuery) ImportQuery() []byte
func (q *ElixirQuery) IsExported(name, sigText string) bool
elixirImportQueryPattern = `
; import statements: import Module
(call
  target: (identifier)
  (arguments
    (alias))) @import_path

; import with options: import Module, only: [...]
(call
  target: (identifier)
  (arguments
    (alias)
    (keywords))) @import_path
`
elixirQueryPattern = `
; Module/protocol/impl definitions: defmodule MyModule do...end
(call
  target: (identifier)
  (arguments
    (alias) @name)
  (do_block)
) @signature @kind

; defimpl with keyword options: defimpl Protocol, for: Module do...end
(call
  target: (identifier)
  (arguments
    (alias) @name
    (keywords))
  (do_block)
) @signature @kind

; Function/macro definitions with arguments: def foo(args) do...end
(call
  target: (identifier)
  (arguments
    (call
      target: (identifier) @name))
  (do_block)
) @signature @kind

; Function/macro definitions with guard clause: def foo(args) when guard do...end
(call
  target: (identifier)
  (arguments
    (binary_operator
      left: (call
        target: (identifier) @name)
      operator: "when"))
  (do_block)
) @signature @kind

; Zero-arity function definitions: def foo do...end
(call
  target: (identifier)
  (arguments
    (identifier) @name)
  (do_block)
) @signature @kind

; Zero-arity function with guard: def foo when guard do...end
(call
  target: (identifier)
  (arguments
    (binary_operator
      left: (identifier) @name
      operator: "when"))
  (do_block)
) @signature @kind

; Guard definitions without do_block: defguard is_positive(x) when ...
(call
  target: (identifier)
  (arguments
    (binary_operator
      left: (call
        target: (identifier) @name)
      operator: "when"))
) @signature @kind

; defdelegate: defdelegate foo(args), to: Bar
(call
  target: (identifier)
  (arguments
    (call
      target: (identifier) @name)
    (keywords))
) @signature @kind

; defstruct with list: defstruct [:field1, :field2]
(call
  target: (identifier) @name
  (arguments
    (list))
) @signature @kind

; defstruct with keywords: defstruct field: default_value
(call
  target: (identifier) @name
  (arguments
    (keywords))
) @signature @kind

; Module attributes: @spec, @type, @typep, @opaque, @callback
(unary_operator
  operator: "@"
  operand: (call
    target: (identifier) @name)
) @signature @kind

; Line comments
(comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/elixir_test.go

```go
import (
	"strings"
	"testing"
	tree_sitter_elixir "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/elixir"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
func extractElixirNames(t *testing.T, code []byte) map[string]bool
func extractElixirSignatures(t *testing.T, code []byte) []string
sigs []string
func extractElixirImportNames(t *testing.T, code []byte) []string
imports []string
func TestElixirQueryLanguage(t *testing.T)
func TestElixirQueryPattern(t *testing.T)
func TestElixirQueryImportPattern(t *testing.T)
func TestElixirQueryExtractFunction(t *testing.T)
func TestElixirQueryExtractModule(t *testing.T)
func TestElixirQueryExtractProtocol(t *testing.T)
func TestElixirQueryExtractMacro(t *testing.T)
func TestElixirQueryExtractGuard(t *testing.T)
func TestElixirQueryExtractDelegate(t *testing.T)
func TestElixirQueryExtractStruct(t *testing.T)
func TestElixirQueryExtractTypeSpec(t *testing.T)
func TestElixirQueryExtractImport(t *testing.T)
func TestElixirQueryExtractZeroArityWithGuard(t *testing.T)
func TestElixirQueryKindMapping(t *testing.T)
func TestElixirQueryCaptures(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/go.go

```go
import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)
captureName      = "name"
captureSignature = "signature"
captureDoc       = "doc"
captureKind      = "kind"
type GoQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewGoQuery() *GoQuery
func (q *GoQuery) Language() *sitter.Language
func (q *GoQuery) Query() []byte
goKindMapping = map[string]string{
	"function_declaration": "function",
	"method_declaration":   "method",
	"type_declaration":     "type",
	"const_declaration":    "variable",
	"var_declaration":      "variable",
	"const_spec":           "variable",
	"var_spec":             "variable",
}
func (q *GoQuery) KindMapping() map[string]string
func (q *GoQuery) ImportQuery() []byte
func (q *GoQuery) CallQuery() []byte
func (q *GoQuery) IsExported(name, _ string) bool
goCallQueryPattern = `
; Direct function calls (e.g., foo())
(call_expression
  function: (identifier) @callee
) @call_node

; Method/package calls (e.g., obj.Method(), fmt.Println())
(call_expression
  function: (selector_expression
    field: (field_identifier) @callee
  )
) @call_node
`
goImportQueryPattern = `
; Import declarations (capture entire declaration)
(import_declaration) @import_path
`
goQueryPattern = `
; Function declarations
(function_declaration
  name: (identifier) @name
) @signature @kind

; Method declarations
(method_declaration
  name: (field_identifier) @name
) @signature @kind

; Type declarations (struct, interface, etc.)
(type_declaration
  (type_spec
    name: (type_identifier) @name
  )
) @signature @kind

; Package-level const specs (captures each const individually)
(const_spec
  name: (identifier) @name
) @signature @kind

; Package-level var specs (captures each var individually)
(var_spec
  name: (identifier) @name
) @signature @kind

; Comments (documentation)
(comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/go_test.go

```go
import (
	"testing"
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)
func TestGoQueryLanguage(t *testing.T)
func TestGoQueryPattern(t *testing.T)
func TestGoQueryExtractFunction(t *testing.T)
funcCaptures map[string]string
funcKindNode *sitter.Node
kindNode *sitter.Node
func TestGoQueryExtractConstAndVar(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/java.go

```go
import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_java "github.com/tree-sitter/tree-sitter-java/bindings/go"
)
type JavaQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewJavaQuery() *JavaQuery
func (q *JavaQuery) Language() *sitter.Language
func (q *JavaQuery) Query() []byte
javaKindMapping = map[string]string{
	"class_declaration":           "class",
	"interface_declaration":       "interface",
	"method_declaration":          "method",
	"constructor_declaration":     "constructor",
	"enum_declaration":            "enum",
	"annotation_type_declaration": "annotation",
	"record_declaration":          "record",
	"field_declaration":           "field",
}
func (q *JavaQuery) KindMapping() map[string]string
func (q *JavaQuery) ImportQuery() []byte
func (q *JavaQuery) CallQuery() []byte
func (q *JavaQuery) IsExported(name, sigText string) bool
javaCallQueryPattern = `
; Method invocations (e.g., obj.method(), method())
(method_invocation
  name: (identifier) @callee
) @call_node
`
javaImportQueryPattern = `
; import statements (capture full declaration)
(import_declaration) @import_path
`
javaQueryPattern = `
; Class declarations (includes inner classes)
(class_declaration
  name: (identifier) @name
) @signature @kind

; Interface declarations
(interface_declaration
  name: (identifier) @name
) @signature @kind

; Method declarations
(method_declaration
  name: (identifier) @name
) @signature @kind

; Constructor declarations
(constructor_declaration
  name: (identifier) @name
) @signature @kind

; Enum declarations
(enum_declaration
  name: (identifier) @name
) @signature @kind

; Annotation type declarations (@interface)
(annotation_type_declaration
  name: (identifier) @name
) @signature @kind

; Record declarations (Java 14+)
(record_declaration
  name: (identifier) @name
) @signature @kind

; Field declarations (static fields filtered in parser.go)
(field_declaration
  (variable_declarator
    name: (identifier) @name
  )
) @signature @kind

; Comments (Javadoc and regular)
(line_comment) @doc
(block_comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/java_test.go

```go
import (
	"testing"
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_java "github.com/tree-sitter/tree-sitter-java/bindings/go"
)
func TestJavaQueryLanguage(t *testing.T)
func TestJavaQueryPattern(t *testing.T)
func TestJavaQueryKindMapping(t *testing.T)
func TestJavaQueryExtractClass(t *testing.T)
foundClass, foundMethod bool
func TestJavaQueryExtractInterface(t *testing.T)
func TestJavaQueryExtractEnum(t *testing.T)
foundEnum bool
func TestJavaQueryExtractAnnotationType(t *testing.T)
foundAnnotation bool
func TestJavaQueryExtractRecord(t *testing.T)
func TestJavaQueryExtractGenerics(t *testing.T)
func TestJavaQueryExtractFieldDeclarations(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/kotlin.go

```go
import (
	tree_sitter_kotlin "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/kotlin"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
type KotlinQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewKotlinQuery() *KotlinQuery
func (q *KotlinQuery) Language() *sitter.Language
func (q *KotlinQuery) Query() []byte
kotlinKindMapping = map[string]string{
	"function_declaration":  "function",
	"class_declaration":     "class",
	"object_declaration":    "class",
	"companion_object":      "class",
	"property_declaration":  "variable",
	"type_alias":            "type",
	"enum_entry":            "variable",
	"secondary_constructor": "constructor",
}
func (q *KotlinQuery) KindMapping() map[string]string
func (q *KotlinQuery) ImportQuery() []byte
func (q *KotlinQuery) IsExported(name, sigText string) bool
kotlinImportQueryPattern = `
; Import statements
(import_header) @import_path
`
kotlinQueryPattern = `
; Function declarations (regular, suspend, inline, extension, operator, infix, tailrec)
(function_declaration
  (simple_identifier) @name
) @signature @kind

; Class declarations (class, data class, sealed class, enum class, interface, annotation class, value class)
(class_declaration
  (type_identifier) @name
) @signature @kind

; Object declarations (singleton)
(object_declaration
  (type_identifier) @name
) @signature @kind

; Companion object with explicit name (e.g., companion object Factory)
(companion_object
  (type_identifier) @name
) @signature @kind

; Property declarations (val/var, const val, lateinit, delegated)
(property_declaration
  (variable_declaration
    (simple_identifier) @name
  )
) @signature @kind

; Type alias
(type_alias
  (type_identifier) @name
) @signature @kind

; Enum entries
(enum_entry
  (simple_identifier) @name
) @signature @kind

; Secondary constructors
(secondary_constructor) @signature @kind

; Line comments
(line_comment) @doc

; Block/multiline comments
(multiline_comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/kotlin_test.go

```go
import (
	"testing"
	tree_sitter_kotlin "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/kotlin"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
func extractKotlinNames(t *testing.T, code []byte) map[string]bool
func TestKotlinQueryLanguage(t *testing.T)
func TestKotlinQueryPattern(t *testing.T)
func TestKotlinQueryImportPattern(t *testing.T)
func TestKotlinQueryExtractFunction(t *testing.T)
func TestKotlinQueryExtractTypes(t *testing.T)
func TestKotlinQueryExtractInterface(t *testing.T)
func TestKotlinQueryExtractObject(t *testing.T)
func TestKotlinQueryExtractProperties(t *testing.T)
func TestKotlinQueryExtractTypeAlias(t *testing.T)
func TestKotlinQueryExtractEnumEntry(t *testing.T)
func TestKotlinQueryExtractImport(t *testing.T)
func TestKotlinQueryExtractGenerics(t *testing.T)
func TestKotlinQueryKindMapping(t *testing.T)
func TestKotlinQueryCaptures(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/lua.go

```go
import (
	tree_sitter_lua "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/lua"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
type LuaQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewLuaQuery() *LuaQuery
func (q *LuaQuery) Language() *sitter.Language
func (q *LuaQuery) Query() []byte
luaKindMapping = map[string]string{
	"function_declaration": "function",
	"variable_declaration": "variable",
	"assignment_statement": "variable",
}
func (q *LuaQuery) KindMapping() map[string]string
func (q *LuaQuery) ImportQuery() []byte
luaImportQueryPattern = `
; local json = require("json")
(variable_declaration
  (assignment_statement
    (expression_list
      value: (function_call
        name: (identifier) @_fn
        arguments: (arguments (string))
      )
    )
  )
) @import_path
(#eq? @_fn "require")
`
luaQueryPattern = `
; Function declarations (global, local, module, method)
; Covers: function foo(), local function foo(), function M.foo(), function M:foo()
(function_declaration
  name: [
    (identifier) @name
    (dot_index_expression field: (identifier) @name)
    (method_index_expression method: (identifier) @name)
  ]
) @signature @kind

; Variable declarations with function assignment: local foo = function() end
(variable_declaration
  (assignment_statement
    (variable_list
      name: (identifier) @name)
    (expression_list
      value: (function_definition))
  )
) @signature @kind

; Variable declarations with table constructor: local M = {}
(variable_declaration
  (assignment_statement
    (variable_list
      name: (identifier) @name)
    (expression_list
      value: (table_constructor))
  )
) @signature @kind

; Comments (LuaDoc --- and regular --)
(comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/lua_test.go

```go
import (
	"testing"
	tree_sitter_lua "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/lua"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
func extractLuaNames(t *testing.T, code []byte) map[string]bool
func TestLuaQueryLanguage(t *testing.T)
func TestLuaQueryPattern(t *testing.T)
func TestLuaQueryImportPattern(t *testing.T)
func TestLuaQueryExtractFunction(t *testing.T)
func TestLuaQueryExtractLocalFunction(t *testing.T)
func TestLuaQueryExtractModuleFunction(t *testing.T)
func TestLuaQueryExtractMethod(t *testing.T)
func TestLuaQueryExtractTableAssignment(t *testing.T)
func TestLuaQueryExtractFunctionAssignment(t *testing.T)
func TestLuaQueryExtractImport(t *testing.T)
func TestLuaQueryNonRequireFalsePositive(t *testing.T)
func TestLuaQueryExtractDoc(t *testing.T)
func TestLuaQueryKindMapping(t *testing.T)
func TestLuaQueryCaptures(t *testing.T)
func TestLuaQueryExtractMixed(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/php.go

```go
import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_php "github.com/tree-sitter/tree-sitter-php/bindings/go"
)
type PHPQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewPHPQuery() *PHPQuery
func (q *PHPQuery) Language() *sitter.Language
func (q *PHPQuery) Query() []byte
phpKindMapping = map[string]string{
	"function_definition":       "function",
	"method_declaration":        "method",
	"class_declaration":         "class",
	"interface_declaration":     "interface",
	"trait_declaration":         "type",
	"enum_declaration":          "enum",
	"const_declaration":         "variable",
	"property_declaration":      "variable",
	"namespace_use_declaration": "import",
}
func (q *PHPQuery) KindMapping() map[string]string
func (q *PHPQuery) ImportQuery() []byte
func (q *PHPQuery) IsExported(name, sigText string) bool
phpImportQueryPattern = `
; use Namespace\\Class;
(namespace_use_declaration) @import_path

; include 'file.php';
(include_expression) @import_path

; require 'vendor/autoload.php';
(require_expression) @import_path

; include_once 'config.php';
(include_once_expression) @import_path

; require_once 'config.php';
(require_once_expression) @import_path
`
phpQueryPattern = `
; Function definitions: function name() {}
(function_definition name: (name) @name) @signature @kind

; Method declarations in classes
(method_declaration name: (name) @name) @signature @kind

; Class declarations
(class_declaration name: (name) @name) @signature @kind

; Interface declarations
(interface_declaration name: (name) @name) @signature @kind

; Trait declarations
(trait_declaration name: (name) @name) @signature @kind

; Enum declarations
(enum_declaration name: (name) @name) @signature @kind

; Const declarations: const NAME = value;
(const_declaration
  (const_element
    (name) @name
  )
) @signature @kind

; Property declarations: public $name;
(property_declaration
  (property_element
    (variable_name (name) @name)
  )
) @signature @kind

; Comments (PHPDoc and regular)
(comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/php_test.go

```go
import (
	"testing"
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_php "github.com/tree-sitter/tree-sitter-php/bindings/go"
)
func extractPHPNames(t *testing.T, code []byte) map[string]bool
func TestPHPQueryLanguage(t *testing.T)
func TestPHPQueryPattern(t *testing.T)
func TestPHPQueryImportPattern(t *testing.T)
func TestPHPQueryExtractFunction(t *testing.T)
func TestPHPQueryExtractClass(t *testing.T)
func TestPHPQueryExtractInterface(t *testing.T)
func TestPHPQueryExtractTrait(t *testing.T)
func TestPHPQueryExtractEnum(t *testing.T)
func TestPHPQueryExtractVariable(t *testing.T)
func TestPHPQueryExtractImport(t *testing.T)
func TestPHPQueryKindMapping(t *testing.T)
func TestPHPQueryCaptures(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/python.go

```go
import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"
)
type PythonQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewPythonQuery() *PythonQuery
func (q *PythonQuery) Language() *sitter.Language
func (q *PythonQuery) Query() []byte
pythonKindMapping = map[string]string{
	"function_definition":  "function",
	"class_definition":     "class",
	"expression_statement": "variable",
	"assignment":           "variable",
}
func (q *PythonQuery) KindMapping() map[string]string
func (q *PythonQuery) ImportQuery() []byte
func (q *PythonQuery) CallQuery() []byte
func (q *PythonQuery) IsExported(name, _ string) bool
pythonCallQueryPattern = `
; Direct function calls (e.g., foo())
(call
  function: (identifier) @callee
) @call_node

; Method/attribute calls (e.g., obj.method())
(call
  function: (attribute
    attribute: (identifier) @callee
  )
) @call_node
`
pythonImportQueryPattern = `
; import module (capture full statement)
(import_statement) @import_path

; from module import ... (capture full statement)
(import_from_statement) @import_path
`
pythonQueryPattern = `
; Function definitions (includes async def, methods)
(function_definition
  name: (identifier) @name
) @signature @kind

; Class definitions
(class_definition
  name: (identifier) @name
) @signature @kind

; Module-level assignments (simple and with type annotations)
(module
  (expression_statement
    (assignment
      left: (identifier) @name
    )
  ) @signature @kind
)

; Comments
(comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/python_test.go

```go
import (
	"testing"
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"
)
func TestPythonQueryLanguage(t *testing.T)
func TestPythonQueryPattern(t *testing.T)
func TestPythonQueryExtractFunction(t *testing.T)
funcCaptures map[string]string
func TestPythonQueryExtractClass(t *testing.T)
func TestPythonQueryExtractAsyncFunction(t *testing.T)
funcCaptures map[string]string
func TestPythonQueryExtractModuleLevelVariables(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/ruby.go

```go
import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_ruby "github.com/tree-sitter/tree-sitter-ruby/bindings/go"
)
type RubyQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewRubyQuery() *RubyQuery
func (q *RubyQuery) Language() *sitter.Language
func (q *RubyQuery) Query() []byte
rubyKindMapping = map[string]string{
	"method":           "method",
	"singleton_method": "method",
	"class":            "class",
	"module":           "namespace",
	"assignment":       "variable",
}
func (q *RubyQuery) KindMapping() map[string]string
func (q *RubyQuery) ImportQuery() []byte
rubyImportQueryPattern = `
; require "library" / require_relative "library"
(call
  method: (identifier) @_fn
  arguments: (argument_list
    (string)
  )
) @import_path
(#match? @_fn "^require")
`
rubyQueryPattern = `
; Instance methods and top-level functions: def foo(args) ... end
(method
  name: (identifier) @name
) @signature @kind

; Class methods: def self.foo(args) ... end
(singleton_method
  name: (identifier) @name
) @signature @kind

; Class definitions: class Foo ... end
(class
  name: (constant) @name
) @signature @kind

; Module definitions: module Foo ... end
(module
  name: (constant) @name
) @signature @kind

; Top-level constant assignments: FOO = value
(program
  (assignment
    left: (constant) @name
  ) @signature @kind
)

; Comments
(comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/ruby_test.go

```go
import (
	"testing"
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_ruby "github.com/tree-sitter/tree-sitter-ruby/bindings/go"
)
func extractRubyNames(t *testing.T, code []byte) map[string]bool
func extractRubyImports(t *testing.T, code []byte) []string
imports []string
func TestRubyQueryLanguage(t *testing.T)
func TestRubyQueryPattern(t *testing.T)
func TestRubyQueryImportPattern(t *testing.T)
func TestRubyQueryExtractFunction(t *testing.T)
func TestRubyQueryExtractTypes(t *testing.T)
func TestRubyQueryExtractClassMethods(t *testing.T)
func TestRubyQueryExtractModuleMethods(t *testing.T)
func TestRubyQueryExtractConstants(t *testing.T)
func TestRubyQueryExtractAttrAccessors(t *testing.T)
func TestRubyQueryExtractImport(t *testing.T)
func TestRubyQueryExtractNestedClasses(t *testing.T)
func TestRubyQueryKindMapping(t *testing.T)
func TestRubyQueryCaptures(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/rust.go

```go
import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_rust "github.com/tree-sitter/tree-sitter-rust/bindings/go"
)
type RustQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewRustQuery() *RustQuery
func (q *RustQuery) Language() *sitter.Language
func (q *RustQuery) Query() []byte
rustKindMapping = map[string]string{
	"function_item":           "function",
	"struct_item":             "struct",
	"enum_item":               "enum",
	"trait_item":              "trait",
	"type_item":               "type",
	"impl_item":               "impl",
	"const_item":              "variable",
	"static_item":             "variable",
	"mod_item":                "namespace",
	"macro_definition":        "macro",
	"foreign_mod_item":        "namespace",
	"union_item":              "struct",
	"associated_type":         "type",
	"function_signature_item": "function",
}
func (q *RustQuery) KindMapping() map[string]string
func (q *RustQuery) ImportQuery() []byte
func (q *RustQuery) CallQuery() []byte
func (q *RustQuery) IsExported(name, _ string) bool
rustCallQueryPattern = `
; Direct function calls (e.g., foo())
(call_expression
  function: (identifier) @callee
) @call_node

; Method/field calls (e.g., obj.method())
(call_expression
  function: (field_expression
    field: (field_identifier) @callee
  )
) @call_node
`
rustImportQueryPattern = `
; Use declarations (capture full statement)
(use_declaration) @import_path

; Extern crate declarations
(extern_crate_declaration) @import_path
`
rustQueryPattern = `
; Functions (including async, unsafe, const, extern)
(function_item
  name: (identifier) @name
) @signature @kind

; Struct declarations
(struct_item
  name: (type_identifier) @name
) @signature @kind

; Enum declarations
(enum_item
  name: (type_identifier) @name
) @signature @kind

; Trait declarations
(trait_item
  name: (type_identifier) @name
) @signature @kind

; Type aliases
(type_item
  name: (type_identifier) @name
) @signature @kind

; Impl blocks (capture the whole impl signature)
(impl_item
  type: (type_identifier) @name
) @signature @kind

; Impl blocks for generic types
(impl_item
  type: (generic_type
    type: (type_identifier) @name
  )
) @signature @kind

; Trait impl blocks (impl Trait for Type)
(impl_item
  trait: (type_identifier)
  type: (type_identifier) @name
) @signature @kind

; Constants
(const_item
  name: (identifier) @name
) @signature @kind

; Statics
(static_item
  name: (identifier) @name
) @signature @kind

; Modules
(mod_item
  name: (identifier) @name
) @signature @kind

; Macro definitions (macro_rules!)
(macro_definition
  name: (identifier) @name
) @signature @kind

; Union declarations
(union_item
  name: (type_identifier) @name
) @signature @kind

; Foreign mod (extern "C" blocks)
(foreign_mod_item) @signature @kind

; Associated types in traits
(associated_type
  name: (type_identifier) @name
) @signature @kind

; Function signatures in traits (without body)
(function_signature_item
  name: (identifier) @name
) @signature @kind

; Doc comments (/// and //!)
(line_comment) @doc

; Block doc comments (/** and /*!)
(block_comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/rust_test.go

```go
import (
	"testing"
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_rust "github.com/tree-sitter/tree-sitter-rust/bindings/go"
)
func TestRustQueryLanguage(t *testing.T)
func TestRustQueryPattern(t *testing.T)
func TestRustQueryImportPattern(t *testing.T)
func TestRustQueryExtractFunction(t *testing.T)
func TestRustQueryExtractTypes(t *testing.T)
func TestRustQueryExtractImplAndMethods(t *testing.T)
func TestRustQueryExtractConstAndStatic(t *testing.T)
func TestRustQueryExtractMacro(t *testing.T)
func TestRustQueryExtractModule(t *testing.T)
func TestRustQueryExtractUse(t *testing.T)
func TestRustQueryExtractGenericsAndLifetimes(t *testing.T)
func TestRustQueryEmptyFile(t *testing.T)
func TestRustQueryKindMapping(t *testing.T)
func TestRustQueryCaptures(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/scala.go

```go
import (
	tree_sitter_scala "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/scala"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
type ScalaQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewScalaQuery() *ScalaQuery
func (q *ScalaQuery) Language() *sitter.Language
func (q *ScalaQuery) Query() []byte
scalaKindMapping = map[string]string{
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
func (q *ScalaQuery) KindMapping() map[string]string
func (q *ScalaQuery) ImportQuery() []byte
func (q *ScalaQuery) IsExported(name, sigText string) bool
scalaImportQueryPattern = `
; Import statements
(import_declaration) @import_path
`
scalaQueryPattern = `
; Function definitions (def with body)
(function_definition
  name: (identifier) @name
) @signature @kind

; Function declarations (abstract methods in traits/classes, no body)
(function_declaration
  name: (identifier) @name
) @signature @kind

; Class definitions (class, abstract class, case class, sealed class, implicit class)
(class_definition
  name: (identifier) @name
) @signature @kind

; Trait definitions (trait, sealed trait)
(trait_definition
  name: (identifier) @name
) @signature @kind

; Object definitions (singleton, companion)
(object_definition
  name: (identifier) @name
) @signature @kind

; Val definitions (val, lazy val, implicit val)
(val_definition
  pattern: (identifier) @name
) @signature @kind

; Val declarations (abstract val in traits)
(val_declaration
  name: (identifier) @name
) @signature @kind

; Var definitions
(var_definition
  pattern: (identifier) @name
) @signature @kind

; Var declarations (abstract var in traits)
(var_declaration
  name: (identifier) @name
) @signature @kind

; Type aliases
(type_definition
  name: (type_identifier) @name
) @signature @kind

; Enum definitions (Scala 3)
(enum_definition
  name: (identifier) @name
) @signature @kind

; Given definitions (Scala 3)
(given_definition
  name: (identifier) @name
) @signature @kind

; Line comments
(comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/scala_test.go

```go
import (
	"testing"
	"unsafe"
	tree_sitter_scala "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/scala"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
func extractScalaNames(t *testing.T, code []byte) map[string]bool
func extractScalaImports(t *testing.T, code []byte) []string
imports []string
func TestScalaQueryLanguage(t *testing.T)
func TestScalaQueryPattern(t *testing.T)
func TestScalaQueryImportPattern(t *testing.T)
func TestScalaQueryExtractFunction(t *testing.T)
func TestScalaQueryExtractTypes(t *testing.T)
func TestScalaQueryExtractClassBody(t *testing.T)
func TestScalaQueryExtractTraitMembers(t *testing.T)
func TestScalaQueryExtractObjectMembers(t *testing.T)
func TestScalaQueryExtractValVar(t *testing.T)
func TestScalaQueryExtractTypeAlias(t *testing.T)
func TestScalaQueryExtractEnum(t *testing.T)
func TestScalaQueryExtractCaseClass(t *testing.T)
func TestScalaQueryExtractExtension(t *testing.T)
func TestScalaQueryExtractImport(t *testing.T)
func TestScalaQueryExtractGenerics(t *testing.T)
func TestScalaQueryKindMapping(t *testing.T)
func TestScalaQueryCaptures(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/shell.go

```go
import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_bash "github.com/tree-sitter/tree-sitter-bash/bindings/go"
)
type ShellQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewShellQuery() *ShellQuery
func (q *ShellQuery) Language() *sitter.Language
func (q *ShellQuery) Query() []byte
shellKindMapping = map[string]string{
	"function_definition": "function",
	"variable_assignment": "variable",
}
func (q *ShellQuery) KindMapping() map[string]string
func (q *ShellQuery) ImportQuery() []byte
shellImportQueryPattern = `
; source /path/to/file and . /path/to/file
; Capture command nodes. Go-side filtering will check if command name is "source" or "."
(command
  name: (command_name) @name
) @import_path
`
shellQueryPattern = `
; Function definitions: function foo { } or function foo() { } or foo() { }
(function_definition
  name: (word) @name
) @signature @kind

; Variable assignments: FOO=bar, FOO="bar", FOO=$(cmd)
(variable_assignment
  name: (variable_name) @name
) @signature @kind

; Comments (# ...)
(comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/shell_test.go

```go
import (
	"testing"
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_bash "github.com/tree-sitter/tree-sitter-bash/bindings/go"
)
func extractShellNames(t *testing.T, code []byte) map[string]bool
func TestShellQueryLanguage(t *testing.T)
func TestShellQueryPattern(t *testing.T)
func TestShellQueryImportPattern(t *testing.T)
func TestShellQueryExtractFunction(t *testing.T)
func TestShellQueryExtractFunctionWithoutKeyword(t *testing.T)
func TestShellQueryExtractVariable(t *testing.T)
func TestShellQueryExtractImport(t *testing.T)
func TestShellQueryExtractDoc(t *testing.T)
func TestShellQueryKindMapping(t *testing.T)
func TestShellQueryCaptures(t *testing.T)
func TestShellQueryExtractMixed(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/sql.go

```go
import (
	tree_sitter_sql "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/sql"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
type SQLQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewSQLQuery() *SQLQuery
func (q *SQLQuery) Language() *sitter.Language
func (q *SQLQuery) Query() []byte
sqlKindMapping = map[string]string{
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
func (q *SQLQuery) KindMapping() map[string]string
func (q *SQLQuery) ImportQuery() []byte
sqlQueryPattern = `
; CREATE TABLE
(create_table
  (object_reference) @name) @signature @kind

; CREATE FUNCTION
(create_function
  (object_reference) @name) @signature @kind

; CREATE VIEW
(create_view
  (object_reference) @name) @signature @kind

; CREATE MATERIALIZED VIEW
(create_materialized_view
  (object_reference) @name) @signature @kind

; CREATE INDEX (name extracted Go-side)
(create_index) @signature @kind

; CREATE TYPE
(create_type
  (object_reference) @name) @signature @kind

; CREATE TRIGGER (first object_reference = trigger name)
(create_trigger
  (object_reference) @name) @signature @kind

; CREATE SCHEMA (bare identifier)
(create_schema
  (identifier) @name) @signature @kind

; CREATE SEQUENCE
(create_sequence
  (object_reference) @name) @signature @kind

; ALTER TABLE
(alter_table
  (object_reference) @name) @signature @kind

; SQL comments (-- single-line)
(comment) @doc

; SQL multi-line comments (/* ... */ are "marginalia" in tree-sitter-sql)
(marginalia) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/sql_test.go

```go
import (
	"testing"
	"unsafe"
	tree_sitter_sql "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/sql"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
func extractSQLNames(t *testing.T, code string) []string
names []string
type sqlCapture struct {
	Name      string
	Signature string
	Kind      string
}
func extractSQLCaptures(t *testing.T, code string) []sqlCapture
captures []sqlCapture
sc sqlCapture
func TestSQLQueryLanguage(t *testing.T)
func TestSQLQueryPattern(t *testing.T)
func TestSQLQueryImportPattern(t *testing.T)
func TestSQLQueryExtractCreateTable(t *testing.T)
func TestSQLQueryExtractCreateFunction(t *testing.T)
func TestSQLQueryExtractCreateView(t *testing.T)
func TestSQLQueryExtractCreateIndex(t *testing.T)
func TestSQLQueryExtractCreateType(t *testing.T)
func TestSQLQueryExtractCreateTrigger(t *testing.T)
func TestSQLQueryExtractCreateSchema(t *testing.T)
func TestSQLQueryExtractMaterializedView(t *testing.T)
func TestSQLQueryExtractAlterTable(t *testing.T)
func TestSQLQueryExtractCreateSequence(t *testing.T)
func TestSQLQueryExtractComments(t *testing.T)
func TestSQLQueryKindMapping(t *testing.T)
func TestSQLQueryCaptures(t *testing.T)
func TestSQLQuerySchemaQualifiedName(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/swift.go

```go
import (
	tree_sitter_swift "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/swift"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
type SwiftQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewSwiftQuery() *SwiftQuery
func (q *SwiftQuery) Language() *sitter.Language
func (q *SwiftQuery) Query() []byte
swiftKindMapping = map[string]string{
	"function_declaration":          "function",
	"class_declaration":             "class",
	"protocol_declaration":          "interface",
	"typealias_declaration":         "type",
	"property_declaration":          "variable",
	"init_declaration":              "constructor",
	"deinit_declaration":            "destructor",
	"subscript_declaration":         "method",
	"operator_declaration":          "function",
	"protocol_function_declaration": "function",
}
func (q *SwiftQuery) KindMapping() map[string]string
func (q *SwiftQuery) ImportQuery() []byte
func (q *SwiftQuery) IsExported(name, sigText string) bool
swiftImportQueryPattern = `
; Import declarations (capture full statement)
(import_declaration) @import_path
`
swiftQueryPattern = `
; Functions
(function_declaration
  name: (simple_identifier) @name
) @signature @kind

; Classes, Structs, Enums (all use class_declaration node type)
(class_declaration
  name: (type_identifier) @name
) @signature @kind

; Extensions (name is in user_type child)
(class_declaration
  name: (user_type
    (type_identifier) @name
  )
) @signature @kind

; Protocol declarations
(protocol_declaration
  name: (type_identifier) @name
) @signature @kind

; Type aliases
(typealias_declaration
  name: (type_identifier) @name
) @signature @kind

; Properties (let/var)
(property_declaration
  name: (pattern
    (simple_identifier) @name
  )
) @signature @kind

; Initializers
(init_declaration) @signature @kind

; Deinitializers
(deinit_declaration) @signature @kind

; Subscript declarations
(subscript_declaration) @signature @kind

; Operator declarations
(operator_declaration
  (custom_operator) @name
) @signature @kind

; Protocol function declarations (methods in protocol body)
(protocol_function_declaration
  name: (simple_identifier) @name
) @signature @kind

; Doc comments (/// style)
(comment) @doc

; Multiline comments (/** style)
(multiline_comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/swift_test.go

```go
import (
	"testing"
	tree_sitter_swift "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/swift"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
func TestSwiftQueryLanguage(t *testing.T)
func TestSwiftQueryPattern(t *testing.T)
func TestSwiftQueryImportPattern(t *testing.T)
func TestSwiftQueryExtractFunction(t *testing.T)
func TestSwiftQueryExtractTypes(t *testing.T)
func TestSwiftQueryExtractProtocol(t *testing.T)
func TestSwiftQueryExtractExtension(t *testing.T)
func TestSwiftQueryExtractProperties(t *testing.T)
func TestSwiftQueryExtractInitDeinit(t *testing.T)
func TestSwiftQueryExtractSubscript(t *testing.T)
func TestSwiftQueryExtractImport(t *testing.T)
func TestSwiftQueryExtractGenerics(t *testing.T)
func TestSwiftQueryKindMapping(t *testing.T)
func TestSwiftQueryCaptures(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/toml.go

```go
import (
	tree_sitter_toml "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/toml"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
type TOMLQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewTOMLQuery() *TOMLQuery
func (q *TOMLQuery) Language() *sitter.Language
func (q *TOMLQuery) Query() []byte
tomlKindMapping = map[string]string{
	"table":               "namespace",
	"table_array_element": "namespace",
	"pair":                "variable",
}
func (q *TOMLQuery) KindMapping() map[string]string
func (q *TOMLQuery) ImportQuery() []byte
tomlQueryPattern = `
; Standard table sections [section]
(table
  (bare_key) @name) @signature @kind

; Standard table sections with quoted key ["section"]
(table
  (quoted_key) @name) @signature @kind

; Standard table sections with dotted key [section.subsection]
(table
  (dotted_key) @name) @signature @kind

; Array of tables [[section]]
(table_array_element
  (bare_key) @name) @signature @kind

; Array of tables with quoted key [["section"]]
(table_array_element
  (quoted_key) @name) @signature @kind

; Array of tables with dotted key [[section.subsection]]
(table_array_element
  (dotted_key) @name) @signature @kind

; Key-value pairs (bare key)
(pair
  (bare_key) @name) @signature @kind

; Key-value pairs (quoted key)
(pair
  (quoted_key) @name) @signature @kind

; Key-value pairs (dotted key)
(pair
  (dotted_key) @name) @signature @kind

; TOML comments
(comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/toml_test.go

```go
import (
	"testing"
	"unsafe"
	tree_sitter_toml "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/toml"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
type tomlCapture struct {
	Name string
	Kind string
}
func extractTOMLCaptures(t *testing.T, code string) []tomlCapture
captures []tomlCapture
tc tomlCapture
func extractTOMLNames(t *testing.T, code string) []string
names []string
func TestTOMLQueryLanguage(t *testing.T)
func TestTOMLQueryPattern(t *testing.T)
func TestTOMLQueryImportPattern(t *testing.T)
func TestTOMLQueryExtractTables(t *testing.T)
func TestTOMLQueryExtractTableArrays(t *testing.T)
func TestTOMLQueryExtractPairs(t *testing.T)
func TestTOMLQueryExtractDottedKey(t *testing.T)
func TestTOMLQueryExtractComments(t *testing.T)
func TestTOMLQueryKindMapping(t *testing.T)
func TestTOMLQueryCaptures(t *testing.T)
func TestTOMLQueryExtractInlineTable(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/typescript.go

```go
import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)
type TypeScriptQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewTypeScriptQuery() *TypeScriptQuery
func (q *TypeScriptQuery) Language() *sitter.Language
func (q *TypeScriptQuery) Query() []byte
tsKindMapping = map[string]string{
	"function_declaration":   "function",
	"method_definition":      "method",
	"class_declaration":      "class",
	"interface_declaration":  "interface",
	"type_alias_declaration": "type",
	"arrow_function":         "function",
	"variable_declaration":   "variable",
	"variable_declarator":    "arrow",
	"lexical_declaration":    "arrow",
	"export_statement":       "export",
}
func (q *TypeScriptQuery) KindMapping() map[string]string
func (q *TypeScriptQuery) ImportQuery() []byte
func (q *TypeScriptQuery) CallQuery() []byte
typeScriptCallQueryPattern = `
; Direct function calls (e.g., foo())
(call_expression
  function: (identifier) @callee
) @call_node

; Method/property calls (e.g., obj.method())
(call_expression
  function: (member_expression
    property: (property_identifier) @callee
  )
) @call_node
`
typeScriptImportQueryPattern = `
; Import statements (capture full statement)
(import_statement) @import_path

; Export statements with source (re-exports)
(export_statement
  source: (string)
) @import_path

; Export clause (barrel exports: export { foo, bar })
(export_statement
  (export_clause) @import_path
)
`
typeScriptQueryPattern = `
; Function declarations
(function_declaration
  name: (identifier) @name
) @signature @kind

; Exported function declarations
(export_statement
  (function_declaration
    name: (identifier) @name
  )
) @signature @kind

; Arrow functions in variable declarations (capture full declaration with const/let/var)
(lexical_declaration
  (variable_declarator
    name: (identifier) @name
    value: (arrow_function)
  )
) @signature @kind

; Module-level const/let declarations with values (captures all module-level)
; Deduplication for arrow functions is handled in parser.go
(program
  (lexical_declaration
    (variable_declarator
      name: (identifier) @name
      value: (_)
    )
  ) @signature @kind
)

; Module-level const/let without initial value (TypeScript declares)
(program
  (lexical_declaration
    (variable_declarator
      name: (identifier) @name
      !value
    )
  ) @signature @kind
)

; Exported module-level const/let declarations with values
(program
  (export_statement
    declaration: (lexical_declaration
      (variable_declarator
        name: (identifier) @name
        value: (_)
      )
    )
  ) @signature @kind
)

; Exported module-level const/let without initial value
(program
  (export_statement
    declaration: (lexical_declaration
      (variable_declarator
        name: (identifier) @name
        !value
      )
    )
  ) @signature @kind
)

; Method definitions
(method_definition
  name: (property_identifier) @name
) @signature @kind

; Class declarations
(class_declaration
  name: (type_identifier) @name
) @signature @kind

; Interface declarations
(interface_declaration
  name: (type_identifier) @name
) @signature @kind

; Type alias declarations
(type_alias_declaration
  name: (type_identifier) @name
) @signature @kind

; Comments (documentation)
(comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/typescript_test.go

```go
import (
	"testing"
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)
func TestTypeScriptQueryLanguage(t *testing.T)
func TestTypeScriptQueryPattern(t *testing.T)
func TestTypeScriptQueryExtractFunction(t *testing.T)
funcCaptures map[string]string
func TestTypeScriptQueryExtractModuleLevelVariables(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/yaml.go

```go
import (
	tree_sitter_yaml "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/yaml"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
type YAMLQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}
func NewYAMLQuery() *YAMLQuery
func (q *YAMLQuery) Language() *sitter.Language
func (q *YAMLQuery) Query() []byte
yamlKindMapping = map[string]string{
	"block_mapping_pair": "variable",
}
func (q *YAMLQuery) KindMapping() map[string]string
func (q *YAMLQuery) ImportQuery() []byte
yamlQueryPattern = `
; Top-level key-value pairs
(block_mapping_pair
  key: (_) @name) @signature @kind

; YAML comments
(comment) @doc
`
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/yaml_test.go

```go
import (
	"testing"
	"unsafe"
	tree_sitter_yaml "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/yaml"
	sitter "github.com/tree-sitter/go-tree-sitter"
)
func extractYAMLNames(t *testing.T, code string) []string
names []string
func TestYAMLQueryLanguage(t *testing.T)
func TestYAMLQueryPattern(t *testing.T)
func TestYAMLQueryImportPattern(t *testing.T)
func TestYAMLQueryExtractSimpleKeys(t *testing.T)
func TestYAMLQueryExtractNestedKeys(t *testing.T)
func TestYAMLQueryExtractComments(t *testing.T)
func TestYAMLQueryKindMapping(t *testing.T)
func TestYAMLQueryCaptures(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/parser.go

```go
import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/indigo-net/Brf.it/pkg/parser"
	"github.com/indigo-net/Brf.it/pkg/parser/treesitter/languages"
)
func init()
type queryType int
queryTypeSignature queryType = iota
queryTypeImport
queryTypeCall
elixirDefPrefixes = [][]byte{
	[]byte("defmodule "),
	[]byte("defprotocol "),
	[]byte("defimpl "),
}
supportedLangs = "go, typescript, tsx, javascript, jsx, python, c, java, cpp, rust, swift, kotlin, csharp, lua, shell, php, ruby, scala, elixir, sql, yaml, toml"
type queryCacheKey struct {
	lang string
	typ  queryType
}
type TreeSitterParser struct {
	queries         map[string]LanguageQuery
	compiledQueries sync.Map // map[queryCacheKey]*sitter.Query
	parserPool      sync.Pool
	cursorPool      sync.Pool
}
func NewTreeSitterParser() *TreeSitterParser
func (p *TreeSitterParser) Close()
func (p *TreeSitterParser) getOrCreateQuery(lang string, langQuery LanguageQuery, typ queryType) (*sitter.Query, error)
queryBytes []byte
func (p *TreeSitterParser) Parse(content []byte, opts *parser.Options) (result *parser.ParseResult, err error)
rawImports []string
calls []parser.FunctionCall
func (p *TreeSitterParser) Languages() []string
func (p *TreeSitterParser) extractSignatures(
	root *sitter.Node,
	content []byte,
	langQuery LanguageQuery,
	opts *parser.Options,
) ([]parser.Signature, error)
type dedupKey struct {
		line int
		name string
	}
kindNode *sitter.Node
func cleanComment(text string) string
func stripBody(text, kind, language string) string
func stripGoBody(text, kind string) string
func stripTypeScriptBody(text, kind string) string
func stripTSFunctionBody(text string) string
func findFunctionBodyStart(text string) int
func findTSClassBodyStart(text string) int
func stripPythonBody(text, kind string) string
func findPythonBodyStart(text string) int
func stripCBody(text, kind string) string
func stripCppBody(text, kind string) string
func findCppBodyStart(text string) int
func isPythonMethod(signature string) bool
func stripJavaBody(text, kind string) string
func findJavaBodyStart(text string) int
func stripRustBody(text, kind string) string
func findRustBodyStart(text string) int
func refineSwiftClassKind(text string) string
func stripSwiftBody(text, kind string) string
func findSwiftBodyStart(text string) int
func stripKotlinBody(text, kind string) string
func findKotlinBodyStart(text string) int
func refineKotlinClassKind(text string) string
func refineLuaFunctionKind(text string) string
func stripLuaBody(text, kind string) string
func stripPHPBody(text, kind string) string
func stripRubyBody(text, kind string) string
func stripScalaBody(text, kind string) string
func findScalaBodyStart(text string) int
func stripShellBody(text, kind string) string
func findPHPBodyStart(text string) int
func stripCSharpBody(text, kind string) string
func findCSharpBodyStart(text string) int
func isExpressionBodied(text string) bool
func findCSharpArrowIndex(text string) int
func extractCSharpOperatorName(text string) string
func extractCSharpConversionOperatorName(text string) string
func (p *TreeSitterParser) extractImports(
	root *sitter.Node,
	content []byte,
	langQuery LanguageQuery,
	opts *parser.Options,
) ([]string, error)
importNode *sitter.Node
func (p *TreeSitterParser) extractCalls(
	root *sitter.Node,
	content []byte,
	langQuery LanguageQuery,
	opts *parser.Options,
	signatures []parser.Signature,
) ([]parser.FunctionCall, error)
callee string
callLine int
func findEnclosingFunction(signatures []parser.Signature, line int) string
func removeBlankLines(text string) string
buf strings.Builder
elixirDefKeywords = map[string]string{
	"defmodule":   "class",
	"defprotocol": "interface",
	"defimpl":     "impl",
	"def":         "function",
	"defp":        "function",
	"defmacro":    "macro",
	"defmacrop":   "macro",
	"defguard":    "function",
	"defguardp":   "function",
	"defdelegate": "function",
	"defstruct":   "struct",
}
func refineElixirCallKind(text string) string
elixirAttrKeywords = map[string]bool{
	"spec":     true,
	"type":     true,
	"typep":    true,
	"opaque":   true,
	"callback": true,
}
func refineElixirAttrKind(text, capturedName string) (string, string)
func stripElixirBody(text, kind string) string
func extractSQLDDLName(text string) string
func extractNextSQLIdentifier(text string) string
func stripSQLBody(text, kind string) string
func stripSQLFunctionBody(text string) string
func stripSQLViewBody(text string) string
func stripYAMLBody(text, _ string) string
func stripTOMLBody(text, kind string) string
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/parser_bench_test.go

```go
import (
	"strings"
	"testing"
	"github.com/indigo-net/Brf.it/pkg/parser"
)
func BenchmarkParseGo(b *testing.B)
code strings.Builder
func BenchmarkParseTypeScript(b *testing.B)
code strings.Builder
func BenchmarkParsePython(b *testing.B)
code strings.Builder
func BenchmarkParseWithImports(b *testing.B)
code strings.Builder
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/parser_fuzz_test.go

```go
import (
	"testing"
	"github.com/indigo-net/Brf.it/pkg/parser"
)
func FuzzParseGo(f *testing.F)
func FuzzParseTypeScript(f *testing.F)
func FuzzParsePython(f *testing.F)
func FuzzParseJava(f *testing.F)
func FuzzParseRust(f *testing.F)
func FuzzParseC(f *testing.F)
func FuzzParseJSON(f *testing.F)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/parser_test.go

```go
import (
	"strings"
	"testing"
	"github.com/indigo-net/Brf.it/pkg/parser"
)
func TestTreeSitterParserImplementsParser(t *testing.T)
_ parser.Parser = (*TreeSitterParser)(nil)
func TestTreeSitterParserLanguages(t *testing.T)
func TestTreeSitterParserParseGo(t *testing.T)
foundAdd bool
func TestTreeSitterParserParseTypeScript(t *testing.T)
foundAdd bool
func TestTreeSitterParserUnsupportedLanguage(t *testing.T)
func TestTreeSitterParserAutoRegistration(t *testing.T)
func TestGoSignatureOnlyExtraction(t *testing.T)
func TestGoIncludeBodyExtraction(t *testing.T)
foundAdd bool
func TestTypeScriptSignatureOnlyExtraction(t *testing.T)
func TestTypeScriptArrowFunctionSignature(t *testing.T)
func contains(s, substr string) bool
func TestTreeSitterParserParseJava(t *testing.T)
foundClass, foundConstructor, foundPublicMethod bool
func TestJavaSignatureOnlyExtraction(t *testing.T)
func TestJavaGenericsExtraction(t *testing.T)
foundClass, foundMethod bool
func TestJavaAutoRegistration(t *testing.T)
func TestTreeSitterParserParseCpp(t *testing.T)
func TestCppSignatureOnlyExtraction(t *testing.T)
func TestCppTemplateExtraction(t *testing.T)
func TestCppAutoRegistration(t *testing.T)
func TestCppImportExtraction(t *testing.T)
func TestGoVariableExtraction(t *testing.T)
func TestTypeScriptVariableExtraction(t *testing.T)
func TestPythonVariableExtraction(t *testing.T)
func TestJavaStaticFieldExtraction(t *testing.T)
func TestCGlobalVariableExtraction(t *testing.T)
func TestVariableSignaturePreservesValue(t *testing.T)
func TestTreeSitterParserParseRust(t *testing.T)
func TestRustSignatureOnlyExtraction(t *testing.T)
func TestRustImportExtraction(t *testing.T)
func TestRustAutoRegistration(t *testing.T)
func TestRustConstAndStaticExtraction(t *testing.T)
func TestRustMacroExtraction(t *testing.T)
func TestRustGenericsAndLifetimes(t *testing.T)
func TestTreeSitterParserParseSwift(t *testing.T)
func TestTreeSitterParserParseKotlin(t *testing.T)
func TestKotlinBodyStripping(t *testing.T)
func TestRefineKotlinClassKind(t *testing.T)
func TestKotlinAutoRegistration(t *testing.T)
func TestKotlinSignatureOnlyExtraction(t *testing.T)
func TestKotlinImportExtraction(t *testing.T)
func TestParsePanicRecovery(t *testing.T)
func TestTreeSitterParserParseCSharp(t *testing.T)
func TestCSharpSignatureOnlyExtraction(t *testing.T)
func TestCSharpOperatorNameSynthesis(t *testing.T)
func TestCSharpStaticFieldExtraction(t *testing.T)
func TestCSharpImportExtraction(t *testing.T)
func TestCSharpAutoRegistration(t *testing.T)
func TestCSharpBodyStripping(t *testing.T)
func TestFindCSharpBodyStart(t *testing.T)
func TestIsExpressionBodied(t *testing.T)
func TestExtractCSharpOperatorName(t *testing.T)
func TestExtractCSharpConversionOperatorName(t *testing.T)
func TestFindFunctionBodyStart(t *testing.T)
func TestParsePanicRecoveryMechanism(t *testing.T)
func TestTreeSitterParserParseLua(t *testing.T)
func TestTreeSitterParserParseLuaImports(t *testing.T)
func TestLuaBodyStripping(t *testing.T)
func TestRefineLuaFunctionKind(t *testing.T)
func TestIsPythonMethod(t *testing.T)
func TestTreeSitterParserParsePHP(t *testing.T)
foundClass, foundMethod, foundInterface, foundTrait, foundFunction, foundConst bool
func TestTreeSitterParserParsePHPImports(t *testing.T)
func TestPHPBodyStripping(t *testing.T)
func TestTreeSitterParserParseRuby(t *testing.T)
func TestTreeSitterParserParseRubyImports(t *testing.T)
func TestTreeSitterParserParseScala(t *testing.T)
func TestTreeSitterParserParseScalaImports(t *testing.T)
func TestTreeSitterParserParseElixir(t *testing.T)
func TestTreeSitterParserParseElixirImports(t *testing.T)
func TestRefineElixirCallKind(t *testing.T)
func TestRefineElixirAttrKind(t *testing.T)
func TestStripElixirBody(t *testing.T)
func TestTreeSitterParserParseSQL(t *testing.T)
names []string
func TestExtractSQLDDLName(t *testing.T)
func TestExtractNextSQLIdentifier(t *testing.T)
func TestStripSQLBody(t *testing.T)
func TestStripSQLFunctionBody(t *testing.T)
func TestTreeSitterParserParseYAML(t *testing.T)
names []string
func TestTreeSitterParserParseTOML(t *testing.T)
names []string
func TestSameLineDifferentNames(t *testing.T)
func TestFindEnclosingFunctionEndLineZero(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/query.go

```go
import (
	sitter "github.com/tree-sitter/go-tree-sitter"
)
type LanguageQuery interface {
	// Language returns the Tree-sitter language for parsing.
	Language() *sitter.Language

	// Query returns the Tree-sitter query pattern for signature extraction.
	Query() []byte

	// ImportQuery returns the Tree-sitter query pattern for import/export extraction.
	// Returns nil if the language doesn't support import extraction.
	ImportQuery() []byte

	// CallQuery returns the Tree-sitter query pattern for function call extraction.
	// Returns nil if the language doesn't support call extraction.
	CallQuery() []byte

	// Captures returns the list of capture names used in the query.
	Captures() []string

	// KindMapping maps Tree-sitter node types to Signature kinds.
	KindMapping() map[string]string

	// IsExported reports whether the symbol identified by name and sigText
	// is public/exported in the given language. sigText is the full
	// signature text which may contain visibility modifiers.
	IsExported(name, sigText string) bool
}
CaptureName      = "name"
CaptureSignature = "signature"
CaptureDoc       = "doc"
CaptureKind      = "kind"
CaptureImportPath = "import_path"
CaptureLuaRequireFn = "_fn"
CaptureCallee = "callee"
DefaultKindMapping = map[string]string{
	"function_declaration": "function",
	"method_declaration":   "method",
	"type_declaration":     "type",
	"struct_type":          "struct",
	"interface_type":       "interface",
	"class_declaration":    "class",
	"arrow_function":       "function",
	"function_expression":  "function",
	"method_definition":    "method",
}
```

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/query_test.go

```go
import (
	"testing"
)
func TestCaptureDefinitions(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/scanner/example_test.go

```go
import (
	"fmt"
	"github.com/indigo-net/Brf.it/pkg/scanner"
)
func ExampleDefaultScanOptions()
func ExampleIsHidden()
```

### /home/runner/work/Brf.it/Brf.it/pkg/scanner/scanner.go

```go
import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"github.com/bmatcuk/doublestar/v4"
	ignore "github.com/sabhiram/go-gitignore"
	"github.com/indigo-net/Brf.it/pkg/parser"
)
type FileEntry struct {
	// Path is the absolute or relative path to the file.
	Path string

	// Language is the detected programming language (e.g., "go", "typescript").
	Language string

	// Size is the file size in bytes.
	Size int64

	// Content holds the file bytes when PreloadContent is enabled.
	// nil when content was not preloaded.
	Content []byte
}
type ScanResult struct {
	// Files is the list of matched files.
	Files []FileEntry

	// TotalSize is the sum of all matched file sizes.
	TotalSize int64

	// SkippedCount is the number of files skipped (too large, unsupported, etc.).
	SkippedCount int

	// Warnings contains non-fatal issues encountered during scanning.
	Warnings []string
}
type ScanOptions struct {
	// RootPath is the directory or file to scan.
	RootPath string

	// SupportedExtensions maps file extensions to language names.
	SupportedExtensions map[string]string

	// IgnoreFiles is the list of ignore file paths (default: [".gitignore"]).
	IgnoreFiles []string

	// IncludePatterns is a list of glob patterns to include.
	// If non-empty, only files matching at least one pattern are included.
	// Supports doublestar (**) patterns.
	IncludePatterns []string

	// ExcludePatterns is a list of glob patterns to exclude.
	// Files matching any pattern are excluded.
	// Supports doublestar (**) patterns.
	ExcludePatterns []string

	// ChangedFiles is an optional whitelist of file paths (relative to RootPath).
	// When non-nil, only files in this list are included in scan results.
	// Used by --changed and --since flags to restrict scanning to git-changed files.
	ChangedFiles map[string]bool

	// IncludeHidden determines whether to include hidden files (dotfiles).
	IncludeHidden bool

	// MaxFileSize is the maximum file size in bytes to include.
	MaxFileSize int64

	// PreloadContent reads file content during scan so downstream consumers
	// (e.g., the extractor) can skip a redundant os.ReadFile call.
	PreloadContent bool

	// MaxTotalPreloadSize limits the total bytes preloaded into memory when
	// PreloadContent is true. Once this budget is exceeded, remaining files
	// are included in the scan results but with Content set to nil (the
	// extractor will fall back to on-demand os.ReadFile). A value of 0 means
	// no limit. Default: 0.
	MaxTotalPreloadSize int64
}
func DefaultScanOptions() *ScanOptions
func copyLanguageMapping() map[string]string
func (o *ScanOptions) GetLanguage(path string) (string, bool)
func IsHidden(name string) bool
func getBaseName(path string) string
type Scanner interface {
	// Scan performs the scan and returns scan results.
	// The context allows cancellation of long-running scans.
	Scan(ctx context.Context) (*ScanResult, error)
}
type FileScanner struct {
	opts              *ScanOptions
	ignorers          []*ignore.GitIgnore
	ignorerErrs       []error
	ignorerErrsWarned bool
	logger            *log.Logger
	rootIsFile        bool
	preloadedSize     int64 // tracks total bytes preloaded so far
}
func NewFileScanner(opts *ScanOptions) (*FileScanner, error)
rootIsFile bool
func (s *FileScanner) Scan(ctx context.Context) (*ScanResult, error)
warning string
func (s *FileScanner) relPath(path string) string
func (s *FileScanner) matchesInclude(path string) bool
func (s *FileScanner) matchesExclude(path string) bool
func (s *FileScanner) matchesExcludeDir(path string) bool
func (s *FileScanner) matchesIgnore(path string) bool
func (s *FileScanner) checkFile(path string, info os.FileInfo) (FileEntry, bool)
```

### /home/runner/work/Brf.it/Brf.it/pkg/scanner/scanner_bench_test.go

```go
import (
	"context"
	"os"
	"path/filepath"
	"testing"
)
func BenchmarkScanDirectory(b *testing.B)
func BenchmarkScanLargeFile(b *testing.B)
func BenchmarkScanWithIgnore(b *testing.B)
```

### /home/runner/work/Brf.it/Brf.it/pkg/scanner/scanner_test.go

```go
import (
	"bytes"
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)
func TestGetBaseName(t *testing.T)
func TestNewFileScanner(t *testing.T)
func TestNewFileScannerNilOptions(t *testing.T)
func TestFileEntryDefaults(t *testing.T)
func TestScanOptionsDefaults(t *testing.T)
expectedMaxSize = 512000
func TestScanOptionsWithExtensions(t *testing.T)
func TestScannerInterface(t *testing.T)
_ Scanner = (*FileScanner)(nil)
func TestScanEmptyDirectory(t *testing.T)
func TestScanSingleFile(t *testing.T)
func TestScanFilterByExtension(t *testing.T)
func TestScanExcludeHidden(t *testing.T)
func TestScanIncludeHidden(t *testing.T)
func TestScanMaxFileSize(t *testing.T)
func TestScanGitignore(t *testing.T)
func TestScanGitignoreLoadFailureWarning(t *testing.T)
buf bytes.Buffer
buf bytes.Buffer
buf bytes.Buffer
func TestScanWalkDirPermissionDenied(t *testing.T)
buf bytes.Buffer
func TestScanSymlinkSkip(t *testing.T)
buf bytes.Buffer
func TestScanNestedDirectories(t *testing.T)
func TestLogOutputNoDoubleNewline(t *testing.T)
buf bytes.Buffer
func TestScanMultipleIgnoreFiles(t *testing.T)
func TestFilepathBaseEdgeCases(t *testing.T)
func TestScanIncludePatterns(t *testing.T)
func TestScanExcludeDirectory(t *testing.T)
func TestScanSingleFileWithIncludePattern(t *testing.T)
func TestScanChangedFilesWhitelist(t *testing.T)
func TestNewFileScannerInvalidPatterns(t *testing.T)
func TestPreloadContent(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/pkg/security/scanner.go

```go
import (
	"fmt"
	"io"
	"regexp"
	"github.com/indigo-net/Brf.it/pkg/extractor"
	"github.com/indigo-net/Brf.it/pkg/parser"
)
type Pattern struct {
	// Name is a human-readable name for the pattern.
	Name string

	// Regex is the compiled regular expression.
	Regex *regexp.Regexp
}
func defaultPatterns() []Pattern
type Finding struct {
	// FilePath is the file where the secret was found.
	FilePath string

	// PatternName is the name of the matched pattern.
	PatternName string
}
type ScanResult struct {
	// Findings is the list of detected secrets.
	Findings []Finding

	// RedactedFiles is the extract result with secrets redacted.
	RedactedFiles []extractor.ExtractedFile
}
type Scanner struct {
	patterns []Pattern
	warnings io.Writer
}
func NewScanner(warnings io.Writer) *Scanner
func (s *Scanner) SetWarnings(w io.Writer)
func (s *Scanner) Scan(result *extractor.ExtractResult) *ScanResult
func (s *Scanner) scanFile(file extractor.ExtractedFile, sr *ScanResult) extractor.ExtractedFile
func (s *Scanner) redactString(filePath, text string, sr *ScanResult) string
```

### /home/runner/work/Brf.it/Brf.it/pkg/security/scanner_test.go

```go
import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"github.com/indigo-net/Brf.it/pkg/extractor"
	"github.com/indigo-net/Brf.it/pkg/parser"
)
func TestScan_NilResult(t *testing.T)
func TestScan_NoSecrets(t *testing.T)
buf bytes.Buffer
func TestScan_AWSAccessKeyDetected(t *testing.T)
buf bytes.Buffer
func TestScan_GitHubTokenDetected(t *testing.T)
buf bytes.Buffer
func TestScan_GenericAPIKeyDetected(t *testing.T)
buf bytes.Buffer
func TestScan_PasswordInDocDetected(t *testing.T)
buf bytes.Buffer
func TestScan_PrivateKeyDetected(t *testing.T)
buf bytes.Buffer
func TestScan_ImportRedacted(t *testing.T)
buf bytes.Buffer
func TestScan_ErrorFileSkipped(t *testing.T)
buf bytes.Buffer
func TestScan_BearerTokenDetected(t *testing.T)
buf bytes.Buffer
func TestScan_MultipleSecretsInOneFile(t *testing.T)
buf bytes.Buffer
```

### /home/runner/work/Brf.it/Brf.it/pkg/tokenizer/example_test.go

```go
import (
	"fmt"
	"github.com/indigo-net/Brf.it/pkg/tokenizer"
)
func ExampleNewNoOpTokenizer()
```

### /home/runner/work/Brf.it/Brf.it/pkg/tokenizer/tiktoken.go

```go
import (
	"unsafe"
	"github.com/pkoukk/tiktoken-go"
)
type TiktokenTokenizer struct {
	encoding string
	tke      *tiktoken.Tiktoken
}
_ Tokenizer = (*TiktokenTokenizer)(nil)
func NewTiktokenTokenizer() (*TiktokenTokenizer, error)
func (t *TiktokenTokenizer) Count(text []byte) (int, error)
func (t *TiktokenTokenizer) Name() string
```

### /home/runner/work/Brf.it/Brf.it/pkg/tokenizer/tokenizer.go

```go
type Tokenizer interface {
	// Count returns the number of tokens in the given text.
	// Returns 0 and error if counting fails.
	Count(text []byte) (int, error)

	// Name returns the tokenizer name (e.g., "tiktoken-cl100k", "noop").
	Name() string
}
type NoOpTokenizer struct{}
_ Tokenizer = (*NoOpTokenizer)(nil)
func NewNoOpTokenizer() *NoOpTokenizer
func (t *NoOpTokenizer) Count(_ []byte) (int, error)
func (t *NoOpTokenizer) Name() string
```

### /home/runner/work/Brf.it/Brf.it/pkg/tokenizer/tokenizer_test.go

```go
import (
	"strings"
	"testing"
)
func TestNoOpTokenizerImplementsTokenizer(t *testing.T)
_ Tokenizer = (*NoOpTokenizer)(nil)
func TestTiktokenTokenizerImplementsTokenizer(t *testing.T)
_ Tokenizer = (*TiktokenTokenizer)(nil)
func TestNoOpTokenizerCount(t *testing.T)
func TestNoOpTokenizerName(t *testing.T)
func TestTiktokenTokenizerCount(t *testing.T)
func TestTiktokenTokenizerName(t *testing.T)
func TestTiktokenTokenizerConsistency(t *testing.T)
func TestTiktokenTokenizerSpecialCharacters(t *testing.T)
```

### /home/runner/work/Brf.it/Brf.it/website/docusaurus.config.ts

```typescript
import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';
const config: Config = {
  title: 'Brf.it',
  tagline: 'Give AI the gist, not the bloat',
  favicon: 'img/favicon.ico',

  future: {
    v4: true,
  },

  url: 'https://indigo-net.github.io',
  baseUrl: '/Brf.it/',

  organizationName: 'indigo-net',
  projectName: 'Brf.it',
  trailingSlash: false,

  onBrokenLinks: 'throw',

  headTags: [
    {
      tagName: 'meta',
      attributes: {
        name: 'keywords',
        content: 'brfit, code extraction, AI, LLM, token optimization, function signatures, Tree-sitter, CLI, code context, MCP',
      },
    },
    {
      tagName: 'meta',
      attributes: {
        name: 'author',
        content: 'indigo-net',
      },
    },
  ],

  i18n: {
    defaultLocale: 'en',
    locales: ['en', 'ko', 'ja', 'de', 'hi'],
    localeConfigs: {
      en: { label: 'English', direction: 'ltr' },
      ko: { label: '한국어', direction: 'ltr' },
      ja: { label: '日本語', direction: 'ltr' },
      de: { label: 'Deutsch', direction: 'ltr' },
      hi: { label: 'हिन्दी', direction: 'ltr' },
    },
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/indigo-net/Brf.it/tree/main/website/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: 'img/docusaurus-social-card.jpg',
    metadata: [
      {name: 'og:type', content: 'website'},
      {name: 'twitter:card', content: 'summary_large_image'},
    ],
    colorMode: {
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'Brf.it',
      logo: {
        alt: 'Brf.it Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          position: 'left',
          label: 'Docs',
        },
        {
          type: 'localeDropdown',
          position: 'right',
          dropdownItemsBefore: [],
          dropdownItemsAfter: [],
        },
        {
          href: 'https://github.com/indigo-net/Brf.it',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Getting Started',
              to: '/docs/',
            },
            {
              label: 'CLI Reference',
              to: '/docs/cli-reference',
            },
            {
              label: 'Supported Languages',
              to: '/docs/languages/',
            },
            {
              label: 'Changelog',
              to: '/docs/changelog',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'GitHub Issues',
              href: 'https://github.com/indigo-net/Brf.it/issues',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/indigo-net/Brf.it',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} indigo-net. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['go', 'typescript', 'python', 'java', 'kotlin', 'rust', 'ruby', 'php', 'swift', 'scala'],
    },
  } satisfies Preset.ThemeConfig,
};
```

### /home/runner/work/Brf.it/Brf.it/website/sidebars.ts

```typescript
import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';
const sidebars: SidebarsConfig = {
  docsSidebar: [
    {
      type: 'doc',
      id: 'getting-started',
      label: 'Getting Started',
    },
    {
      type: 'doc',
      id: 'cli-reference',
      label: 'CLI Reference',
    },
    {
      type: 'category',
      label: 'Languages',
      items: [
        'languages/index',
        'languages/go',
        'languages/typescript',
        'languages/python',
        'languages/java',
        'languages/kotlin',
        'languages/rust',
        'languages/ruby',
        'languages/php',
        'languages/swift',
        'languages/scala',
        'languages/c-cpp',
      ],
    },
    {
      type: 'doc',
      id: 'changelog',
      label: 'Changelog',
    },
  ],
};
```

### /home/runner/work/Brf.it/Brf.it/website/src/components/FeaturesSection.tsx

```typescript
import React from 'react';
import Translate from '@docusaurus/Translate';
interface Feature {
  title: string;
  description: string;
  icon: string;
}
```

### /home/runner/work/Brf.it/Brf.it/website/src/components/Hero.tsx

```typescript
import React from 'react';
import Link from '@docusaurus/Link';
import Translate, {translate} from '@docusaurus/Translate';
export default function Hero(): JSX.Element
```

### /home/runner/work/Brf.it/Brf.it/website/src/components/InstallSection.tsx

```typescript
import React, {useState} from 'react';
import Translate from '@docusaurus/Translate';
type Platform = keyof typeof installCommands;
const copyToClipboard = (text: string)
```

### /home/runner/work/Brf.it/Brf.it/website/src/components/LanguageGrid.tsx

```typescript
import React from 'react';
import Translate from '@docusaurus/Translate';
interface Language {
  name: string;
  icon: string;
}
const languages: Language[] = [
  { name: 'Go', icon: '🐹' },
  { name: 'TypeScript', icon: '📘' },
  { name: 'JavaScript', icon: '📒' },
  { name: 'Python', icon: '🐍' },
  { name: 'Java', icon: '☕' },
  { name: 'Kotlin', icon: '🟣' },
  { name: 'Rust', icon: '🦀' },
  { name: 'Ruby', icon: '💎' },
  { name: 'PHP', icon: '🐘' },
  { name: 'Swift', icon: '🍎' },
  { name: 'Scala', icon: '🔴' },
  { name: 'C/C++', icon: '⚙️' },
];
export default function LanguageGrid(): JSX.Element {
  return (
    <section className="section" style={{ background: 'var(--ifm-color-emphasis-100)' }}>
      <div className="container">
        <h2 className="section-title">
          <Translate id="languages.title">Supported Languages</Translate>
        </h2>
        <p className="section-subtitle">
          <Translate id="languages.subtitle">
            Tree-sitter powered parsing for accurate signature extraction
          </Translate>
        </p>

        <div className="language-grid">
          {languages.map(lang
```

### /home/runner/work/Brf.it/Brf.it/website/src/components/TokenComparison.tsx

```typescript
import React, {useState} from 'react';
import Translate, {translate} from '@docusaurus/Translate';
interface CodeExample {
  language: string;
  label: string;
  before: string;
  after: string;
  beforeTokens: number;
  afterTokens: number;
  beforeLines: number;
  afterLines: number;
}
```

### /home/runner/work/Brf.it/Brf.it/website/src/pages/index.tsx

```typescript
import type {ReactNode} from 'react';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Hero from '@site/src/components/Hero';
import TokenComparison from '@site/src/components/TokenComparison';
import FeaturesSection from '@site/src/components/FeaturesSection';
import LanguageGrid from '@site/src/components/LanguageGrid';
import InstallSection from '@site/src/components/InstallSection';
export default function Home(): ReactNode
```

