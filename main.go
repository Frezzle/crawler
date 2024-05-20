package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/frezzle/web-crawler/fetcher"
	"github.com/frezzle/web-crawler/parser"
)

// possible TODOs:
// - include seed URLs in diagram even if they don't link to anywhere else
// - optimise mermaid file output by not repeating the name of same web pages; if nodes can be defined separately from connections then we get the above^ TODO for free!
// - concurrent crawling for better performance
// - interactive crawling e.g. interactive mermaid diagram
// - page title in diagram nodes, alongside
// - check if a domain has a sitemap defined for new URLs to crawl
// - check for robots.txt file to adhere to (to be a polite crawler = less likely to be blocked)
// - make this crawler into a library so anyone can use with their own seed urls, config, interface implementations, etc.

func hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

func truncateUrl(url string, n int) string {
	if len(url) > n {
		if n > 3 {
			return url[:n-3] + "..."
		}
		return url[:n]
	}
	return url
}

// Most of these are because they spoil the mermaid diagram,
// but RL, for example, is to not mess up their visitor metrics.
var domainsToNotCrawl = []string{
	"aseprite.org",
	"github.com",
	"hollowknight.com",
	"linkedin.com",
	"mercadolibre.com.ar",
	"uptimerobot.com",
	"wikipedia.org",
	"youtube.com",
	// don't mess up metrics for these peeps:
	"https://www.arkamitra.com",
	"https://dimitris.dev",
	"riskledger.com",
}

// Returns true if we are allowed to crawl domain, otherwise false.
// It's overly aggressive right now e.g. if we block bla.com it'll also block blabla.com.
func canCrawlUrl(url string) bool {
	for _, domain := range domainsToNotCrawl {
		if strings.Contains(url, domain) {
			return false
		}
	}
	return true
}

func main() {
	crawlLimit := 20
	urlsToCrawl := []string{
		"https://banch.io",
		"https://brunocalogero.dev",
		"https://frederico.dev",
	}
	crawledUrls := make(map[string]bool)
	f := fetcher.NewWebFetcher()
	p := parser.NewWebParser(f)
	urlConnections := make([][2]string, 0, 1000)

	for len(urlsToCrawl) > 0 && len(crawledUrls) < crawlLimit {
		// dequeue URL
		url := urlsToCrawl[0]
		urlsToCrawl = urlsToCrawl[1:]
		// skip if already crawled
		if crawledUrls[url] {
			log.Printf("Skipping already-crawled URL %s\n", url)
			continue
		}
		// skip if i don't want to crawl it for some reason
		if !canCrawlUrl(url) {
			log.Printf("Not crawling URL %s\n", url)
			continue
		}
		// record it as crawled
		crawledUrls[url] = true
		// crawl it
		urls, err := p.Parse(url)
		if err != nil {
			log.Printf("failed to crawl %s: %s\n", url, fmt.Errorf("%w", err))
		}
		log.Printf("Crawled %s and found %d unique links.\n", url, len(urls))

		// queue URLs to crawl,
		// and record the source->destination links for mermaid diagram
		// newUrls := 0
		for _, u := range urls {
			urlConnections = append(urlConnections, [2]string{url, u})

			// queue any new url that needs crawling
			// if crawledUrls[u] {
			// 	continue
			// }
			urlsToCrawl = append(urlsToCrawl, u)
			// newUrls++
			// if newUrls == 3 {
			// 	break // limit number of links from each page, to try get a nicer mermaid graph
			// }

		}
	}

	err := saveMermaidFlowchart(urlConnections, "flowchart.mermaid")
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to save mermaid graph: %w", err))
	}
}

// Outputs many URL->URL "connections" to a mermaid flowchart file,
// which can be rendered to an image e.g. using mermaid.js, mermaid CLI, mermaid.live.
func saveMermaidFlowchart(connections [][2]string, filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("unable to open or create file for mermaid graph: %w", err)
	}
	_, err = file.WriteString("flowchart LR\n")
	if err != nil {
		return fmt.Errorf("unable to write header to mermaid file: %w", err)
	}
	for _, conn := range connections {
		_, err = file.WriteString(fmt.Sprintf(
			"%s[%s]-->%s[%s]\n",
			truncate(hash(conn[0]), 10), // truncate to stay below char limit on mermaid.live
			truncateUrl(conn[0], 150),
			truncate(hash(conn[1]), 10), // truncate to stay below char limit on mermaid.live
			truncateUrl(conn[1], 150),
			// TODO: don't truncate now that i know how to increase char limit on mermaid.live / am using mermaid cli ?
		))
		if err != nil {
			return fmt.Errorf("unable to write url connection to mermaid file: %w", err)
		}
	}

	return nil
}
