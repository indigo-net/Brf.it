package parser_test

import (
	"fmt"

	"github.com/indigo-net/Brf.it/pkg/parser"
)

func ExampleDetectLanguage() {
	fmt.Println(parser.DetectLanguage("main.go"))
	fmt.Println(parser.DetectLanguage("app.tsx"))
	fmt.Println(parser.DetectLanguage("script.py"))
	fmt.Println(parser.DetectLanguage("unknown.xyz"))
	// Output:
	// go
	// typescript
	// python
	//
}

func ExampleNewRegistry() {
	reg := parser.NewRegistry()
	fmt.Println(len(reg.Languages()))
	// Output:
	// 0
}
