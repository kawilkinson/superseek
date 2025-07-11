package indexerutil

import (
	"regexp"
	"strings"
)

var nameSplitPattern = regexp.MustCompile(`[-_./\s]+`)
var urlSplitPattern = regexp.MustCompile(`[.,\-_\/+:\(\)]+`)

func SplitFilename(filename string) []string {
	parts := nameSplitPattern.Split(filename, -1)

	var splitParts []string
	for _, part := range parts {
		part = strings.ToLower(part)
		if !existsInSet(part, fileTypes) && !strings.Contains(part, "px") {
			splitParts = append(splitParts, part)
		}
	}

	return splitParts
}

func SplitURL(url string) []string {
	parts := urlSplitPattern.Split(url, -1)

	var splitParts []string
	for _, part := range parts {
		if !existsInSet(part, popularDomains) && !existsInSet(strings.ToLower(part), popularDomains) {
			splitParts = append(splitParts, strings.ToLower(part))
		}
	}

	return splitParts
}

// helper function that does a constant time check to see if a word exists in a set
func existsInSet(word string, set map[string]struct{}) bool {
	_, exists := set[word]

	return exists
}
