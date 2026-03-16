package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
)

func main() {
	path := filepath.Join("./content/")
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	for _, f := range files {
		filePath := filepath.Join(path, f.Name())
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
			fmt.Printf("Found content: %s\n", filePath)
			htmlContent := convertToHtml(filePath)
			htmlPath := strings.Replace(filePath, ".md", ".html", 1)
			os.WriteFile(htmlPath, htmlContent, 0o777)
			fmt.Printf("Wrote file: %s\n", htmlPath)
		}
	}
}

func convertToHtml(path string) []byte {
	mds, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	md := []byte(mds)
	html := markdown.ToHTML(md, nil, nil)

	// fmt.Printf("--- Markdown:\n%s\n\n--- HTML:\n%s\n", md, html)
	return html
}
