package middleware

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
)

type EnsureUserIdPresentFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func EnsureUserIdPresent(h func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return EnsureUserIdPresentFunc(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		userId, err := getUserIdFromAuth(request)
		if err != nil {
			fmt.Printf("Err: %s\n", err)
			return utils.BuildResponse("Forbidden", 403, nil), nil
		}
		if userId == nil {
			fmt.Printf("Err: No user_id found in Authorisation from Authorizer\n")
			return utils.BuildResponse("Forbidden", 403, nil), nil
		}
		return h(ctx, request)
	})
}
