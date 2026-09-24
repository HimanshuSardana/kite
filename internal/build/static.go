package build

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyStatic mirrors every file under staticDir into outputDir,
// preserving relative paths and file modes. If staticDir does not
// exist it is a silent no-op so existing sites keep working.
// It returns the number of files copied.
func CopyStatic(staticDir, outputDir string) (int, error) {
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		return 0, nil
	}

	count := 0
	err := filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(staticDir, path)
		if err != nil {
			return fmt.Errorf("resolving %s: %w", path, err)
		}
		dst := filepath.Join(outputDir, rel)

		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return fmt.Errorf("creating dir for %s: %w", dst, err)
		}

		if err := copyStaticFile(path, dst, info.Mode()); err != nil {
			return fmt.Errorf("copying %s: %w", rel, err)
		}
		count++
		return nil
	})
	if err != nil {
		return count, err
	}
	return count, nil
}

func copyStaticFile(src, dst string, mode os.FileMode) error {
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
