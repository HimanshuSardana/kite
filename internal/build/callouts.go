package build

import (
	"regexp"
	"strings"
)

// Callouts are question/answer blocks written with ordinary Markdown
// blockquote syntax:
//
//	> [!question] Optional title
//	> Body, in full Markdown.

//	> [!solution] Optional title
//	> The answer.
//
// The Markdown parser merges consecutive `>` blocks into a single
// <blockquote>, so one blockquote may contain any number of markers. This file
// rewrites the rendered HTML: each marker starts a new callout, and everything
// up to the next marker is that callout's body. Question blocks render as a
// styled <div>; solution blocks render as a <details> so they can be revealed
// on demand.

// calloutMarker matches the marker plus any spaces before the title.
var calloutMarker = regexp.MustCompile(`^\[!(question|solution|answer)\][ \t]*`)

type calloutKind int

const (
	calloutQuestion calloutKind = iota
	calloutSolution
)

// RenderCallouts rewrites callout blockquotes in rendered HTML, recursing into
// nested blockquotes first.
func RenderCallouts(html []byte) []byte {
	s := string(html)
	var out strings.Builder
	i := 0

	for {
		k := strings.Index(s[i:], "<blockquote")
		if k < 0 {
			out.WriteString(s[i:])
			break
		}
		k += i
		out.WriteString(s[i:k])

		end := matchBlockquote(s, k)
		if end < 0 {
			// Unbalanced; leave the rest untouched.
			out.WriteString(s[k:])
			break
		}

		openEnd := strings.IndexByte(s[k:], '>') + k
		inner := s[openEnd+1 : end]
		inner = string(RenderCallouts([]byte(inner))) // nested callouts first

		transformed, changed := transformCalloutBlock(inner)
		if changed {
			out.WriteString(transformed)
		} else {
			// No marker here, but nested blockquotes may still have been rewritten.
			out.WriteString(s[k : openEnd+1])
			out.WriteString(inner)
			out.WriteString("</blockquote>")
		}
		i = end + len("</blockquote>")
	}

	return []byte(out.String())
}

// matchBlockquote returns the index of the </blockquote> matching the
// <blockquote ...> that starts at start, or -1 if it is unbalanced.
func matchBlockquote(s string, start int) int {
	gt := strings.IndexByte(s[start:], '>')
	if gt < 0 {
		return -1
	}
	i := start + gt + 1
	depth := 1

	for i < len(s) {
		nextOpen := strings.Index(s[i:], "<blockquote")
		nextClose := strings.Index(s[i:], "</blockquote>")
		if nextClose < 0 {
			return -1
		}
		if nextOpen >= 0 && nextOpen < nextClose {
			depth++
			i += nextOpen + len("<blockquote")
			gt := strings.IndexByte(s[i:], '>')
			if gt < 0 {
				return -1
			}
			i += gt + 1
			continue
		}
		depth--
		if depth == 0 {
			return i + nextClose
		}
		i += nextClose + len("</blockquote>")
	}
	return -1
}

// transformCalloutBlock splits a blockquote's inner HTML at callout markers and
// wraps each run in a callout. It returns ok=false when there is no marker.
func transformCalloutBlock(inner string) (string, bool) {
	markers := findCalloutMarkers(inner)
	if len(markers) == 0 {
		return inner, false
	}

	var out strings.Builder
	out.WriteString(inner[:markers[0].start]) // anything before the first marker

	for i, m := range markers {
		segEnd := len(inner)
		if i+1 < len(markers) {
			segEnd = markers[i+1].start
		}
		kind, title, body := parseCalloutSegment(inner[m.start:segEnd])
		writeCallout(&out, kind, title, body)
	}

	return out.String(), true
}

type calloutMarkerPos struct {
	start int
}

func findCalloutMarkers(inner string) []calloutMarkerPos {
	var found []calloutMarkerPos
	search := 0
	for {
		k := strings.Index(inner[search:], "<p>[!")
		if k < 0 {
			return found
		}
		k += search
		if calloutMarker.MatchString(inner[k+len("<p>"):]) {
			found = append(found, calloutMarkerPos{start: k})
		}
		search = k + len("<p>[!")
	}
}

// parseCalloutSegment turns one marker run into (kind, title, body HTML).
func parseCalloutSegment(segment string) (calloutKind, string, string) {
	content := segment[len("<p>"):]
	m := calloutMarker.FindStringSubmatch(content)

	kind := calloutQuestion
	if m[1] == "solution" || m[1] == "answer" {
		kind = calloutSolution
	}

	after := content[len(m[0]):]
	inline, rest := after, ""
	if closeP := strings.Index(after, "</p>"); closeP >= 0 {
		inline = after[:closeP]
		rest = after[closeP+len("</p>"):]
	}

	var title, bodyInline string
	if nl := strings.IndexByte(inline, '\n'); nl >= 0 {
		title = strings.TrimSpace(inline[:nl])
		bodyInline = inline[nl+1:]
	} else {
		title = strings.TrimSpace(inline)
	}
	bodyInline = strings.TrimLeft(bodyInline, " \t")

	var body strings.Builder
	if bodyInline != "" {
		body.WriteString("<p>")
		body.WriteString(bodyInline)
		body.WriteString("</p>")
	}
	body.WriteString(rest)

	return kind, title, body.String()
}

func writeCallout(w *strings.Builder, kind calloutKind, title, body string) {
	titleHTML := ""
	if title != "" {
		titleHTML = `<span class="callout-title">` + title + `</span>`
	}
	hasBody := strings.TrimSpace(body) != ""

	if kind == calloutSolution {
		w.WriteString(`<details class="callout callout-solution"><summary class="callout-head">`)
		w.WriteString(`<span class="callout-badge">A</span><span class="callout-label">Solution</span>`)
		w.WriteString(titleHTML)
		w.WriteString(`</summary><div class="callout-body">`)
		w.WriteString(body)
		w.WriteString(`</div></details>`)
		return
	}

	w.WriteString(`<div class="callout callout-question"><div class="callout-head">`)
	w.WriteString(`<span class="callout-badge">Q</span><span class="callout-label">Question</span>`)
	w.WriteString(titleHTML)
	w.WriteString(`</div>`)
	if hasBody {
		w.WriteString(`<div class="callout-body">`)
		w.WriteString(body)
		w.WriteString(`</div>`)
	}
	w.WriteString(`</div>`)
}
