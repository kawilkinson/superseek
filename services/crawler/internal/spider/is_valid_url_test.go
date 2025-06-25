package spider

import "testing"

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected bool
	}{
		{
			name:     "valid url",
			inputURL: "https://test.domain.com/test/tester",
			expected: true,
		},
		{
			name:     "valid normalized URL",
			inputURL: "test.domain.com/test/tester",
			expected: true,
		},
		{
			name:     "invalid URL",
			inputURL: "https://test.domain.com/test/このリンクが駄目ですね",
			expected: false,
		},
		{
			name:     "invalid URL with %",
			inputURL: "https://test.domain.com/test/hello%there%",
			expected: false,
		},
	}

	for i, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actual := IsValidURL(testCase.inputURL)
			if actual != testCase.expected {
				t.Errorf("Test %v - '%s' FAIL: expected URL: %v, actual: %v", i, testCase.name, testCase.expected, actual)
			}
		})
	}
}
