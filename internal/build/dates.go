package build

import (
	"fmt"
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

// RelativeDate renders t as "X minutes/hours/days/weeks ago" relative to now,
// falling back to "Jan 2006" once it is a month or more old.
func RelativeDate(t, now time.Time) string {
	d := now.Sub(t)
	if d < 0 {
		d = 0
	}
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m <= 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	case d < 7*24*time.Hour:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	case d < 30*24*time.Hour:
		weeks := int(d.Hours() / (24 * 7))
		if weeks <= 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	default:
		return t.Format("Jan 2006")
	}
}

// displayPostDate formats a frontmatter date for cards and tag pages:
// relative while it is fresh, the plain month once it is older than a month.
func displayPostDate(raw string) string {
	if t := parsePostDate(raw); !t.IsZero() {
		return RelativeDate(t, time.Now())
	}
	return raw
}
