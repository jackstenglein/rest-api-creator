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
			// TODO: Remove this and dynamically create projects
			"Projects": {
				M: map[string]*dynamodb.AttributeValue{
					"default": {
						M: map[string]*dynamodb.AttributeValue{
							"Id":      {S: aws.String("default")},
							"Name":    {S: aws.String("Default Project")},
							"Objects": {M: map[string]*dynamodb.AttributeValue{}},
						},
					},
				},
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

func (dynamo) getUser(email string, expression string, attributeNames map[string]*string) (*User, error) {
	input := &dynamodb.GetItemInput{
		ExpressionAttributeNames: attributeNames,
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
		ProjectionExpression: aws.String(expression),
		TableName:            aws.String(os.Getenv("TABLE_NAME")),
	}

	result, err := getSvc.GetItem(input)
	if err != nil {
		return nil, errors.Wrap(err, "Failed DynamoDB GetItem call")
	}
	if result.Item == nil {
		return nil, errors.NewClient(fmt.Sprintf("Email '%s' not found", email))
	}

	user := User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal GetItem result")
	}
	return &user, nil
}

// GetUserInfo returns the basic User info associated with the given email. Projects are not included.
// If the email does not exist, the returned user will be nil and the returned error will be a new client
// error.
func (dynamo) GetUserInfo(email string) (*User, error) {
	expression := "Email, Password, SessionToken"
	return Dynamo.getUser(email, expression, nil)
}

// GetProject returns the Project object associated with the given email and projectID.
// If the projectID does not exist for the specified email, the returned project will be
// nil and the returned error will be a new client error.
func (dynamo) GetProject(email string, projectID string) (*Project, error) {
	expression := "Projects.#pid"
	attributeNames := map[string]*string{
		"#pid": aws.String(projectID),
	}

	user, err := Dynamo.getUser(email, expression, attributeNames)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to get user with email '%s'", email))
	}

	project := user.Projects[projectID]
	if project == nil {
		return nil, errors.NewClient(fmt.Sprintf("Project '%s' not found", projectID))
	}
	return project, nil
}

// updateUser updates the properties of the user given in expression with the given items. If expression is not a valid
// property path for the given user, a client error is returned. If something else goes wrong, a server
// error is returned.
//
// TODO: Return client error if path does not exist
func (dynamo) updateUser(email string, expression string, attributeNames map[string]*string, items map[string]interface{}) error {

	expressionAttributeValues := make(map[string]*dynamodb.AttributeValue)
	for key, item := range items {
		itemAV, err := dynamodbattribute.Marshal(item)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to marshal item %s", key))
		}
		expressionAttributeValues[key] = itemAV
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  attributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
		TableName:        aws.String(os.Getenv("TABLE_NAME")),
		UpdateExpression: aws.String(expression),
	}

	_, err := updateSvc.UpdateItem(input)
	return errors.Wrap(err, "Failed DynamoDB UpdateItem call")
}

// UpdateObject either creates or replaces the given object within the given project. If an error occurs,
// it is returned.
func (dynamo) UpdateObject(email string, projectID string, object *Object) error {
	expression := "SET Projects.#pid.Objects.#oid = :obj"
	attributeNames := map[string]*string{
		"#pid": aws.String(projectID),
		"#oid": aws.String(object.ID),
	}
	items := map[string]interface{}{
		":obj": object,
	}
	return Dynamo.updateUser(email, expression, attributeNames, items)
}

// UpdateUserToken sets the auth token on the User object associated with the given email
// in the database.
func (dynamo) UpdateUserToken(email string, token string) error {
	expression := "SET SessionToken = :tok"
	items := map[string]interface{}{
		":tok": token,
	}
	return Dynamo.updateUser(email, expression, nil, items)
}
