package parser

import (
	"bytes"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/frezzle/crawler/utils"
	"golang.org/x/net/html"
)

// WebParser can parse a HTML document to find more URLs that it links to.
// Single-page apps (SPAs) are not supported... would be cool to do later!
type WebParser struct {
}

func NewWebParser() *WebParser {
	return &WebParser{}
}

// Parse retrieves all unique anchor links in the web page specified by the URL.
// Links to the same URL specified (i.e. links to self) are ignored, though the check is simple so some may slip by.
// All URLs are normalised.
// Returns an error if the page is not HTML content.
// May return error for other reasons, e.g. failing to fetch the page.
func (wc *WebParser) Parse(body []byte, pageUrl string) ([]string, error) {
	pageUrl = utils.NormaliseUrl(pageUrl)

	baseUrl, err := url.Parse(pageUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base url: %s", err)
	}

	uniqueLinks := map[string]struct{}{}
	r := bytes.NewReader(body)
	tokenizer := html.NewTokenizer(r)
Loop:
	for {
		tokenType := tokenizer.Next()
		// log.Println("token type:", tokenType)
		switch tokenType {
		case html.ErrorToken:
			break Loop
		case html.StartTagToken, html.EndTagToken:
			tagName, hasAttrs := tokenizer.TagName()
			if string(tagName) != "a" || !hasAttrs {
				break
			}
			// log.Println("token tag name:", string(tagName), "hasAttrs:", hasAttrs)
			var key, val []byte
			moreAttrs := hasAttrs
			for moreAttrs {
				key, val, moreAttrs = tokenizer.TagAttr()
				if string(key) != "href" {
					continue
				}
				href := strings.TrimSpace(string(val))
				if href == "" || href == "/" || href[0] == '#' {
					continue // skip relative links to self
				}
				// log.Println("tag href:     ", string(val))

				hrefAbsolute, err := baseUrl.Parse(href)
				if err != nil {
					return nil, fmt.Errorf("failed to parse absolute url: %w", err)
				}

				url := utils.NormaliseUrl(hrefAbsolute.String())

				if url == pageUrl {
					continue // skip absolute link to itself
				}

				uniqueLinks[url] = struct{}{}
				// log.Println("absolute href:", hrefAbsolute)
			}
			// case html.TextToken, html.DoctypeToken, html.CommentToken:
			// 	text := tokenizer.Text()
			// 	log.Println("token text:", string(text))
		}
	}
	// TODO try html.Parse instead^ ?

	links := make([]string, len(uniqueLinks))
	i := 0
	for link := range uniqueLinks {
		links[i] = link
		i++
	}

	// sort it to make results deterministic (assuming web page content stays the same!)
	// TODO: could be an option of the crawler? i.e. move this sorting out of the parser
	slices.Sort(links)

	return links, nil
}
