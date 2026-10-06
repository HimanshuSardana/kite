package build

// Scrollytelling blocks (pomb.us "build your own react" style).
//
// Authoring (in Markdown):
//
//	::: scrolly
//	```js focus=1
//	const element = <h1>Hello</h1>
//	```
//	```js focus=1:3
//	const element = React.createElement("h1", null, "Hello")
//	```
//	```tex
//	element = \{type, props\}
//	```
//
//	---
//	prose for step 0, in **Markdown**
//	---
//	prose for step 1
//	:::
//
// Fenced code/tex blocks become animated steps (in order); prose sections
// split on a `---` line become the scrolling column. Step i shows code/tex
// step i (clamped). Rendering is static HTML + a small vanilla runtime in
// the scrolly theme: sticky stage, IntersectionObserver picks the active
// step, code lines morph with FLIP (stable data-lid across steps via LCS),
// tex steps crossfade and are typeset by MathJax.

import (
	"encoding/json"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type scrollyLine struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Focus bool   `json:"focus"`
}

type scrollyStep struct {
	Lang  string        `json:"lang"` // "js", "tex", ...
	Lines []scrollyLine `json:"lines,omitempty"`
	Tex   string        `json:"tex,omitempty"`
	Title string        `json:"title,omitempty"`
}

type scrollyBlock struct {
	Steps []scrollyStep
	Prose []string // markdown per step
}

var scrollyFenceRe = regexp.MustCompile("(?m)^```([a-zA-Z0-9+_-]*)[ \t]*([^\n]*)\n([\\s\\S]*?)^```[ \t]*$")

// extractScrollyBlocks pulls ::: scrolly blocks out of the markdown source,
// replacing each with a block-HTML placeholder the markdown renderer passes
// through untouched.
func extractScrollyBlocks(src string) (string, []scrollyBlock) {
	var blocks []scrollyBlock
	lines := strings.Split(src, "\n")
	var out []string
	i := 0
	for i < len(lines) {
		if strings.TrimSpace(lines[i]) == "::: scrolly" {
			j := i + 1
			depth := 1
			for j < len(lines) {
				t := strings.TrimSpace(lines[j])
				if t == "::: scrolly" {
					depth++
				} else if t == ":::" {
					depth--
					if depth == 0 {
						break
					}
				}
				j++
			}
			if j >= len(lines) {
				out = append(out, lines[i:]...)
				break
			}
			inner := strings.Join(lines[i+1:j], "\n")
			blocks = append(blocks, parseScrollyBlock(inner))
			idx := len(blocks) - 1
			out = append(out, "", fmt.Sprintf("<div class=\"kite-scrolly-ph\" data-i=\"%d\"></div>", idx), "")
			i = j + 1
			continue
		}
		out = append(out, lines[i])
		i++
	}
	return strings.Join(out, "\n"), blocks
}

func parseScrollyBlock(inner string) scrollyBlock {
	var b scrollyBlock
	// Split off prose: everything after the last fence.
	fences := scrollyFenceRe.FindAllStringSubmatchIndex(inner, -1)
	end := 0
	if len(fences) > 0 {
		end = fences[len(fences)-1][1]
	}
	head := inner
	proseRaw := ""
	if end > 0 {
		head = inner[:end]
		proseRaw = inner[end:]
	}
	for _, m := range scrollyFenceRe.FindAllStringSubmatch(head, -1) {
		lang := strings.ToLower(strings.TrimSpace(m[1]))
		meta := strings.TrimSpace(m[2])
		code := strings.Trim(m[3], "\n")
		if lang == "tex" || lang == "latex" || lang == "math" {
			b.Steps = append(b.Steps, scrollyStep{Lang: "tex", Tex: code})
			continue
		}
		if lang == "" {
			lang = "text"
		}
		focus := parseFocus(meta, code)
		rawLines := strings.Split(code, "\n")
		stepLines := make([]scrollyLine, len(rawLines))
		for k, ln := range rawLines {
			stepLines[k] = scrollyLine{Text: ln, Focus: focus[k+1]}
		}
		b.Steps = append(b.Steps, scrollyStep{Lang: lang, Lines: stepLines})
	}
	// Prose split on --- lines.
	var prose []string
	if strings.TrimSpace(proseRaw) == "" {
		prose = []string{"Scroll to continue."}
	} else {
		cur := []string{}
		for _, ln := range strings.Split(proseRaw, "\n") {
			if strings.TrimSpace(ln) == "---" {
				prose = append(prose, strings.Join(cur, "\n"))
				cur = []string{}
				continue
			}
			cur = append(cur, ln)
		}
		prose = append(prose, strings.Join(cur, "\n"))
		// drop leading empty chunk before first ---
		if len(prose) > 1 && strings.TrimSpace(prose[0]) == "" {
			prose = prose[1:]
		}
	}
	for k := range prose {
		prose[k] = strings.Trim(prose[k], "\n")
		if strings.TrimSpace(prose[k]) == "" {
			prose[k] = "Continue."
		}
	}
	// Pad prose to steps.
	for len(prose) < len(b.Steps) {
		prose = append(prose, "Continue.")
	}
	b.Prose = prose
	assignScrollyLineIDs(&b)
	return b
}

// parseFocus understands "focus=1,3:5" (1-indexed, ranges with colon).
// Empty meta => everything focused.
func parseFocus(meta, code string) map[int]bool {
	n := len(strings.Split(code, "\n"))
	all := map[int]bool{}
	for k := 1; k <= n; k++ {
		all[k] = true
	}
	m := regexp.MustCompile(`focus=([0-9,:-]+)`).FindStringSubmatch(meta)
	if m == nil {
		return all
	}
	out := map[int]bool{}
	for _, part := range strings.Split(m[1], ",") {
		if strings.Contains(part, ":") {
			se := strings.SplitN(part, ":", 2)
			a, _ := strconv.Atoi(se[0])
			e, _ := strconv.Atoi(se[1])
			for k := a; k <= e; k++ {
				out[k] = true
			}
		} else if v, err := strconv.Atoi(part); err == nil {
			out[v] = true
		}
	}
	return out
}

// assignScrollyLineIDs gives stable IDs to identical lines across consecutive
// code steps (LCS on trimmed text), so the client can FLIP-morph them.
func assignScrollyLineIDs(b *scrollyBlock) {
	counter := 0
	newID := func() string {
		counter++
		return fmt.Sprintf("L%d", counter)
	}
	var prev []scrollyLine
	for si := range b.Steps {
		st := &b.Steps[si]
		if st.Lang == "tex" {
			continue
		}
		if prev == nil {
			for k := range st.Lines {
				st.Lines[k].ID = newID()
			}
		} else {
			a := trimmedTexts(prev)
			bt := trimmedTexts(st.Lines)
			match := lcsMatch(a, bt)
			used := map[string]bool{}
			for k := range st.Lines {
				if match[k] >= 0 && !used[prev[match[k]].ID] {
					st.Lines[k].ID = prev[match[k]].ID
					used[prev[match[k]].ID] = true
				} else {
					st.Lines[k].ID = newID()
				}
			}
		}
		prev = st.Lines
	}
}

func trimmedTexts(lines []scrollyLine) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = strings.TrimSpace(l.Text)
	}
	return out
}

// lcsMatch returns for each index in b the matched index in a (or -1).
func lcsMatch(a, b []string) []int {
	n, m := len(a), len(b)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if a[i] != "" && a[i] == b[j] {
				dp[i][j] = dp[i+1][j+1] + 1
			} else if dp[i+1][j] >= dp[i][j+1] {
				dp[i][j] = dp[i+1][j]
			} else {
				dp[i][j] = dp[i][j+1]
			}
		}
	}
	match := make([]int, m)
	for j := range match {
		match[j] = -1
	}
	i, j := 0, 0
	for i < n && j < m {
		if a[i] != "" && a[i] == b[j] {
			match[j] = i
			i++
			j++
		} else if dp[i+1][j] >= dp[i][j+1] {
			i++
		} else {
			j++
		}
	}
	return match
}

func renderMarkdownInline(md string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))
	r := mdhtml.NewRenderer(mdhtml.RendererOptions{Flags: mdhtml.CommonFlags})
	return string(markdown.Render(doc, r))
}

// renderScrollyBlock emits the static HTML for one block (stage shows step 0;
// full step data embedded as JSON for the runtime).
func renderScrollyBlock(b scrollyBlock, idx int) string {
	data, _ := json.Marshal(b.Steps)
	var sb strings.Builder
	fmt.Fprintf(&sb, `<div class="scrolly" data-scrolly="%d">`, idx)
	fmt.Fprintf(&sb, `<script type="application/json" class="scrolly-data">%s</script>`, html.EscapeString(string(data)))
	sb.WriteString(`<div class="scrolly-stage-col"><div class="scrolly-sticky"><div class="scrolly-stage">`)
	sb.WriteString(renderScrollyStepHTML(b.Steps, 0))
	sb.WriteString(`</div><div class="scrolly-dots">`)
	for k := range b.Steps {
		cls := "scrolly-dot"
		if k == 0 {
			cls += " on"
		}
		fmt.Fprintf(&sb, `<span class="%s" data-dot="%d"></span>`, cls, k)
	}
	sb.WriteString(`</div></div></div>`)
	sb.WriteString(`<div class="scrolly-prose-col">`)
	for k, prose := range b.Prose {
		step := k
		if step >= len(b.Steps) {
			step = len(b.Steps) - 1
		}
		fmt.Fprintf(&sb, `<section class="scrolly-step" data-step="%d">%s</section>`, step, renderMarkdownInline(prose))
	}
	sb.WriteString(`</div></div>`)
	return sb.String()
}

func renderScrollyStepHTML(steps []scrollyStep, n int) string {
	if len(steps) == 0 {
		return `<div class="scrolly-empty">Empty scrolly block.</div>`
	}
	if n < 0 {
		n = 0
	}
	if n >= len(steps) {
		n = len(steps) - 1
	}
	st := steps[n]
	if st.Lang == "tex" {
		return fmt.Sprintf(`<div class="scrolly-tex">$$%s$$</div>`, html.EscapeString(st.Tex))
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, `<pre class="scrolly-code lang-%s"><code>`, html.EscapeString(st.Lang))
	for _, ln := range st.Lines {
		text := ln.Text
		if strings.TrimSpace(text) == "" {
			text = " "
		}
		cls := "s-line"
		if !ln.Focus {
			cls += " dim"
		}
		fmt.Fprintf(&sb, `<div class="%s" data-lid="%s"><span>%s</span></div>`,
			cls, html.EscapeString(ln.ID), html.EscapeString(text))
	}
	sb.WriteString(`</code></pre>`)
	return sb.String()
}

// RenderScrollyPlaceholders swaps placeholders left by extractScrollyBlocks
// with final HTML. Call on rendered page HTML.
func RenderScrollyPlaceholders(rendered []byte, blocks []scrollyBlock) []byte {
	s := string(rendered)
	for i, b := range blocks {
		ph := fmt.Sprintf("<div class=\"kite-scrolly-ph\" data-i=\"%d\"></div>", i)
		// gomarkdown may wrap raw block html in <p> when not blank-separated;
		// handle both forms.
		s = strings.ReplaceAll(s, "<p>"+ph+"</p>", renderScrollyBlock(b, i))
		s = strings.ReplaceAll(s, ph, renderScrollyBlock(b, i))
	}
	return []byte(s)
}
