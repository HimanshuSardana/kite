package build

import (
	"testing"

	"github.com/HimanshuSardana/kite/pkg/content"
)

func TestParsePostDate(t *testing.T) {
	cases := []struct {
		in   string
		y, m int
		d, h int
		zero bool
	}{
		{"2026-09-27", 2026, 9, 27, 0, false},
		{"2026-09-27T12:38:00Z", 2026, 9, 27, 12, false},
		{"2026-09-27T01:28:00Z", 2026, 9, 27, 1, false},
		{"24/3/26", 2026, 3, 24, 0, false},
		{"Sep 2026", 2026, 9, 1, 0, false},
		{"not-a-date", 0, 0, 0, 0, true},
		{"", 0, 0, 0, 0, true},
	}

	for _, c := range cases {
		got := parsePostDate(c.in)
		if c.zero {
			if !got.IsZero() {
				t.Errorf("parsePostDate(%q) = %v, want zero", c.in, got)
			}
			continue
		}
		if got.IsZero() {
			t.Errorf("parsePostDate(%q) = zero, want a date", c.in)
			continue
		}
		if got.Year() != c.y || int(got.Month()) != c.m || got.Day() != c.d || got.Hour() != c.h {
			t.Errorf("parsePostDate(%q) = %v, want %04d-%02d-%02d %02d:00", c.in, got, c.y, c.m, c.d, c.h)
		}
	}
}

func TestSortPostSummariesNewestFirst(t *testing.T) {
	posts := []content.PostSummary{
		{Title: "Old", Slug: "old", Date: "2026-09-20"},
		{Title: "New", Slug: "new", Date: "2026-09-27T12:38:00Z"},
		{Title: "Mid", Slug: "mid", Date: "2026-09-27T01:28:00Z"},
		{Title: "Same A", Slug: "b-same", Date: "2026-09-27"},
		{Title: "Same B", Slug: "a-same", Date: "2026-09-27"},
	}

	sortPostSummaries(posts)

	want := []string{"new", "mid", "a-same", "b-same", "old"}
	for i, slug := range want {
		if posts[i].Slug != slug {
			t.Fatalf("position %d = %s, want %s (full order: %v)", i, posts[i].Slug, slug, slugs(posts))
		}
	}
}

func TestDisplayPostDate(t *testing.T) {
	if got := displayPostDate("2026-09-27T12:38:00Z"); got != "Sep 2026" {
		t.Errorf("displayPostDate(timestamp) = %q, want %q", got, "Sep 2026")
	}
	if got := displayPostDate("2026-09-27"); got != "Sep 2026" {
		t.Errorf("displayPostDate(date) = %q, want %q", got, "Sep 2026")
	}
	if got := displayPostDate("garbage"); got != "garbage" {
		t.Errorf("displayPostDate(garbage) = %q, want it unchanged", got)
	}
}

func slugs(posts []content.PostSummary) []string {
	out := make([]string, len(posts))
	for i, p := range posts {
		out[i] = p.Slug
	}
	return out
}
