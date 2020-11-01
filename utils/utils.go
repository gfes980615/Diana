package utils

import (
	"unicode"
)

func RemoveExtraChar(title string) string {
	var s []int32
	for _, t := range title {
		if unicode.Is(unicode.Han, t) || unicode.IsDigit(t) || unicode.IsLetter(t) {
			s = append(s, t)
		}
	}
	return string(s)
}
