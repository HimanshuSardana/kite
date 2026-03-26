package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	step         int
	blogName     string
	siteTitle    string
	authorName   string
	authorRole   string
	authorBio    string
	theme        string
	themes       []string
	cursor       int
	inputBuffer  string
	quitting     bool
	focusedInput bool
}

func (m *InitModel) Init() tea.Cmd {
	m.themes = []string{
		"modern-light",
		"modern-dark",
		"modern-dark-2",
		"modern-dark-catppuccin",
		"everforest",
		"gruvbox",
		"rose-pine",
		"terminal-gruvbox",
		"tufte",
	}
	return nil
}

func (m *InitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, nil
		case "enter":
			return m.handleEnter()
		case "up":
			if m.step == 6 {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.themes) - 1
				}
			}
		case "down":
			if m.step == 6 {
				m.cursor++
				if m.cursor >= len(m.themes) {
					m.cursor = 0
				}
			}
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
		m.theme = m.inputBuffer
		m.inputBuffer = ""
		if m.theme != "" {
			m.step = 7
		} else {
			m.step++
		}
	case 6:
		m.theme = m.themes[m.cursor]
		m.step++
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
		s = headerStyle.Render("╭─── Kite Setup") + "\n" +
			"\n" + "Preferred theme (or press enter to skip):" + "\n\n" +
			inputStyle.Render(m.inputBuffer+"_") + "\n\n" +
			helpStyle.Render("type to enter · enter to skip · esc to cancel")
	case 6:
		s = headerStyle.Render("╭─── Kite Setup") + "\n\n" +
			"Select a theme:\n\n"
		for i, theme := range m.themes {
			if i == m.cursor {
				s += "  " + buttonActiveStyle.Render("● "+theme) + "\n"
			} else {
				s += "  " + buttonStyle.Render("○ "+theme) + "\n"
			}
		}
		s += "\n" + helpStyle.Render("↑↓ to select · enter to confirm")
	}

	return s
}

func RunInit() error {
	m := &InitModel{}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}

	if m.quitting {
		fmt.Println("\nInit cancelled.")
		return nil
	}

	fmt.Println("\n" + headerStyle.Render("Setting up your blog..."))

	theme := m.theme
	if theme == "" {
		theme = "modern-light"
	}

	dirs := []string{"content", "output", "themes"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("creating %s directory: %w", dir, err)
		}
	}

	siteTitle := m.siteTitle
	if siteTitle == "" {
		siteTitle = m.blogName
	}

	configContent := fmt.Sprintf(`siteTitle: "%s"
authorName: "%s"
authorRole: "%s"
authorBio: "%s"
defaultTheme: "%s"
`, siteTitle, m.authorName, m.authorRole, m.authorBio, theme)

	if err := os.WriteFile("config.yaml", []byte(configContent), 0o644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	sampleContent := "---\n" +
		"title: Welcome to Kite\n" +
		"date: 2026-01-01\n" +
		"tags: [getting-started]\n" +
		"---\n\n" +
		"# Welcome\n\n" +
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
	fmt.Println("  ✓ Created themes/ directory")
	fmt.Println("  ✓ Created sample post (content/1.md)")
	fmt.Println("\nRun `kite build` to generate your site!")
	fmt.Println("Run `kite serve` to preview it locally.")

	return nil
}
