package indexerutil

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/abadojack/whatlanggo"
	"golang.org/x/net/html"
)

var nameSplitPattern = regexp.MustCompile(`[-_./\s]+`)
var urlSplitPattern = regexp.MustCompile(`[.,\-_\/+:\(\)]+`)
var newlinePattern = regexp.MustCompile(`\n+`)
var whitespacePattern = regexp.MustCompile(`\s+`)
var bracketsPattern = regexp.MustCompile(`\[[^\]]*]`)

func GetHTMLData(htmlData string) (map[string]interface{}, error) {
	htmlReader := strings.NewReader(htmlData)
	htmlRootNode, err := html.Parse(htmlReader)
	if err != nil {
		return nil, fmt.Errorf("unable to parse the HTML data given: %v", err)
	}

	metaTags := map[string]string{}
	var paragraphs []string
	var title string

	recurseHTMLTree(htmlRootNode, &metaTags, &paragraphs, &title)

	fullText := strings.Join(paragraphs, " ")
	cleanText := bracketsPattern.ReplaceAllString(fullText, "")
	words := strings.Fields(cleanText)

	var summary string
	if len(words) > 500 {
		summary = strings.Join(words[:500], " ")
	} else {
		summary = strings.Join(words, " ")
	}

	tokens := tokenizeLargeText(summary, 10000)

	var filteredText []string
	for _, word := range tokens {
		lowerWord := strings.ToLower(word)
		if !existsInSet(lowerWord, stopWordsSet) && isAlnum(lowerWord) {
			filteredText = append(filteredText, lowerWord)
		}
	}

	result := map[string]interface{}{
		"title":        fallback(metaTags["og:title"], metaTags["title"], title),
		"description":  fallback(metaTags["og:description"], metaTags["description"]),
		"summary_text": summary,
		"text":         filteredText,
		"language":     detectLanguage(summary, 1000),
	}

	return result, nil
}

func recurseHTMLTree(node *html.Node, metaTags *map[string]string, paragraphs *[]string, title *string) {
	if node.Type == html.ElementNode {
		switch node.Data {
		case "meta":
			var key, content string
			for _, attr := range node.Attr {
				if attr.Key == "property" || attr.Key == "name" {
					key = attr.Val
				}
				if attr.Key == "content" {
					content = attr.Val
				}
			}
			if key != "" && content != "" {
				(*metaTags)[key] = content
			}
		case "title":
			if node.FirstChild != nil && *title == "" {
				*title = strings.TrimSpace(node.FirstChild.Data)
			}
		case "p":
			*paragraphs = append(*paragraphs, extractText(node))
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		recurseHTMLTree(child, metaTags, paragraphs, title)
	}
}

// helper function to extract text out of a html node
func extractText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}

	var strBuilder strings.Builder
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		strBuilder.WriteString(extractText((child)))
	}

	return strBuilder.String()
}

// helper function to have a fallback value for the return in the getHtmlData function
func fallback(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}

// helper function using whatlanggo to grab the language of a sample of text
func detectLanguage(text string, sampleSize int) string {
	var end int
	if len(text) >= sampleSize {
		end = sampleSize
	} else {
		end = len(text)
	}

	sample := text[:end]
	info := whatlanggo.Detect(sample)
	language := info.Lang.String()

	return language
}

// helper function to process large chunks of text
func tokenizeLargeText(text string, chunkSize int) []string {
	var tokens []string

	for i := 0; i < len(text); i += chunkSize {
		end := i + chunkSize
		if end > len(text) {
			end = len(text)
		}

		chunk := text[i:end]
		tokens = append(tokens, chunk)
	}

	return tokens
}

// helper function to make sure a string is all alpha numeric
func isAlnum(str string) bool {
	for _, r := range str {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return false
		}
	}

	return true
}
