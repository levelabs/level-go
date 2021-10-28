package common

import (
	"strconv"
	"strings"
	"unicode"
)

func BuildUrl(url string, i int) string {
	return strings.Join([]string{strings.TrimRightFunc(url, func(r rune) bool {
		return unicode.IsNumber(r)
	}), strconv.Itoa(i)}, "")
}
