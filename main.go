package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	// "github.com/gomarkdown/markdown/html"
	// "github.com/gomarkdown/markdown/parser"
)

func main() {
	path := filepath.Join("./test.md")
	mds, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error %s", err)
	}
	md := []byte(mds)
	html := markdown.ToHTML(md, nil, nil)

	fmt.Printf("--- Markdown:\n%s\n\n--- HTML:\n%s\n", md, html)
}
