package spider

import (
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		// test 0
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
		<html>
			<body>
				<a href="/path/one">
					<span>Boot.dev</span>
				</a>
				<a href="https://other.com/path/one">
					<span>Boot.dev</span>
				</a>
			</body>
		</html>
		`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		//test 1
		{
			name:     "no anchor tags",
			inputURL: "https://blog.boot.dev",
			inputBody: `
				<html>
					<body>
						<p>
							No links here
						</p>
					</body>
				</html>
			`,
			expected: []string{},
		},
		// test 2
		{
			name:     "empty href",
			inputURL: "https://blog.boot.dev",
			inputBody: `
				<html>
					<body>
						<a href="">
							Empty link
						</a>
					</body>
				</html>
			`,
			expected: []string{"https://blog.boot.dev"},
		},
		// test 3
		{
			name:     "malformed href",
			inputURL: "https://blog.boot.dev",
			inputBody: `
				<html>
					<body>
						<a href=":://broken-url">
							Broken url
						</a>
					</body>
				</html>
			`,
			expected: []string{},
		},
		// test 4
		{
			name:     "protocol-relative href",
			inputURL: "https://blog.boot.dev",
			inputBody: `
				<html>
					<body>
						<a href="//cdn.boot.dev/script.js">
							CDN
						</a>
					</body>
				</html>
			`,
			expected: []string{"https://cdn.boot.dev/script.js"},
		},
		// test 5
		{
			name:     "relative path traversal",
			inputURL: "https://blog.boot.dev/tutorials",
			inputBody: `
				<html>
					<body>
						<a href="../about">
							About
						</a>
					</body>
				</html>
			`,
			expected: []string{"https://blog.boot.dev/about"},
		},
	}

	for i, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(testCase.inputBody, testCase.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, testCase.name, err)
				return
			}

			if !reflect.DeepEqual(testCase.expected, actual) {
				t.Errorf("Test %v - '%s' FAIL: expected parsed URLs: %v, actual: %v", i, testCase.name, testCase.expected, actual)
			}
		})
	}
}
