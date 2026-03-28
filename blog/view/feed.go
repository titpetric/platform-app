package view

import (
	"context"
	"fmt"
	"html"
	"io"
	"io/fs"
	"time"

	"github.com/titpetric/platform-app/blog/model"
)

// FeedConfig holds configuration for generating Atom feeds.
type FeedConfig struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Language string `json:"language"`
	Author   Author `json:"author"`
}

// Author holds author information for the feed.
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// DefaultFeedConfig returns a default feed configuration.
func DefaultFeedConfig() *FeedConfig {
	return &FeedConfig{
		URL:      "https://blog.localhost",
		Title:    "Blog",
		Subtitle: "Articles and thoughts",
		Language: "en",
		Author: Author{
			Name:  "Author",
			Email: "author@example.com",
		},
	}
}

// AtomFeed generates an Atom XML feed for articles.
func (v *Views) AtomFeed(ctx context.Context, w io.Writer, articles []model.Article, contentFS fs.FS) error {
	config := DefaultFeedConfig()
	return v.atomFeed(ctx, w, articles, config, contentFS)
}

// AtomFeedWithConfig generates an Atom XML feed with custom configuration.
func (v *Views) AtomFeedWithConfig(ctx context.Context, w io.Writer, articles []model.Article, config *FeedConfig, contentFS fs.FS) error {
	return v.atomFeed(ctx, w, articles, config, contentFS)
}

func (v *Views) atomFeed(ctx context.Context, w io.Writer, articles []model.Article, config *FeedConfig, contentFS fs.FS) error {
	// Find the most recent article date
	newestDate := time.Now()
	for _, a := range articles {
		if a.Date != nil && a.Date.After(newestDate) {
			newestDate = *a.Date
		}
	}

	// Write feed header
	fmt.Fprintf(w, `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xml:base="%s">
  <title>%s</title>
  <subtitle>%s</subtitle>
  <link href="%s/feed.xml" rel="self"/>
  <link href="%s"/>
  <updated>%s</updated>
  <id>%s</id>
  <author>
    <name>%s</name>
    <email>%s</email>
  </author>
`,
		escapeXML(config.URL),
		escapeXML(config.Title),
		escapeXML(config.Subtitle),
		escapeXML(config.URL),
		escapeXML(config.URL),
		newestDate.Format(time.RFC3339),
		escapeXML(config.URL),
		escapeXML(config.Author.Name),
		escapeXML(config.Author.Email),
	)

	// Add entries for each article
	for _, article := range articles {
		if article.Date == nil {
			continue
		}

		var contentStr string
		if contentFS != nil {
			content, err := fs.ReadFile(contentFS, article.Filename)
			if err == nil {
				content = StripFrontMatter(content)
				contentStr = string(content)
			}
		}

		fmt.Fprintf(w, `  <entry>
    <title>%s</title>
    <link href="%s%s"/>
    <updated>%s</updated>
    <id>%s%s</id>
    <summary>%s</summary>
    <content xml:lang="%s" type="html">%s</content>
  </entry>
`,
			escapeXML(article.Title),
			escapeXML(config.URL),
			escapeXML(article.URL),
			article.Date.Format(time.RFC3339),
			escapeXML(config.URL),
			escapeXML(article.URL),
			escapeXML(article.Description),
			escapeXML(config.Language),
			escapeXML(contentStr),
		)
	}

	io.WriteString(w, `</feed>`)
	return nil
}

// escapeXML escapes special XML characters.
func escapeXML(s string) string {
	return html.EscapeString(s)
}
