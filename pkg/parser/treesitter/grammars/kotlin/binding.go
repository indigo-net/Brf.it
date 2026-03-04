// Package kotlin provides the tree-sitter grammar for Kotlin.
//
// The C source files (parser.c, scanner.c) are vendored from
// https://github.com/fwcd/tree-sitter-kotlin (v0.3.8).
// CGO automatically compiles all .c files in this directory.
package kotlin

// #cgo CFLAGS: -std=c11 -fPIC
// extern const void *tree_sitter_kotlin(void);
import "C"

import "unsafe"

// Language returns the tree-sitter Language for Kotlin.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_kotlin())
}
