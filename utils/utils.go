package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func Truncate(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

func TruncateUrl(url string, n int) string {
	if len(url) > n {
		if n > 3 {
			return url[:n-3] + "..."
		}
		return url[:n]
	}
	return url
}

// Very basic normalisation of URL so we don't have duplicates of same page.
// For now it simply removes a trailing slash if appears at the end,
// however this won't work in many cases...
//
// - bla.com/?q=hello will not become bla.com?q=hello
//
// - bla.com/#title will not become bla.com#title
//
// - bla.com// will not become bla.com
//
// ...it's fine for now, hopefully.
func NormaliseUrl(url string) string {
	s, _ := strings.CutSuffix(url, "/")
	return s
}
