package indexerutil

import (
	"regexp"
	"strings"
)

var nameSplitPattern = regexp.MustCompile(`[-_./\s]+`)

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

// helper function that does a constant time check to see if a word exists in a set
func existsInSet(word string, set map[string]struct{}) bool {
	_, exists := set[word]

	return exists
}
