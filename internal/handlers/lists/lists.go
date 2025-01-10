package lists

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/calza27/Gift-Registry/GR-API/internal/repositories"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
)

type ListHandler interface {
	HandleAddList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	HandleGetListList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	HandleGetList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	HandleUpdateList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	HandleRemoveList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
}

type ListHandlerImpl struct {
	ListRepository repositories.ListRepository
}

func NewListHandler(listRepo repositories.ListRepository) ListHandler {
	return &ListHandlerImpl{
		ListRepository: listRepo,
	}
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

func (h *ListHandlerImpl) HandleUpdateList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return utils.BuildResponse("", 200, nil)
}

func (h *ListHandlerImpl) HandleRemoveList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return utils.BuildResponse("", 200, nil)
}
