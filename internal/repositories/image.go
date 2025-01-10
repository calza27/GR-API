package repositories

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/calza27/Gift-Registry/GR-API/internal/aws/awsclient"
	"github.com/calza27/Gift-Registry/GR-API/internal/models"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
)

type ImageRepository interface {
	PutImage(image models.Image) (*string, error)
	GetImageUrl(key string) (*string, error)
}

type S3ImageRepository struct {
	s3          *s3.Client
	bucketname  string
	urlLifespan time.Duration
}

func NewImageRepository(bucketname string, urlLifespan time.Duration) (ImageRepository, error) {
	s3, err := awsclient.GetS3Client()
	if err != nil {
		return nil, fmt.Errorf("Error when initialising connection to S3: %w", err)
	}
	return &S3ImageRepository{
		s3:          s3,
		bucketname:  bucketname,
		urlLifespan: urlLifespan,
	}, nil
}

func (r *S3ImageRepository) PutImage(image models.Image) (*string, error) {
	extension := filepath.Ext(image.FileName)
	newFileName := fmt.Sprintf("%s%s", utils.GenerateUUID(), extension)
	fileBytes, err := base64.StdEncoding.DecodeString(image.FileData)
	if err != nil {
		return nil, fmt.Errorf("Error decoding base64 string: %w", err)
	}
	params := &s3.PutObjectInput{
		Bucket: aws.String(r.bucketname),
		Key:    aws.String(newFileName),
		Body:   bytes.NewReader(fileBytes),
	}
	if _, err := r.s3.PutObject(context.Background(), params); err != nil {
		return nil, fmt.Errorf("Error when writing object to S3: %w", err)
	}
	return &newFileName, nil
}

func (r *S3ImageRepository) GetImageUrl(fileName string) (*string, error) {
	_, err := r.s3.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(r.bucketname),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, fmt.Errorf("file %s not found: %w", fileName, err)
	}
	presignClient := s3.NewPresignClient(r.s3)
	req, err := presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(r.bucketname),
		Key:    aws.String(fileName),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = r.urlLifespan
	})
	if err != nil {
		return nil, fmt.Errorf("Error when geenerating presigned URL for object: %w", err)
	}
	return &req.URL, nil
}
