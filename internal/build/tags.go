package build

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/HimanshuSardana/kite/pkg/content"
)

// TagInfo is one tag's index page: display name, URL slug, newest-first posts.
type TagInfo struct {
	Name  string
	Slug  string
	Posts []content.PostSummary
}

// CollectTags groups posts by tag slug ("Go" and "go" share a page).
// Tags that slugify to nothing are skipped. Result sorted by name.
func CollectTags(posts []content.PostSummary) []TagInfo {
	bySlug := map[string]*TagInfo{}
	for _, p := range posts {
		for _, tag := range p.Tags {
			slug := content.SlugifyTag(tag)
			if slug == "" {
				continue
			}
			t, ok := bySlug[slug]
			if !ok {
				t = &TagInfo{Name: tag, Slug: slug}
				bySlug[slug] = t
			}
			t.Posts = append(t.Posts, p)
		}
	}

	tags := make([]TagInfo, 0, len(bySlug))
	for _, t := range bySlug {
		sort.Slice(t.Posts, func(i, j int) bool {
			return t.Posts[i].Date > t.Posts[j].Date
		})
		tags = append(tags, *t)
	}
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})
	return tags
}

// RenderTagPages writes output/tag/<slug>/index.html per tag plus a
// output/tag/index.html directory of all tags, reusing the theme's
// layout.html so no theme changes are required.
func RenderTagPages(themePath, outputDir string, tags []TagInfo) error {
	if len(tags) == 0 {
		return nil
	}

	tmpl, err := LoadTemplate(themePath, "layout.html")
	if err != nil {
		return err
	}

	for _, t := range tags {
		list := "<p><a href=\"/\">\u2190 All posts</a> \u00b7 <a href=\"/tag/\">All tags</a></p>\n<ul>\n"
		for _, p := range t.Posts {
			list += fmt.Sprintf("  <li><a href=\"/%s/\">%s</a> — %s</li>\n",
				template.HTMLEscapeString(p.Slug),
				template.HTMLEscapeString(p.Title),
				template.HTMLEscapeString(displayDate(p.Date)))
		}
		list += "</ul>\n"

		page := Page{
			Title:   "Tag: " + t.Name,
			Content: template.HTML(list),
			Year:    time.Now().Year(),
		}
		if err := renderTagPage(tmpl, outputDir+"/tag/"+t.Slug+"/index.html", page); err != nil {
			return fmt.Errorf("rendering tag %s: %w", t.Name, err)
		}
	}

	index := "<p><a href=\"/\">\u2190 All posts</a></p>\n<ul>\n"
	for _, t := range tags {
		noun := "posts"
		if len(t.Posts) == 1 {
			noun = "post"
		}
		index += fmt.Sprintf("  <li><a href=\"/tag/%s/\">%s</a> (%d %s)</li>\n",
			t.Slug, template.HTMLEscapeString(t.Name), len(t.Posts), noun)
	}
	index += "</ul>\n"

	if err := RenderPage(tmpl, outputDir+"/tag/index.html", Page{
		Title:   "Tags",
		Content: template.HTML(index),
		Year:    time.Now().Year(),
	}); err != nil {
		return fmt.Errorf("rendering tag index: %w", err)
	}

	fmt.Printf("Tag pages written: %d tags\n", len(tags))
	return nil
}

func displayDate(raw string) string {
	if t, err := time.Parse("2006-01-02", raw); err == nil {
		return t.Format("Jan 2006")
	}
	return raw
}

// tagSlugs returns tag slugs for sitemap inclusion.
func tagSlugs(tags []TagInfo) []string {
	slugs := make([]string, 0, len(tags))
	for _, t := range tags {
		slugs = append(slugs, t.Slug)
	}
	return slugs
}

// renderTagPage renders like RenderPage but pins relative theme assets
// (e.g. ../style.css) to the site root with a <base> tag: tag pages live
// two levels deep (/tag/<slug>/), where theme-relative paths would 404.
// Tag content uses absolute links and carries no TOC, so this is safe.
func renderTagPage(tmpl *template.Template, outputPath string, page Page) error {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, page); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	html := buf.String()
	if !strings.Contains(html, "<base") {
		html = strings.Replace(html, "<head>", "<head>\n<base href=\"/\">", 1)
	}
	if err := os.MkdirAll(dirOf(outputPath), 0o755); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}
	if err := os.WriteFile(outputPath, []byte(html), 0o644); err != nil {
		return fmt.Errorf("writing tag page: %w", err)
	}
	return nil
}

func dirOf(path string) string {
	i := strings.LastIndex(path, "/")
	if i == -1 {
		return "."
	}
	return path[:i]
}
