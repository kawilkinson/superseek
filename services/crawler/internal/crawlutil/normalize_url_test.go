package crawlutil

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
		err      bool
	}{
		{
			name:     "remove scheme",
			inputURL: "https://test.domain.com/path",
			expected: "test.domain.com/path",
			err:      false,
		},
		{
			name:     "remove trailing forward slash",
			inputURL: "https://test.domain.com/path/",
			expected: "test.domain.com/path",
			err:      false,
		},
		{
			name:     "remove both scheme and end forward slash",
			inputURL: "https://test.domain.com/path/",
			expected: "test.domain.com/path",
			err:      false,
		},
		{
			name:     "remove uppercase characters",
			inputURL: "https://Test.DOMAIN.com/path",
			expected: "test.domain.com/path",
			err:      false,
		},
		{
			name:     "remove queries",
			inputURL: "https://test.domain.com/path?param1=value1&param2=value2",
			expected: "test.domain.com/path",
			err:      false,
		},
		{
			name:     "remove www.",
			inputURL: "https://www.test.domain.com/path",
			expected: "test.domain.com/path",
			err:      false,
		},
		{
			name:     "invalid scheme",
			inputURL: "htp://test.domain.com/path",
			expected: "",
			err:      true,
		},
	}

	for i, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := NormalizeURL(testCase.inputURL)
			if err != nil && !testCase.err {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, testCase.name, err)
				return
			}
			if actual != testCase.expected {
				t.Errorf("Test %v - '%s' FAIL: expected URL: %v, actual: %v", i, testCase.name, testCase.expected, actual)
			}
		})
	}
}
