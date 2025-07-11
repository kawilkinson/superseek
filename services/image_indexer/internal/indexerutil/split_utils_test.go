package indexerutil

import (
	"reflect"
	"testing"
)

func TestSplitFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple image file",
			input:    "cat-cute-photo.jpg",
			expected: []string{"cat", "cute", "photo"},
		},
		{
			name:     "filename with px and extension",
			input:    "icon-128px.svg",
			expected: []string{"icon"},
		},
		{
			name:     "complex filename with underscores and dots",
			input:    "holiday_pics.2024.final_version.PNG",
			expected: []string{"holiday", "pics", "2024", "final", "version"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := SplitFilename(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %d - '%s' FAIL: expected %v, got %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
