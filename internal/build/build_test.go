package build

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testBuildDirs(t *testing.T, withPost bool) (contentDir, outputDir, themesDir, staticDir, configPath string) {
	t.Helper()
	base := t.TempDir()
	contentDir = filepath.Join(base, "content")
	outputDir = filepath.Join(base, "output")
	themesDir = filepath.Join(base, "themes")
	staticDir = filepath.Join(base, "static")
	for _, d := range []string{contentDir, outputDir, themesDir, staticDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	configPath = filepath.Join(base, "missing-config.yaml")
	if withPost {
		post := "---\ntitle: Hello\ndate: 2026-01-02\n---\n\nHi there\n"
		if err := os.WriteFile(filepath.Join(contentDir, "hello.md"), []byte(post), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return contentDir, outputDir, themesDir, staticDir, configPath
}

func TestBuildMissingThemeReturnsError(t *testing.T) {
	contentDir, outputDir, themesDir, staticDir, configPath := testBuildDirs(t, false)

	err := Build(BuildOptions{
		ThemeName:  "no-such-theme",
		ContentDir: contentDir,
		OutputDir:  outputDir,
		ThemesDir:  themesDir,
		StaticDir:  staticDir,
		ConfigPath: configPath,
	})
	if err == nil {
		t.Fatal("Build with missing theme = nil, want error")
	}
	if !strings.Contains(err.Error(), "no-such-theme") {
		t.Errorf("error %q does not mention theme name", err.Error())
	}
}

func TestBuildThemeMissingLayoutReturnsError(t *testing.T) {
	contentDir, outputDir, themesDir, staticDir, configPath := testBuildDirs(t, true)

	// Theme directory exists (so the upfront ThemeExists check passes)
	// but layout.html is missing, so loading the template must fail
	// with a returned error instead of exiting the process.
	if err := os.MkdirAll(filepath.Join(themesDir, "broken"), 0o755); err != nil {
		t.Fatal(err)
	}

	err := Build(BuildOptions{
		ThemeName:  "broken",
		ContentDir: contentDir,
		OutputDir:  outputDir,
		ThemesDir:  themesDir,
		StaticDir:  staticDir,
		ConfigPath: configPath,
	})
	if err == nil {
		t.Fatal("Build with missing layout.html = nil, want error")
	}
	if !strings.Contains(err.Error(), "broken") {
		t.Errorf("error %q does not mention theme name", err.Error())
	}
}

func TestListThemesBadDirReturnsError(t *testing.T) {
	badDir := filepath.Join(t.TempDir(), "does-not-exist")
	if _, err := ListThemes(badDir); err == nil {
		t.Fatal("ListThemes with bad dir = nil, want error")
	}
}

func TestListThemesSuccess(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"alpha", "beta"} {
		if err := os.MkdirAll(filepath.Join(dir, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	got, err := ListThemes(dir)
	if err != nil {
		t.Fatalf("ListThemes: %v", err)
	}
	if len(got) != 2 || got[0] != "alpha" || got[1] != "beta" {
		t.Errorf("ListThemes = %v, want [alpha beta]", got)
	}
}
