package crawler

import (
	"fmt"
	"log"
	"strings"

	"github.com/frezzle/web-crawler/fetcher"
	"github.com/frezzle/web-crawler/parser"
)

type Crawler struct {
	fetcher             fetcher.Fetcher
	parser              parser.Parser
	locationsToNotCrawl []string
}

type Option func(wc *Crawler)

func NewCrawler(options ...Option) *Crawler {
	wc := &Crawler{}
	for _, opt := range options {
		opt(wc)
	}
	return wc
}

func WithFetcher(f fetcher.Fetcher) func(*Crawler) {
	return func(wc *Crawler) {
		wc.fetcher = f
	}
}

func WithParser(p parser.Parser) func(*Crawler) {
	return func(wc *Crawler) {
		wc.parser = p
	}
}

func WithLocationsToNotCrawl(locations []string) func(*Crawler) {
	return func(wc *Crawler) {
		wc.locationsToNotCrawl = locations
	}
}

// Pages are crawled on a best-effort basis.
// Returns location pairs, where each represents a source location that links to another location.
// TODO: decouple logger/logging from here? errors channel? logs channel? errors array returned?
func (wc *Crawler) Crawl(seedLocations []string, crawlLimit int) ([][2]string, error) {
	if len(seedLocations) == 0 {
		return nil, fmt.Errorf("must specify seed location(s)")
	}
	if crawlLimit < 1 {
		return nil, fmt.Errorf("must crawl at least 1 page")
	}

	locationsToCrawl := make([]string, len(seedLocations))
	copy(locationsToCrawl, seedLocations)
	crawledLocations := make(map[string]bool)
	f := fetcher.NewWebFetcher()
	p := parser.NewWebParser()
	locationLinks := make([][2]string, 0, 1000)

	for len(locationsToCrawl) > 0 && len(crawledLocations) < crawlLimit {
		// dequeue
		loc := locationsToCrawl[0]
		locationsToCrawl = locationsToCrawl[1:]
		// skip if already crawled
		if crawledLocations[loc] {
			log.Printf("Skipping already-crawled location %s\n", loc)
			continue
		}
		// skip if i don't want to crawl it for some reason
		if !wc.canCrawlLocation(loc) {
			log.Printf("Not crawling location %s\n", loc)
			continue
		}

		// record it as crawled
		crawledLocations[loc] = true

		// fetch it
		content, err := f.Fetch(loc)
		if err != nil {
			log.Println(fmt.Errorf("failed to fetch content from location %s: %w", loc, err))
			continue
		}

		// parse it
		locs, err := p.Parse(content, loc)
		if err != nil {
			log.Println(fmt.Errorf("failed to parse content from location %s: %w", loc, err))
			continue
		}

		log.Printf("Crawled %s and found %d other unique links.\n", loc, len(locs))

		// queue locations to crawl and record links
		for _, l := range locs {
			locationLinks = append(locationLinks, [2]string{loc, l})
			locationsToCrawl = append(locationsToCrawl, l)
		}
	}

	return locationLinks, nil
}

// Returns true if we are allowed to crawl the location.
// TODO: not customisable e.g. it specifically does a strings.Contains for the use-case
// of not crawling domains that appear in a link :/ there's room for over-engineering this! xP
func (wc *Crawler) canCrawlLocation(loc string) bool {
	for _, not := range wc.locationsToNotCrawl {
		if strings.Contains(loc, not) {
			return false
		}
	}
	return true
}
