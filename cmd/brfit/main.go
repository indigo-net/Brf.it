package main

import "fmt"

// Version is set at build time via -ldflags
var Version = "0.1.0"

func main() {
	fmt.Println("brf.it - Don't just feed raw code. Brief it for your AI.")
}
