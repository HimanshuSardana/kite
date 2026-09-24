package build

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/HimanshuSardana/kite/pkg/config"
	"github.com/HimanshuSardana/kite/pkg/content"
	"github.com/HimanshuSardana/kite/pkg/themes"
)

const (
	DefaultContentDir = "./content"
	DefaultOutputDir  = "./output"
	DefaultThemesDir  = "./themes"
	DefaultStaticDir  = "./static"
	DefaultConfigPath = "./config.yaml"
	DefaultThemeName  = "modern-light"
)

type BuildOptions struct {
	ThemeName     string
	ContentDir    string
	OutputDir     string
	ThemesDir     string
	StaticDir     string
	ConfigPath    string
	IncludeDrafts bool
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
	if opts.StaticDir == "" {
		opts.StaticDir = DefaultStaticDir
	}
	if opts.ConfigPath == "" {
		opts.ConfigPath = DefaultConfigPath
	}

	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		log.Printf("Warning: Could not load config: %v", err)
	}

	themePath := themes.GetThemePath(opts.ThemesDir, opts.ThemeName)

	files, err := content.ListContentFiles(opts.ContentDir)
	if err != nil {
		return fmt.Errorf("listing content files: %w", err)
	}

	summaries := make([]content.PostSummary, 0, len(files))
	skippedDrafts := 0

	for _, file := range files {
		start := time.Now()

		parsed, err := ParseMarkdown(file.Path)
		if err != nil {
			log.Printf("Error parsing %s: %v", file.Path, err)
			continue
		}

		if parsed.Frontmatter.Draft && !opts.IncludeDrafts {
			skippedDrafts++
			// Remove previously built output so drafts never leak
			// into production from an earlier --drafts build.
			if outputPath, err := content.GetOutputPath(opts.ContentDir, file.Path, opts.OutputDir); err == nil {
				os.RemoveAll(outputPath)
				os.Remove(filepath.Dir(outputPath)) // best-effort: drops the dir if now empty
			}
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

		elapsed := time.Since(start).Milliseconds()
		fmt.Printf("Compiling %s (%dms)\n", parsed.Frontmatter.Title, elapsed)
	}

	fmt.Println("All files processed!")
	if skippedDrafts > 0 {
		fmt.Printf("Skipped %d draft(s) (build with --drafts to include)\n", skippedDrafts)
	}

	if err := RenderHomePage(themePath, opts.OutputDir, opts.ConfigPath, summaries); err != nil {
		log.Printf("Error rendering home page: %v", err)
	}

	siteURL := "https://your-site.com"
	if cfg != nil && cfg.SiteURL != "" {
		siteURL = cfg.SiteURL
	}
	if err := GenerateRSS(opts.OutputDir, opts.ConfigPath, siteURL, summaries); err != nil {
		log.Printf("Error generating RSS feed: %v", err)
	}

	if err := GenerateSitemap(opts.OutputDir, siteURL, summaries); err != nil {
		log.Printf("Error generating sitemap: %v", err)
	}

	feedTitle := "RSS"
	if cfg != nil && cfg.SiteTitle != "" {
		feedTitle = cfg.SiteTitle
	}
	if n, err := InjectFeedDiscovery(opts.OutputDir, siteURL, feedTitle); err != nil {
		log.Printf("Error injecting feed discovery links: %v", err)
	} else if n > 0 {
		fmt.Printf("Feed discovery link added to %d pages\n", n)
	}

	if n, err := CopyStatic(opts.StaticDir, opts.OutputDir); err != nil {
		log.Printf("Error copying static files: %v", err)
	} else if n > 0 {
		fmt.Printf("Copied %d static files\n", n)
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
  new           Create a new post with frontmatter

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
