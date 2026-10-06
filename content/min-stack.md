---
title: "Min stack: O(1) minimum with two stacks"
date: 2026-10-06
tags: [data-structures, demo]
---

A stack with `push`, `pop` and `getMin` — all in O(1). The trick is a second stack that remembers the minimum. Scroll: the illustration and the code morph together.

::: scrolly
```scene
stack S: []
stack M: []
```
```tex
\text{push, pop, getMin — all in } O(1)
```
```js focus=1:4
class MinStack {
  constructor() {
    this.s = [];  // main stack
    this.m = [];  // min stack
  }
  // push, pop, top, getMin ...
}

//> MinStack()
//  S = []   M = []
```
```scene
stack S: [3] highlight=top
stack M: [3] highlight=top
```
```tex
S = [3],\quad M = [3],\quad \min = 3
```
```js focus=1:4,7:8
push(x) {
  this.s.push(x);
  if (!this.m.length || x <= this.getMin())
    this.m.push(x);
}

//> push(3)
//  S = [3]   M = [3]
```
```scene
stack S: [3, 5] highlight=top
stack M: [3] highlight=top
```
```tex
S = [3, 5],\quad M = [3],\quad \min = 3
```
```js focus=1:4,7:8
push(x) {
  this.s.push(x);
  if (!this.m.length || x <= this.getMin())
    this.m.push(x);
}

//> push(5)
//  S = [3, 5]   M = [3]
```
```scene
stack S: [3, 5, 2] highlight=top
stack M: [3, 2] highlight=top
```
```tex
S = [3, 5, 2],\quad M = [3, 2],\quad \min = 2
```
```js focus=1:4,7:8
push(x) {
  this.s.push(x);
  if (!this.m.length || x <= this.getMin())
    this.m.push(x);
}

//> push(2)
//  S = [3, 5, 2]   M = [3, 2]
```
```scene
stack S: [3, 5] highlight=top
stack M: [3] highlight=top
```
```tex
S = [3, 5],\quad M = [3],\quad \min = 3
```
```js focus=1:3,6:7
pop() {
  if (this.top() === this.getMin()) this.m.pop();
  return this.s.pop();
}

//> pop()
//  S = [3, 5]   M = [3]
```
```scene
stack S: [3, 5] highlight=top
stack M: [3] highlight=top
```
```tex
\min = 3
```
```js focus=1:2,5:6
top()    { return this.s[this.s.length - 1]; }
getMin() { return this.m[this.m.length - 1]; }

//> getMin() -> 3
//  S = [3, 5]   M = [3]
```

---
A plain stack finds its minimum by scanning — O(n). The min stack keeps a second stack `M` alongside the main stack `S`. The top of `M` is always the current minimum, so `getMin` is a single peek.

---
`push(3)` onto an empty stack. `S` becomes `[3]`. `M` is empty, so `3` is trivially the smallest — it goes onto `M` too. Both tops agree: the minimum is `3`.

---
`push(5)`. `S` becomes `[3, 5]`, but `5` is bigger than the current minimum `3`, so `M` ignores it. The guard `x <= getMin()` is the whole trick: `M` only records a value when it lowers (or ties) the minimum.

---
`push(2)`. New low — `2` lands on both stacks. `M` now reads `[3, 2]` bottom-to-top: the minimum of every prefix of pushes, recoverable after pops. Minimum is `2`.

---
`pop()` removes `2` from `S`. Since the popped value equals the top of `M`, it is popped from `M` as well — exposing the previous minimum `3` underneath. Popping `5` instead would have left `M` untouched.

---
`getMin()` peeks at the top of `M`: `3`. No scan, no counter — just two stacks and one invariant: `M` holds the minimums in decreasing order of arrival.
:::
