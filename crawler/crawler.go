package crawler

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/frezzle/crawler/fetcher"
	"github.com/frezzle/crawler/parser"
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

	if len(seedLocations) > crawlLimit {
		seedLocations = seedLocations[:crawlLimit]
	}

	queued := map[string]bool{}
	queue := make(chan string, len(seedLocations))
	for _, loc := range seedLocations {
		queued[loc] = true
		queue <- loc
	}

	// start worker pool, where each worker fetches and parses a location
	workerCount := 10
	crawlResults := make(chan crawlResult, workerCount)
	var crawlers sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		crawlers.Add(1)
		go func(workerId int) {
			log.Println("worker", i, "started")
			defer func() {
				log.Println("worker", i, "ending")
				crawlers.Done()
			}()
			for loc := range queue {
				log.Println("worker", i, "crawling", loc)

				// fetch it
				content, err := wc.fetcher.Fetch(loc)
				if err != nil {
					log.Println(fmt.Errorf("worker %d: failed to fetch content from location %s: %w", workerId, loc, err))
					continue
				}

				// parse it
				locs, err := wc.parser.Parse(content, loc)
				if err != nil {
					log.Println(fmt.Errorf("worker %d: failed to parse content from location %s: %w", workerId, loc, err))
					continue
				}

				log.Printf("worker %d: Crawled %s and found %d other unique links.\n", workerId, loc, len(locs))
				crawlResults <- crawlResult{SourceLocation: loc, LinkedLocations: locs}
			}
		}(i)
	}

	// Separate routine that will wait for all workers to be done before finally closing the results channel,
	// which notifies the code below that we have all results.
	go func() {
		crawlers.Wait()
		log.Println("closing crawlResults channel")
		close(crawlResults)
	}()

	// gather crawl results,
	// keep queueing more crawls until we've reached crawl limit.
	// TODO: what if we run out of locations to crawl before reaching crawl limit? detectable with channels?
	locationLinks := make([][2]string, 0, 1000)
	queueClosed := false
outer:
	for {
		select {
		case result, ok := <-crawlResults:
			if !ok {
				break outer
			}
			log.Println("top")
			log.Println("received crawl results", result.SourceLocation, "-->", result.LinkedLocations)
			for _, loc := range result.LinkedLocations {
				// always record links
				locationLinks = append(locationLinks, [2]string{result.SourceLocation, loc})

				// possibly don't crawl this location
				if queued[loc] {
					continue // skip location that's already been crawled or will be crawled
				} else if !wc.canCrawlLocation(loc) {
					continue
				} else if len(queued) < crawlLimit {
					// schedule to crawl this location
					log.Println("queuing location", loc)
					queued[loc] = true
					queue <- loc
					log.Println("queued length is", len(queued))
				}
			}

			if len(queued) == crawlLimit && !queueClosed {
				log.Println("closing queue channel")
				close(queue)
				queueClosed = true
			}
		case <-time.After(time.Duration(3) * time.Second): // TODO better way to detect that queue is empty AND no more results to process
			log.Println("Not received results in a while. Assuming done, exiting...")
			break outer
		}
	}

	return locationLinks, nil
}

type crawlResult struct {
	SourceLocation  string
	LinkedLocations []string
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
