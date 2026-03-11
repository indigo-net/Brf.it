package extractor_test

import (
	"fmt"

	"github.com/indigo-net/Brf.it/pkg/extractor"
)

func ExampleDefaultExtractOptions() {
	opts := extractor.DefaultExtractOptions()
	fmt.Println(opts.Concurrency)
	fmt.Println(opts.IncludePrivate)
	fmt.Println(opts.IncludeBody)
	// Output:
	// 0
	// false
	// false
}
