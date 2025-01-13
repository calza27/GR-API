package user

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/calza27/Gift-Registry/GR-API/internal/models"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
)

type UserHandler interface {
	HandleGetUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	HandleUpdateUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
}

type UserHandlerImpl struct {
}

func NewUserHandler() UserHandler {
	return &UserHandlerImpl{}
}

func (h *UserHandlerImpl) HandleUpdateUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var user models.User
	err := utils.DecodeRequestBody(request, &user)
	if err != nil {
		fmt.Printf("Error unmarshalling request body: %w\n", err)
		return utils.BuildResponse("Error unmarshalling request body", 400, nil)
	}
	authUserId := utils.GetUserIdFromRequest(request)
	if authUserId == "" {
		return utils.BuildResponse("Missing required auth parameter: user_id", 400, nil)
	}
	userId := request.PathParameters["user_id"]
	if userId == "" {
		return utils.BuildResponse("Missing required path parameter: user_id", 400, nil)
	}
	if authUserId != user.Id {
		return utils.BuildResponse("Authentication User ID in path does not match User ID in request body", 400, nil)
	}
	if userId != user.Id {
		return utils.BuildResponse("User ID in path does not match User ID in request body", 400, nil)
	}
	return utils.BuildResponse("", 201, nil)
}

func (h *UserHandlerImpl) HandleGetUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var user models.User
	err := utils.DecodeRequestBody(request, &user)
	if err != nil {
		fmt.Printf("Error unmarshalling request body: %w\n", err)
		return utils.BuildResponse("Error unmarshalling request body", 400, nil)
	}
	authUserId := utils.GetUserIdFromRequest(request)
	if authUserId == "" {
		return utils.BuildResponse("Missing required auth parameter: user_id", 400, nil)
	}
	userId := request.PathParameters["user_id"]
	if userId == "" {
		return utils.BuildResponse("Missing required path parameter: user_id", 400, nil)
	}
	if authUserId != user.Id {
		return utils.BuildResponse("Authentication User ID in path does not match User ID in request body", 400, nil)
	}
	if userId != user.Id {
		return utils.BuildResponse("User ID in path does not match User ID in request body", 400, nil)
	}
	return utils.BuildResponse("", 200, nil)
}
