package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
)

type LoggingHandlerFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func BoundaryLogging(h func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return LoggingHandlerFunc(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		jReq, _ := json.Marshal(buildRequestForLog(req))
		fmt.Printf("Incoming Request: %+v\n", string(jReq))
		resp, err := h(ctx, req)
		jResp, _ := json.Marshal(buildResponseForLog(resp))
		fmt.Printf("Outgoing Response: %+v\n", string(jResp))
		return resp, err
	})
}

func buildRequestForLog(r events.APIGatewayProxyRequest) map[string]string {
	request := map[string]string{
		"Resource":                    r.Resource,
		"Path":                        r.Path,
		"HTTPMethod":                  r.HTTPMethod,
		"Referer":                     r.Headers["Referer"],
		"User-Agent":                  r.Headers["User-Agent"],
		"X-Amzn-Trace-Id":             r.Headers["X-Amzn-Trace-Id"],
		"X-Forwarded-For":             r.Headers["X-Forwarded-For"],
		"Authorizer":                  fmt.Sprintf("%+v", r.RequestContext.Authorizer),
		"Authorizer.cognito:username": utils.GetUserIdFromRequest(r),
		"PathParameters":              preparePathParams(r),
		"QueryStringParameters":       prepareQueryParams(r),
		"Body":                        r.Body,
		"RequestTime":                 r.RequestContext.RequestTime,
	}
	return request
}

func buildResponseForLog(r events.APIGatewayProxyResponse) map[string]string {
	response := map[string]string{
		"StatusCode":      fmt.Sprintf("%d", r.StatusCode),
		"Body":            truncateText(r.Body, 200),
		"IsBase64Encoded": fmt.Sprintf("%t", r.IsBase64Encoded),
		"ResponseTime":    time.Now().Format(time.RFC3339),
	}
	return response
}

func prepareQueryParams(r events.APIGatewayProxyRequest) string {
	queryParts := []string{}
	for k, v := range r.QueryStringParameters {
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", k, v))
	}
	return utils.BuildQueryString(queryParts)
}

func preparePathParams(r events.APIGatewayProxyRequest) string {
	pathParts := []string{}
	for k, v := range r.PathParameters {
		pathParts = append(pathParts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(pathParts, ", ")
}

func truncateText(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "... [truncated]"
}
