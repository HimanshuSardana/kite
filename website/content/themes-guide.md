---
title: Building themes
date: 2026-09-03
tags: [guide, themes]
---

A kite theme is two plain HTML files using Go's `html/template` syntax. No new templating language, no build step for the theme itself:

```sh
themes/mytheme/
  home.html     # the index page
  layout.html   # every post page
```

## home.html

Receives the site config plus the post list:

```html
<h1>{{ .AuthorName }}</h1>
<p>{{ .AuthorBio }}</p>

<ul>
{{ range .Posts }}
  <li>
    <a href="/{{ .Slug }}/">{{ .Title }}</a>
    <span>{{ .Date }}</span>
  </li>
{{ end }}
</ul>
```

Available fields: `SiteTitle`, `AuthorName`, `AuthorRole`, `AuthorBio`, `Year`, and `Posts` (each with `Title`, `Slug`, `Date`, `Tags`).

Template helper: `{{ . | tagSlug }}` slugifies a tag for its index URL
(`My Tag` → `my-tag`), matching the generated `/tag/<slug>/` pages.

## layout.html

Receives the rendered post:

```html
<h1>{{ .Title }}</h1>

<nav>
{{ range .TOC }}
  <a href="#{{ .ID }}">{{ .Text }}</a>
{{ end }}
</nav>

{{ .Content }}
```

Fields: `Title`, `Content` (already-rendered HTML), `TOC` (heading `Text` + `ID` pairs), `Year`.

## Styling

Inline a `<style>` block (like this site's theme does) and there is nothing else to deploy — one self-contained HTML file per page. `kite serve` also copies `themes/<name>/style.css` to `output/style.css` for previewing; for `kite build`, copy any extra assets yourself.

## Ship it

```sh
kite build mytheme
```

Ten themes ship with kite to steal ideas from: `modern-light`, `modern-dark`, `everforest`, `gruvbox`, `rose-pine`, `tufte` and more. Run `kite list-themes` to see them all.
