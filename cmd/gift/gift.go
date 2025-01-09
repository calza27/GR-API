package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/calza27/Gift-Registry/GR-API/internal/handlers/gifts"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
	"github.com/calza27/Gift-Registry/GR-API/middleware"
)

type Handler struct {
	GiftHandler gifts.GiftHandler
}

func main() {
	handler := Handler{
		GiftHandler: gifts.NewGiftHandler(),
	}
	lambda.Start(middleware.BoundaryLogging(middleware.EnsureUserIdPresent(handler.HandleRequest)))
}

func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	operationName := request.RequestContext.OperationName
	var handlerFunc func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	switch operationName {
	case "addGift":
		handlerFunc = h.GiftHandler.HandleAddGift
	case "getGift":
		handlerFunc = h.GiftHandler.HandleGetGift
	case "updateGift":
		handlerFunc = h.GiftHandler.HandleUpdateGift
	case "deleteGift":
		handlerFunc = h.GiftHandler.HandleRemoveGift
	default:
		return utils.BuildResponse(fmt.Sprintf("unknown operation %s", operationName), 400, nil), nil
	}
	return handlerFunc(request), nil
}
