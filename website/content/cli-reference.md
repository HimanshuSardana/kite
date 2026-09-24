---
title: CLI reference
date: 2026-09-04
tags: [reference]
---

Everything kite does, in one place.

## Commands

| Command | Description |
|---|---|
| `kite init` | Interactively scaffold a new site (content, themes, config, sample post) |
| `kite new hello.md` | Create a post with stamped frontmatter (`--draft` starts it as a draft) |
| `kite build` | Build with the `defaultTheme` from `config.yaml` |
| `kite build <theme>` | Build with a specific theme, e.g. `kite build gruvbox` |
| `kite serve` | Build + serve `output/` at `http://localhost:8000` with live reload |
| `kite serve --port 8080` | Serve on a custom port |
| `kite serve <theme>` | Preview a specific theme, e.g. `kite serve rose-pine` |
| `kite list-themes` | List themes in `./themes` |

Theme resolution order: explicit argument → `defaultTheme` in `config.yaml` → `modern-light`.

## Configuration

All keys in `config.yaml`:

| Key | Purpose |
|---|---|
| `siteTitle` | `<title>` of the home page, RSS feed title |
| `authorName` | Hero name, footer, RSS author |
| `authorRole` | Hero subtitle |
| `authorBio` | Hero paragraph, RSS description |
| `defaultTheme` | Theme used when none is passed on the CLI |
| `siteUrl` | Base URL used for absolute links in `feed.xml` |

## Output

```sh
output/
  index.html        # home page
  feed.xml          # RSS, regenerated every build
  <slug>/index.html
```

Point any static file server at `output/` and you're live.
