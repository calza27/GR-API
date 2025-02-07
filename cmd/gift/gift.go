package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/calza27/Gift-Registry/GR-API/internal/aws/awsclient"
	"github.com/calza27/Gift-Registry/GR-API/internal/handlers/gifts"
	"github.com/calza27/Gift-Registry/GR-API/internal/repositories"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
	"github.com/calza27/Gift-Registry/GR-API/middleware"
)

type Handler struct {
	GiftHandler gifts.GiftHandler
}

func main() {
	ssmClient, err := awsclient.GetSsmClient()
	if err != nil {
		panic(err)
	}
	ddbName, err := ssmClient.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name: aws.String("/gift-registry/ddb/gifts/name"),
	})
	if err != nil {
		panic(err)
	}
	tableName := *ddbName.Parameter.Value
	if tableName == "" {
		panic("gift table name is blank")
	}

	indexName, err := ssmClient.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name: aws.String("/gift-registry/ddb/gifts/list-index/name"),
	})
	if err != nil {
		panic(err)
	}
	listIdIndexName := *indexName.Parameter.Value
	if tableName == "" {
		panic("listId index name is blank")
	}

	giftRepository, err := repositories.NewGiftRepository(tableName, listIdIndexName)
	if err != nil {
		panic(err)
	}

	handler := Handler{
		GiftHandler: gifts.NewGiftHandler(giftRepository),
	}
	lambda.Start(middleware.BoundaryLogging(handler.HandleRequest))
}

func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	operationName := request.RequestContext.OperationName
	var handlerFunc func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	switch operationName {
	case "addGift":
		handlerFunc = middleware.EnsureUserIdPresent(h.GiftHandler.HandleAddGift)
	case "getGiftList":
		handlerFunc = h.GiftHandler.HandleGetGiftList
	case "getGift":
		handlerFunc = h.GiftHandler.HandleGetGift
	case "updateGift":
		handlerFunc = middleware.EnsureUserIdPresent(h.GiftHandler.HandleUpdateGift)
	case "deleteGift":
		handlerFunc = middleware.EnsureUserIdPresent(h.GiftHandler.HandleRemoveGift)
	case "purchaseGift":
		handlerFunc = h.GiftHandler.HandlePurchaseGift
	case "unpurchaseGift":
		handlerFunc = middleware.EnsureUserIdPresent(h.GiftHandler.HandleUnpurchaseGift)
	default:
		return utils.BuildResponse(fmt.Sprintf("unknown operation %s", operationName), 400, nil), nil
	}
	return handlerFunc(request), nil
}
