package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestAlreadyInitialized(t *testing.T) {
	mkDir := func(t *testing.T) string {
		t.Helper()
		return t.TempDir()
	}

	t.Run("empty dir is not initialized", func(t *testing.T) {
		dir := mkDir(t)
		if alreadyInitialized(filepath.Join(dir, "config.yaml"), filepath.Join(dir, "content")) {
			t.Fatal("expected empty dir to not look initialized")
		}
	})

	t.Run("existing config.yaml counts", func(t *testing.T) {
		dir := mkDir(t)
		if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("siteTitle: x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if !alreadyInitialized(filepath.Join(dir, "config.yaml"), filepath.Join(dir, "content")) {
			t.Fatal("expected existing config.yaml to count as initialized")
		}
	})

	t.Run("non-empty content dir counts", func(t *testing.T) {
		dir := mkDir(t)
		content := filepath.Join(dir, "content")
		if err := os.MkdirAll(content, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(content, "1.md"), []byte("hi"), 0o644); err != nil {
			t.Fatal(err)
		}
		if !alreadyInitialized(filepath.Join(dir, "config.yaml"), content) {
			t.Fatal("expected non-empty content/ to count as initialized")
		}
	})

	t.Run("empty content dir alone does not count", func(t *testing.T) {
		dir := mkDir(t)
		content := filepath.Join(dir, "content")
		if err := os.MkdirAll(content, 0o755); err != nil {
			t.Fatal(err)
		}
		if alreadyInitialized(filepath.Join(dir, "config.yaml"), content) {
			t.Fatal("expected empty content/ without config.yaml to not count as initialized")
		}
	})
}

func TestMarshalSiteConfigRoundTrip(t *testing.T) {
	cases := map[string]struct {
		siteTitle, authorName, authorRole, authorBio, theme string
	}{
		"plain":       {"My Blog", "Jane", "Writer", "Hello world", "modern-light"},
		"quotes":      {`My "quoted" blog`, `Jane "JD" Doe`, `Dev "x"`, `Bio with "quotes" inside`, "gruvbox"},
		"newlines":    {"line1\nline2", "Jane\nDoe", "Role\nTitle", "bio line1\nbio line2\nbio line3", "rose-pine"},
		"yaml-tricky": {`a: b # c`, "'single'", "[brackets] {braces}", "a\n  indented\n\t tab", "tufte"},
		"empty":       {"", "", "", "", "modern-light"},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			data, err := marshalSiteConfig(tc.siteTitle, tc.authorName, tc.authorRole, tc.authorBio, tc.theme)
			if err != nil {
				t.Fatalf("marshalSiteConfig: %v", err)
			}
			var got siteConfig
			if err := yaml.Unmarshal(data, &got); err != nil {
				t.Fatalf("marshaled YAML does not parse: %v\n---\n%s", err, data)
			}
			if got.SiteTitle != tc.siteTitle || got.AuthorName != tc.authorName ||
				got.AuthorRole != tc.authorRole || got.AuthorBio != tc.authorBio ||
				got.DefaultTheme != tc.theme {
				t.Fatalf("round-trip mismatch: got %+v", got)
			}
		})
	}
}
