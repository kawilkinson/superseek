package spider

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func getHTML(rawURL string) (string, int, string, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", 0, "", fmt.Errorf("invalid URL: %w", err)
	}

	// double check and make sure no internal IPs or schemes are passed through here
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", 0, "", fmt.Errorf("parsed URL is invalid, no https or http detected: %v", parsedURL)
	}

	ip := net.ParseIP(parsedURL.Hostname())
	if ip != nil && ip.IsPrivate() {
		return "", 0, "", fmt.Errorf("refusing to crawl private IP address: %s", ip)
	}

	// #nosec G107
	response, err := http.Get(rawURL)
	if err != nil {
		return "", 0, "", fmt.Errorf("error making GET request to provided URL: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return "", response.StatusCode, "", fmt.Errorf("bad status code returned from provided URL: %d\n%s", response.StatusCode, http.StatusText(response.StatusCode))
	}

	contentType := response.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", 0, contentType, fmt.Errorf("response is not in text/html content-type: %v", response.Header)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", response.StatusCode, contentType, fmt.Errorf("error reading the response body: %v", err)
	}

	return string(body), response.StatusCode, contentType, nil
}
