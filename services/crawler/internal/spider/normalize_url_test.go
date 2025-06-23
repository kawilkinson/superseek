package spider

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove end forward slash",
			inputURL: "blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove both scheme and end forward slash",
			inputURL: "http://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove uppercase characters",
			inputURL: "Blog.BOOT.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove queries",
			inputURL: "blog.boot.dev/path?param1=value1&param2=value2",
			expected: "blog.boot.dev/path",
		},
	}

	for i, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := normalizeURL(testCase.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, testCase.name, err)
				return
			}
			if actual != testCase.expected {
				t.Errorf("Test %v - '%s' FAIL: expected URL: %v, actual: %v", i, testCase.name, testCase.expected, actual)
			}
		})
	}
}
