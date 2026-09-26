package build

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// samplePost is a realistic Markdown post: frontmatter, headings, prose,
// a list, a fenced code block and a table. It is roughly the size of a
// short blog article.
const samplePost = `---
title: "Benchmark Post %d"
date: "2024-01-%02d"
tags: ["go", "performance"]
---

Kite turns Markdown into HTML. This paragraph exists to give the parser a
realistic amount of prose to tokenise, so the benchmark reflects an actual
post rather than an empty file.

## Why speed matters

A static site generator runs on every save while you are writing. If the
build takes seconds, you stop previewing. Kite keeps the whole pipeline in
memory and writes plain HTML, so even large sites rebuild in milliseconds.

- List item one
- List item two with some *emphasis* and a [link](https://example.com)
- List item three

### Code

` + "```go\n" + `func main() {
    fmt.Println("hello, kite")
}
` + "```\n" + `

### Data

| Command      | What it does            |
|--------------|-------------------------|
| kite build   | Build the static site   |
| kite serve   | Serve with live reload  |
| kite check   | Find broken links       |

## Wrapping up

That is the end of the post. It contains enough structure for the heading
walker, the renderer and the template pass to do real work.
`

func writeSamplePost(tb testing.TB, dir string, i int) {
	tb.Helper()
	name := filepath.Join(dir, fmt.Sprintf("post-%03d.md", i))
	body := fmt.Sprintf(samplePost, i, (i%28)+1)
	if err := os.WriteFile(name, []byte(body), 0o644); err != nil {
		tb.Fatal(err)
	}
}

// quiet redirects stdout to os.DevNull for the duration of fn so that the
// per-file "Compiling …" output of Build does not swamp the results.
func quiet(fn func()) {
	old := os.Stdout
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		fn()
		return
	}
	os.Stdout = devNull
	defer func() { os.Stdout = old; devNull.Close() }()
	fn()
}

// setupWorkspace creates a throwaway kite project with n posts, a config and
// a theme, returning the BuildOptions that point at it.
func setupWorkspace(tb testing.TB, posts int) BuildOptions {
	tb.Helper()
	root := tb.TempDir()
	contentDir := filepath.Join(root, "content")
	outputDir := filepath.Join(root, "output")
	staticDir := filepath.Join(root, "static")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		tb.Fatal(err)
	}
	if err := os.MkdirAll(staticDir, 0o755); err != nil {
		tb.Fatal(err)
	}
	for i := 0; i < posts; i++ {
		writeSamplePost(tb, contentDir, i)
	}

	configPath := filepath.Join(root, "config.yaml")
	config := "siteTitle: \"Bench\"\nauthorName: \"Bench\"\nsiteUrl: \"https://bench.local\"\n"
	if err := os.WriteFile(configPath, []byte(config), 0o644); err != nil {
		tb.Fatal(err)
	}

	// Use the repository's real theme so template rendering is measured too.
	themesDir, err := filepath.Abs(filepath.Join("..", "..", "themes"))
	if err != nil {
		tb.Fatal(err)
	}

	return BuildOptions{
		ThemeName:  "modern-light",
		ContentDir: contentDir,
		OutputDir:  outputDir,
		ThemesDir:  themesDir,
		StaticDir:  staticDir,
		ConfigPath: configPath,
	}
}

func BenchmarkParseMarkdown(b *testing.B) {
	path := filepath.Join(b.TempDir(), "post.md")
	if err := os.WriteFile(path, []byte(fmt.Sprintf(samplePost, 1, 1)), 0o644); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ParseMarkdown(path); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuild10Posts(b *testing.B)  { benchmarkBuild(b, 10) }
func BenchmarkBuild100Posts(b *testing.B) { benchmarkBuild(b, 100) }
func BenchmarkBuild500Posts(b *testing.B) { benchmarkBuild(b, 500) }

func benchmarkBuild(b *testing.B, posts int) {
	opts := setupWorkspace(b, posts)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		quiet(func() {
			if err := Build(opts); err != nil {
				b.Fatal(err)
			}
		})
	}
}

// BenchmarkBuildWebsite measures a real-world build: the docs site that this
// repository publishes, using its own theme and content.
func BenchmarkBuildWebsite(b *testing.B) {
	websiteDir, err := filepath.Abs(filepath.Join("..", "..", "website"))
	if err != nil {
		b.Fatal(err)
	}
	opts := BuildOptions{
		ThemeName:  "showcase",
		ContentDir: filepath.Join(websiteDir, "content"),
		OutputDir:  filepath.Join(b.TempDir(), "output"),
		ThemesDir:  filepath.Join(websiteDir, "themes"),
		StaticDir:  filepath.Join(websiteDir, "static"),
		ConfigPath: filepath.Join(websiteDir, "config.yaml"),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		quiet(func() {
			if err := Build(opts); err != nil {
				b.Fatal(err)
			}
		})
	}
}
