package content

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func slugsByPath(files []ContentFile) map[string]string {
	m := make(map[string]string, len(files))
	for _, f := range files {
		m[f.Path] = f.Slug
	}
	return m
}

func TestListContentFilesNestedDuplicateBasenames(t *testing.T) {
	contentDir := filepath.Join(t.TempDir(), "content")
	writeFile(t, filepath.Join(contentDir, "a", "hello.md"), "# A\n")
	writeFile(t, filepath.Join(contentDir, "b", "hello.md"), "# B\n")
	writeFile(t, filepath.Join(contentDir, "top.md"), "# Top\n")

	files, err := ListContentFiles(contentDir)
	if err != nil {
		t.Fatalf("ListContentFiles: %v", err)
	}
	if len(files) != 3 {
		t.Fatalf("got %d files, want 3: %v", len(files), files)
	}

	got := slugsByPath(files)
	want := map[string]string{
		filepath.Join(contentDir, "a", "hello.md"): "a/hello",
		filepath.Join(contentDir, "b", "hello.md"): "b/hello",
		filepath.Join(contentDir, "top.md"):        "top",
	}
	for path, wantSlug := range want {
		if got[path] != wantSlug {
			t.Errorf("path %s: slug = %q, want %q", path, got[path], wantSlug)
		}
	}

	// Slugs must be distinct so RSS/sitemap/home links don't collide.
	seen := map[string]bool{}
	for _, f := range files {
		if seen[f.Slug] {
			t.Errorf("duplicate slug %q", f.Slug)
		}
		seen[f.Slug] = true
	}
}

func TestGetOutputPathNested(t *testing.T) {
	contentDir := filepath.Join(t.TempDir(), "content")
	outputDir := filepath.Join(t.TempDir(), "output")

	got, err := GetOutputPath(contentDir, filepath.Join(contentDir, "a", "hello.md"), outputDir)
	if err != nil {
		t.Fatalf("GetOutputPath: %v", err)
	}
	want := filepath.Join(outputDir, "a", "hello", "index.html")
	if got != want {
		t.Errorf("GetOutputPath = %q, want %q", got, want)
	}
}

func TestGetOutputPathDotDirectory(t *testing.T) {
	contentDir := filepath.Join(t.TempDir(), "content")
	outputDir := filepath.Join(t.TempDir(), "output")

	got, err := GetOutputPath(contentDir, filepath.Join(contentDir, "my.md", "file.md"), outputDir)
	if err != nil {
		t.Fatalf("GetOutputPath: %v", err)
	}
	want := filepath.Join(outputDir, "my.md", "file", "index.html")
	if got != want {
		t.Errorf("GetOutputPath = %q, want %q", got, want)
	}
}

func TestSlugMatchesOutputPath(t *testing.T) {
	contentDir := filepath.Join(t.TempDir(), "content")
	outputDir := filepath.Join(t.TempDir(), "output")
	writeFile(t, filepath.Join(contentDir, "a", "hello.md"), "# A\n")
	writeFile(t, filepath.Join(contentDir, "my.md", "file.md"), "# Dots\n")

	files, err := ListContentFiles(contentDir)
	if err != nil {
		t.Fatalf("ListContentFiles: %v", err)
	}
	for _, f := range files {
		out, err := GetOutputPath(contentDir, f.Path, outputDir)
		if err != nil {
			t.Fatalf("GetOutputPath(%s): %v", f.Path, err)
		}
		want := filepath.Join(outputDir, filepath.FromSlash(f.Slug), "index.html")
		if out != want {
			t.Errorf("slug %q: output = %q, want %q", f.Slug, out, want)
		}
	}
}
