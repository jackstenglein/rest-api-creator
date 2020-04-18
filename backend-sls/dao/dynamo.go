package dao

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	apierrors "github.com/rest_api_creator/backend-sls/errors"
)

type DynamoService interface {
	GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	UpdateItem(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
}

type DynamoStore struct {
	service DynamoService
}

func DefaultDynamoStore() *DynamoStore {
	return &DynamoStore{dynamodb.New(session.New())}
}

func NewDynamoStore(service DynamoService) *DynamoStore {
	return &DynamoStore{service}
}

func (store *DynamoStore) CreateUser(email string, password string, token string) error {
	input := &dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(email)"),
		Item: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
			"Password": {
				S: aws.String(password),
			},
			"SessionToken": {
				S: aws.String(token),
			},
		},
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	_, err := store.service.PutItem(input)
	if err != nil {
		var aerr awserr.Error
		if errors.As(err, &aerr) {
			if aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
				return apierrors.NewUserError("Email already in use")
			}
		}
		return errors.Wrap(err, "Failed DynamoDB PutItem call")
	}
	return nil
}

func (store *DynamoStore) GetUser(email string) (User, error) {
	user := User{}

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	result, err := store.service.GetItem(input)
	if err != nil {
		return user, errors.Wrap(err, "Failed DynamoDB GetItem call")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return user, errors.Wrap(err, "Failed to unmarshal GetItem result")
	}

	return user, nil
}

func (store *DynamoStore) GetProject(email string, projectId string) (*Project, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
		ProjectionExpression: aws.String(fmt.Sprintf("Project.%s", projectId)),
		TableName:            aws.String(os.Getenv("TABLE_NAME")),
	}
	result, err := store.service.GetItem(input)
	if err != nil {
		return nil, errors.Wrap(err, "Failed DynamoDB GetItem call")
	}

	project := Project{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &project)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal GetItem result")
	}
	return &project, nil
}

func (store *DynamoStore) UpdateUserToken(email string, token string) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(token),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
		TableName:        aws.String(os.Getenv("TABLE_NAME")),
		UpdateExpression: aws.String("SET SessionToken=:t"),
	}

	_, err := store.service.UpdateItem(input)
	return errors.Wrap(err, "Failed DynamoDB UpdateItem call")
}
