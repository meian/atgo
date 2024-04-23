package csrf

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	csrfPrefix = `var csrfToken = "`
)

func FromDocument(doc *goquery.Document) string {
	script := doc.Find("head script").FilterFunction(func(i int, s *goquery.Selection) bool {
		return strings.Contains(s.Text(), csrfPrefix)
	}).First().Text()
	if len(script) == 0 {
		return ""
	}
	script = script[strings.Index(script, csrfPrefix)+len(csrfPrefix):]
	if idx := strings.Index(script, `"`); idx != -1 {
		return script[:idx]
	}
	return ""
}
