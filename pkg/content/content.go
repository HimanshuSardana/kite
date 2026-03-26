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
}

type PostSummary struct {
	Title string
	Slug  string
	Date  string
	Tags  []string
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
