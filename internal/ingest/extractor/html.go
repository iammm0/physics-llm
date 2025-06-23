package extractor

import (
	"github.com/PuerkitoBio/goquery"
	"os"
	"strings"
)

type htmlExt struct{}

func (htmlExt) Extract(p string) (string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", err
	}
	doc, _ := goquery.NewDocumentFromReader(f)
	text := strings.TrimSpace(doc.Text())
	return text, nil
}

func init() { Register(".html", htmlExt{}); Register(".htm", htmlExt{}) }
