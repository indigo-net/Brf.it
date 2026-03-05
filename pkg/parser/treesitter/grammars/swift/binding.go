// Package swift provides the tree-sitter grammar for Swift.
//
// The C source files (parser.c, scanner.c) are vendored from
// https://github.com/alex-pinkus/tree-sitter-swift (v0.7.1).
// CGO automatically compiles all .c files in this directory.
package swift

// #cgo CFLAGS: -std=c11 -fPIC
// extern const void *tree_sitter_swift(void);
import "C"

import "unsafe"

// Language returns the tree-sitter Language for Swift.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_swift())
}
