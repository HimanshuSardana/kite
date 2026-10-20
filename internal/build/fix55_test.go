package build

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckSiteSingleQuotedBrokenLink(t *testing.T) {
	dir := t.TempDir()
	page := `<html><body>` +
		`<a href='/missing-page/'>broken single</a>` +
		`<a href="/missing-double/">broken double</a>` +
		`<img src='/missing-img.png'>` +
		`<a href='/exists/'>ok</a>` +
		`</body></html>`
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(page), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "exists"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "exists", "index.html"), []byte("<html></html>"), 0o644); err != nil {
		t.Fatal(err)
	}

	broken, err := CheckSite(dir)
	if err != nil {
		t.Fatalf("CheckSite: %v", err)
	}
	targets := map[string]bool{}
	for _, b := range broken {
		targets[b.Target] = true
	}
	for _, want := range []string{"/missing-page/", "/missing-double/", "/missing-img.png"} {
		if !targets[want] {
			t.Errorf("CheckSite missed broken target %q (got %v)", want, broken)
		}
	}
	if targets["/exists/"] {
		t.Errorf("CheckSite flagged valid /exists/ as broken: %v", broken)
	}
}

func TestEscapeXMLMultibyte(t *testing.T) {
	in := "café 日本語 🎉 & <tag> \"q\" 'a'"
	got := escapeXML(in)
	// Multibyte runes must survive intact.
	for _, want := range []string{"café", "日本語", "🎉"} {
		if !strings.Contains(got, want) {
			t.Errorf("escapeXML(%q) = %q, missing %q", in, got, want)
		}
	}
	// ASCII specials must still be escaped.
	for _, want := range []string{"&amp;", "&lt;tag&gt;", "&quot;q&quot;", "&apos;a&apos;"} {
		if !strings.Contains(got, want) {
			t.Errorf("escapeXML(%q) = %q, missing %q", in, got, want)
		}
	}
}

func TestRenderRSSMultibyteTitle(t *testing.T) {
	feed := RSSFeed{
		Title:       "café 日本語 🎉",
		Link:        "https://example.com",
		Description: "desc",
		Items: []RSSItem{
			{Title: "naïve résumé — 東京 🍰", Link: "https://example.com/a/", Date: "Mon, 02 Jan 2006 15:04:05 MST", Content: "body café"},
		},
	}
	out := renderRSS(feed)
	for _, want := range []string{"café 日本語 🎉", "naïve résumé — 東京 🍰", "body café"} {
		if !strings.Contains(out, want) {
			t.Errorf("renderRSS missing %q in output", want)
		}
	}
	if strings.Contains(out, "�") {
		t.Errorf("renderRSS output contains replacement char: %q", out)
	}
}
