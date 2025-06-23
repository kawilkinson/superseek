package spider

import (
	"reflect"
	"testing"
)

func TestSortReport(t *testing.T) {
	tests := []struct {
		name      string
		inputURLs map[string]int
		expected  []pagesStruct
	}{
		{
			name: "sort by descending count, then URL ascending",
			inputURLs: map[string]int{
				"https://wagslane.dev":            5,
				"https://wagslane.dev/about":      7,
				"https://wagslane.dev/tags":       3,
				"https://wagslane.dev/index.html": 3,
			},
			expected: []pagesStruct{
				{URL: "https://wagslane.dev/about", count: 7},
				{URL: "https://wagslane.dev", count: 5},
				{URL: "https://wagslane.dev/index.html", count: 3},
				{URL: "https://wagslane.dev/tags", count: 3},
			},
		},
		{
			name: "sort with all the same counts",
			inputURLs: map[string]int{
				"https://wagslane.dev":       5,
				"https://wagslane.dev/about": 5,
				"https://wagslane.dev/tags":  5,
			},
			expected: []pagesStruct{
				{URL: "https://wagslane.dev", count: 5},
				{URL: "https://wagslane.dev/about", count: 5},
				{URL: "https://wagslane.dev/tags", count: 5},
			},
		},
		{
			name:      "empty map",
			inputURLs: map[string]int{},
			expected:  []pagesStruct{},
		},
	}

	for i, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actual := sortReport(testCase.inputURLs)

			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Test %v - '%s' FAIL: expected %v, actual %v", i, testCase.name, testCase.expected, actual)
			}
		})
	}
}
