package build

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var linkRe = regexp.MustCompile(`(?:href|src)="([^"]+)"`)

// BrokenLink is one internal reference with no file behind it.
type BrokenLink struct {
	Page   string // output-relative page containing the link, e.g. "hello/index.html"
	Target string // the link target as written
}

// CheckSite walks built HTML pages and verifies every internal href/src.
// External URLs, fragments, and non-HTML schemes are skipped. It returns
// the broken links found, sorted for stable output.
func CheckSite(outputDir string) ([]BrokenLink, error) {
	if _, err := os.Stat(outputDir); err != nil {
		return nil, fmt.Errorf("no %s directory (run kite build first)", outputDir)
	}

	var broken []BrokenLink
	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		page, _ := filepath.Rel(outputDir, path)

		for _, m := range linkRe.FindAllStringSubmatch(string(raw), -1) {
			target := m[1]
			if !isInternal(target) {
				continue
			}
			if !resolves(outputDir, page, target) {
				broken = append(broken, BrokenLink{Page: page, Target: target})
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(broken, func(i, j int) bool {
		if broken[i].Page == broken[j].Page {
			return broken[i].Target < broken[j].Target
		}
		return broken[i].Page < broken[j].Page
	})
	return broken, nil
}

func isInternal(target string) bool {
	if target == "" || strings.HasPrefix(target, "#") {
		return false
	}
	lower := strings.ToLower(target)
	for _, scheme := range []string{"http://", "https://", "mailto:", "tel:", "data:", "javascript:"} {
		if strings.HasPrefix(lower, scheme) {
			return false
		}
	}
	if strings.HasPrefix(target, "//") {
		return false
	}
	return true
}

// resolves reports whether target lands on a real output file, following
// pretty URLs (/hello/ -> hello/index.html) and ignoring query/fragment.
func resolves(outputDir, page, target string) bool {
	clean := target
	if i := strings.IndexAny(clean, "?#"); i != -1 {
		clean = clean[:i]
	}
	if clean == "" {
		return true
	}

	var abs string
	if strings.HasPrefix(clean, "/") {
		abs = filepath.Join(outputDir, filepath.FromSlash(clean))
	} else {
		abs = filepath.Join(outputDir, filepath.Dir(page), clean)
	}

	if st, err := os.Stat(abs); err == nil {
		if !st.IsDir() {
			return true
		}
		if _, err := os.Stat(filepath.Join(abs, "index.html")); err == nil {
			return true
		}
		return false
	}
	// Extensionless pretty link: /hello -> hello/index.html
	if _, err := os.Stat(abs + ".html"); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(abs, "index.html")); err == nil {
		return true
	}
	return false
}
