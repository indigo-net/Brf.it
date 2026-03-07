// Package lua provides the tree-sitter grammar for Lua.
//
// The C source files (parser.c, scanner.c) are vendored from
// https://github.com/tree-sitter-grammars/tree-sitter-lua (v0.5.0).
// CGO automatically compiles all .c files in this directory.
package lua

// #cgo CFLAGS: -std=c11 -fPIC
// extern const void *tree_sitter_lua(void);
import "C"

import "unsafe"

// Language returns the tree-sitter Language for Lua.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_lua())
}
