package indexerutil

import (
	"reflect"
	"testing"
)

// COMMENTED OUT FOR NOW, NOT USED
// func TestSplitName(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		input    string
// 		expected []string
// 	}{
// 		{
// 			name:     "simple image file",
// 			input:    "cat-cute-photo.jpg",
// 			expected: []string{"cat", "cute", "photo"},
// 		},
// 		{
// 			name:     "filename with px and extension",
// 			input:    "icon-128px.svg",
// 			expected: []string{"icon"},
// 		},
// 		{
// 			name:     "complex filename with underscores and dots",
// 			input:    "holiday_pics.2024.final_version.PNG",
// 			expected: []string{"holiday", "pics", "2024", "final", "version"},
// 		},
// 	}

// 	for i, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			actual := SplitName(tc.input)
// 			if !reflect.DeepEqual(actual, tc.expected) {
// 				t.Errorf("Test %d - '%s' FAIL: expected %v, got %v", i, tc.name, tc.expected, actual)
// 			}
// 		})
// 	}
// }

func TestSplitURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "basic image URL",
			input:    "https://images.unsplash.com/photo-1627846798698-1f4e59c9d9a4",
			expected: []string{"https", "images", "unsplash", "photo", "1627846798698", "1f4e59c9d9a4"},
		},
		{
			name:     "cdn URL with dimensions",
			input:    "https://cdn.example.com/assets/icons-256px-v1.png",
			expected: []string{"https", "cdn", "example", "assets", "icons", "256px", "v1", "png"},
		},
		{
			name:     "URL with brand and TLD filtered",
			input:    "https://www.google.com/search/images/cute-dog.jpg",
			expected: []string{"https", "www", "search", "images", "cute", "dog", "jpg"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := SplitURL(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %d - '%s' FAIL: expected %v, got %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}