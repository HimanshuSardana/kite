package build

import (
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/HimanshuSardana/kite/pkg/content"
	"github.com/HimanshuSardana/kite/pkg/themes"
)

const (
	DefaultContentDir = "./content"
	DefaultOutputDir  = "./output"
	DefaultThemesDir  = "./themes"
	DefaultConfigPath = "./config.yaml"
	DefaultThemeName  = "modern-light"
)

type BuildOptions struct {
	ThemeName  string
	ContentDir string
	OutputDir  string
	ThemesDir  string
	ConfigPath string
}

func Build(opts BuildOptions) error {
	if opts.ThemeName == "" {
		opts.ThemeName = DefaultThemeName
	}
	if opts.ContentDir == "" {
		opts.ContentDir = DefaultContentDir
	}
	if opts.OutputDir == "" {
		opts.OutputDir = DefaultOutputDir
	}
	if opts.ThemesDir == "" {
		opts.ThemesDir = DefaultThemesDir
	}
	if opts.ConfigPath == "" {
		opts.ConfigPath = DefaultConfigPath
	}

	themePath := themes.GetThemePath(opts.ThemesDir, opts.ThemeName)

	files, err := content.ListContentFiles(opts.ContentDir)
	if err != nil {
		return fmt.Errorf("listing content files: %w", err)
	}

	summaries := make([]content.PostSummary, 0, len(files))

	for _, file := range files {
		fmt.Println("Processing:", file.Path)

		parsed, err := ParseMarkdown(file.Path)
		if err != nil {
			log.Printf("Error parsing %s: %v", file.Path, err)
			continue
		}

		summaries = append(summaries, content.PostSummary{
			Title: parsed.Frontmatter.Title,
			Slug:  file.Slug,
			Date:  parsed.Frontmatter.Date,
			Tags:  parsed.Frontmatter.Tags,
		})

		outputPath, err := content.GetOutputPath(opts.ContentDir, file.Path, opts.OutputDir)
		if err != nil {
			log.Printf("Error computing output path: %v", err)
			continue
		}

		tmpl, err := LoadTemplate(themePath, "layout.html")
		if err != nil {
			log.Fatalf("Error loading template: %v", err)
		}

		page := Page{
			Title:   parsed.Frontmatter.Title,
			Content: template.HTML(parsed.Content),
			TOC:     parsed.TOC,
			Year:    time.Now().Year(),
		}

		if err := RenderPage(tmpl, outputPath, page); err != nil {
			log.Printf("Error rendering page: %v", err)
		}
	}

	fmt.Println("All files processed!")

	if err := RenderHomePage(themePath, opts.OutputDir, opts.ConfigPath, summaries); err != nil {
		log.Printf("Error rendering home page: %v", err)
	}

	return nil
}

func ListThemes(themesDir string) []string {
	if themesDir == "" {
		themesDir = DefaultThemesDir
	}

	themeList, err := themes.List(themesDir)
	if err != nil {
		log.Fatal("Error:", err)
	}

	result := make([]string, len(themeList))
	for i, t := range themeList {
		result[i] = t.Name
	}
	return result
}

func ShowHelpMessage() {
	fmt.Println(`
Kite — A lightweight static site generator

USAGE:
  kite <command> [options]

COMMANDS:
  build         Build the static site into the output directory
  serve         Start a local development server with live reload
  list-themes   List all available themes

OPTIONS:
  -h, --help    Show this help message

EXAMPLES:
  kite build
  kite serve
  kite serve --port 8080
  kite list-themes

DESCRIPTION:
  Kite converts your content into a static website using themes and templates.
  Use 'build' for production output and 'serve' for local development.
`)
}
