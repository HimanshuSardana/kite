package build

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/HimanshuSardana/kite/pkg/config"
	"github.com/HimanshuSardana/kite/pkg/content"
)

type Page struct {
	Title   string
	Content template.HTML
	TOC     []TOCItem
	Year    int
}

type HomePageData struct {
	SiteTitle    string
	AuthorName   string
	AuthorRole   string
	AuthorBio    string
	DefaultTheme string
	Year         int
	Posts        []content.PostSummary
}

func RenderPage(tmpl *template.Template, outputPath string, data Page) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer outputFile.Close()

	if err := tmpl.Execute(outputFile, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	return nil
}

func LoadTemplate(themePath, templateFile string) (*template.Template, error) {
	tmpl, err := template.ParseFiles(filepath.Join(themePath, templateFile))
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}
	return tmpl, nil
}

func RenderHomePage(themePath, outputDir, configPath string, summaries []content.PostSummary) error {
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Date > summaries[j].Date
	})

	for i, p := range summaries {
		if t, err := time.Parse("2006-01-02", p.Date); err == nil {
			summaries[i].Date = t.Format("Jan 2006")
		}
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	data := HomePageData{
		SiteTitle:    cfg.SiteTitle,
		AuthorName:   cfg.AuthorName,
		AuthorRole:   cfg.AuthorRole,
		AuthorBio:    cfg.AuthorBio,
		DefaultTheme: cfg.DefaultTheme,
		Year:         time.Now().Year(),
		Posts:        summaries,
	}

	tmpl, err := LoadTemplate(themePath, "home.html")
	if err != nil {
		return err
	}

	outPath := filepath.Join(outputDir, "index.html")
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating index.html: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("rendering home page: %w", err)
	}
	fmt.Println("Home page written to", outPath)

	return nil
}
