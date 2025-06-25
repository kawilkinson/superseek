package spider

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// double check and make sure no internal IPs or schemes are passed through here
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("parsed URL is invalid, no https or http detected: %v", parsedURL)
	}

	ip := net.ParseIP(parsedURL.Hostname())
	if ip != nil && ip.IsPrivate() {
		return "", fmt.Errorf("refusing to crawl private IP address: %s", ip)
	}

	// #nosec G107
	response, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("error making GET request to provided URL: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return "", fmt.Errorf("bad status code returned from provided URL: %d\n%s", response.StatusCode, http.StatusText(response.StatusCode))
	}
	if !strings.Contains(response.Header.Get("Content-Type"), "text/html") {
		return "", fmt.Errorf("response is not in text/html content-type: %v", response.Header)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the response body: %v", err)
	}

	return string(body), nil
}
