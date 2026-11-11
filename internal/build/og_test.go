package build

import (
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/HimanshuSardana/kite/pkg/content"
)

func writePNG(t *testing.T, path string, width, height int) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}

func writePostPage(t *testing.T, outputDir, slug, title string) string {
	t.Helper()
	path := filepath.Join(outputDir, slug, "index.html")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	html := "<html><head><title>" + title + "</title></head><body>post</body></html>"
	if err := os.WriteFile(path, []byte(html), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestInjectOGMetaAddsTagsOnce(t *testing.T) {
	outputDir := t.TempDir()
	posts := []content.PostSummary{{Title: "Hello & <World>", Slug: "hello"}}
	pagePath := writePostPage(t, outputDir, "hello", "Hello")
	writePNG(t, filepath.Join(outputDir, "hello", "og.png"), 1200, 630)

	n, err := InjectOGMeta(outputDir, "https://example.com/", posts)
	if err != nil {
		t.Fatalf("InjectOGMeta: %v", err)
	}
	if n != 1 {
		t.Fatalf("injected %d pages, want 1", n)
	}

	raw, err := os.ReadFile(pagePath)
	if err != nil {
		t.Fatal(err)
	}
	page := string(raw)
	for _, want := range []string{
		`property="og:image" content="https://example.com/hello/og.png"`,
		`property="og:image:width" content="1200"`,
		`property="og:image:height" content="630"`,
		`property="og:url" content="https://example.com/hello/"`,
		`property="og:title" content="Hello &amp; &lt;World&gt;"`,
		`name="twitter:card" content="summary_large_image"`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("page missing %q\n%s", want, page)
		}
	}
	if strings.Count(page, `property="og:image"`) != 1 {
		t.Errorf("og:image injected more than once:\n%s", page)
	}

	// Rebuilds must not stack tags on top of each other.
	if n, err := InjectOGMeta(outputDir, "https://example.com", posts); err != nil || n != 0 {
		t.Fatalf("second run: injected %d, err %v; want 0, nil", n, err)
	}
}

func TestInjectOGMetaSkipsPostsWithoutCard(t *testing.T) {
	outputDir := t.TempDir()
	writePostPage(t, outputDir, "no-card", "No Card")

	n, err := InjectOGMeta(outputDir, "https://example.com", []content.PostSummary{{Slug: "no-card"}})
	if err != nil {
		t.Fatalf("InjectOGMeta: %v", err)
	}
	if n != 0 {
		t.Fatalf("injected %d pages, want 0", n)
	}
}

func TestGenerateOGImagesSkipsMissingTemplate(t *testing.T) {
	outputDir := t.TempDir()
	posts := []content.PostSummary{{Slug: "hello"}}

	n, err := GenerateOGImages(outputDir, filepath.Join(t.TempDir(), "og.typ"), "https://example.com", "Me", "Blog", posts)
	if err != nil {
		t.Fatalf("missing template should not error, got %v", err)
	}
	if n != 0 {
		t.Fatalf("generated %d images, want 0", n)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "hello", "og.png")); !os.IsNotExist(err) {
		t.Fatalf("expected no og.png, stat err = %v", err)
	}
}

func TestGenerateOGImagesWithTypst(t *testing.T) {
	if _, err := exec.LookPath("typst"); err != nil {
		t.Skip("typst not installed")
	}

	template := filepath.Join(t.TempDir(), "og.typ")
	src := `#set page(width: 1200pt, height: 630pt, margin: 0pt)
#let title = sys.inputs.at("title")
#let date = sys.inputs.at("date")
#let author = sys.inputs.at("author")
#text(size: 48pt)[#title]
#text(size: 20pt)[#date — #author]
`
	if err := os.WriteFile(template, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	outputDir := t.TempDir()
	posts := []content.PostSummary{{
		Title:    "Typst & Go",
		Slug:     "typst-and-go",
		Date:     "2026-03-24",
		DateTime: "2026-03-24T00:00:00Z",
		Tags:     []string{"typst", "go"},
	}}

	n, err := GenerateOGImages(outputDir, template, "https://example.com", "Me", "My Blog", posts)
	if err != nil {
		t.Fatalf("GenerateOGImages: %v", err)
	}
	if n != 1 {
		t.Fatalf("generated %d images, want 1", n)
	}

	width, height, err := pngSize(filepath.Join(outputDir, "typst-and-go", "og.png"))
	if err != nil {
		t.Fatalf("reading card: %v", err)
	}
	if width != 1200 || height != 630 {
		t.Fatalf("card is %dx%d, want 1200x630", width, height)
	}
}

func TestOGDateIsStable(t *testing.T) {
	got := ogDate(content.PostSummary{Date: "2026-03-24"})
	if got != "Mar 24, 2026" {
		t.Fatalf("ogDate = %q, want %q", got, "Mar 24, 2026")
	}
	if got := ogDate(content.PostSummary{Date: "sometime"}); got != "sometime" {
		t.Fatalf("ogDate fallback = %q, want raw value", got)
	}
}
