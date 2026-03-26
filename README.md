# Kite

A fast, minimal static site generator written in Go. Transform Markdown files into beautiful, themed websites with zero dependencies at runtime.

<p>
  <img src="https://img.shields.io/badge/version-1.0.0-blue.svg" alt="Version">
  <img src="https://img.shields.io/badge/Go-1.25+-00ADD8.svg" alt="Go Version">
</p>

## Installation

```bash
go install github.com/HimanshuSardana/kite@latest
```

## Usage

### Initialize a New Blog

```bash
kite init
```

This interactive command walks you through:
- Blog name and site title
- Author information
- Theme selection
- Creates `content/`, `output/`, `themes/` directories
- Generates config and a sample post

### Build Your Site

```bash
kite build
```

Or specify a theme:

```bash
kite build gruvbox
```

### Preview Locally

```bash
kite serve
```

Visit `http://localhost:8000` to see your site.

## Commands

| Command | Description |
|---------|-------------|
| `kite init` | Initialize a new blog project |
| `kite build` | Build the static site |
| `kite build <theme>` | Build with a specific theme |
| `kite serve` | Start local development server |
| `kite serve --port 8080` | Serve on custom port |
| `kite list-themes` | Show available themes |

## Configuration

Edit `config.yaml` to customize your site:

```yaml
siteTitle: "Your Blog Name"
authorName: "Your Name"
authorRole: "Writer & Developer"
authorBio: "A short bio about yourself"
defaultTheme: "modern-light"
siteUrl: "https://your-domain.com"
```

## Themes

Kite comes with 9 built-in themes:
- modern-light
- modern-dark
- modern-dark-2
- modern-dark-catppuccin
- everforest
- gruvbox
- rose-pine
- terminal-gruvbox
- tufte

