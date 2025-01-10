package gifts

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/calza27/Gift-Registry/GR-API/internal/models"
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
	var gift models.Gift
	err := utils.DecodeRequestBody(request, &gift)
	if err != nil {
		fmt.Printf("Error unmarshalling request body: %w\n", err)
		return utils.BuildResponse("Error unmarshalling request body", 400, nil)
	}
	listId := request.PathParameters["list_id"]
	if listId == "" {
		return utils.BuildResponse("Missing required path parameter: list_id", 400, nil)
	}
	if listId != gift.ListId {
		return utils.BuildResponse("List ID in path does not match List ID in request body", 400, nil)
	}

	err = h.GiftRepository.CreateGift(gift)
	if err != nil {
		fmt.Printf("Error writing gift to DDB: %w\n", err)
		return utils.BuildResponse("Error adding gift", 500, nil)
	}
	return utils.BuildResponse("", 201, nil)
}

func (h *GiftHandlerImpl) HandleGetGiftList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	listId := request.PathParameters["list_id"]
	if listId == "" {
		return utils.BuildResponse("Missing required path parameter: list_id", 400, nil)
	}
	gifts, err := h.GiftRepository.GetGiftsByListId(listId)
	if err != nil {
		fmt.Printf("Error getting gift list: %w\n", err)
		return utils.BuildResponse("Error getting gift list", 500, nil)
	}
	body, err := utils.EncodeResponseBody(gifts)
	if err != nil {
		fmt.Printf("Error marshalling gift list: %w\n", err)
		return utils.BuildResponse("Error getting gift list", 500, nil)
	}
	return utils.BuildResponse(body, 200, nil)
}

func (h *GiftHandlerImpl) HandleGetGift(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	listId := request.PathParameters["list_id"]
	if listId == "" {
		return utils.BuildResponse("Missing required path parameter: list_id", 400, nil)
	}
	giftId := request.PathParameters["gift_id"]
	if giftId == "" {
		return utils.BuildResponse("Missing required path parameter: gift_id", 400, nil)
	}

	gift, err := h.GiftRepository.GetGiftById(giftId)
	if err != nil {
		fmt.Printf("Error getting gift: %w\n", err)
		return utils.BuildResponse("Error getting gift", 500, nil)
	}
	if gift.ListId != listId {
		return utils.BuildResponse("Gift not found", 404, nil)
	}

	body, err := utils.EncodeResponseBody(gift)
	if err != nil {
		fmt.Printf("Error marshalling gift: %w\n", err)
		return utils.BuildResponse("Error getting gift", 500, nil)
	}
	return utils.BuildResponse(body, 200, nil)
}

func (h *GiftHandlerImpl) HandleUpdateGift(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var gift models.Gift
	err := utils.DecodeRequestBody(request, &gift)
	if err != nil {
		fmt.Printf("Error unmarshalling request body: %w\n", err)
		return utils.BuildResponse("Error unmarshalling request body", 400, nil)
	}
	giftId := request.PathParameters["gift_id"]
	if giftId == "" {
		return utils.BuildResponse("Missing required path parameter: gift_id", 400, nil)
	}
	if giftId != gift.Id {
		return utils.BuildResponse("Gift ID in path does not match Gift ID in request body", 400, nil)
	}
	listId := request.PathParameters["list_id"]
	if listId == "" {
		return utils.BuildResponse("Missing required path parameter: list_id", 400, nil)
	}
	if listId != gift.ListId {
		return utils.BuildResponse("List ID in path does not match List ID in request body", 400, nil)
	}

	err = h.GiftRepository.UpdateGift(gift)
	if err != nil {
		fmt.Printf("Error updating gift: %w\n", err)
		return utils.BuildResponse("Error updating gift", 500, nil)
	}
	return utils.BuildResponse("", 201, nil)
}

func (h *GiftHandlerImpl) HandleRemoveGift(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	listId := request.PathParameters["list_id"]
	if listId == "" {
		return utils.BuildResponse("Missing required path parameter: list_id", 400, nil)
	}
	giftId := request.PathParameters["gift_id"]
	if giftId == "" {
		return utils.BuildResponse("Missing required path parameter: gift_id", 400, nil)
	}

	gift, err := h.GiftRepository.GetGiftById(giftId)
	if err != nil {
		fmt.Printf("Error getting gift: %w\n", err)
		return utils.BuildResponse("Error removing gift", 500, nil)
	}
	if gift.ListId != listId {
		return utils.BuildResponse("Gift not found", 404, nil)
	}

	err = h.GiftRepository.DeleteGift(giftId)
	if err != nil {
		fmt.Printf("Error removing gift: %w\n", err)
		return utils.BuildResponse("Error removing gift", 500, nil)
	}
	return utils.BuildResponse("", 201, nil)
}
