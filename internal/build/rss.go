package build

import (
	"fmt"
	"os"
	"time"

	"github.com/HimanshuSardana/kite/pkg/config"
	"github.com/HimanshuSardana/kite/pkg/content"
)

type RSSItem struct {
	Title   string
	Link    string
	Date    string
	Content string
}

type RSSFeed struct {
	Title       string
	Link        string
	Description string
	Items       []RSSItem
}

func GenerateRSS(outputDir, configPath, siteURL string, posts []content.PostSummary) error {
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	feed := RSSFeed{
		Title:       cfg.SiteTitle,
		Link:        siteURL,
		Description: cfg.AuthorBio,
		Items:       make([]RSSItem, 0, len(posts)),
	}

	for _, post := range posts {
		dateStr := post.Date
		var pubDate time.Time

		if t, err := time.Parse("2006-01-02", post.Date); err == nil {
			pubDate = t
			dateStr = pubDate.Format(time.RFC1123)
		} else if t, err := time.Parse("Jan 2006", post.Date); err == nil {
			pubDate = t
			dateStr = pubDate.Format(time.RFC1123)
		}

		feed.Items = append(feed.Items, RSSItem{
			Title:   post.Title,
			Link:    fmt.Sprintf("%s/%s/", siteURL, post.Slug),
			Date:    dateStr,
			Content: fmt.Sprintf("Read more at %s/%s/", siteURL, post.Slug),
		})
	}

	rssContent := renderRSS(feed)

	feedPath := "feed.xml"
	if outputDir != "" && outputDir != "." {
		feedPath = outputDir + "/feed.xml"
	}

	if err := os.WriteFile(feedPath, []byte(rssContent), 0644); err != nil {
		return fmt.Errorf("writing feed: %w", err)
	}

	fmt.Println("Feed written to", feedPath)
	return nil
}

func renderRSS(feed RSSFeed) string {
	updated := time.Now().Format(time.RFC1123)

	s := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>` + escapeXML(feed.Title) + `</title>
    <link>` + escapeXML(feed.Link) + `</link>
    <description>` + escapeXML(feed.Description) + `</description>
    <language>en-us</language>
    <lastBuildDate>` + updated + `</lastBuildDate>
    <atom:link href="` + escapeXML(feed.Link) + `/feed.xml" rel="self" type="application/rss+xml"/>
`

	for _, item := range feed.Items {
		s += `    <item>
      <title>` + escapeXML(item.Title) + `</title>
      <link>` + escapeXML(item.Link) + `</link>
      <guid>` + escapeXML(item.Link) + `</guid>
      <pubDate>` + item.Date + `</pubDate>
      <description>` + escapeXML(item.Content) + `</description>
    </item>
`
	}

	s += `  </channel>
</rss>`

	return s
}

func escapeXML(s string) string {
	s = replaceAll(s, "&", "&amp;")
	s = replaceAll(s, "<", "&lt;")
	s = replaceAll(s, ">", "&gt;")
	s = replaceAll(s, "\"", "&quot;")
	s = replaceAll(s, "'", "&apos;")
	return s
}

func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); i++ {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old) - 1
		} else {
			result += string(s[i])
		}
	}
	return result
}
