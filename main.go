package main

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
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

	themeName := "modern-light"

	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "serve":

			copyFile("./themes/"+themeName+".css", "./output/style.css")

			fs := http.FileServer(http.Dir("./output/"))
			http.Handle("/", fs)

			log.Println("Serving on http://localhost:8000")

			err := http.ListenAndServe(":8000", nil)
			if err != nil {
				log.Fatalf("Error occured %s\n", err)
			}
		}
	}

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
				Title:   title,
				Content: template.HTML(htmlContent),
			}

			tmpl, err := template.ParseFiles("./layout.html")
			if err != nil {
				log.Fatalf("Error parsing template: %s", err)
			}

			relPath, err := filepath.Rel(contentDir, path)
			if err != nil {
				log.Fatalf("Error computing relative path: %s", err)
			}
			// test.md -> test.html
			// test.md -> test/index.html
			outputFilePath := filepath.Join(outputDir, strings.Replace(relPath, ".md", "/index.html", 1))

			err = os.MkdirAll(filepath.Dir(outputFilePath), 0o755)
			if err != nil {
				log.Fatalf("Error creating directories: %s", err)
			}

			outputFile, err := os.Create(outputFilePath)
			if err != nil {
				log.Fatalf("Error creating output file: %s", err)
			}
			defer outputFile.Close()

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

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
