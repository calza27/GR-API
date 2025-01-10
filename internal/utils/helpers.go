package utils

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
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

func DecodeRequestBody(request events.APIGatewayProxyRequest, v interface{}) error {
	return json.Unmarshal([]byte(request.Body), &v)
}

func EncodeResponseBody(v interface{}) (string, error) {
	body, err := json.Marshal(v)
	return string(body), err
}
