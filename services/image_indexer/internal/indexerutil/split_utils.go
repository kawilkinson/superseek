package indexerutil

import (
	"image"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var nameSplitPattern = regexp.MustCompile(`[-_./\s]+`)
var urlSplitPattern = regexp.MustCompile(`[.,\-_\/+:\(\)]+`)
var newlinePattern = regexp.MustCompile(`\n+`)
var whitespacePattern = regexp.MustCompile(`\s+`)
var bracketsPattern = regexp.MustCompile(`\[[^\]]*]`)

func SplitName(filename string) []string {
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

func IsValidImage(url string, minWidth, minHeight int) bool {
	var absoluteURL string
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		absoluteURL = "https://" + url
	}

	client := http.Client{
		Timeout: Timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("unable to get image %s: %v\n", absoluteURL, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("non 200 response for image %s: %d\n", absoluteURL, resp.StatusCode)
		return false
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Printf("unable to decode image %s: %v\n", absoluteURL, err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	return width >= minWidth && height >= minHeight
}

// helper function that does a constant time check to see if a word exists in a set
func existsInSet(word string, set map[string]struct{}) bool {
	_, exists := set[word]

	return exists
}
