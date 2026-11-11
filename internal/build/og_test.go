package build

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/HimanshuSardana/kite/pkg/content"
)

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
