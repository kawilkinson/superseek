package indexerutil

import (
	"reflect"
	"testing"
)

func TestCountWords(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string]int
	}{
		{
			name:     "basic word count",
			input:    []string{"apple", "banana", "apple", "orange", "banana", "apple"},
			expected: map[string]int{"apple": 3, "banana": 2, "orange": 1},
		},
		{
			name:     "single word repeated",
			input:    []string{"test", "test", "test"},
			expected: map[string]int{"test": 3},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: map[string]int{},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := CountWords(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %d - '%s' FAIL: expected %v, got %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestMostCommonWords(t *testing.T) {
	tests := []struct {
		name     string
		inputMap map[string]int
		n        int
		expected map[string]int
	}{
		{
			name:     "return top 2 most frequent words",
			inputMap: map[string]int{"apple": 3, "banana": 5, "orange": 2},
			n:        2,
			expected: map[string]int{"banana": 5, "apple": 3},
		},
		{
			name:     "handle ties by alphabetical order",
			inputMap: map[string]int{"apple": 2, "banana": 2, "cherry": 2},
			n:        2,
			expected: map[string]int{"apple": 2, "banana": 2},
		},
		{
			name:     "n greater than number of elements",
			inputMap: map[string]int{"apple": 1},
			n:        5,
			expected: map[string]int{"apple": 1},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := MostCommonWords(tc.inputMap, tc.n)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %d - '%s' FAIL: expected %v, got %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
