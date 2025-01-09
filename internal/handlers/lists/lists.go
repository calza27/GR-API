package lists

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
)

type ListHandler interface {
	HandleAddList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	HandleGetListList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	HandleGetList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	HandleGetGiftList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	HandleUpdateList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	HandleRemoveList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
}

type ListHandlerImpl struct {
}

func NewListHandler() ListHandler {
	return &ListHandlerImpl{}
}

func (h *ListHandlerImpl) HandleAddList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return utils.BuildResponse("", 200, nil)
}

func (h *ListHandlerImpl) HandleGetListList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return utils.BuildResponse("", 200, nil)
}

func (h *ListHandlerImpl) HandleGetList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return utils.BuildResponse("", 200, nil)
}

func (h *ListHandlerImpl) HandleGetGiftList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return utils.BuildResponse("", 200, nil)
}

func (h *ListHandlerImpl) HandleUpdateList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return utils.BuildResponse("", 200, nil)
}

func (h *ListHandlerImpl) HandleRemoveList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return utils.BuildResponse("", 200, nil)
}
