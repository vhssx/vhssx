package utils

import "strings"

var require = NotEmpty

func NotEmpty(value string) bool {
	return strings.TrimSpace(value) != ""
}
