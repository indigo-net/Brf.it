// Package sql provides the tree-sitter grammar for SQL.
//
// The C source files (parser.c, scanner.c) are vendored from
// https://github.com/DerekStride/tree-sitter-sql (v0.3.11).
// CGO automatically compiles all .c files in this directory.
package sql

// #cgo CFLAGS: -std=c11 -fPIC
// extern const void *tree_sitter_sql(void);
import "C"

import "unsafe"

// Language returns the tree-sitter Language for SQL.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_sql())
}
