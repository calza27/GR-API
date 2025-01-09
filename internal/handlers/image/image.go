package image

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
)

type ImageHandler interface {
	HandleUploadImage(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
}

type ImageHandlerImpl struct {
}

func NewImageHandler() ImageHandler {
	return &ImageHandlerImpl{}
}

func (h *ImageHandlerImpl) HandleUploadImage(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return utils.BuildResponse("", 200, nil)
}
