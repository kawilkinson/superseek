package crawler_utilities

import (
	"fmt"
	"net/url"
	"strings"
)

func StripURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("could not parse raw url: %v", err)
	}

	if parsedURL.Scheme == "" {
		return "", fmt.Errorf("URL has no filled in 'Scheme' field")
	}

	if parsedURL.Host == "" {
		return "", fmt.Errorf("URL has no filled in 'Host' field")
	}

	strippedURL := parsedURL.Scheme + "://" + parsedURL.Host

	if parsedURL.Path != "" {
		trimmedPath := strings.TrimSuffix(parsedURL.Path, "/")
		strippedURL += trimmedPath
	}

	return strippedURL, nil
}
