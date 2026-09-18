---
title: Why kite exists
date: 2026-09-05
tags: [philosophy]
---

Tired of bloated frameworks eating 100MBs of RAM on a small VPS, I wanted the smallest possible tool that turns Markdown into a website. Remembering Luke Smith's old videos about Hugo, I almost wrote a Hugo theme — then had the "bright" idea of writing the SSG itself.

> Enter **kite**.

## The philosophy

- **One binary.** `go install` and you're done. No Node, no bundler, no lockfile.
- **Millisecond builds.** A handful of Markdown files compile before your finger leaves the Enter key.
- **Themes are just HTML.** Two files, Go template syntax you already know, inline CSS. Nothing to learn, nothing to eject from.
- **Batteries that matter.** Live-reload dev server, RSS feed generation, table of contents data — built in, not plugins.

## When not to use kite

Honest answer: if you need taxonomies, i18n, image pipelines or 500 themes, use Hugo. If you need interactivity, use Astro or Next.js. Kite is for blogs and docs that should build instantly and host anywhere — this philosophy page and the whole [showcase](../) are kite output.
