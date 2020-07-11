package dao

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
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
					"defaultProject": {
						M: map[string]*dynamodb.AttributeValue{
							"Id":          {S: aws.String("defaultProject")},
							"Name":        {S: aws.String("Default Project")},
							"Description": {S: aws.String(defaultProjectDesc)},
							"Objects":     {M: map[string]*dynamodb.AttributeValue{}},
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

func (dynamo) DeleteObject(email string, projectID string, objectID string) error {
	expression := "REMOVE Projects.#pid.Objects.#oid"
	attributeNames := map[string]*string{
		"#pid": aws.String(projectID),
		"#oid": aws.String(objectID),
	}
	return Dynamo.updateUser(email, expression, attributeNames, nil)
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

// GetUser returns the entire user object associated with the given email. All projects are included in
// their entirety. If an error occurs, the returned user will be nil.
func (dynamo) GetUser(email string) (*User, error) {
	expression := "Email, Projects"
	return Dynamo.getUser(email, expression, nil)
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

	var expressionAttributeValues map[string]*dynamodb.AttributeValue
	if len(items) > 0 {
		expressionAttributeValues = make(map[string]*dynamodb.AttributeValue)
		for key, item := range items {
			itemAV, err := dynamodbattribute.Marshal(item)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to marshal item %s", key))
			}
			expressionAttributeValues[key] = itemAV
		}
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

func (dynamo) UpdateDeployment(email string, projectID string, instanceID string, instanceURL string) error {
	expression := "SET Projects.#pid.InstanceId = :id, Projects.#pid.DeployUrl = :url"
	attributeNames := map[string]*string{
		"#pid": aws.String(projectID),
	}
	items := map[string]interface{}{
		":id":  instanceID,
		":url": instanceURL,
	}
	return Dynamo.updateUser(email, expression, attributeNames, items)
}

// UpdateObject either creates or replaces the given object within the given project. If an error occurs,
// it is returned.
func (dynamo) UpdateObject(email string, projectID string, object *Object, originalID string) error {
	expression := ""
	attributeNames := map[string]*string{
		"#pid": aws.String(projectID),
	}
	items := map[string]interface{}{
		":obj": object,
	}

	if originalID == object.ID || originalID == "" {
		// We are updating an existing object or creating a new object
		expression = "SET Projects.#pid.Objects.#oid = :obj"
		attributeNames["#oid"] = aws.String(object.ID)
	} else {
		// We are changing the ID of an existing object and need to delete the old ID
		expression = "REMOVE Projects.#pid.Objects.#oid1 SET Projects.#pid.Objects.#oid2 = :obj"
		attributeNames["#oid1"] = aws.String(originalID)
		attributeNames["#oid2"] = aws.String(object.ID)
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
