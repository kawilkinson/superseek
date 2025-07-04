package indexerutil

import (
	"strings"
)

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

// helper function that does a constant time check to see if a word exists in a set
func existsInSet(word string, set map[string]struct{}) bool {
	_, exists := set[word]

	return exists
}
