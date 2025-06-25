package spider

import (
	"reflect"
	"sort"
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
			inputURL: "https://test.domain.com",
			inputBody: `
		<html>
			<body>
				<a href="/path/one">
					<span>Test.domain.com</span>
				</a>
				<a href="https://other.com/path/one">
					<span>Test.domain.com</span>
				</a>
			</body>
		</html>
		`,
			expected: []string{"https://test.domain.com/path/one", "https://other.com/path/one"},
		},
		//test 1
		{
			name:     "no anchor tags",
			inputURL: "https://test.domain.com",
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
			inputURL: "https://test.domain.com",
			inputBody: `
				<html>
					<body>
						<a href="">
							Empty link
						</a>
					</body>
				</html>
			`,
			expected: []string{"https://test.domain.com"},
		},
		// test 3
		{
			name:     "malformed href",
			inputURL: "https://test.domain.com",
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
			inputURL: "https://test.domain.com",
			inputBody: `
				<html>
					<body>
						<a href="//cdn.test.domain.com/script.js">
							CDN
						</a>
					</body>
				</html>
			`,
			expected: []string{"https://cdn.test.domain.com/script.js"},
		},
		// test 5
		{
			name:     "relative path traversal",
			inputURL: "https://test.domain.com/tutorials",
			inputBody: `
				<html>
					<body>
						<a href="../about">
							About
						</a>
					</body>
				</html>
			`,
			expected: []string{"https://test.domain.com/about"},
		},
		// test 6
		{
			name:     "remove duplicate links",
			inputURL: "https://test.domain.com",
			inputBody: `
			<html>
                <body>
                    <a href="/valid-link"><span>Valid</span></a>
                    <a href="<invalid></a>"><span>Broken</span></a>
                    <a href="https://valid.com/path"></a>
                    <a href="/valid-link"><span>Valid</span></a>
                    <a href="<invalid></a>"><span>Broken</span></a>
                    <a href="https://valid.com/path"></a>
                </body>
            </html>
			`,
			expected: []string{"https://test.domain.com/valid-link", "https://valid.com/path"},
		},
		// test 7
		{
			name:     "ignore non-ASCII links",
			inputURL: "https://test.domain.com",
			inputBody: `
			<html>
                <body>
                    <a href="/valid-link"><span>Valid</span></a>
                    <a href="https://valid.com/path"></a>
                    <a href="https://пример.рф">Cyrillic</a>
                    <a href="https://例子.com">Chinese</a>
                    <a href="https://テスト.jp">Japanese</a>
                    <a href="/another-valid"></a>
                </body>
            </html>
			`,
			expected: []string{"https://test.domain.com/valid-link", "https://valid.com/path", "https://test.domain.com/another-valid"},
		},
	}

	for i, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actual, _, err := getURLsFromHTML(testCase.inputBody, testCase.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, testCase.name, err)
				return
			}

			sort.Strings(testCase.expected)
			sort.Strings(actual)

			if !reflect.DeepEqual(testCase.expected, actual) {
				t.Errorf("Test %v - '%s' FAIL: expected parsed URLs: %v, actual: %v", i, testCase.name, testCase.expected, actual)
			}
		})
	}
}
