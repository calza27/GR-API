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

type ListRepository interface {
	CreateList(list models.List) error
	UpdateList(list models.List) error
	GetListById(listId string) (models.List, error)
	GetListsByUserId(userId string) ([]models.List, error)
	DeleteList(listId string) error
	persistList(item map[string]types.AttributeValue) error
}

type DynamoDbListRepository struct {
	db             *dynamodb.Client
	tablename      string
	userIdIndex    string
	sharingIdIndex string
}

func NewListRepository(tableName, userIdIndex, sharingIdIndex string) (ListRepository, error) {
	db, err := awsclient.GetDynamodbClient()
	if err != nil {
		return nil, fmt.Errorf("Error when initialising connection to DDB: %w", err)
	}

	return &DynamoDbListRepository{
		db:             db,
		tablename:      tableName,
		userIdIndex:    userIdIndex,
		sharingIdIndex: sharingIdIndex,
	}, nil
}

type ListEntity struct {
	Id            string `dynamodbav:"id"`
	UserId        string `dynamodbav:"userId"`
	CreatedAt     string `dynamodbav:"createdAt"`
	Name          string `dynamodbav:"listName,omitempty"`
	SharingId     string `dynamodbav:"sharingId,omitempty"`
	ImageFileName string `dynamodbav:"imageFileName,omitempty"`
}

func (r *DynamoDbListRepository) CreateList(list models.List) error {
	listItem := convertToListModel(list)
	now := time.Now()
	listItem.CreatedAt = utils.DateTimeToString(now)
	listItem.Id = utils.GenerateUUID()
	sharingId, err := r.uniqueSharingId()
	if err != nil {
		return fmt.Errorf("Error when trying to generate unique sharing id: %w", err)
	}
	listItem.SharingId = sharingId
	item, err := attributevalue.MarshalMap(listItem)
	if err != nil {
		return fmt.Errorf("error when trying to convert list data to dynamodbattribute: %w", err)
	}
	return r.persistList(item)
}

func (r *DynamoDbListRepository) UpdateList(list models.List) error {
	listItem := convertToListModel(list)
	item, err := attributevalue.MarshalMap(listItem)
	if err != nil {
		return fmt.Errorf("error when trying to convert list data to dynamodbattribute: %w", err)
	}
	return r.persistList(item)
}

func (r *DynamoDbListRepository) persistList(item map[string]types.AttributeValue) error {
	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(r.tablename),
	}
	if _, err := r.db.PutItem(context.Background(), params); err != nil {
		return fmt.Errorf("Error when trying to persist: %w", err)
	}
	return nil
}

// GetListById retrieves a list by its id or sharing id
func (r *DynamoDbListRepository) GetListById(listId string) (models.List, error) {
	params := &dynamodb.QueryInput{
		TableName:              aws.String(r.tablename),
		KeyConditionExpression: aws.String("id = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: listId},
		},
	}
	result, err := r.db.Query(context.Background(), params)
	if err != nil {
		return models.List{}, fmt.Errorf("Error when trying to get list by id: %w", err)
	}
	if len(result.Items) == 1 {
		var list ListEntity
		if err := attributevalue.UnmarshalMap(result.Items[0], &list); err != nil {
			return models.List{}, fmt.Errorf("Error when trying to unmarshal list data: %w", err)
		}
		return convertToList(list), nil
	}
	params = &dynamodb.QueryInput{
		TableName:              aws.String(r.tablename),
		IndexName:              aws.String(r.sharingIdIndex),
		KeyConditionExpression: aws.String("sharingId = :sharingId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sharingId": &types.AttributeValueMemberS{Value: listId},
		},
	}
	result, err = r.db.Query(context.Background(), params)
	if err != nil {
		return models.List{}, fmt.Errorf("Error when trying to get list by sharing id: %w", err)
	}
	if len(result.Items) == 1 {
		var list ListEntity
		if err := attributevalue.UnmarshalMap(result.Items[0], &list); err != nil {
			return models.List{}, fmt.Errorf("Error when trying to unmarshal list data: %w", err)
		}
		return convertToList(list), nil
	}
	return models.List{}, fmt.Errorf("No list found with id or sharingId matching: %s", listId)
}

func (r *DynamoDbListRepository) GetListsByUserId(userId string) ([]models.List, error) {
	params := &dynamodb.QueryInput{
		TableName:              aws.String(r.tablename),
		IndexName:              aws.String(r.userIdIndex),
		KeyConditionExpression: aws.String("userId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userId},
		},
	}
	result, err := r.db.Query(context.Background(), params)
	if err != nil {
		return []models.List{}, fmt.Errorf("Error when trying to get lists by user id: %w", err)
	}
	var lists []models.List
	for _, item := range result.Items {
		var listItem ListEntity
		if err := attributevalue.UnmarshalMap(item, &listItem); err != nil {
			return []models.List{}, fmt.Errorf("Error when trying to unmarshal list data: %w", err)
		}
		list := convertToList(listItem)
		lists = append(lists, list)
	}
	return lists, nil
}

func (r *DynamoDbListRepository) DeleteList(listId string) error {
	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tablename),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: listId},
		},
	}
	_, err := r.db.DeleteItem(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error when trying to delete list: %w", err)
	}
	return nil
}

func convertToListModel(list models.List) ListEntity {

	listItem := ListEntity{
		Id:            list.Id,
		UserId:        list.UserId,
		CreatedAt:     list.CreatedAt,
		Name:          list.Name,
		SharingId:     list.SharingId,
		ImageFileName: list.ImageFileName,
	}
	return listItem
}

func convertToList(list ListEntity) models.List {
	return models.List{
		Id:            list.Id,
		UserId:        list.UserId,
		CreatedAt:     list.CreatedAt,
		Name:          list.Name,
		SharingId:     list.SharingId,
		ImageFileName: list.ImageFileName,
	}
}

func (r *DynamoDbListRepository) uniqueSharingId() (string, error) {
	limit := 10
	for {
		sharingId := generateSharingId()
		params := &dynamodb.QueryInput{
			TableName:              aws.String(r.tablename),
			IndexName:              aws.String(r.sharingIdIndex),
			KeyConditionExpression: aws.String("sharingId = :sharingId"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":sharingId": &types.AttributeValueMemberS{Value: sharingId},
			},
		}
		result, err := r.db.Query(context.Background(), params)
		if err != nil {
			return "", fmt.Errorf("Error when trying to validate potential sharing id: %w", err)
		}
		if len(result.Items) == 0 {
			return sharingId, nil
		}
		limit -= 1
		if limit == 0 {
			return "", fmt.Errorf("Could not generate unique sharing id")
		}
	}
}

func generateSharingId() string {
	length := 12
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	sharingId := make([]byte, length)
	for i := range sharingId {
		sharingId[i] = chars[utils.RandomInt(0, len(chars))]
	}
	return string(sharingId)
}
