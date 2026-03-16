# Kite

**Kite** is a *blazingly fast* static site generator written in Go. It supports
multiple themes, boasts near-perfect Lighthouse Scores (99 on Desktop, 94 on
Mobile)

## Usage

Place your markdown content in the `content/` directory, the output will be in the `output/` directory
```txt
.
├── content
│   └── blog
│       └── test.md
└── output
    └── blog
        └── test.html
 ```

 Run the following command:
 ```bash
 kite
 ```

