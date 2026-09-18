---
title: Quickstart
date: 2026-09-01
tags: [guide, install]
---

Get from zero to a live site in about a minute.

## Install

You need Go 1.25 or newer, then:

```sh
go install github.com/HimanshuSardana/kite@latest
```

That drops a single `kite` binary on your `PATH`. No runtime dependencies, no Node, nothing else.

## Scaffold a site

```sh
mkdir mysite && cd mysite
kite init
```

`init` is interactive — it asks for your blog name, author info and a theme, then creates `content/`, `themes/`, `output/` and a `config.yaml` like this:

```yaml
siteTitle: "My Blog"
authorName: "Your Name"
authorRole: "Writer & Developer"
authorBio: "A short bio about yourself"
defaultTheme: "modern-light"
siteUrl: "https://your-domain.com"
```

## Write, build, serve

Drop Markdown files in `content/` (see [Writing posts](writing/)), then:

```sh
kite build            # build with the default theme -> output/
kite build gruvbox    # build with a specific theme
kite serve            # preview at http://localhost:8000 with live reload
kite serve --port 8080
kite list-themes      # show available themes
```

`output/` is plain static HTML — serve it with anything (Caddy, nginx, GitHub Pages). It even includes a generated `feed.xml`.
