package spider

import (
	"fmt"
	"sort"
)

type pagesStruct struct {
	URL   string
	count int
}

func PrintReport(pages map[string]int, baseURL string) {
	fmt.Println("")
	fmt.Println("=============================")
	fmt.Printf("REPORT for %s\n", baseURL)
	fmt.Println("=============================")
	fmt.Println("")

	sortedPages := sortReport(pages)

	for _, page := range sortedPages {
		fmt.Printf("Found %d internal links to %s\n", page.count, page.URL)
	}
}

func sortReport(pages map[string]int) []pagesStruct {
	sortedPages := []pagesStruct{}
	for URL, count := range pages {
		sortedPages = append(sortedPages, pagesStruct{URL, count})
	}

	sort.Slice(sortedPages, func(i, j int) bool {
		if sortedPages[i].count != sortedPages[j].count {
			return sortedPages[i].count > sortedPages[j].count
		}
		return sortedPages[i].URL < sortedPages[j].URL
	})

	return sortedPages
}
