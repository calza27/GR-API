package gifts

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/calza27/Gift-Registry/GR-API/internal/repositories"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
)

type GiftHandler interface {
	HandleAddGift(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	HandleGetGiftList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	HandleGetGift(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	HandleUpdateGift(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	HandleRemoveGift(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
}

type GiftHandlerImpl struct {
	GiftRepository repositories.GiftRepository
}

func NewGiftHandler(giftRepo repositories.GiftRepository) GiftHandler {
	return &GiftHandlerImpl{
		GiftRepository: giftRepo,
	}
}

func (h *GiftHandlerImpl) HandleAddGift(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return utils.BuildResponse("", 200, nil)
}

func (h *GiftHandlerImpl) HandleGetGiftList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return utils.BuildResponse("", 200, nil)
}

func (h *GiftHandlerImpl) HandleGetGift(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return utils.BuildResponse("", 200, nil)
}

func (h *GiftHandlerImpl) HandleUpdateGift(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return utils.BuildResponse("", 200, nil)
}

func (h *GiftHandlerImpl) HandleRemoveGift(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return utils.BuildResponse("", 200, nil)
}
