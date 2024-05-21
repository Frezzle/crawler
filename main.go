package main

import (
	"fmt"
	"log"

	"github.com/frezzle/web-crawler/crawler"
	"github.com/frezzle/web-crawler/graph"
)

func main() {
	wc := crawler.NewWebCrawler(crawler.WithLocationsToNotCrawl([]string{
		// These ones cause too many uninteresting links, and make the graph uglier:
		"aseprite.org",
		"github.com",
		"hollowknight.com",
		"linkedin.com",
		"mercadolibre.com.ar",
		"twitter.com",
		"uptimerobot.com",
		"wikipedia.org",
		"youtube.com",

		// temporary, for making nice example image for readme
		"https://banch.io/astrocappy",
		"https://banch.io/blog",

		// Don't mess up metrics for these peeps:
		"arkamitra.com",
		"dimitris.dev",
		"riskledger.com",
	}))

	urlConnections, err := wc.Crawl([]string{
		// These have blessed me with permission to crawl:
		"https://banch.io",
		"https://brunocalogero.dev",
		"https://frederico.dev",
	}, 20)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed crawling: %s", err))
	}

	err = graph.SaveMermaidFlowchart(urlConnections, "flowchart.mermaid")
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to save mermaid graph: %w", err))
	}
}
