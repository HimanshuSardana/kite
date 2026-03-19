---
title: "Writing a Static Site Generator in Golang"
---

I recently decided to write a static site generator as a small toy project to get more familiar with the Go standard library. I’ve always liked learning languages by building something practical, and a static site generator seemed like a perfect mix of file handling, templating, and content processing.

Another reason was infrastructure. My blog used to run on a small VPS from Contabo, but after they raised their prices I decided to move everything to a Raspberry Pi sitting on my desk. Since the Pi isn’t particularly powerful, I wanted the site to be as lean as possible — no runtime rendering, no databases, just plain static files served over HTTP.

> NOTE
> I decided to call it **Kite**.

The idea behind Kite is simple: take Markdown files, convert them to HTML, wrap them in a template, and write them to an output directory. That’s it.

---

## Project Structure

The generator assumes a simple structure:

```
.
├── content/
│   ├── index.md
│   └── blog/
│       └── first-post.md
├── output/
├── layout.html
└── main.go
```

* **content/** contains Markdown files
* **layout.html** is the HTML template
* **output/** is where generated pages go

Each Markdown file gets converted into an HTML file with the same relative path.

For example:

```
content/blog/post.md
```

becomes:

```
output/blog/post.html
```

---

## The Page Structure

First we define a simple struct that represents the data passed to our template.

```go
type Page struct {
	Title   string
	Content template.HTML
}
```

* **Title** will eventually hold the page title (hardcoded for now).
* **Content** contains the rendered HTML.

Notice the use of `template.HTML`. This tells Go’s template engine that the content is already trusted HTML and should not be escaped.

---

## Walking the Content Directory

The core of the generator uses `filepath.WalkDir` to recursively traverse the `content/` directory.

```go
err := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
```

For every file encountered, we:

1. Ignore directories
2. Check if the file ends with `.md`
3. Convert it to HTML
4. Render it into a template
5. Write the result to the output directory

Filtering Markdown files is straightforward:

```go
if strings.HasSuffix(d.Name(), ".md") {
	fmt.Println("Processing:", path)
```

---

## Converting Markdown to HTML

For Markdown parsing I used the `gomarkdown/markdown` package.

```go
func convertToHtml(path string) []byte {
	md, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading %s: %s", path, err)
	}

	html := markdown.ToHTML(md, nil, nil)
	return html
}
```

This function simply:

1. Reads the Markdown file
2. Converts it to HTML
3. Returns the result

No configuration, no extensions — just the default behavior.

---

## Creating the Page Object

After conversion we construct the page object:

```go
htmlContent := convertToHtml(path)

newPage := Page{
	Title:   "hello", // you can customize this later
	Content: template.HTML(htmlContent),
}
```

Right now the title is hardcoded, but later this could be extracted from:

* frontmatter
* the first heading
* metadata in the Markdown file

---

## Parsing the Layout Template

Next we load the HTML template that wraps the content.

```go
tmpl, err := template.ParseFiles("./layout.html")
if err != nil {
	log.Fatalf("Error parsing template: %s", err)
}
```

This happens once per file in the current version. It could be optimized by parsing the template only once before the directory walk, but for a small site this overhead is negligible.

---

## Computing the Output Path

One important step is preserving the directory structure.

```go
relPath, err := filepath.Rel(contentDir, path)
```

This converts something like:

```
content/blog/post.md
```

into:

```
blog/post.md
```

Then we swap the extension:

```go
outputFilePath := filepath.Join(
	outputDir,
	strings.Replace(relPath, ".md", ".html", 1),
)
```

Result:

```
output/blog/post.html
```

---

## Ensuring Directories Exist

Before writing the file we make sure the directory exists:

```go
err = os.MkdirAll(filepath.Dir(outputFilePath), 0o755)
```

This creates any missing folders in the path.

---

## Writing the Final HTML

Finally we create the output file and render the template:

```go
outputFile, err := os.Create(outputFilePath)
if err != nil {
	log.Fatalf("Error creating output file: %s", err)
}
defer outputFile.Close()

err = tmpl.Execute(outputFile, newPage)
if err != nil {
	log.Fatalf("Error generating output content: %s", err)
}
```

The template receives the `Page` struct and generates the final HTML.

---

## Example Layout Template

A minimal `layout.html` might look like this:

```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>{{ .Title }}</title>
</head>
<body>
  <main>
    {{ .Content }}
  </main>
</body>
</html>
```

When the generator runs, the Markdown content gets injected into `{{ .Content }}`.

---

## Running the Generator

Just run:

```
go run main.go
```

Example output:

```
Processing: content/index.md
Processing: content/blog/first-post.md
All files processed!
```

And the generated files appear in `output/`.

---

## Why This Is Nice

Even this tiny generator already has a few advantages:

* **Extremely fast builds**
* **No runtime dependencies**
* **Very small hosting requirements**
* **Full control over the pipeline**

Perfect for something like a Raspberry Pi server.

---

## Future Improvements

Kite is intentionally minimal right now, but there are plenty of directions to take it:

* Frontmatter support (YAML/TOML)
* Automatic title extraction
* Template caching
* Asset copying (CSS, images)
* Incremental builds
* Live reload during development
* RSS feed generation
* Tag and category pages

The nice part is that Go’s standard library already provides most of the building blocks needed.

---

## Final Thoughts

Writing a static site generator is one of those projects that’s small enough to finish but complex enough to teach useful concepts:

* filesystem traversal
* templating
* content pipelines
* project structure

And the result is something genuinely useful.

Kite might stay a simple tool, or it might slowly grow features as I need them. Either way, it’s already doing exactly what I wanted: generating a lightweight blog that my Raspberry Pi can serve effortlessly.
