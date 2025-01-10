package image

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/calza27/Gift-Registry/GR-API/internal/models"
	"github.com/calza27/Gift-Registry/GR-API/internal/repositories"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
)

type ImageHandler interface {
	HandleUploadImage(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	HandleGetImageUrl(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
}

type ImageHandlerImpl struct {
	ImageRepository repositories.ImageRepository
}

func NewImageHandler(imageRepo repositories.ImageRepository) ImageHandler {
	return &ImageHandlerImpl{
		ImageRepository: imageRepo,
	}
}

func (h *ImageHandlerImpl) HandleUploadImage(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var image models.Image
	err := utils.DecodeRequestBody(request, &image)
	if err != nil {
		fmt.Printf("Error unmarshalling request body: %w\n", err)
		return utils.BuildResponse("Error unmarshalling request body", 400, nil)
	}
	if image.FileName == "" {
		return utils.BuildResponse("Missing required field: file_name", 400, nil)
	}
	if image.FileData == "" {
		return utils.BuildResponse("Missing required field: file_data", 400, nil)
	}
	newFileName, err := h.ImageRepository.PutImage(image)
	if err != nil {
		fmt.Printf("Error uploading image: %w\n", err)
		return utils.BuildResponse("Error uploading image", 500, nil)
	}
	if newFileName == nil {
		return utils.BuildResponse("Error uploading image, file rename blank", 500, nil)
	}
	return utils.BuildResponse(*newFileName, 200, nil)
}

func (h *ImageHandlerImpl) HandleGetImageUrl(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	fileName := request.PathParameters["file_name"]
	if fileName == "" {
		return utils.BuildResponse("Missing required path parameter: file_name", 400, nil)
	}
	url, err := h.ImageRepository.GetImageUrl(fileName)
	if err != nil {
		fmt.Printf("Error getting image url: %w\n", err)
		return utils.BuildResponse("Error getting image url", 500, nil)
	}
	if url == nil {
		return utils.BuildResponse("Image not found", 404, nil)
	}
	return utils.BuildResponse(*url, 200, nil)
}
