package middleware

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
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
	userId := utils.GetUserIdFromRequest(request)
	if userId == "" {
		return nil, fmt.Errorf("err: No cognito:username found in Authorisation from Authorizer")
	}
	return &userId, nil
}
