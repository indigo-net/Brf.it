// Package yaml provides the tree-sitter grammar for YAML.
//
// The C source files (parser.c, scanner.c) are vendored from
// https://github.com/tree-sitter-grammars/tree-sitter-yaml (v0.7.2).
// CGO automatically compiles all .c files in this directory.
package yaml

// #cgo CFLAGS: -std=c11 -fPIC
// extern const void *tree_sitter_yaml(void);
import "C"

import "unsafe"

// Language returns the tree-sitter Language for YAML.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_yaml())
}
