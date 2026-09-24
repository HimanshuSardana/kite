package content

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type Frontmatter struct {
	Title string   `yaml:"title"`
	Date  string   `yaml:"date"`
	Tags  []string `yaml:"tags"`
	Draft bool     `yaml:"draft"`
}

type PostSummary struct {
	Title       string
	Slug        string
	Date        string
	Tags        []string
	WordCount   int
	ReadingTime int
}

// SlugifyTag maps a display tag ("My Tag") to its URL slug ("my-tag").
// The same function backs generated /tag/ URLs and the tagSlug template
// helper, so pills and pages can never disagree.
func SlugifyTag(tag string) string {
	var b strings.Builder
	prevDash := false
	for _, r := range strings.ToLower(tag) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		case r == ' ' || r == '-' || r == '_':
			if !prevDash && b.Len() > 0 {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

type ContentFile struct {
	Path        string
	Slug        string
	Frontmatter Frontmatter
}

func ListContentFiles(contentDir string) ([]ContentFile, error) {
	var files []ContentFile

	err := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".md") {
			slug := strings.TrimSuffix(d.Name(), ".md")
			files = append(files, ContentFile{
				Path: path,
				Slug: slug,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func GetOutputPath(contentDir, contentPath, outputDir string) (string, error) {
	relPath, err := filepath.Rel(contentDir, contentPath)
	if err != nil {
		return "", err
	}

	outputFilePath := filepath.Join(outputDir, strings.Replace(relPath, ".md", "/index.html", 1))
	return outputFilePath, nil
}
