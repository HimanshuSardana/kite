---
title: Markdown Kitchen Sink
date: 2026-09-25
tags: [demo, markdown]
---

Every component kite's Markdown renderer supports, in one page.

## Headings

### Third level

#### Fourth level

##### Fifth level

###### Sixth level

## Text styles

**Bold**, *italic*, ***bold italic***, ~~strikethrough~~, `inline code`, and a [link to the docs](https://go.dev). Also an external [GitHub link](https://github.com/HimanshuSardana/kite).

---

## Lists

- Unordered item one
- Unordered item two
  - Nested item
    - Another nested level
- Unordered item three

1. Ordered item one
2. Ordered item two
   1. Nested ordered item
3. Ordered item three

- [x] Completed task
- [ ] Pending task

## Blockquote

> Kite is a fast, minimal static site generator written in Go.
> Markdown in, website out.

> A second blockquote with **bold** and `code`.

## Code blocks

Inline: `kite build`.

```go
func main() {
    fmt.Println("Hello from kite")
}
```

```bash
go install github.com/HimanshuSardana/kite@latest
kite init
kite build
kite serve --port 8000
```

```js
const greeting = "kite";
console.log(`Markdown in, ${greeting} out`);
```

```text
A plain text code block.
Nothing highlighted here.
```

## Table

| Command | Description | Runtime |
|---|---|---|
| `kite init` | Scaffold a new site | Interactive |
| `kite build` | Build static output | Instant |
| `kite serve` | Preview with live reload | Dev only |
| `kite check` | Find broken internal links | Instant |

## Images

![A local image](/logo.svg)

![External image](https://raw.githubusercontent.com/facebook/react/main/fixtures/dom/public/react-logo.svg)

## Horizontal rule

Above and below.

---

## HTML passthrough

<p align="center">
  <strong>Raw HTML works too.</strong>
</p>

## Footnotes

This post was created with `kite new` and filled in with Markdown.
