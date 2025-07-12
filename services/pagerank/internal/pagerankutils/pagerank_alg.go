package pagerankutils

import (
	"log"
	"sort"
)

type SortedPageRanks struct {
	URL  string
	Rank float64
}

func PageRankSort(pageRank map[string]float64, backlinks map[string][]string, outlinksCount map[string]int, count int64) []SortedPageRanks {
	iterations := 10
	damping := 0.85
	for i := 0; i < iterations; i++ {
		newPageRank := make(map[string]float64)

		for url := range pageRank {
			var newCumulativeRank float64

			if backlinksForURL, exists := backlinks[url]; exists {
				for _, backlink := range backlinksForURL {
					outlinkCount, outlinkExists := outlinksCount[backlink]
					backlinkRank, backlinkExists := pageRank[backlink]

					if backlinkExists && outlinkExists {
						newCumulativeRank += (backlinkRank / float64(outlinkCount))
					}
				}
			}
			newPageRank[url] = (1-damping)/float64(count) + (damping * newCumulativeRank)
		}
		pageRank = newPageRank
	}

	sortedPageRanks := []SortedPageRanks{}
	for url, rank := range pageRank {
		sortedPageRanks = append(sortedPageRanks, SortedPageRanks{
			URL:  url,
			Rank: rank,
		})
	}

	sort.Slice(sortedPageRanks, func(i, j int) bool {
		return sortedPageRanks[i].Rank > sortedPageRanks[j].Rank
	})

	log.Println("sorted page rank values:")

	for _, pageRank := range sortedPageRanks {
		log.Printf("page URL: %s, page rank: %f\n", pageRank.URL, pageRank.Rank)
	}

	return sortedPageRanks
}
