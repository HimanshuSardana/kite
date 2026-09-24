package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

func runNew(args []string) {
	if len(args) < 3 || args[2] == "--help" || args[2] == "-h" {
		fmt.Fprintln(os.Stderr, "Usage: kite new <name.md> [--draft]")
		os.Exit(1)
	}

	name := args[2]
	draft := false
	for _, arg := range args[3:] {
		if arg == "--draft" {
			draft = true
		}
	}

	// Bare names land in content/; paths with a slash are used as given.
	path := name
	if !strings.Contains(path, "/") {
		path = filepath.Join("content", path)
	}
	if !strings.HasSuffix(path, ".md") {
		path += ".md"
	}

	if _, err := os.Stat(path); err == nil {
		fmt.Fprintf(os.Stderr, "Error: %s already exists (refusing to overwrite)\n", path)
		os.Exit(1)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	title := titleFromSlug(strings.TrimSuffix(filepath.Base(path), ".md"))
	post := "---\n" +
		"title: " + title + "\n" +
		"date: " + time.Now().Format("2006-01-02") + "\n" +
		"tags: []\n"
	if draft {
		post += "draft: true\n"
	}
	post += "---\n\nWrite here.\n"

	if err := os.WriteFile(path, []byte(post), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Created", path)
}

// titleFromSlug turns hello-world.md's stem into "Hello World".
func titleFromSlug(slug string) string {
	words := strings.FieldsFunc(slug, func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})
	for i, w := range words {
		runes := []rune(w)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	if len(words) == 0 {
		return "Untitled"
	}
	return strings.Join(words, " ")
}
