package build

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadingTimeMinutes(t *testing.T) {
	if got := ReadingTimeMinutes(0); got != 1 {
		t.Errorf("0 words = %d min, want 1", got)
	}
	if got := ReadingTimeMinutes(200); got != 1 {
		t.Errorf("200 words = %d min, want 1", got)
	}
	if got := ReadingTimeMinutes(201); got != 2 {
		t.Errorf("201 words = %d min, want 2", got)
	}
	if got := ReadingTimeMinutes(1000); got != 5 {
		t.Errorf("1000 words = %d min, want 5", got)
	}
}

func TestParseMarkdownReadingTime(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "post.md")
	body := "word " + strings.Repeat("word ", 400)
	src := "---\ntitle: T\ndate: 2026-01-02\n---\n\n" + body + "\n"
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	page, err := ParseMarkdown(path)
	if err != nil {
		t.Fatalf("ParseMarkdown: %v", err)
	}
	if page.ReadingTime < 2 {
		t.Errorf("ReadingTime = %d, want >= 2 for a 400-word post", page.ReadingTime)
	}
}
