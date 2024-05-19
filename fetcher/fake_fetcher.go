package fetcher

import "fmt"

type FakeFetcher struct {
	// Maps URLs to the response body bytes.
	content map[string][]byte
}

// NewFakeFetcher creates a new FakeFetcher,
// with the ability to fetch from the map of URLs to response bodies.
// Fetch will error if trying to fetch a URL that's not in this map.
func NewFakeFetcher(content map[string][]byte) *FakeFetcher {
	return &FakeFetcher{content}
}

func (ff *FakeFetcher) Fetch(url string) ([]byte, error) {
	if bytes, ok := ff.content[url]; ok {
		return bytes, nil
	}
	return nil, fmt.Errorf("page does not exist: %s", url)
}
