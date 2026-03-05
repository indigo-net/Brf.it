// Package csharp provides the tree-sitter grammar for C#.
//
// The C source files (parser.c, scanner.c) are vendored from
// https://github.com/tree-sitter/tree-sitter-c-sharp (v0.23.1).
// CGO automatically compiles all .c files in this directory.
package csharp

// #cgo CFLAGS: -std=c11 -fPIC
// extern const void *tree_sitter_c_sharp(void);
import "C"

import "unsafe"

// Language returns the tree-sitter Language for C#.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_c_sharp())
}
