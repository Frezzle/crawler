package fetcher

type Fetcher interface {
	// Fetch GETs the page specified by the provided HTTP URL,
	// returning the bytes of the response body, if there are any.
	Fetch(url string) ([]byte, error)
}
