package dao

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

// -------------- GetItem Mock -----------------

type getItemFunc func(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)

func (f getItemFunc) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return f(input)
}

func getItemMock(mockInput *dynamodb.GetItemInput, mockOutput *dynamodb.GetItemOutput, mockErr error) getItemFunc {
	return func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
		if reflect.DeepEqual(input, mockInput) {
			return mockOutput, mockErr
		}
		return nil, errors.NewServer("Incorrect GetItemInput to mock")
	}
}

// -------------- PutItem Mock -----------------

type putItemFunc func(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)

func (f putItemFunc) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return f(input)
}

func putItemMock(mockInput *dynamodb.PutItemInput, mockOutput *dynamodb.PutItemOutput, mockErr error) putItemFunc {
	return func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
		if reflect.DeepEqual(input, mockInput) {
			return mockOutput, mockErr
		}
		return nil, errors.NewServer("Incorrect PutItemInput to mock")
	}
}

// -------------- UpdateItem Mock -----------------

type updateItemFunc func(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)

func (f updateItemFunc) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	return f(input)
}

func updateItemMock(mockInput *dynamodb.UpdateItemInput, mockOutput *dynamodb.UpdateItemOutput, mockErr error) updateItemFunc {
	return func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
		if reflect.DeepEqual(input, mockInput) {
			return mockOutput, mockErr
		}
		fmt.Println("Actual input:", input)
		fmt.Println("Expected input:", mockInput)
		return nil, errors.NewServer("Incorrect UpdateItemInput to mock")
	}
}

// ------------- GetUser Tests ------------------

var getUserTests = []struct {
	name string

	// Input
	email          string
	expression     string
	attributeNames map[string]*string

	// Mock data
	mockInput  *dynamodb.GetItemInput
	mockOutput *dynamodb.GetItemOutput
	mockErr    error

	// Expected output
	wantUser *User
	wantErr  error
}{
	{
		name:       "ServiceError",
		email:      "email",
		expression: "Email, Password, SessionToken",
		mockInput: &dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"Email": {
					S: aws.String("email"),
				},
			},
			ProjectionExpression: aws.String("Email, Password, SessionToken"),
			TableName:            aws.String(os.Getenv("TABLE_NAME")),
		},
		mockErr: errors.NewServer("DynamoDB failed"),
		wantErr: errors.Wrap(errors.NewServer("DynamoDB failed"), "Failed DynamoDB GetItem call"),
	},
	{
		name:       "NonexistentUser",
		email:      "email",
		expression: "Projects.#pid",
		attributeNames: map[string]*string{
			"#pid": aws.String("projectID"),
		},
		mockInput: &dynamodb.GetItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#pid": aws.String("projectID"),
			},
			Key: map[string]*dynamodb.AttributeValue{
				"Email": {
					S: aws.String("email"),
				},
			},
			ProjectionExpression: aws.String("Projects.#pid"),
			TableName:            aws.String(os.Getenv("TABLE_NAME")),
		},
		mockOutput: &dynamodb.GetItemOutput{},
		wantErr:    errors.NewClient("Email 'email' not found"),
	},
	{
		name:       "SuccessfulInvocation",
		email:      "email",
		expression: "Email, Password, SessionToken",
		mockInput: &dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"Email": {
					S: aws.String("email"),
				},
			},
			ProjectionExpression: aws.String("Email, Password, SessionToken"),
			TableName:            aws.String(os.Getenv("TABLE_NAME")),
		},
		mockOutput: &dynamodb.GetItemOutput{
			Item: map[string]*dynamodb.AttributeValue{
				"Email":        {S: aws.String("email")},
				"Password":     {S: aws.String("password")},
				"SessionToken": {S: aws.String("sessionToken")},
			},
		},
		wantUser: &User{Email: "email", Password: "password", Token: "sessionToken"},
	},
}

func TestGetUser(t *testing.T) {
	for _, test := range getUserTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			getSvc = getItemMock(test.mockInput, test.mockOutput, test.mockErr)
			defer func() {
				getSvc = defaultSvc
			}()

			// Execute
			gotUser, gotErr := Dynamo.getUser(test.email, test.expression, test.attributeNames)

			// Verify
			if !reflect.DeepEqual(gotUser, test.wantUser) {
				t.Errorf("Got user %v; want %v", gotUser, test.wantUser)
			}
			if !errors.Equal(gotErr, test.wantErr) {
				t.Errorf("Got error '%s'; want '%s'", gotErr, test.wantErr)
			}
		})
	}
}

// ---------------- GetProject Tests ----------------

var getProjectTests = []struct {
	name        string
	email       string
	projectID   string
	mockInput   *dynamodb.GetItemInput
	mockOutput  *dynamodb.GetItemOutput
	mockErr     error
	wantProject *Project
	wantErr     error
}{
	{
		name:      "ServiceError",
		email:     "email",
		projectID: "projectID",
		mockInput: &dynamodb.GetItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#pid": aws.String("projectID"),
			},
			Key: map[string]*dynamodb.AttributeValue{
				"Email": {
					S: aws.String("email"),
				},
			},
			ProjectionExpression: aws.String("Projects.#pid"),
			TableName:            aws.String(os.Getenv("TABLE_NAME")),
		},
		mockErr: errors.NewServer("DynamoDB error"),
		wantErr: errors.Wrap(errors.Wrap(errors.NewServer("DynamoDB error"), "Failed DynamoDB GetItem call"), "Failed to get user with email 'email'"),
	},
	{
		name:      "NonexistentProject",
		email:     "email",
		projectID: "projectID",
		mockInput: &dynamodb.GetItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#pid": aws.String("projectID"),
			},
			Key: map[string]*dynamodb.AttributeValue{
				"Email": {
					S: aws.String("email"),
				},
			},
			ProjectionExpression: aws.String("Projects.#pid"),
			TableName:            aws.String(os.Getenv("TABLE_NAME")),
		},
		mockOutput: &dynamodb.GetItemOutput{
			Item: map[string]*dynamodb.AttributeValue{
				"Projects": {M: map[string]*dynamodb.AttributeValue{}},
			},
		},
		wantErr: errors.NewClient("Project 'projectID' not found"),
	},
	{
		name:      "SuccessfulInvocation",
		email:     "email",
		projectID: "projectID",
		mockInput: &dynamodb.GetItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#pid": aws.String("projectID"),
			},
			Key: map[string]*dynamodb.AttributeValue{
				"Email": {
					S: aws.String("email"),
				},
			},
			ProjectionExpression: aws.String("Projects.#pid"),
			TableName:            aws.String(os.Getenv("TABLE_NAME")),
		},
		mockOutput: &dynamodb.GetItemOutput{
			Item: map[string]*dynamodb.AttributeValue{
				"Projects": {
					M: map[string]*dynamodb.AttributeValue{
						"projectID": {
							M: map[string]*dynamodb.AttributeValue{
								"Id":   {S: aws.String("projectID")},
								"Name": {S: aws.String("default")},
							},
						},
					},
				},
			},
		},
		wantProject: &Project{ID: "projectID", Name: "default"},
	},
}

func TestGetProject(t *testing.T) {
	for _, test := range getProjectTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			getSvc = getItemMock(test.mockInput, test.mockOutput, test.mockErr)
			defer func() {
				getSvc = defaultSvc
			}()

			// Execute
			project, err := Dynamo.GetProject(test.email, test.projectID)

			// Verify
			if !reflect.DeepEqual(project, test.wantProject) {
				t.Errorf("Got project %v; want %v", project, test.wantProject)
			}
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got error '%s'; want '%s'", err, test.wantErr)
			}
		})
	}
}

// ----------- CreateUser Tests --------------

var createUserTests = []struct {
	name      string
	email     string
	password  string
	token     string
	mockInput *dynamodb.PutItemInput
	mockErr   error
	wantErr   error
}{
	{
		name:      "ServiceError",
		email:     "email",
		password:  "password",
		token:     "token",
		mockInput: createUserMockInput("email", "password", "token"),
		mockErr:   errors.NewServer("DynamoDB failure"),
		wantErr:   errors.Wrap(errors.NewServer("DynamoDB failure"), "Failed DynamoDB PutItem call"),
	},
	{
		name:      "EmailAlreadyExists",
		email:     "email",
		password:  "password",
		token:     "token",
		mockInput: createUserMockInput("email", "password", "token"),
		mockErr:   awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "Email exists", nil),
		wantErr:   errors.NewClient("Email already in use"),
	},
	{
		name:      "SuccessfulInvocation",
		email:     "email",
		password:  "password",
		token:     "token",
		mockInput: createUserMockInput("email", "password", "token"),
	},
}

func createUserMockInput(email string, password string, token string) *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
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
}

func TestCreateUser(t *testing.T) {
	for _, test := range createUserTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			putSvc = putItemMock(test.mockInput, nil, test.mockErr)
			defer func() {
				putSvc = defaultSvc
			}()

			// Execute
			gotErr := Dynamo.CreateUser(test.email, test.password, test.token)

			// Verify
			if !errors.Equal(gotErr, test.wantErr) {
				t.Errorf("Got error '%s'; want '%s'", gotErr, test.wantErr)
			}
		})
	}
}

// -------------- Update Tests -----------------

var updateUserTests = []struct {
	name string

	// Input
	email          string
	expression     string
	attributeNames map[string]*string
	items          map[string]interface{}

	// Mock data
	mockInput *dynamodb.UpdateItemInput
	mockErr   error

	// Expected output
	wantErr error
}{
	{
		name:       "ServiceError",
		email:      "error@test.com",
		expression: "TEST update expression",
		items: map[string]interface{}{
			"key": "value",
		},
		mockInput: &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				"key": {
					S: aws.String("value"),
				},
			},
			Key: map[string]*dynamodb.AttributeValue{
				"Email": {
					S: aws.String("error@test.com"),
				},
			},
			TableName:        aws.String(os.Getenv("TABLE_NAME")),
			UpdateExpression: aws.String("TEST update expression"),
		},
		mockErr: errors.NewServer("DynamoDB failure"),
		wantErr: errors.Wrap(errors.NewServer("DynamoDB failure"), "Failed DynamoDB UpdateItem call"),
	},
	{
		name:       "SuccessfulInvocation",
		email:      "success@test.com",
		expression: "TEST update expression 2",
		attributeNames: map[string]*string{
			"#pid": aws.String("attributeName"),
		},
		items: map[string]interface{}{
			"key2": "value2",
		},
		mockInput: &dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#pid": aws.String("attributeName"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				"key2": {
					S: aws.String("value2"),
				},
			},
			Key: map[string]*dynamodb.AttributeValue{
				"Email": {
					S: aws.String("success@test.com"),
				},
			},
			TableName:        aws.String(os.Getenv("TABLE_NAME")),
			UpdateExpression: aws.String("TEST update expression 2"),
		},
	},
}

func TestUpdateItem(t *testing.T) {
	for _, test := range updateUserTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			updateSvc = updateItemMock(test.mockInput, nil, test.mockErr)
			defer func() {
				updateSvc = defaultSvc
			}()

			// Execute
			gotErr := Dynamo.updateUser(test.email, test.expression, test.attributeNames, test.items)

			// Verify
			if !errors.Equal(gotErr, test.wantErr) {
				t.Errorf("Got error '%s'; want '%s'", gotErr, test.wantErr)
			}
		})
	}
}

var updateObjectTests = []struct {
	name string

	// Input
	email     string
	projectID string
	object    *Object

	// Mock data
	mockInput *dynamodb.UpdateItemInput
	mockErr   error

	// Expected output
	wantErr error
}{
	{
		name:      "ServiceError",
		email:     "error@test.com",
		projectID: "projectID",
		object:    &Object{ID: "objectID", Name: "objectName", CodeName: "ObjectName", Description: "objectDesc"},
		mockInput: &dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#pid": aws.String("projectID"),
				"#oid": aws.String("objectID"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":obj": {
					M: map[string]*dynamodb.AttributeValue{
						"Id":          {S: aws.String("objectID")},
						"Name":        {S: aws.String("objectName")},
						"CodeName":    {S: aws.String("ObjectName")},
						"Description": {S: aws.String("objectDesc")},
					},
				},
			},
			Key: map[string]*dynamodb.AttributeValue{
				"Email": {
					S: aws.String("error@test.com"),
				},
			},
			TableName:        aws.String(os.Getenv("TABLE_NAME")),
			UpdateExpression: aws.String("SET Projects.#pid.Objects.#oid = :obj"),
		},
		mockErr: errors.NewServer("DynamoDB failure"),
		wantErr: errors.Wrap(errors.NewServer("DynamoDB failure"), "Failed DynamoDB UpdateItem call"),
	},
}

func TestUpdateObject(t *testing.T) {
	for _, test := range updateObjectTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			updateSvc = updateItemMock(test.mockInput, nil, test.mockErr)
			defer func() {
				updateSvc = defaultSvc
			}()

			// Execute
			gotErr := Dynamo.UpdateObject(test.email, test.projectID, test.object)

			// Verify
			if !errors.Equal(gotErr, test.wantErr) {
				t.Errorf("Got error '%s'; want '%s'", gotErr, test.wantErr)
			}
		})
	}
}

var updateUserTokenTests = []struct {
	name string

	// Input
	email string
	token string

	// Mock data
	mockInput *dynamodb.UpdateItemInput
	mockErr   error

	// Expected output
	wantErr error
}{
	{
		name:  "SuccessfulInvocation",
		email: "success@test.com",
		token: "tokenValue",
		mockInput: &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":tok": {S: aws.String("tokenValue")},
			},
			Key: map[string]*dynamodb.AttributeValue{
				"Email": {
					S: aws.String("success@test.com"),
				},
			},
			TableName:        aws.String(os.Getenv("TABLE_NAME")),
			UpdateExpression: aws.String("SET SessionToken = :tok"),
		},
	},
}

func TestUpdateUserToken(t *testing.T) {
	for _, test := range updateUserTokenTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			updateSvc = updateItemMock(test.mockInput, nil, test.mockErr)
			defer func() {
				updateSvc = defaultSvc
			}()

			// Execute
			gotErr := Dynamo.UpdateUserToken(test.email, test.token)

			// Verify
			if !errors.Equal(gotErr, test.wantErr) {
				t.Errorf("Got error '%s'; want '%s'", gotErr, test.wantErr)
			}
		})
	}
}
