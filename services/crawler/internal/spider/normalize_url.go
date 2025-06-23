package spider

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("error when parsing URL given: %v", err)
	}

	cleanedPath := cleanPath(parsedURL.Path)
	normalizedURL := strings.ToLower(parsedURL.Host) + strings.ToLower(cleanedPath)

	return normalizedURL, nil
}

func cleanPath(pathURL string) string {
	return strings.TrimRight(pathURL, "/")
}
