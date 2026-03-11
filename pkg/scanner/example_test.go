package scanner_test

import (
	"fmt"

	"github.com/indigo-net/Brf.it/pkg/scanner"
)

func ExampleDefaultScanOptions() {
	opts := scanner.DefaultScanOptions()
	fmt.Println(opts.MaxFileSize)
	fmt.Println(opts.IncludeHidden)
	// Output:
	// 512000
	// false
}

func ExampleIsHidden() {
	fmt.Println(scanner.IsHidden(".gitignore"))
	fmt.Println(scanner.IsHidden("main.go"))
	// Output:
	// true
	// false
}
