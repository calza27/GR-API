package middleware

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type HandlerFunc func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse

func ChainHandlers(h ...HandlerFunc) HandlerFunc {
	return HandlerFunc(func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		var resp events.APIGatewayProxyResponse
		for _, handler := range h {
			resp = handler(request)
			if resp.StatusCode >= 300 {
				return resp
			}
		}
		return resp
	})
}

func getUserIdFromAuth(request events.APIGatewayProxyRequest) (*string, error) {
	if request.RequestContext.Authorizer == nil {
		return nil, fmt.Errorf("err: No authorizer found in request context")
	}
	if _, ok := request.RequestContext.Authorizer["user_id"]; !ok {
		return nil, fmt.Errorf("err: No user_id field found in Authorisation from Authorizer")
	}
	userId := request.RequestContext.Authorizer["user_id"].(string)
	if len(userId) == 0 {
		return nil, fmt.Errorf("err: No user_id found in Authorisation from Authorizer")
	}
	return &userId, nil
}
