package main

import (
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
)

type Page struct {
	Title   string
	Content template.HTML
}

func main() {
	path := filepath.Join("./content/")
	// outputPath := filepath.Join("./output/")
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".md") {
			if err != nil {
				log.Fatalf("Error: %s", err)
			}
			htmlContent := convertToHtml(path)
			newPage := Page{
				Title:   "hello",
				Content: template.HTML(htmlContent),
			}

			tmpl, err := template.ParseFiles("./layout.html")
			if err != nil {
				log.Fatalf("Error parsing template: %s", err)
			}

			outputFile, err := os.Create("./output.html")
			if err != nil {
				log.Fatalf("Error creating output file: %s", err)
			}

			defer outputFile.Close()

			err = tmpl.Execute(outputFile, newPage)
			if err != nil {
				log.Fatalf("Error generating output content %s", err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error here: %s", err)
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
