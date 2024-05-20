package crawler

import (
	"github.com/frezzle/web-crawler/fetcher"
	"github.com/frezzle/web-crawler/parser"
)

// NewWebCrawler creates a crawler that fetches web pages and parses their HTML content.
// This is simply a convenience function for the most popular type of crawler.
func NewWebCrawler(options ...Option) *Crawler {
	wc := NewCrawler(
		WithFetcher(fetcher.NewWebFetcher()),
		WithParser(parser.NewWebParser()),
	)
	for _, opt := range options {
		opt(wc)
	}
	return wc
}
