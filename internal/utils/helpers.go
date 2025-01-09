package utils

import "strings"

func BuildQueryString(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	return "?" + strings.Join(parts, "&")
}
