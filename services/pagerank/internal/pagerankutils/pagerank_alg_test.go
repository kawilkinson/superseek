package pagerankutils

import "testing"

func TestPageRankSort(t *testing.T) {
	tests := []struct {
		name          string
		pageRank      map[string]float64
		backlinks     map[string][]string
		outlinksCount map[string]int
		count         int64
		expectedTop   string
	}{
		{
			name: "hub and spokes",
			pageRank: map[string]float64{
				"hub":   1.0,
				"page1": 1.0,
				"page2": 1.0,
				"page3": 1.0,
			},
			backlinks: map[string][]string{
				"page1": {"hub"},
				"page2": {"hub", "page1"}, // have page2 backlink to an extra page to make it stand out in rank
				"page3": {"hub"},
			},
			outlinksCount: map[string]int{
				"hub":   3,
				"page1": 1,
			},
			count:       4,
			expectedTop: "page2",
		},
		{
			name: "page with no backlinks",
			pageRank: map[string]float64{
				"orphan": 1.0,
				"hub":    1.0,
			},
			backlinks: map[string][]string{
				"hub": {"orphan"},
			},
			outlinksCount: map[string]int{
				"orphan": 1,
			},
			count:       2,
			expectedTop: "hub",
		},
		{
			name: "three nodes with one link count being higher",
			pageRank: map[string]float64{
				"a": 1.0,
				"b": 1.0,
				"c": 1.0,
			},
			backlinks: map[string][]string{
				"a": {"b"},
				"b": {"a", "c"},
				"c": {"b"},
			},
			outlinksCount: map[string]int{
				"a": 1,
				"b": 2,
				"c": 1,
			},
			count:       3,
			expectedTop: "b",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sortedPageRanks := PageRankSort(tc.pageRank, tc.backlinks, tc.outlinksCount, tc.count)

			if len(sortedPageRanks) == 0 {
				t.Errorf("Test %d - '%s' FAIL: got empty result", i, tc.name)
			}

			topURL := sortedPageRanks[0].URL
			if topURL != tc.expectedTop {
				t.Errorf("Test %d - '%s' FAIL: expected top URL to be '%s', got '%s'", i, tc.name, tc.expectedTop, topURL)
			}
		})
	}
}
