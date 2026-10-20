package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
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

func typeText(t *testing.T, m *InitModel, s string) {
	t.Helper()
	for _, r := range s {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
		// Space arrives as KeySpace in bubbletea; mirror that.
		if r == ' ' {
			msg = tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{r}}
		}
		if _, err := m.Update(msg); err != nil {
			_ = err
		}
	}
}

func pressBackspace(m *InitModel) {
	msg := tea.KeyMsg{Type: tea.KeyBackspace}
	m.Update(msg)
}

func TestInitListsAllThemes(t *testing.T) {
	m := &InitModel{}
	m.Init()

	items := m.themeList.Items()
	if len(items) != len(initThemes) {
		t.Fatalf("picker has %d items, want %d (%v)", len(items), len(initThemes), initThemes)
	}

	seen := map[string]bool{}
	for _, it := range items {
		li, ok := it.(listItem)
		if !ok {
			t.Fatalf("picker item %T is not a listItem", it)
		}
		seen[li.title] = true
	}
	for _, want := range initThemes {
		if !seen[want] {
			t.Errorf("theme %q missing from picker", want)
		}
	}
	if !seen["magical"] {
		t.Errorf(`theme "magical" missing from picker`)
	}
}

func TestUpdateAcceptsPunctuation(t *testing.T) {
	m := &InitModel{}
	input := `@'""/:!?-_.,;()#&*=+<>[]{}|~$%^\`
	typeText(t, m, input)
	if m.inputBuffer != input {
		t.Fatalf("buffer = %q, want %q", m.inputBuffer, input)
	}
	if !utf8.ValidString(m.inputBuffer) {
		t.Fatalf("buffer is not valid UTF-8: %q", m.inputBuffer)
	}
}

func TestUpdateAcceptsUnicode(t *testing.T) {
	m := &InitModel{}
	input := "Héllo 日本語 🎉 test-user_2"
	typeText(t, m, input)
	if m.inputBuffer != input {
		t.Fatalf("buffer = %q, want %q", m.inputBuffer, input)
	}
	if !utf8.ValidString(m.inputBuffer) {
		t.Fatalf("buffer is not valid UTF-8: %q", m.inputBuffer)
	}
}

func TestBackspaceDeletesRuneNotByte(t *testing.T) {
	m := &InitModel{}
	typeText(t, m, "aé日🎉")
	if m.inputBuffer != "aé日🎉" {
		t.Fatalf("setup buffer = %q, want %q", m.inputBuffer, "aé日🎉")
	}

	pressBackspace(m) // removes 🎉
	if m.inputBuffer != "aé日" {
		t.Fatalf("after 1 backspace buffer = %q, want %q", m.inputBuffer, "aé日")
	}
	pressBackspace(m) // removes 日
	if m.inputBuffer != "aé" {
		t.Fatalf("after 2 backspaces buffer = %q, want %q", m.inputBuffer, "aé")
	}
	pressBackspace(m) // removes é (2 bytes)
	if m.inputBuffer != "a" {
		t.Fatalf("after 3 backspaces buffer = %q, want %q", m.inputBuffer, "a")
	}
	if !utf8.ValidString(m.inputBuffer) {
		t.Fatalf("buffer is not valid UTF-8: %q", m.inputBuffer)
	}

	pressBackspace(m)
	if m.inputBuffer != "" {
		t.Fatalf("after 4 backspaces buffer = %q, want empty", m.inputBuffer)
	}
	// Backspace on empty buffer must not panic.
	pressBackspace(m)
	if m.inputBuffer != "" {
		t.Fatalf("backspace on empty buffer = %q, want empty", m.inputBuffer)
	}
}
