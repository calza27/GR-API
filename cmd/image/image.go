package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/calza27/Gift-Registry/GR-API/internal/aws/awsclient"
	"github.com/calza27/Gift-Registry/GR-API/internal/handlers/image"
	"github.com/calza27/Gift-Registry/GR-API/internal/repositories"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
	"github.com/calza27/Gift-Registry/GR-API/middleware"
)

type Handler struct {
	ImageHandler image.ImageHandler
}

func main() {
	ssmClient, err := awsclient.GetSsmClient()
	if err != nil {
		panic(err)
	}
	s3Name, err := ssmClient.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name: aws.String("/gift-registry/s3/images/name"),
	})
	if err != nil {
		panic(err)
	}
	bucketName := *s3Name.Parameter.Value
	if bucketName == "" {
		panic("bucket name is blank")
	}

	urlLifespanParam, err := ssmClient.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name: aws.String("/gift-registry/s3/images/url-lifespan"),
	})
	if err != nil {
		panic(err)
	}
	lifespan, err := time.ParseDuration(*urlLifespanParam.Parameter.Value)
	if err != nil {
		panic(err)
	}

	imageRepository, err := repositories.NewImageRepository(bucketName, lifespan)
	if err != nil {
		panic(err)
	}
	handler := Handler{
		ImageHandler: image.NewImageHandler(imageRepository),
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
