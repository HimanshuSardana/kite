package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v2"
)

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			Padding(0, 0, 1, 0)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("235")).
			Padding(0, 1)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Background(lipgloss.Color("235")).
			Padding(0, 2).
			Margin(0, 1)

	buttonActiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("235")).
				Background(lipgloss.Color("86")).
				Padding(0, 2).
				Margin(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)

type InitModel struct {
	step        int
	blogName    string
	siteTitle   string
	authorName  string
	authorRole  string
	authorBio   string
	theme       string
	themeList   list.Model
	inputBuffer string
	quitting    bool
	finished    bool
}

func (m *InitModel) Init() tea.Cmd {
	themes := []string{
		"modern-light",
		"modern-dark",
		"modern-dark-2",
		"modern-dark-catppuccin",
		"everforest",
		"gruvbox",
		"rose-pine",
		"terminal-gruvbox",
		"tufte",
		"notes",
		"minimal-notes",
		"flectar",
	}

	items := make([]list.Item, len(themes))
	for i, t := range themes {
		items[i] = listItem{t}
	}

	m.themeList = list.New(items, list.NewDefaultDelegate(), 10, 20)
	m.themeList.Title = "Select a theme:"
	m.themeList.SetShowStatusBar(false)
	m.themeList.SetFilteringEnabled(false)
	m.themeList.SetWidth(40)

	return nil
}

type listItem struct {
	title string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return "" }
func (i listItem) FilterValue() string { return i.title }

func (m *InitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.step == 5 {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "enter" {
				selected := m.themeList.SelectedItem().(listItem)
				m.theme = selected.title
				m.step++
				m.finished = true
				return m, tea.Quit
			}
			if keyMsg.String() == "esc" {
				m.quitting = true
				return m, tea.Quit
			}
		}
		var cmd tea.Cmd
		m.themeList, cmd = m.themeList.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			return m.handleEnter()
		case "backspace":
			if len(m.inputBuffer) > 0 {
				m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
			}
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
			"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
			"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
			"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
			"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
			"-", "_", " ", ".":
			m.inputBuffer += msg.String()
		}
	}
	return m, nil
}

func (m *InitModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case 0:
		m.blogName = m.inputBuffer
		m.inputBuffer = ""
		m.step++
	case 1:
		m.siteTitle = m.inputBuffer
		if m.siteTitle == "" {
			m.siteTitle = m.blogName
		}
		m.inputBuffer = ""
		m.step++
	case 2:
		m.authorName = m.inputBuffer
		m.inputBuffer = ""
		m.step++
	case 3:
		m.authorRole = m.inputBuffer
		m.inputBuffer = ""
		m.step++
	case 4:
		m.authorBio = m.inputBuffer
		m.inputBuffer = ""
		m.step++
	case 5:
		selected := m.themeList.SelectedItem().(listItem)
		m.theme = selected.title
		m.step++
		m.finished = true
		return m, tea.Quit
	}
	return m, nil
}

func (m *InitModel) View() string {
	var s string

	switch m.step {
	case 0:
		s = headerStyle.Render("╭─── Kite Setup") + "\n" +
			"\n" + "What's the name of your blog?" + "\n" +
			"(e.g. my-tech-blog, dev-diary)" + "\n\n" +
			inputStyle.Render(m.inputBuffer+"_") + "\n\n" +
			helpStyle.Render("type to enter · enter to continue · esc to cancel")
	case 1:
		s = headerStyle.Render("╭─── Kite Setup") + "\n" +
			"\n" + "Site title (for the header):" + "\n\n" +
			inputStyle.Render(m.inputBuffer+"_") + "\n\n" +
			helpStyle.Render("type to enter · enter to continue · esc to cancel")
	case 2:
		s = headerStyle.Render("╭─── Kite Setup") + "\n" +
			"\n" + "Your name:" + "\n\n" +
			inputStyle.Render(m.inputBuffer+"_") + "\n\n" +
			helpStyle.Render("type to enter · enter to continue · esc to cancel")
	case 3:
		s = headerStyle.Render("╭─── Kite Setup") + "\n" +
			"\n" + "Your role (e.g. Developer, Writer):" + "\n\n" +
			inputStyle.Render(m.inputBuffer+"_") + "\n\n" +
			helpStyle.Render("type to enter · enter to continue · esc to cancel")
	case 4:
		s = headerStyle.Render("╭─── Kite Setup") + "\n" +
			"\n" + "Short bio:" + "\n\n" +
			inputStyle.Render(m.inputBuffer+"_") + "\n\n" +
			helpStyle.Render("type to enter · enter to continue · esc to cancel")
	case 5:
		s = headerStyle.Render("╭─── Kite Setup") + "\n\n" +
			m.themeList.View() + "\n" +
			helpStyle.Render("↑↓ to select · enter to confirm")
	default:
		s = "Setting up your blog..."
	}

	return s
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(dstPath, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, data, info.Mode())
	})
}

const themesURL = "https://github.com/HimanshuSardana/kite/archive/refs/heads/main.zip"

// siteConfig mirrors the on-disk config.yaml schema. It is marshaled with
// yaml.v2 so user input containing quotes or newlines stays valid YAML.
type siteConfig struct {
	SiteTitle    string `yaml:"siteTitle"`
	AuthorName   string `yaml:"authorName"`
	AuthorRole   string `yaml:"authorRole"`
	AuthorBio    string `yaml:"authorBio"`
	DefaultTheme string `yaml:"defaultTheme"`
}

func marshalSiteConfig(siteTitle, authorName, authorRole, authorBio, theme string) ([]byte, error) {
	cfg := siteConfig{
		SiteTitle:    siteTitle,
		AuthorName:   authorName,
		AuthorRole:   authorRole,
		AuthorBio:    authorBio,
		DefaultTheme: theme,
	}
	return yaml.Marshal(&cfg)
}

// alreadyInitialized reports whether configPath exists or contentDir exists
// and is non-empty — i.e. the directory looks like an initialized site.
func alreadyInitialized(configPath, contentDir string) bool {
	if st, err := os.Stat(configPath); err == nil && !st.IsDir() {
		return true
	}
	entries, err := os.ReadDir(contentDir)
	if err != nil {
		return false
	}
	return len(entries) > 0
}

func downloadThemes() error {
	fmt.Println("  Downloading themes from GitHub...")

	tmpDir, err := os.MkdirTemp("", "kite-themes-*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	zipPath := filepath.Join(tmpDir, "kite-main.zip")
	if err := downloadURLToFile(themesURL, zipPath); err != nil {
		return err
	}

	extractDir := filepath.Join(tmpDir, "extracted")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		return fmt.Errorf("creating extract dir: %w", err)
	}
	if err := extractZipDir(zipPath, extractDir); err != nil {
		return err
	}

	srcThemes, err := findThemesDir(extractDir)
	if err != nil {
		return err
	}

	if err := os.MkdirAll("themes", 0o755); err != nil {
		return err
	}

	// Same semantics as `cp -r <src>/themes/. themes/`: copy each entry.
	entries, err := os.ReadDir(srcThemes)
	if err != nil {
		return fmt.Errorf("reading extracted themes: %w", err)
	}
	for _, e := range entries {
		src := filepath.Join(srcThemes, e.Name())
		dst := filepath.Join("themes", e.Name())
		if e.IsDir() {
			if err := copyDir(src, dst); err != nil {
				return fmt.Errorf("copying theme %s: %w", e.Name(), err)
			}
			continue
		}
		data, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("copying theme file %s: %w", e.Name(), err)
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return fmt.Errorf("copying theme file %s: %w", e.Name(), err)
		}
	}

	fmt.Println("  ✓ Themes downloaded successfully")
	return nil
}

func downloadURLToFile(url, dest string) error {
	resp, err := http.Get(url) //nolint:gosec,noctx
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: unexpected status %s", resp.Status)
	}
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		return fmt.Errorf("failed to download: %w", err)
	}
	if err := out.Close(); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	return nil
}

func extractZipDir(zipPath, dest string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to extract: %w", err)
	}
	defer r.Close()
	for _, f := range r.File {
		// Guard against Zip Slip.
		clean := filepath.Clean(f.Name)
		if strings.Contains(clean, "..") {
			return fmt.Errorf("failed to extract: unsafe path %q", f.Name)
		}
		target := filepath.Join(dest, clean)
		if !strings.HasPrefix(target, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("failed to extract: unsafe path %q", f.Name)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("failed to extract: %w", err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("failed to extract: %w", err)
		}
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to extract: %w", err)
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return fmt.Errorf("failed to extract: %w", err)
		}
		if err := os.WriteFile(target, data, f.Mode()); err != nil {
			return fmt.Errorf("failed to extract: %w", err)
		}
	}
	return nil
}

// findThemesDir locates the extracted "themes" directory. The archive is
// expected to unpack to <extractDir>/kite-main/themes, but fall back to a
// walk so a changed top-level directory name does not break extraction.
func findThemesDir(extractDir string) (string, error) {
	preferred := filepath.Join(extractDir, "kite-main", "themes")
	if st, err := os.Stat(preferred); err == nil && st.IsDir() {
		return preferred, nil
	}
	var found string
	err := filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "themes" {
			found = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("locating themes: %w", err)
	}
	if found == "" {
		return "", fmt.Errorf("locating themes: no themes directory in archive")
	}
	return found, nil
}

func RunInit(force bool) error {
	m := &InitModel{}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}

	if m.quitting {
		fmt.Println("\nInit cancelled.")
		return nil
	}

	if !m.finished {
		fmt.Println("\nInit not complete.")
		return nil
	}

	fmt.Println("\n" + headerStyle.Render("Setting up your blog..."))

	theme := m.theme
	if theme == "" {
		theme = "modern-light"
	}

	if !force && alreadyInitialized("config.yaml", "content") {
		return fmt.Errorf("current directory looks already initialized (found existing config.yaml or non-empty content/); re-running `kite init` would overwrite your site — pass --force to proceed")
	}

	dirs := []string{"content", "output", "static", "themes"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("creating %s directory: %w", dir, err)
		}
	}

	if _, err := os.Stat("themes/modern-light/layout.html"); err == nil {
		fmt.Println("  ✓ Found themes in current directory")
	} else if entries, err := os.ReadDir("../themes"); err == nil && len(entries) > 0 {
		fmt.Println("  ✓ Copying themes from parent directory...")
		for _, e := range entries {
			if e.IsDir() {
				src := filepath.Join("../themes", e.Name())
				dst := filepath.Join("themes", e.Name())
				if err := copyDir(src, dst); err != nil {
					fmt.Printf("  ⚠ Could not copy theme %s: %v\n", e.Name(), err)
				}
			}
		}
		fmt.Println("  ✓ Themes copied successfully")
	} else {
		fmt.Println("  ⚠ No themes found - downloading from GitHub...")
		if err := downloadThemes(); err != nil {
			return fmt.Errorf("downloading themes: %w", err)
		}
	}

	siteTitle := m.siteTitle
	if siteTitle == "" {
		siteTitle = m.blogName
	}

	configContent, err := marshalSiteConfig(siteTitle, m.authorName, m.authorRole, m.authorBio, theme)
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	if err := os.WriteFile("config.yaml", configContent, 0o644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	sampleContent := "---\n" +
		"title: Welcome to Kite\n" +
		"date: 2026-01-01\n" +
		"tags: [getting-started]\n" +
		"---\n\n" +
		"This is your first post! Write your content in Markdown here.\n\n" +
		"## Getting Started\n\n" +
		"- Add more posts to the content/ directory\n" +
		"- Run `kite build` to generate your site\n" +
		"- Run `kite serve` to preview locally\n\n" +
		"Enjoy blogging!\n"

	if err := os.WriteFile("content/1.md", []byte(sampleContent), 0o644); err != nil {
		return fmt.Errorf("writing sample content: %w", err)
	}

	fmt.Println("  ✓ Created config.yaml")
	fmt.Println("  ✓ Created content/ directory")
	fmt.Println("  ✓ Created output/ directory")
	fmt.Println("  ✓ Created static/ directory (copied as-is to output/)")
	fmt.Println("  ✓ Created themes/ directory")
	fmt.Println("  ✓ Created sample post (content/1.md)")
	fmt.Println("\nRun `kite build` to generate your site!")
	fmt.Println("Run `kite serve` to preview it locally.")

	return nil
}
