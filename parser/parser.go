package parser

type Parser interface {
	// Parse content and return other unique locations that it links to, if any.
	// Relative location link are made absolute, determined by the given location.
	Parse(content []byte, location string) ([]string, error)
}
