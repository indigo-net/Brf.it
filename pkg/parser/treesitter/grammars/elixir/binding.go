// Package elixir provides the tree-sitter grammar for Elixir.
//
// The C source files (parser.c, scanner.c) are vendored from
// https://github.com/elixir-lang/tree-sitter-elixir (v0.3.5).
// CGO automatically compiles all .c files in this directory.
package elixir

// #cgo CFLAGS: -std=c11 -fPIC
// extern const void *tree_sitter_elixir(void);
import "C"

import "unsafe"

// Language returns the tree-sitter Language for Elixir.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_elixir())
}
