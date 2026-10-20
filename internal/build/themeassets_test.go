package build

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testLayoutHTML = `<!doctype html>
<html><head><title>{{ .Title }}</title><link rel="stylesheet" href="/style.css" /></head>
<body>{{ .Content }}</body></html>
`

const testHomeHTML = `<!doctype html>
<html><head><title>home</title></head><body>home</body></html>
`

func writeFixtureSite(t *testing.T) BuildOptions {
	t.Helper()
	root := t.TempDir()

	contentDir := filepath.Join(root, "content")
	outputDir := filepath.Join(root, "output")
	staticDir := filepath.Join(root, "static")
	themesDir := filepath.Join(root, "themes")
	themeDir := filepath.Join(themesDir, "fixture")
	for _, d := range []string{contentDir, staticDir, themeDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	post := `---
title: "Hello"
date: "2024-01-02"
tags: ["go"]
---

Hello world.
`
	if err := os.WriteFile(filepath.Join(contentDir, "hello.md"), []byte(post), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(themeDir, "layout.html"), []byte(testLayoutHTML), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(themeDir, "home.html"), []byte(testHomeHTML), 0o644); err != nil {
		t.Fatal(err)
	}
	css := "body { color: red; }\n"
	if err := os.WriteFile(filepath.Join(themeDir, "style.css"), []byte(css), 0o644); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(root, "config.yaml")
	cfg := "siteTitle: \"Test\"\nauthorName: \"Tester\"\nsiteUrl: \"https://example.com\"\n"
	if err := os.WriteFile(configPath, []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}

	return BuildOptions{
		ThemeName:  "fixture",
		ContentDir: contentDir,
		OutputDir:  outputDir,
		ThemesDir:  themesDir,
		StaticDir:  staticDir,
		ConfigPath: configPath,
	}
}

func quietBuild(t *testing.T, opts BuildOptions) {
	t.Helper()
	old := os.Stdout
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		if err := Build(opts); err != nil {
			t.Fatalf("Build: %v", err)
		}
		return
	}
	os.Stdout = devNull
	defer func() { os.Stdout = old; devNull.Close() }()
	if err := Build(opts); err != nil {
		t.Fatalf("Build: %v", err)
	}
}

// TestBuildCopiesThemeAssets is the regression test for issue #48:
// a plain Build must copy the theme stylesheet next to the rendered pages.
func TestBuildCopiesThemeAssets(t *testing.T) {
	opts := writeFixtureSite(t)
	quietBuild(t, opts)

	postHTML := filepath.Join(opts.OutputDir, "hello", "index.html")
	if _, err := os.Stat(postHTML); err != nil {
		t.Fatalf("post HTML missing at %s: %v", postHTML, err)
	}

	cssOut := filepath.Join(opts.OutputDir, "style.css")
	got, err := os.ReadFile(cssOut)
	if err != nil {
		t.Fatalf("theme style.css missing in output: %v", err)
	}
	if string(got) != "body { color: red; }\n" {
		t.Errorf("output style.css = %q, want theme stylesheet", string(got))
	}
}

// TestBuildTagPagesHaveNoBase ensures the <base href="/"> hack is gone:
// tag pages must render without rewriting every relative URL.
func TestBuildTagPagesHaveNoBase(t *testing.T) {
	opts := writeFixtureSite(t)
	quietBuild(t, opts)

	tagPage := filepath.Join(opts.OutputDir, "tag", "go", "index.html")
	data, err := os.ReadFile(tagPage)
	if err != nil {
		t.Fatalf("tag page missing at %s: %v", tagPage, err)
	}
	if strings.Contains(strings.ToLower(string(data)), "<base") {
		t.Errorf("tag page contains <base> tag:\n%s", string(data))
	}
}

// TestBuiltinLayoutsUseRootAbsoluteCSS guards all shipped themes: no
// layout.html may reference the stylesheet via a depth-relative URL.
func TestBuiltinLayoutsUseRootAbsoluteCSS(t *testing.T) {
	themesDir, err := filepath.Abs(filepath.Join("..", "..", "themes"))
	if err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("no themes found")
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(themesDir, e.Name(), "layout.html"))
		if err != nil {
			t.Errorf("%s: reading layout.html: %v", e.Name(), err)
			continue
		}
		if strings.Contains(string(data), "../style.css") {
			t.Errorf("%s/layout.html still references ../style.css, want /style.css", e.Name())
		}
	}
}

// TestCopyThemeAssetsSkipsTemplates ensures the theme copier ships assets
// without leaking raw templates into the public output.
func TestCopyThemeAssetsSkipsTemplates(t *testing.T) {
	root := t.TempDir()
	themeDir := filepath.Join(root, "themes", "fx")
	outputDir := filepath.Join(root, "output")
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(themeDir, "layout.html"), []byte("tmpl"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(themeDir, "style.css"), []byte("css"), 0o644); err != nil {
		t.Fatal(err)
	}

	CopyThemeAssets(themeDir, outputDir)

	if _, err := os.Stat(filepath.Join(outputDir, "style.css")); err != nil {
		t.Errorf("style.css not copied: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "layout.html")); err == nil {
		t.Errorf("layout.html leaked into output dir")
	}
}
