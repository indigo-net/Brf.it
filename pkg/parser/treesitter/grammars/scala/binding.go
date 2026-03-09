// Package scala provides the tree-sitter grammar for Scala.
//
// The C source files (parser.c, scanner.c) are vendored from
// https://github.com/tree-sitter/tree-sitter-scala (latest).
// CGO automatically compiles all .c files in this directory.
package scala

// #cgo CFLAGS: -std=c11 -fPIC
// extern const void *tree_sitter_scala(void);
import "C"

import "unsafe"

// Language returns the tree-sitter Language for Scala.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_scala())
}
