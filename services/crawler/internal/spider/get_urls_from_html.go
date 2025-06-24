package spider

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Initialize a compiled regex to make sure characters are in ASCII
var nonASCIIRegex = regexp.MustCompile(`[^\x20-\x7E]`)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	htmlReader := strings.NewReader(htmlBody)
	htmlRootNode, err := html.Parse(htmlReader)
	if err != nil {
		return nil, fmt.Errorf("error when parsing HTML given: %v", err)
	}

	parsedURLs := []string{}
	parsedURLs = recurseHTMLTree(htmlRootNode, parsedURLs, rawBaseURL)

	return parsedURLs, nil
}

func recurseHTMLTree(node *html.Node, parsedURLs []string, rawBaseURL string) []string {
	if node == nil {
		return parsedURLs
	}

	if node.Type == html.ElementNode && node.Data == "a" {
		for _, anchor := range node.Attr {
			if anchor.Key == "href" {
				rawHref := anchor.Val

				// Skip any malformed URLs and any non-ASCII URLs
				if strings.ContainsAny(rawHref, " <>\"") || nonASCIIRegex.MatchString(rawHref) {
					continue
				}

				hrefURL, err := url.Parse(rawHref)
				if err != nil {
					break
				}
				baseURL, err := url.Parse(rawBaseURL)
				if err != nil {
					break
				}
				
				resolvedURL := baseURL.ResolveReference(hrefURL)
				parsedURLs = append(parsedURLs, resolvedURL.String())
				break
			}
		}
	}

	parsedURLs = recurseHTMLTree(node.FirstChild, parsedURLs, rawBaseURL)
	parsedURLs = recurseHTMLTree(node.NextSibling, parsedURLs, rawBaseURL)

	return parsedURLs
}
