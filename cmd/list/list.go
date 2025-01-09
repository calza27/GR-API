package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/calza27/Gift-Registry/GR-API/internal/handlers/lists"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
	"github.com/calza27/Gift-Registry/GR-API/middleware"
)

type Handler struct {
	ListHandler lists.ListHandler
}

func main() {
	handler := Handler{
		ListHandler: lists.NewListHandler(),
	}
	lambda.Start(middleware.BoundaryLogging(handler.HandleRequest))
}

func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	operationName := request.RequestContext.OperationName
	var handlerFunc func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	switch operationName {
	case "addList":
		handlerFunc = middleware.EnsureUserIdPresent(h.ListHandler.HandleAddList)
	case "getListList":
		handlerFunc = middleware.EnsureUserIdPresent(h.ListHandler.HandleGetListList)
	case "getList":
		handlerFunc = h.ListHandler.HandleGetList
	case "getGiftList":
		handlerFunc = h.ListHandler.HandleGetGiftList
	case "updateList":
		handlerFunc = middleware.EnsureUserIdPresent(h.ListHandler.HandleUpdateList)
	case "removeList":
		handlerFunc = middleware.EnsureUserIdPresent(h.ListHandler.HandleRemoveList)
	default:
		return utils.BuildResponse(fmt.Sprintf("unknown operation %s", operationName), 400, nil), nil
	}
	return handlerFunc(request), nil
}
