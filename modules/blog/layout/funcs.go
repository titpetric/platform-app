package layout

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/titpetric/vuego"
)

var Funcs = vuego.FuncMap{
	"postDate": func(val any) any {
		layoutStr := "2006/01/02 15:04"

		switch t := val.(type) {
		case time.Time:
			return t.Format(layoutStr)
		case string:
			// Try to parse as RFC3339 or Unix timestamp
			if parsed, err := time.Parse(time.RFC3339, t); err == nil {
				return parsed.Format(layoutStr)
			}
		}
		return val
	},
	"readingTime": func(content string) string {
		// average reading speed ~200 words/minute
		words := len(strings.Fields(content))
		minutes := words / 200
		if words%200 != 0 {
			minutes++
		}

		if minutes <= 2 {
			return "a few minutes"
		}
		return fmt.Sprintf("%d minutes", minutes)
	},
	"metaTitle": func(title string) string {
		return title
	},
	"metaDescription": func(description string) string {
		return description
	},
	"metaOGImage": func(in string) string {
		return in
	},
	"getCss": func(string) string {
		return ""
	},
	"getJs": func(string) string {
		return ""
	},
	"json": func(in any) (string, error) {
		b, err := json.MarshalIndent(in, "", "  ")
		return string(b), err
	},
}
