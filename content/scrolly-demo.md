---
title: "Build your own React — a scrollytelling demo"
date: 2026-10-06
tags: [demo]
---

We are going to rewrite React from scratch. Step by step. Scroll — the code on the left stays pinned and morphs as the story advances, just like pomb.us.

::: scrolly
```js focus=1
const element = <h1 title="foo">Hello</h1>
const container = document.getElementById("root")
ReactDOM.render(element, container)
```
```js focus=1:6
const element = React.createElement(
  "h1",
  { title: "foo" },
  "Hello"
)
const container = document.getElementById("root")
ReactDOM.render(element, container)
```
```js focus=1:7
const element = {
  type: "h1",
  props: {
    title: "foo",
    children: "Hello",
  },
}
const container = document.getElementById("root")
ReactDOM.render(element, container)
```
```tex
element = \{ type,\ props \}
```

---
On the first line we have the element, defined with JSX. It isn't even valid JavaScript, so in order to replace it with vanilla JS, first we need to replace it with valid JS.

---
JSX is transformed by build tools like Babel. The transformation is simple: replace the tags with a call to `createElement`, passing tag name, props and children.

---
`createElement` just builds an object from its arguments. So we can safely replace the call with its output. And this is what an element is — an object with `type` and `props`.

---
And the same engine animates LaTeX. An element is a pair:

$$
element = \{type, props\}
$$

Unchanged code lines keep their identity and slide; tex steps crossfade, typeset live by MathJax.
:::
