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
	Images      []string
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

	images := extractImages(doc)

	return &ParsedPage{
		Frontmatter: matter,
		Content:     output,
		TOC:         toc,
		Images:      images,
	}, nil
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

func extractImages(doc ast.Node) []string {
	var images []string
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if img, ok := node.(*ast.Image); ok && entering {
			if len(img.Destination) > 0 {
				images = append(images, string(img.Destination))
			}
		}
		return ast.GoToNext
	})
	return images
}
