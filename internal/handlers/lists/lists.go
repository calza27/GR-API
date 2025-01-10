package lists

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/calza27/Gift-Registry/GR-API/internal/models"
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
	var list models.List
	err := utils.DecodeRequestBody(request, &list)
	if err != nil {
		fmt.Printf("Error unmarshalling request body: %w\n", err)
		return utils.BuildResponse("Error unmarshalling request body", 400, nil)
	}
	userId := utils.GetUserIdFromRequest(request)
	if userId == "" {
		return utils.BuildResponse("Missing required auth parameter: user_id", 400, nil)
	}
	list.UserId = userId
	err = h.ListRepository.CreateList(list)
	if err != nil {
		fmt.Printf("Error writing list to DDB: %w\n", err)
		return utils.BuildResponse("Error adding list", 500, nil)
	}
	return utils.BuildResponse("", 201, nil)
}

func (h *ListHandlerImpl) HandleGetListList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	userId := utils.GetUserIdFromRequest(request)
	if userId == "" {
		return utils.BuildResponse("Missing required auth parameter: user_id", 400, nil)
	}
	lists, err := h.ListRepository.GetListsByUserId(userId)
	if err != nil {
		fmt.Printf("Error getting lists for user: %w\n", err)
		return utils.BuildResponse("Error getting list of gift lists", 500, nil)
	}
	if len(lists) == 0 {
		return utils.BuildResponse("No lists found", 404, nil)
	}
	body, err := utils.EncodeResponseBody(lists)
	if err != nil {
		fmt.Printf("Error marshalling lists for user: %w\n", err)
		return utils.BuildResponse("Error getting list of gift lists", 500, nil)
	}
	return utils.BuildResponse(body, 200, nil)
}

func (h *ListHandlerImpl) HandleGetList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	listId := request.PathParameters["list_id"]
	if listId == "" {
		return utils.BuildResponse("Missing required path parameter: list_id", 400, nil)
	}
	list, err := h.ListRepository.GetListById(listId)
	if err != nil {
		fmt.Printf("Error getting list: %w\n", err)
		return utils.BuildResponse("Error getting list", 500, nil)
	}

	body, err := utils.EncodeResponseBody(list)
	if err != nil {
		fmt.Printf("Error marshalling list: %w\n", err)
		return utils.BuildResponse("Error getting list", 500, nil)
	}
	return utils.BuildResponse(body, 200, nil)
}

func (h *ListHandlerImpl) HandleUpdateList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	userId := utils.GetUserIdFromRequest(request)
	if userId == "" {
		return utils.BuildResponse("Missing required auth parameter: user_id", 400, nil)
	}
	listId := request.PathParameters["list_id"]
	if listId == "" {
		return utils.BuildResponse("Missing required path parameter: list_id", 400, nil)
	}
	var list models.List
	err := utils.DecodeRequestBody(request, &list)
	if err != nil {
		fmt.Printf("Error unmarshalling request body: %w\n", err)
		return utils.BuildResponse("Error unmarshalling request body", 400, nil)
	}
	if listId != list.Id {
		return utils.BuildResponse("List ID in path does not match List ID in request body", 400, nil)
	}
	if list.UserId != userId {
		return utils.BuildResponse("List does not belong to user", 403, nil)
	}
	err = h.ListRepository.UpdateList(list)
	if err != nil {
		fmt.Printf("Error updating list: %w\n", err)
		return utils.BuildResponse("Error updating list", 500, nil)
	}
	return utils.BuildResponse("", 201, nil)
}

func (h *ListHandlerImpl) HandleRemoveList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	userId := utils.GetUserIdFromRequest(request)
	if userId == "" {
		return utils.BuildResponse("Missing required auth parameter: user_id", 400, nil)
	}
	listId := request.PathParameters["list_id"]
	if listId == "" {
		return utils.BuildResponse("Missing required path parameter: list_id", 400, nil)
	}
	list, err := h.ListRepository.GetListById(listId)
	if err != nil {
		fmt.Printf("Error getting list: %w\n", err)
		return utils.BuildResponse("Error removing gift", 500, nil)
	}
	if list.UserId != userId {
		return utils.BuildResponse("List does not belong to user", 403, nil)
	}
	err = h.ListRepository.DeleteList(listId)
	if err != nil {
		fmt.Printf("Error deleting list: %w\n", err)
		return utils.BuildResponse("Error removing gift", 500, nil)
	}
	return utils.BuildResponse("", 201, nil)
}
