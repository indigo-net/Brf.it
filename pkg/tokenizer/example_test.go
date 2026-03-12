package tokenizer_test

import (
	"fmt"

	"github.com/indigo-net/Brf.it/pkg/tokenizer"
)

func ExampleNewNoOpTokenizer() {
	t := tokenizer.NewNoOpTokenizer()
	fmt.Println(t.Name())

	count, err := t.Count([]byte("hello world"))
	fmt.Println(count, err)
	// Output:
	// noop
	// 0 <nil>
}
