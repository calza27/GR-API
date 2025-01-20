package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/calza27/Gift-Registry/GR-API/internal/aws/awsclient"
	"github.com/calza27/Gift-Registry/GR-API/internal/models"
	"github.com/calza27/Gift-Registry/GR-API/internal/utils"
)

type GiftRepository interface {
	CreateGift(gift models.Gift) error
	UpdateGift(gift models.Gift) error
	GetGiftById(giftId string) (models.Gift, error)
	GetGiftsByListId(listId string) ([]models.Gift, error)
	DeleteGift(giftId string) error
	persistGift(item map[string]types.AttributeValue) error
}

type DynamoDbGiftRepository struct {
	db          *dynamodb.Client
	tablename   string
	listIdIndex string
}

func NewGiftRepository(tableName, listIdIndex string) (GiftRepository, error) {
	db, err := awsclient.GetDynamodbClient()
	if err != nil {
		return nil, fmt.Errorf("Error when initialising connection to DDB: %w", err)
	}

	return &DynamoDbGiftRepository{
		db:          db,
		tablename:   tableName,
		listIdIndex: listIdIndex,
	}, nil
}

type GiftEntity struct {
	Id              string `dynamodbav:"id"`
	ListId          string `dynamodbav:"listId"`
	CreatedAt       string `dynamodbav:"createdAt"`
	Title           string `dynamodbav:"title,omitempty"`
	Desription      string `dynamodbav:"description,omitempty"`
	PlaceOfPurchase string `dynamodbav:"placeOfPurchase,omitempty"`
	ImageFileName   string `dynamodbav:"imageFileName,omitempty"`
	Url             string `dynamodbav:"url,omitempty"`
	Price           int    `dynamodbav:"price,omitempty"`
	Rank            int    `dynamodbav:"rank,omitempty"`
}

func (r *DynamoDbGiftRepository) CreateGift(gift models.Gift) error {
	giftItem := convertToGiftModel(gift)
	now := time.Now()
	giftItem.CreatedAt = utils.DateTimeToString(now)
	giftItem.Id = utils.GenerateUUID()
	item, err := attributevalue.MarshalMap(giftItem)
	if err != nil {
		return fmt.Errorf("error when trying to convert gift data to dynamodbattribute: %w", err)
	}
	return r.persistGift(item)
}

func (r *DynamoDbGiftRepository) UpdateGift(gift models.Gift) error {
	giftItem := convertToGiftModel(gift)
	item, err := attributevalue.MarshalMap(giftItem)
	if err != nil {
		return fmt.Errorf("error when trying to convert gift data to dynamodbattribute: %w", err)
	}
	return r.persistGift(item)
}

func (r *DynamoDbGiftRepository) persistGift(item map[string]types.AttributeValue) error {
	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(r.tablename),
	}
	if _, err := r.db.PutItem(context.Background(), params); err != nil {
		return fmt.Errorf("Error when trying to persist: %w", err)
	}
	return nil
}

func (r *DynamoDbGiftRepository) GetGiftById(giftId string) (models.Gift, error) {
	params := &dynamodb.QueryInput{
		TableName:              aws.String(r.tablename),
		KeyConditionExpression: aws.String("id = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: giftId},
		},
	}
	result, err := r.db.Query(context.Background(), params)
	if err != nil {
		return models.Gift{}, fmt.Errorf("Error when trying to get gift by id: %w", err)
	}
	if len(result.Items) == 0 {
		return models.Gift{}, fmt.Errorf("No gift found with id: %s", giftId)
	}
	if len(result.Items) > 1 {
		return models.Gift{}, fmt.Errorf("More than one gift found with id: %s", giftId)
	}
	var gift GiftEntity
	if err := attributevalue.UnmarshalMap(result.Items[0], &gift); err != nil {
		return models.Gift{}, fmt.Errorf("Error when trying to unmarshal gift data: %w", err)
	}
	return convertToGift(gift), nil
}

func (r *DynamoDbGiftRepository) GetGiftsByListId(listId string) ([]models.Gift, error) {
	params := &dynamodb.QueryInput{
		TableName:              aws.String(r.tablename),
		IndexName:              aws.String(r.listIdIndex),
		KeyConditionExpression: aws.String("listId = :listId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":listId": &types.AttributeValueMemberS{Value: listId},
		},
	}
	result, err := r.db.Query(context.Background(), params)
	if err != nil {
		return []models.Gift{}, fmt.Errorf("Error when trying to get gift by id: %w", err)
	}
	var gifts []models.Gift
	for _, item := range result.Items {
		var giftItem GiftEntity
		if err := attributevalue.UnmarshalMap(item, &giftItem); err != nil {
			return []models.Gift{}, fmt.Errorf("Error when trying to unmarshal gift data: %w", err)
		}
		gift := convertToGift(giftItem)
		gifts = append(gifts, gift)
	}
	return gifts, nil
}

func (r *DynamoDbGiftRepository) DeleteGift(giftId string) error {
	gift, err := r.GetGiftById(giftId)
	if err != nil {
		return fmt.Errorf("Error when trying to verify gift by id: %w", err)
	}
	if gift.Id == "" {
		return fmt.Errorf("No gift found with id to delete: %s", giftId)
	}
	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tablename),
		Key: map[string]types.AttributeValue{
			"id":        &types.AttributeValueMemberS{Value: gift.Id},
			"createdAt": &types.AttributeValueMemberS{Value: gift.CreatedAt},
		},
	}
	_, err = r.db.DeleteItem(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error when trying to delete gift: %w", err)
	}
	return nil
}

func convertToGiftModel(gift models.Gift) GiftEntity {

	giftItem := GiftEntity{
		Id:              gift.Id,
		ListId:          gift.ListId,
		CreatedAt:       gift.CreatedAt,
		Title:           gift.Title,
		Desription:      gift.Desription,
		PlaceOfPurchase: gift.PlaceOfPurchase,
		ImageFileName:   gift.ImageFileName,
		Url:             gift.Url,
		Price:           gift.Price,
		Rank:            gift.Rank,
	}
	return giftItem
}

func convertToGift(gift GiftEntity) models.Gift {
	return models.Gift{
		Id:              gift.Id,
		ListId:          gift.ListId,
		CreatedAt:       gift.CreatedAt,
		Title:           gift.Title,
		Desription:      gift.Desription,
		PlaceOfPurchase: gift.PlaceOfPurchase,
		ImageFileName:   gift.ImageFileName,
		Url:             gift.Url,
		Price:           gift.Price,
		Rank:            gift.Rank,
	}
}
