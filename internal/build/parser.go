package build

import (
	"fmt"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/HimanshuSardana/kite/pkg/content"
)

type TOCItem struct {
	Level int
	Text  string
	ID    string
}

type ParsedPage struct {
	Frontmatter content.Frontmatter
	Content     []byte
	TOC         []TOCItem
	ReadingTime int
}

func ParseMarkdown(path string) (*ParsedPage, error) {
	md, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var matter content.Frontmatter
	rest, err := frontmatter.Parse(strings.NewReader(string(md)), &matter)
	if err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	stripped, scrollyBlocks := extractScrollyBlocks(string(rest))
	rest = []byte(stripped)

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(rest)

	var toc []TOCItem
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if heading, ok := node.(*ast.Heading); ok && entering {
			text := extractText(heading)
			id := string(heading.HeadingID)

			toc = append(toc, TOCItem{
				Level: heading.Level,
				Text:  text,
				ID:    id,
			})
		}
		return ast.GoToNext
	})

	renderer := html.NewRenderer(html.RendererOptions{
		Flags: html.CommonFlags,
	})

	output := markdown.Render(doc, renderer)
	output = RenderCallouts(output)
	output = RenderScrollyPlaceholders(output, scrollyBlocks)

	words := len(strings.Fields(string(rest)))

	return &ParsedPage{
		Frontmatter: matter,
		Content:     output,
		TOC:         toc,
		ReadingTime: ReadingTimeMinutes(words),
	}, nil
}

// ReadingTimeMinutes estimates reading time at 200 words per minute,
// with a one-minute floor.
func ReadingTimeMinutes(words int) int {
	m := (words + 199) / 200
	if m < 1 {
		m = 1
	}
	return m
}

func extractText(h *ast.Heading) string {
	var text string
	ast.WalkFunc(h, func(node ast.Node, entering bool) ast.WalkStatus {
		if leaf, ok := node.(*ast.Text); ok && entering {
			text += string(leaf.Literal)
		}
		return ast.GoToNext
	})
	return text
}
