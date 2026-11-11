package build

import (
	"fmt"
	"image/png"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/HimanshuSardana/kite/pkg/content"
)

// DefaultOGTemplate names the Typst source for open graph cards, resolved
// relative to the project root. config.yaml can point elsewhere with
// ogTemplate.
const DefaultOGTemplate = "og.typ"

// ogPPI keeps the render 1:1 with the template's point dimensions, so a
// 1200pt x 630pt page compiles to exactly 1200x630 pixels.
const ogPPI = 72

// OGInputs is the data handed to the Typst template as --input key=value
// pairs. Inside the template read them with sys.inputs.at("title"),
// sys.inputs.at("date"), and so on.
type OGInputs struct {
	Title  string
	Date   string
	Author string
	Site   string
	Slug   string
	URL    string
	Tags   string
}

// GenerateOGImages compiles templatePath once per post into
// <outputDir>/<slug>/og.png. A missing template or a missing typst binary is
// a warning, not a failure: cards stay opt-in and never break a build.
func GenerateOGImages(outputDir, templatePath, siteURL, authorName, siteTitle string, posts []content.PostSummary) (int, error) {
	if templatePath == "" {
		templatePath = DefaultOGTemplate
	}
	if _, err := os.Stat(templatePath); err != nil {
		log.Printf("Skipping OG images: template %q not found", templatePath)
		return 0, nil
	}
	if _, err := exec.LookPath("typst"); err != nil {
		log.Printf("Skipping OG images: typst not installed (see https://typst.app)")
		return 0, nil
	}

	baseURL := strings.TrimSuffix(siteURL, "/")
	generated := 0
	for _, post := range posts {
		outPath := filepath.Join(outputDir, post.Slug, "og.png")
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return generated, fmt.Errorf("creating og image directory: %w", err)
		}
		inputs := OGInputs{
			Title:  post.Title,
			Date:   ogDate(post),
			Author: authorName,
			Site:   siteTitle,
			Slug:   post.Slug,
			URL:    baseURL + "/" + post.Slug + "/",
			Tags:   strings.Join(post.Tags, ", "),
		}
		if err := compileOG(templatePath, outPath, inputs); err != nil {
			return generated, fmt.Errorf("og image for %q: %w", post.Slug, err)
		}
		generated++
	}
	return generated, nil
}

func compileOG(templatePath, outPath string, in OGInputs) error {
	cmd := exec.Command("typst",
		"compile", templatePath, outPath,
		"--ppi", strconv.Itoa(ogPPI),
		"--input", "title="+in.Title,
		"--input", "date="+in.Date,
		"--input", "author="+in.Author,
		"--input", "site="+in.Site,
		"--input", "slug="+in.Slug,
		"--input", "url="+in.URL,
		"--input", "tags="+in.Tags,
	)
	var stderr strings.Builder
	cmd.Stdout = io.Discard
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// ogDate renders a stable, human date for the card. Relative dates would go
// stale the moment the card is baked into a PNG.
func ogDate(post content.PostSummary) string {
	raw := post.Date
	if post.DateTime != "" {
		raw = post.DateTime
	}
	if t := parsePostDate(raw); !t.IsZero() {
		return t.Format("Jan 2, 2006")
	}
	return post.Date
}

func pngSize(path string) (int, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	cfg, err := png.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}
	return cfg.Width, cfg.Height, nil
}
