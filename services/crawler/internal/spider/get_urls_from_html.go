package spider

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, map[string]map[string]string, error) {
	htmlReader := strings.NewReader(htmlBody)
	htmlRootNode, err := html.Parse(htmlReader)
	if err != nil {
		return nil, nil, fmt.Errorf("error when parsing HTML given: %v", err)
	}

	URLsSet := make(map[string]struct{})
	images := make(map[string]map[string]string)

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error when parsing base URL: %v", err)
	}

	recurseHTMLTree(htmlRootNode, URLsSet, images, baseURL)

	parsedURLs := setToSlice(URLsSet)

	return parsedURLs, images, nil
}

func recurseHTMLTree(node *html.Node, URLSet map[string]struct{}, images map[string]map[string]string, baseURL *url.URL) {
	if node == nil {
		return
	}

	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				rawHref := attr.Val

				// Skip any malformed URLs and any non-ASCII URLs
				if strings.ContainsAny(rawHref, " <>\"") || !isASCII(rawHref) {
					break
				}

				hrefURL, err := url.Parse(rawHref)
				if err != nil {
					break
				}

				resolvedURL := baseURL.ResolveReference(hrefURL)
				URLSet[resolvedURL.String()] = struct{}{}
				break
			}
		}
	} else if node.Type == html.ElementNode && node.Data == "img" {
		imgInfo := map[string]string{}
		for _, attr := range node.Attr {
			if attr.Key == "src" {
				rawSrc := attr.Val

				// Skip any malformed URLs and any non-ASCII URLs
				if strings.ContainsAny(rawSrc, " <>\"") || !isASCII(rawSrc) {
					break
				}

				srcURL, err := url.Parse(rawSrc)
				if err != nil {
					break
				}

				resolvedURL := baseURL.ResolveReference(srcURL)
				imgInfo["src"] = resolvedURL.String()

			} else if attr.Key == "alt" {
				imgInfo["alt"] = attr.Val
			}
		}

		if src, exists := imgInfo["src"]; exists {
			images[src] = imgInfo
		}
	}

	recurseHTMLTree(node.FirstChild, URLSet, images, baseURL)
	recurseHTMLTree(node.NextSibling, URLSet, images, baseURL)
}

// Helper function for creating a slice
func setToSlice(set map[string]struct{}) []string {
	slice := make([]string, 0, len(set))
	for item := range set {
		slice = append(slice, item)
	}

	return slice
}

// Helper function for ensuring a URL is all in ASCII
func isASCII(str string) bool {
	for _, char := range str {
		if char < 0x20 || char > 0x7E {
			return false
		}
	}

	return true
}
