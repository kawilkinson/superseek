package indexerutil

import (
	"regexp"
	"strings"
)

var nameSplitPattern = regexp.MustCompile(`[-_./\s]+`)
var urlSplitPattern = regexp.MustCompile(`[.,\-_\/+:\(\)]+`)
var newlinePattern = regexp.MustCompile(`\n+`)
var whitespacePattern = regexp.MustCompile(`\s+`)
var bracketsPattern = regexp.MustCompile(`\[[^\]]*]`)

func splitName(filename string) []string {
	parts := nameSplitPattern.Split(filename, -1)

	var splitParts []string
	for _, part := range parts {
		if !existsInSet(part, fileTypes) &&
			!existsInSet(strings.ToLower(part), fileTypes) &&
			!strings.Contains(part, "px") {
			splitParts = append(splitParts, strings.ToLower(part))
		}
	}

	return splitParts
}

func splitURL(url string) []string {
	parts := urlSplitPattern.Split(url, -1)

	var splitParts []string
	for _, part := range parts {
		if !existsInSet(part, popularDomains) && !existsInSet(strings.ToLower(part), popularDomains) {
			splitParts = append(splitParts, strings.ToLower(part))
		}
	}

	return splitParts
}

// use html.parse here to grab metadata, paragraphs. Tokenize, and filter it.
func getHTMLData(html string) {

}

// helper function for converting a slice to a set for constant time lookups
func sliceToSet(slice []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, item := range slice {
		set[item] = struct{}{}
	}

	return set
}

// helper function that does a constant time check to see if a word exists in a set
func existsInSet(word string, set map[string]struct{}) bool {
	_, exists := set[word]

	return exists
}
