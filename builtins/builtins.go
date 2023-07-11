package builtins

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Title(input string) string {
	return cases.Title(language.Make("en")).String(input)
}

func SnakeToCamelCase(input string) string {
	var output string
	var prevChar rune

	for _, ch := range input {
		if ch != '_' {
			if prevChar == '_' || int32(prevChar) == 0 {
				output = output + strings.ToUpper(string(ch))
			} else {
				output = output + string(ch)
			}
		}
		prevChar = ch

	}
	return output
}
