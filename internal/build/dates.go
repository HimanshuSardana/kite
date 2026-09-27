package build

import (
	"sort"
	"strings"
	"time"

	"github.com/HimanshuSardana/kite/pkg/content"
)

// postDateLayouts are tried in order when interpreting a frontmatter date.
// ISO day first (what `kite new` writes), then full timestamps so posts
// written on the same day can still be ordered, then legacy formats.
var postDateLayouts = []string{
	"2006-01-02",
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006/01/02",
	"02/01/2006",
	"2/1/06",
	"Jan 2006",
}

// parsePostDate interprets a frontmatter date, returning the zero time when
// nothing matches so callers can fall back to slug order.
func parsePostDate(s string) time.Time {
	s = strings.TrimSpace(s)
	for _, layout := range postDateLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

// sortPostSummaries orders newest-first by parsed date, breaking ties by slug
// so same-day posts never shuffle between builds.
func sortPostSummaries(posts []content.PostSummary) {
	sort.SliceStable(posts, func(i, j int) bool {
		ti, tj := parsePostDate(posts[i].Date), parsePostDate(posts[j].Date)
		if !ti.Equal(tj) {
			return ti.After(tj)
		}
		return posts[i].Slug < posts[j].Slug
	})
}

// displayPostDate formats a frontmatter date for cards and tag pages,
// keeping the raw string when it cannot be parsed.
func displayPostDate(raw string) string {
	if t := parsePostDate(raw); !t.IsZero() {
		return t.Format("Jan 2006")
	}
	return raw
}
