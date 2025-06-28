package crawlutil

import (
	"testing"
)

func TestStripURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
		err      bool
	}{
		{
			name:     "remove trailing forward slash",
			inputURL: "https://test.domain.com/path/",
			expected: "https://test.domain.com/path",
			err:      false,
		},
		{
			name:     "keep www.",
			inputURL: "https://www.test.domain.com/path",
			expected: "https://www.test.domain.com/path",
			err:      false,
		},
		{
			name:     "remove fragments",
			inputURL: "https://test.domain.com/path#Test",
			expected: "https://test.domain.com/path",
			err:      false,
		},
		{
			name:     "remove query parameters",
			inputURL: "https://test.domain.com/path?version=1.5&test=tester",
			expected: "https://test.domain.com/path",
			err:      false,
		},
	}

	for i, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := StripURL(testCase.inputURL)
			if err != nil && !testCase.err {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, testCase.name, err)
			}

			if actual != testCase.expected {
				t.Errorf("Test %v - '%s' FAIL: expected URL: %v, actual: %v", i, testCase.name, testCase.expected, actual)
			}
		})
	}
}
