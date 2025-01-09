package middleware

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
)

type EnsureUserIdPresentFunc func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse

func EnsureUserIdPresent(h func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse) func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return EnsureUserIdPresentFunc(func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		userId, err := getUserIdFromAuth(request)
		if err != nil {
			fmt.Printf("Err: %s\n", err)
			return utils.BuildResponse("Forbidden", 403, nil)
		}
		if userId == nil {
			fmt.Printf("Err: No user_id found in Authorisation from Authorizer\n")
			return utils.BuildResponse("Forbidden", 403, nil)
		}
		return h(request)
	})
}
