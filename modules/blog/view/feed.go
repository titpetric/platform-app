package view

import (
	"context"
	"fmt"
	"html"
	"io"
	"os"
	"time"

	"github.com/titpetric/platform-app/modules/blog/model"
)

// AtomFeed generates an Atom XML feed for articles
func (v *Views) AtomFeed(ctx context.Context, w io.Writer, articles []model.Article, meta map[string]any) error {
	var author map[string]any
	if m, ok := meta["author"].(map[string]any); ok {
		author = m
	}

	// Find the most recent article date
	newestDate := time.Now()
	if len(articles) > 0 {
		newestDate = *articles[0].Date
		for _, a := range articles {
			if a.Date == nil {
				continue
			}
			articleDate := *a.Date
			if articleDate.After(newestDate) {
				newestDate = articleDate
			}
		}
	}

	var language string
	if lang, ok := meta["language"].(string); ok {
		language = lang
	}

	io.WriteString(w, fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xml:base="%s">
  <title>%s</title>
  <subtitle>Blogging general thoughts and rambles, code snippets, and front-end web dev discoveries</subtitle>
  <link href="%s/feed.xml" rel="self"/>
  <link href="%s"/>
  <updated>%s</updated>
  <id>%s</id>
  <author>
    <name>%s</name>
    <email>%s</email>
  </author>
`, meta["url"], escapeXML(meta["title"]), meta["url"], meta["url"], newestDate.Format(time.RFC3339), meta["url"], escapeXML(author["name"]), author["email"]))

	// Add entries for each article
	for _, article := range articles {

		content, err := os.ReadFile(article.Filename)
		if err != nil {
			return err
		}

		// Strip front matter from content
		content = StripFrontMatter(content)

		entryXML := fmt.Sprintf(`  <entry>
    <title>%s</title>
    <link href="%s/blog/%s"/>
    <updated>%s</updated>
    <id>%s/blog/%s</id>
    <content xml:lang="%s" type="html">%s</content>
  </entry>
`, escapeXML(article.Title), meta["url"], article.Slug, article.Date.Format(time.RFC3339), meta["url"], article.Slug, language, escapeXML(content))

		io.WriteString(w, entryXML)
	}

	io.WriteString(w, `</feed>`)

	return nil
}

// escapeXML escapes special XML characters
func escapeXML(s any) string {
	switch v := s.(type) {
	case string:
		return html.EscapeString(v)
	case []byte:
		return html.EscapeString(string(v))
	default:
		return html.EscapeString(fmt.Sprintf("%v", v))
	}
}
