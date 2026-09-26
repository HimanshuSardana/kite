package build

import (
	"os"
	"path/filepath"
	"testing"
)

// TestBuiltinThemesParse guards every shipped theme: each one must provide
// both templates and both must parse with the site FuncMap. It catches
// syntax errors in a theme long before a build does.
func TestBuiltinThemesParse(t *testing.T) {
	themesDir := filepath.Join("..", "..", "themes")
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		t.Fatalf("reading themes dir: %v", err)
	}

	checked := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		checked++
		themePath := filepath.Join(themesDir, entry.Name())
		for _, file := range []string{"home.html", "layout.html"} {
			if _, err := os.Stat(filepath.Join(themePath, file)); err != nil {
				t.Errorf("%s: missing %s", entry.Name(), file)
				continue
			}
			if _, err := LoadTemplate(themePath, file); err != nil {
				t.Errorf("%s/%s does not parse: %v", entry.Name(), file, err)
			}
		}
	}

	if checked == 0 {
		t.Fatal("no themes found; is the themes directory in the wrong place?")
	}
}
