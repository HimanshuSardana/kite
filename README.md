# Kite — Minimal Go Static Site Generator

Kite is a lightweight (2.8MB) static site generator written in Go. 

## Features

- Markdown to HTML conversion
- Multiple built-in CSS themes
- Simple layout templating
- Fast builds with Go
- Clean output structure



## Usage

1. Clone the repository
```bash
git clone https://github.com/HimanshuSardana/kite
cd kite
```

2. Install dependencies
```bash
go mod tidy
```

3. Run the generator
```bash
go run main.go
```

Or use the compiled binary:
```bash
./kite-release
```

---

To write new posts
- Add Markdown files inside the `content/` directory.
- Each file will be converted into its own page.
- Folder structure is preserved in output.

Example:
```
content/test.md → output/test/index.html
```

---

## Inbuilt Themes

Themes are located in the `themes/` directory.

Available themes include:

* `modern-light.css`
* `modern-dark.css`
* `everforest.css`
* `rose-pine.css`
* `terminal-gruvbox.css`
* `tufte.css`

To change a theme, update your layout or configuration to reference the desired CSS file.
