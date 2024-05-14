package helpers

import (
	"strings"
)

func ToSnakeCase(str string) string {
	var result string

	for i, char := range str {
		if i > 0 && 'A' <= char && char <= 'Z' {
			result += "_"
		}

		result += strings.ToLower(string(char))
	}

	return result
}