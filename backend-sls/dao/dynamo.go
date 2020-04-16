package dao

import (
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rest_api_creator/backend-sls/errors"
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
			"email": {
				S: aws.String(email),
			},
			"password": {
				S: aws.String(password),
			},
			"token": {
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
