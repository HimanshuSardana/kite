package themes

import (
	"os"
	"path/filepath"
)

type Theme struct {
	Name string
	Path string
}

func List(themesDir string) ([]Theme, error) {
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		return nil, err
	}

	var themes []Theme
	for _, entry := range entries {
		if entry.IsDir() {
			themes = append(themes, Theme{
				Name: entry.Name(),
				Path: filepath.Join(themesDir, entry.Name()),
			})
		}
	}

	return themes, nil
}

func GetThemePath(themesDir, themeName string) string {
	return filepath.Join(themesDir, themeName)
}

func ThemeExists(themesDir, themeName string) bool {
	path := GetThemePath(themesDir, themeName)
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
