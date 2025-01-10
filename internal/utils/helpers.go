package utils

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

func BuildQueryString(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	return "?" + strings.Join(parts, "&")
}

func GenerateUUID() string {
	id := uuid.New()
	return id.String()
}

func DecodeRequestBody(body string, v interface{}) error {
	return json.Unmarshal([]byte(body), &v)
}
