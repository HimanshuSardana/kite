package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
)

type Page struct {
	Title   string
	Content template.HTML
}

type Frontmatter struct {
	Title string `yaml:"title"`
}

func main() {
	contentDir := "./content"
	outputDir := "./output"

	err := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".md") {
			fmt.Println("Processing:", path)

			title, htmlContent := convertToHtml(path)
			newPage := Page{
				Title:   title, // you can customize this later
				Content: template.HTML(htmlContent),
			}

			// Parse template once per file (could be optimized)
			tmpl, err := template.ParseFiles("./layout.html")
			if err != nil {
				log.Fatalf("Error parsing template: %s", err)
			}

			// Compute relative path and output file path
			relPath, err := filepath.Rel(contentDir, path)
			if err != nil {
				log.Fatalf("Error computing relative path: %s", err)
			}
			outputFilePath := filepath.Join(outputDir, strings.Replace(relPath, ".md", ".html", 1))

			// Make sure the output directory exists
			err = os.MkdirAll(filepath.Dir(outputFilePath), 0o755)
			if err != nil {
				log.Fatalf("Error creating directories: %s", err)
			}

			// Create the output file
			outputFile, err := os.Create(outputFilePath)
			if err != nil {
				log.Fatalf("Error creating output file: %s", err)
			}
			defer outputFile.Close()

			// Execute template
			err = tmpl.Execute(outputFile, newPage)
			if err != nil {
				log.Fatalf("Error generating output content: %s", err)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Error walking directory: %s", err)
	}

	fmt.Println("All files processed!")
}

func convertToHtml(path string) (string, []byte) {
	md, err := os.ReadFile(path)
	var matter Frontmatter
	rest, err := frontmatter.Parse(strings.NewReader(string(md)), &matter)
	if err != nil {
		log.Fatalf("Error reading %s: %s", path, err)
	}

	fmt.Printf("Found post %s\n", matter.Title)

	html := markdown.ToHTML(rest, nil, nil)
	return matter.Title, html
}
