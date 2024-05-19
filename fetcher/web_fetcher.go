package fetcher

import (
	"fmt"
	"io"
	"net/http"
)

type WebFetcher struct{}

func NewWebFetcher() *WebFetcher {
	return &WebFetcher{}
}

func (wf *WebFetcher) Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch url %s: %w", url, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body of url %s: %w", url, err)
	}

	return body, nil
}
