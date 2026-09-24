package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/HimanshuSardana/kite/pkg/content"
)

// GenerateSitemap writes sitemap.xml (home + one entry per post) into
// outputDir. It mirrors the GenerateRSS post loop, so no new config needed.
func GenerateSitemap(outputDir, siteURL string, posts []content.PostSummary, tagSlugs []string) error {
	today := time.Now().Format("2006-01-02")

	s := `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
		`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n" +
		`  <url><loc>` + escapeXML(siteURL) + `/</loc><lastmod>` + today + `</lastmod></url>` + "\n"

	for _, post := range posts {
		lastmod := today
		if t, err := time.Parse("2006-01-02", post.Date); err == nil {
			lastmod = t.Format("2006-01-02")
		}
		s += `  <url><loc>` + escapeXML(siteURL) + `/` + escapeXML(post.Slug) +
			`/</loc><lastmod>` + lastmod + `</lastmod></url>` + "\n"
	}

	for _, slug := range tagSlugs {
		s += `  <url><loc>` + escapeXML(siteURL) + `/tag/` + escapeXML(slug) +
			`/</loc><lastmod>` + today + `</lastmod></url>` + "\n"
	}

	s += `</urlset>`

	sitemapPath := "sitemap.xml"
	if outputDir != "" && outputDir != "." {
		sitemapPath = outputDir + "/sitemap.xml"
	}

	if err := os.WriteFile(sitemapPath, []byte(s), 0o644); err != nil {
		return fmt.Errorf("writing sitemap: %w", err)
	}

	fmt.Println("Sitemap written to", sitemapPath)
	return nil
}

// InjectFeedDiscovery adds an RSS autodiscovery <link> to every built HTML
// page that does not already have one. This keeps theme templates free of
// SEO plumbing — including custom themes. Rebuilds are idempotent.
func InjectFeedDiscovery(outputDir, siteURL, feedTitle string) (int, error) {
	link := `<link rel="alternate" type="application/rss+xml" title="` +
		escapeXML(feedTitle) + `" href="` + escapeXML(siteURL) + `/feed.xml">`

	count := 0
	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		html := string(raw)

		if strings.Contains(html, `type="application/rss+xml"`) {
			return nil
		}
		idx := strings.Index(html, "</head>")
		if idx == -1 {
			return nil
		}

		html = html[:idx] + "  " + link + "\n" + html[idx:]
		if err := os.WriteFile(path, []byte(html), info.Mode()); err != nil {
			return err
		}
		count++
		return nil
	})
	if err != nil {
		return count, err
	}
	return count, nil
}
