package tfidfutils

import "testing"

func TestGetHTMLData(t *testing.T) {
	tests := []struct {
		name           string
		htmlInput      string
		expectError    bool
		expectedFields map[string]string
	}{
		{
			name: "basic HTML with title and paragraph",
			htmlInput: `
				<html>
					<head><title>Hello World This Is a Title</title></head>
					<body><p>This is a test paragraph for testing purposes.</p></body>
				</html>
			`,
			expectError: false,
			expectedFields: map[string]string{
				"title":       "Hello World This Is a Title",
				"description": "",
				"language":    "English",
			},
		},
		{
			name: "HTML with og:title and description metadata",
			htmlInput: `
				<html>
					<head>
						<title>Ignored Title</title>
						<meta property="og:title" content="Meta Title For Meta Purposes">
						<meta name="description" content="This is the description that I'm testing with.">
					</head>
					<body><p>Testing meta tag parsing for the amazing superseek</p></body>
				</html>
			`,
			expectError: false,
			expectedFields: map[string]string{
				"title":       "Meta Title For Meta Purposes",
				"description": "This is the description that I'm testing with.",
				"language":    "English",
			},
		},
		{
			name:        "invalid HTML input",
			htmlInput:   "<<>>",
			expectError: true,
		},
		{
			name: "HTML in Japanese",
			htmlInput: `
				<html>
					<head><title>こんにちは世界</title></head>
					<body><p>これは日本語のテキストです。言語検出のテストです。</p></body>
				</html>
			`,
			expectError: false,
			expectedFields: map[string]string{
				"title":       "こんにちは世界",
				"description": "",
				"language":    "Japanese",
			},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GetHTMLData(tc.htmlInput)
			if tc.expectError {
				if err == nil {
					t.Errorf("Test %d - '%s' FAIL: expected error but got nil", i, tc.name)
				}
				return
			}

			if err != nil {
				t.Errorf("Test %d - '%s' FAIL: unexpected error: %v", i, tc.name, err)
			}

			for key, expectedValue := range tc.expectedFields {
				actualValue, exists := result[key].(string)
				if !exists {
					t.Errorf("Test %d - '%s' FAIL: expected %s to be a string, got %T", i, tc.name, key, result[key])
					continue
				}
				if actualValue != expectedValue {
					t.Errorf("Test %d - '%s' FAIL: expected %s = '%s', got '%s'", i, tc.name, key, expectedValue, actualValue)
				}
			}
		})
	}
}
