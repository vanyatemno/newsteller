package templates

import (
	"html/template"
	"strings"
	"time"
)

type Template interface {
	GeneratePage() (string, error)
}

// Template functions
var funcMap = template.FuncMap{
	"truncateContent": func(content string, maxCharacters int) string {
		if len(content) < maxCharacters {
			return content
		}

		words := strings.Fields(content)
		if len(words) > 6 && len(strings.Join(words[:6], " ")) < maxCharacters {
			return strings.Join(words[:6], " ") + "..."
		}

		return content[:maxCharacters] + "..."
	},
	"formatDate": func(t time.Time) string {
		return t.Format("Jan 02, 2006")
	},
	"formatDateTime": func(t time.Time) string {
		return t.Format("January 2, 2006 at 3:04 PM")
	},
}
