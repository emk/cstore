package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s HOST:POST\n", os.Args[0])
		os.Exit(1)
	}
	addr := os.Args[1]
	fmt.Printf("Serving on http://%s/\n", addr)
	ListenAndServe(addr)
}
