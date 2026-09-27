package build

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderCalloutsQuestionAndSolution(t *testing.T) {
	in := `<blockquote>
<p>[!question] Why does a process need a stack?
Because the call stack tracks active frames.</p>

<p>[!solution]
It stores return addresses and locals.</p>
</blockquote>`

	out := string(RenderCallouts([]byte(in)))

	for _, want := range []string{
		`<div class="callout callout-question">`,
		`<span class="callout-label">Question</span>`,
		`<span class="callout-title">Why does a process need a stack?</span>`,
		`<p>Because the call stack tracks active frames.</p>`,
		`<details class="callout callout-solution">`,
		`<summary class="callout-head">`,
		`<span class="callout-label">Solution</span>`,
		`<p>It stores return addresses and locals.</p>`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n---\n%s", want, out)
		}
	}
	if strings.Contains(out, "[!question]") || strings.Contains(out, "[!solution]") {
		t.Errorf("marker text leaked into output:\n%s", out)
	}
	if strings.Contains(out, "<blockquote>") {
		t.Errorf("callout blockquote should have been replaced:\n%s", out)
	}
}

func TestRenderCalloutsMarkerWithoutTitle(t *testing.T) {
	in := "<blockquote>\n<p>[!question]\nWhat is a thread?</p>\n</blockquote>"
	out := string(RenderCallouts([]byte(in)))

	if !strings.Contains(out, `<p>What is a thread?</p>`) {
		t.Errorf("body not extracted:\n%s", out)
	}
	if strings.Contains(out, `class="callout-title"`) {
		t.Errorf("empty title should be omitted:\n%s", out)
	}
}

func TestRenderCalloutsLeavesPlainBlockquote(t *testing.T) {
	in := "<blockquote>\n<p>Just a normal quote.</p>\n</blockquote>"
	out := string(RenderCallouts([]byte(in)))
	if out != in {
		t.Errorf("plain blockquote changed:\nwant: %s\ngot:  %s", in, out)
	}
}

func TestRenderCalloutsNested(t *testing.T) {
	in := `<blockquote>
<p>outer</p>
<blockquote>
<p>[!question] inner?</p>
</blockquote>
</blockquote>`
	out := string(RenderCallouts([]byte(in)))

	if !strings.Contains(out, "<p>outer</p>") {
		t.Errorf("outer content lost:\n%s", out)
	}
	if !strings.Contains(out, `<div class="callout callout-question">`) {
		t.Errorf("nested callout not rendered:\n%s", out)
	}
	if !strings.Contains(out, "</blockquote>") {
		t.Errorf("outer blockquote should be preserved:\n%s", out)
	}
}

func TestParseMarkdownCallouts(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "post.md")
	src := `---
title: Callouts
date: 2026-01-02
---

> [!question] What is a process?
> A program in execution.

> [!answer] And a thread?
> A unit of scheduling within a process.
`
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	page, err := ParseMarkdown(path)
	if err != nil {
		t.Fatalf("ParseMarkdown: %v", err)
	}
	content := string(page.Content)

	if !strings.Contains(content, `callout-question`) || !strings.Contains(content, `callout-solution`) {
		t.Errorf("callouts not rendered:\n%s", content)
	}
	if !strings.Contains(content, `<span class="callout-title">What is a process?</span>`) {
		t.Errorf("question title missing:\n%s", content)
	}
	if !strings.Contains(content, `<p>A unit of scheduling within a process.</p>`) {
		t.Errorf("answer body missing:\n%s", content)
	}
}
