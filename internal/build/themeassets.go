package build

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// CopyThemeAssets mirrors non-template files from a theme directory
// (e.g. style.css) into the output directory, preserving relative paths.
// Template files (*.html) are skipped. A missing theme directory is a
// silent no-op so existing sites keep working.
func CopyThemeAssets(themePath, outputDir string) {
	info, err := os.Stat(themePath)
	if err != nil || !info.IsDir() {
		return
	}

	count := 0
	err = filepath.Walk(themePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ".html") {
			return nil
		}

		rel, err := filepath.Rel(themePath, path)
		if err != nil {
			return fmt.Errorf("resolving %s: %w", path, err)
		}
		dst := filepath.Join(outputDir, rel)

		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return fmt.Errorf("creating dir for %s: %w", dst, err)
		}

		if err := copyThemeAssetFile(path, dst, info.Mode()); err != nil {
			return fmt.Errorf("copying %s: %w", rel, err)
		}
		count++
		return nil
	})
	if err != nil {
		log.Printf("Error copying theme assets: %v", err)
	} else if count > 0 {
		fmt.Printf("Copied %d theme asset(s)\n", count)
	}
}

func copyThemeAssetFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
