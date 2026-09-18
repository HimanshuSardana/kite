---
title: Writing posts
date: 2026-09-02
tags: [guide, markdown]
---

Every `.md` file in `content/` becomes a page at `/<filename>/`. `hello.md` → `/hello/`.

## Frontmatter

Each file starts with a small block describing the post:

```md
---
title: My first post
date: 2026-09-18
tags: [go, ssg]
---

Your markdown here...
```

Dates in `YYYY-MM-DD` format are displayed nicely on the home page (e.g. "Sep 2026"). Tags render as pills under each post title.

## Markdown

Standard Markdown, rendered with gomarkdown: headings, lists, quotes, tables, fenced code blocks with language hints, links and images.

Headings automatically build the page's table of contents — themes receive it as structured data and can render a sidebar, a dropdown, or nothing at all.

## Home page + RSS

On every build, kite collects all posts (newest first) for the home page and regenerates `feed.xml`, so readers can subscribe without any plugins or config.
