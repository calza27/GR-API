package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/calza27/Gift-Registry/GR-API/internal/handlers/image"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
	"github.com/calza27/Gift-Registry/GR-API/middleware"
)

type Handler struct {
	ImageHandler image.ImageHandler
}

func main() {
	handler := Handler{
		ImageHandler: image.NewImageHandler(),
	}
	lambda.Start(middleware.BoundaryLogging(handler.HandleRequest))
}

func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	operationName := request.RequestContext.OperationName
	var handlerFunc func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	switch operationName {
	case "uploadImage":
		handlerFunc = middleware.EnsureUserIdPresent(h.ImageHandler.HandleUploadImage)
	case "getImageUrl":
		handlerFunc = h.ImageHandler.HandleGetImageUrl
	default:
		return utils.BuildResponse(fmt.Sprintf("unknown operation %s", operationName), 400, nil), nil
	}
	return handlerFunc(request), nil
}
