package parser

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GodotSearchOptions struct {
	Version  string
	Slug     string
	Platform string
	Mono     bool
}

const (
	GodotArchivePageUrl = "https://godotengine.org/download/archive/"
)

func GetDocumentFromURL(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Failed to close response body: %v", closeErr)
		}
	}()

	return goquery.NewDocumentFromReader(resp.Body)
}

func GetLatestVersion() string {
	doc, err := GetDocumentFromURL(GodotArchivePageUrl)
	if err != nil {
		return ""
	}

	return doc.Find(".archive-version:first-of-type>h4").First().Text()
}

func GetLatestExperimentalVersion() string {
	doc, err := GetDocumentFromURL(GodotArchivePageUrl)
	if err != nil {
		return ""
	}

	return doc.Find(".archive-version>h4").First().Text()
}

func GetGodotDownloadURL(options GodotSearchOptions) string {
	doc, err := GetDocumentFromURL(GodotArchivePageUrl + options.Version + "-" + options.Slug)
	if err != nil {
		return ""
	}

	return doc.Find(".btn-download").FilterFunction(func(index int, s *goquery.Selection) bool {
		href, exists := s.Attr("href")
		if !exists {
			return false
		}

		if !strings.Contains(href, options.Version) {
			return false
		}

		if !strings.Contains(href, options.Slug) {
			return false
		}

		if !strings.Contains(href, options.Platform) {
			return false
		}

		if options.Mono != strings.Contains(href, "mono") {
			return false
		}

		return true
	}).First().AttrOr("href", "")
}
