package dao

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

// updater wraps the UpdateItem method in order to perform dependency injection
// in the dynamo tests.
type updater interface {
	UpdateItem(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
}

// getter wraps the GetItem method in order to perform dependency injection
// in the dynamo tests.
type getter interface {
	GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
}

// putter wraps the PutItem method in order to perform dependency injection
// in the dynamo tests.
type putter interface {
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
}

// defaultSvc actually queries DynamoDB. The other Svc variables exist only for
// dependency injection and should only be changed inside a test.
var defaultSvc = dynamodb.New(session.New())
var getSvc getter = defaultSvc
var putSvc putter = defaultSvc
var updateSvc updater = defaultSvc

// dynamo is an empty struct that acts as a collection of database methods.
type dynamo struct{}

// Dynamo provides a high-level interface to perform database queries against AWS DynamoDB.
var Dynamo = dynamo{}

// CreateUser adds a User object to the database with the given email, password and session token.
// If the email already exists in the database, CreateUser makes no changes to the databse and returns
// a client error.
func (dynamo) CreateUser(email string, password string, token string) error {
	input := &dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(Email)"),
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

	_, err := putSvc.PutItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
				return errors.NewClient("Email already in use")
			}
		}
		return errors.Wrap(err, "Failed DynamoDB PutItem call")
	}
	return nil
}

// GetUser returns the User object associated with the given email. If the email
// does not exist, the returned user will be nil and the returned error will be
// a new client error.
func (dynamo) GetUser(email string) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	result, err := getSvc.GetItem(input)
	if err != nil {
		return nil, errors.Wrap(err, "Failed DynamoDB GetItem call")
	}
	if result.Item == nil {
		return nil, errors.NewClient(fmt.Sprintf("Email `%s` not found", email))
	}

	user := User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal GetItem result")
	}

	return &user, nil
}

// GetProject returns the Project object associated with the given email and projectID.
// If the projectID does not exist for the specified email, the returned project will be
// nil and the returned error will be a new client error.
func (dynamo) GetProject(email string, projectID string) (*Project, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
		ProjectionExpression: aws.String(fmt.Sprintf("Projects.%s", projectID)),
		TableName:            aws.String(os.Getenv("TABLE_NAME")),
	}
	result, err := getSvc.GetItem(input)
	if err != nil {
		return nil, errors.Wrap(err, "Failed DynamoDB GetItem call")
	}

	if result.Item == nil {
		return nil, errors.NewClient(fmt.Sprintf("Project `%s` not found", projectID))
	}

	project := Project{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &project)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal GetItem result")
	}
	return &project, nil
}

// UpdateUserToken sets the auth token on the User object associated with the given email
// in the database.
// TODO: Check that this doesn't create a new user obejct if email doesn't exist in DB.
func (dynamo) UpdateUserToken(email string, token string) error {
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

	_, err := updateSvc.UpdateItem(input)
	return errors.Wrap(err, "Failed DynamoDB UpdateItem call")
}
