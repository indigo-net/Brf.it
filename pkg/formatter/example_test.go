package formatter_test

import (
	"fmt"

	"github.com/indigo-net/Brf.it/pkg/formatter"
)

func ExampleNewXMLFormatter() {
	f := formatter.NewXMLFormatter()
	fmt.Println(f.Name())
	// Output:
	// xml
}

func ExampleNewMarkdownFormatter() {
	f := formatter.NewMarkdownFormatter()
	fmt.Println(f.Name())
	// Output:
	// markdown
}

func ExampleNewJSONFormatter() {
	f := formatter.NewJSONFormatter()
	fmt.Println(f.Name())
	// Output:
	// json
}

func ExampleXMLFormatter_Format() {
	f := formatter.NewXMLFormatter()
	data := &formatter.PackageData{
		RootPath: "/project",
		Version:  "1.0.0",
		Files: []formatter.FileData{
			{
				Path:     "main.go",
				Language: "go",
			},
		},
	}
	output, err := f.Format(data)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(len(output) > 0)
	// Output:
	// true
}
