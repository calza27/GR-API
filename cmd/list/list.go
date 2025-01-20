package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/calza27/Gift-Registry/GR-API/internal/aws/awsclient"
	"github.com/calza27/Gift-Registry/GR-API/internal/handlers/lists"
	"github.com/calza27/Gift-Registry/GR-API/internal/repositories"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
	"github.com/calza27/Gift-Registry/GR-API/middleware"
)

type Handler struct {
	ListHandler lists.ListHandler
}

func main() {
	ssmClient, err := awsclient.GetSsmClient()
	if err != nil {
		panic(err)
	}
	ddbName, err := ssmClient.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name: aws.String("/gift-registry/ddb/lists/name"),
	})
	if err != nil {
		panic(err)
	}
	tableName := *ddbName.Parameter.Value
	if tableName == "" {
		panic("list table name is blank")
	}

	userIndexName, err := ssmClient.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name: aws.String("/gift-registry/ddb/lists/user-index/name"),
	})
	if err != nil {
		panic(err)
	}
	userIdIndexName := *userIndexName.Parameter.Value
	if tableName == "" {
		panic("userId index name is blank")
	}

	sharingIndexName, err := ssmClient.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name: aws.String("/gift-registry/ddb/lists/sharing-id-index/name"),
	})
	if err != nil {
		panic(err)
	}
	sharingIdIndexName := *sharingIndexName.Parameter.Value
	if tableName == "" {
		panic("sharingId index name is blank")
	}

	listRepository, err := repositories.NewListRepository(tableName, userIdIndexName, sharingIdIndexName)
	if err != nil {
		panic(err)
	}

	handler := Handler{
		ListHandler: lists.NewListHandler(listRepository),
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
	case "updateList":
		handlerFunc = middleware.EnsureUserIdPresent(h.ListHandler.HandleUpdateList)
	case "deleteList":
		handlerFunc = middleware.EnsureUserIdPresent(h.ListHandler.HandleRemoveList)
	default:
		return utils.BuildResponse(fmt.Sprintf("unknown operation %s", operationName), 400, nil), nil
	}
	return handlerFunc(request), nil
}
