package builtins

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Title(input string) string {
	return cases.Title(language.Make("en")).String(input)
}
