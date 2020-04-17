package dao

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rest_api_creator/backend-sls/errors"
	"os"
)

type DynamoStore struct {
	service *dynamodb.DynamoDB
}

func NewDynamoStore() *DynamoStore {
	return &DynamoStore{dynamodb.New(session.New())}
}

func (store *DynamoStore) CreateUser(email string, password string, token string) errors.ApiError {
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
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return errors.NewUserError("Email already in use")
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return errors.NewServerError(err.Error())
	}
	return nil
}

func (store *DynamoStore) GetUser(email string) (User, errors.ApiError) {
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
		fmt.Println(err)
		return user, errors.NewServerError(err.Error())
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		fmt.Println(err)
		return user, errors.NewServerError(err.Error())
	}

	return user, nil
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
	return err
}
