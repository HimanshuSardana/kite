---
title: "LaTeX that renders itself: math and diagrams"
date: 2026-09-26
tags: [mathematics, latex, diagrams]
---

Notes are more useful when the maths survives the trip from your head to the
page. This post is a tour of what the **notes** theme renders natively —
inline and display equations, aligned derivations, matrices, and diagrams
drawn in LaTeX itself (no screenshots, no exported images).

## Inline math

Wrap an expression in single dollar signs and it typesets in the flow of a
sentence: the identity $e^{i\pi} + 1 = 0$ ties five constants together, while
$a^2 + b^2 = c^2$ is the Pythagorean theorem. Subscripts, superscripts and
fractions all work: $\frac{1}{n}\sum_{k=1}^{n} k = \frac{n+1}{2}$.

## Display equations

A paragraph of its own, delimited by double dollar signs, is centred and
given room to breathe:

$$
x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}
$$

The Gaussian integral is a classic:

$$
\int_{-\infty}^{\infty} e^{-x^2}\,dx = \sqrt{\pi}
$$

## Aligned derivations

Use `aligned` to line up a chain of equalities on the `&` markers:

$$
\begin{aligned}
(a+b)^2 &= (a+b)(a+b) \\
        &= a^2 + ab + ba + b^2 \\
        &= a^2 + 2ab + b^2
\end{aligned}
$$

Maxwell's equations fit neatly into a two-column alignment:

$$
\begin{aligned}
\nabla \cdot \mathbf{E} &= \frac{\rho}{\varepsilon_0}, &
\nabla \cdot \mathbf{B} &= 0, \\
\nabla \times \mathbf{E} &= -\frac{\partial \mathbf{B}}{\partial t}, &
\nabla \times \mathbf{B} &= \mu_0 \mathbf{J} + \mu_0 \varepsilon_0 \frac{\partial \mathbf{E}}{\partial t}.
\end{aligned}
$$

## Matrices and vectors

$$
\begin{pmatrix} a & b \\ c & d \end{pmatrix}
\begin{pmatrix} x \\ y \end{pmatrix}
=
\begin{pmatrix} ax + by \\ cx + dy \end{pmatrix}
$$

## Commutative diagrams

Diagrams authored with the `AMScd` package render through MathJax — this is a
short exact sequence:

$$
\begin{CD}
0 @>>> A @>f>> B @>g>> C @>>> 0\\
@. @VV{\alpha}V @VV{\beta}V @VV{\gamma}V @.\\
0 @>>> A' @>>{f'}> B' @>>{g'}> C' @>>> 0
\end{CD}
$$

## TikZ diagrams

Arbitrary drawings go in a `tikz` fenced block. It is compiled to inline SVG
in the browser by TikZJax, so the figure is vector, themeable and generated
from the same LaTeX you would put in a paper. Here is a small directed graph:

```tikz
\begin{tikzpicture}[
  every node/.style={circle, draw, thick, inner sep=2pt, minimum size=7mm},
  every edge/.style={draw, thick, ->}
]
  \node (a) at (0,0) {$a$};
  \node (b) at (2.4,1.1) {$b$};
  \node (c) at (2.4,-1.1) {$c$};
  \node (d) at (4.8,0) {$d$};
  \draw (a) edge (b);
  \draw (a) edge (c);
  \draw (b) edge (d);
  \draw (c) edge (d);
  \draw[->, thick, bend left=22] (a) to node[draw=none, fill=none, midway, above, font=\small] {$f$} (d);
\end{tikzpicture}
```

A labelled right triangle, with the right angle marked:

```tikz
\begin{tikzpicture}[scale=1.1]
  \coordinate (A) at (0,0);
  \coordinate (B) at (3,0);
  \coordinate (C) at (3,2);
  \draw[thick] (A) -- (B) -- (C) -- cycle;
  \draw[thick] (2.75,0) -- (2.75,0.25) -- (3,0.25);
  \node[below left] at (A) {$A$};
  \node[below right] at (B) {$B$};
  \node[above right] at (C) {$C$};
  \node[below] at (1.5,0) {$b$};
  \node[right] at (3,1) {$a$};
  \node[above left] at (1.5,1) {$c$};
\end{tikzpicture}
```

And a plot of two trigonometric functions, drawn with TikZ's own `plot`
operation:

```tikz
\begin{tikzpicture}
  \draw[->, thick] (-0.4,0) -- (6.9,0) node[right] {$x$};
  \draw[->, thick] (0,-1.5) -- (0,1.6) node[above] {$y$};
  \draw[domain=0:6.283, samples=120, smooth, variable=\x, thick]
    plot (\x,{sin(\x r)});
  \draw[domain=0:6.283, samples=120, smooth, variable=\x, thick, dashed]
    plot (\x,{cos(\x r)});
  \draw[dashed] (1.571,-1.5) -- (1.571,1.6);
  \node[below] at (1.571,-1.5) {$\frac{\pi}{2}$};
\end{tikzpicture}
```

## Writing your own

The pattern is deliberately small: `$ … $` for inline, `$$ … $$` for display,
and a fenced `tikz` block for diagrams. Everything else is standard Markdown.
