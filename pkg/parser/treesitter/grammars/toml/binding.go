// Package toml provides the tree-sitter grammar for TOML.
//
// The C source files (parser.c, scanner.c) are vendored from
// https://github.com/tree-sitter-grammars/tree-sitter-toml (v0.7.0).
// CGO automatically compiles all .c files in this directory.
package toml

// #cgo CFLAGS: -std=c11 -fPIC
// extern const void *tree_sitter_toml(void);
import "C"

import "unsafe"

// Language returns the tree-sitter Language for TOML.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_toml())
}
