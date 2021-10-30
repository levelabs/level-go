package common

import (
	"strconv"
	"strings"
	"unicode"
)

func BuildUrl(url string, i int) string {
	return strings.Join([]string{url, strconv.Itoa(i)}, "")
}

func TrimRightNumber(url string) string {
	return strings.TrimRightFunc(url, func(r rune) bool {
		return unicode.IsNumber(r)
	})
}
