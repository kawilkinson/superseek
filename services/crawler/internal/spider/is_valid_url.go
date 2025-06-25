package spider

import (
	"strings"
	"unicode"
)

func IsValidURL(link string) bool {
	if strings.Contains(link, "w/index.php") {
		return false
	}

	for _, char := range link {
		if char > 127 || (!unicode.IsLetter(char) && !unicode.IsDigit(char) && !isAllowedSymbol(char)) {
			return false
		}
	}

	return !strings.Contains(link, "%")
}

func isAllowedSymbol(char rune) bool {
	allowed := "-._~:/?#[]@!$&'()*+,;="
	return (char < 127 && unicode.IsPrint(char)) || containsRune(allowed, char)
}

func containsRune(str string, char rune) bool {
	for _, strChar := range str {
		if strChar == char {
			return true
		}
	}

	return false
}
