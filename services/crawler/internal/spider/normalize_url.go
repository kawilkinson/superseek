package spider

import (
	"fmt"
	"net/url"
	"strings"
)

func NormalizeURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("error when parsing URL given: %v", err)
	}

	if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
		return "", fmt.Errorf("parsed URL is invalid, no https or http detected: %v", parsedURL)
	}

	if parsedURL.Host == "" {
		return "", fmt.Errorf("URL has no filled in 'Host' field")
	}

	parsedURL.Host = strings.TrimPrefix(parsedURL.Host, "www.")
	cleanedPath := cleanPath(parsedURL.Path)
	normalizedURL := strings.ToLower(parsedURL.Host) + strings.ToLower(cleanedPath)

	return normalizedURL, nil
}

func cleanPath(pathURL string) string {
	return strings.TrimSuffix(pathURL, "/")
}
