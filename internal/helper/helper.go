// Package helper contains helper functions that used in different packages
package helper

import (
	"regexp"
	"strings"
	"unicode"
)

// ToSnakeCase used to convert strings from camel to a snake case
func ToSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	rez := strings.ToLower(snake)
	return rez
}

// LcFirst convert first string symbol to a lower case
func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
