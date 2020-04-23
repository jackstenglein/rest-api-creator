package dao

import (
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
		return nil, errors.NewServer("Incorrect UpdateItemInput to mock")
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
		name:       "ServiceError",
		email:      "email",
		projectID:  "projectID",
		mockInput:  getProjectMockInput("email", "Projects.projectID"),
		mockOutput: getProjectMockOutput(Project{ID: "projectID", Name: "default"}),
		mockErr:    errors.NewServer("DynamoDB error"),
		wantErr:    errors.Wrap(errors.NewServer("DynamoDB error"), "Failed DynamoDB GetItem call"),
	},
	{
		name:        "SuccessfulInvocation",
		email:       "email",
		projectID:   "projectID",
		mockInput:   getProjectMockInput("email", "Projects.projectID"),
		mockOutput:  getProjectMockOutput(Project{ID: "projectID", Name: "default"}),
		wantProject: &Project{ID: "projectID", Name: "default"},
	},
	{
		name:       "NonexistentProject",
		email:      "email",
		projectID:  "projectID",
		mockInput:  getProjectMockInput("email", "Projects.projectID"),
		mockOutput: &dynamodb.GetItemOutput{},
		wantErr:    errors.NewClient("Project `projectID` not found"),
	},
}

func getProjectMockInput(email string, projectionExpression string) *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
		ProjectionExpression: aws.String(projectionExpression),
		TableName:            aws.String(os.Getenv("TABLE_NAME")),
	}
}

func getProjectMockOutput(project Project) *dynamodb.GetItemOutput {
	av, _ := dynamodbattribute.MarshalMap(project)
	return &dynamodb.GetItemOutput{Item: av}
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

// ------------- GetUser Tests ------------------

var getUserTests = []struct {
	name       string
	email      string
	mockInput  *dynamodb.GetItemInput
	mockOutput *dynamodb.GetItemOutput
	mockErr    error
	wantUser   *User
	wantErr    error
}{
	{
		name:      "ServiceError",
		email:     "email",
		mockInput: getUserMockInput("email"),
		mockErr:   errors.NewServer("DynamoDB failed"),
		wantErr:   errors.Wrap(errors.NewServer("DynamoDB failed"), "Failed DynamoDB GetItem call"),
	},
	{
		name:       "NonexistentUser",
		email:      "email",
		mockInput:  getUserMockInput("email"),
		mockOutput: &dynamodb.GetItemOutput{},
		wantErr:    errors.NewClient("Email `email` not found"),
	},
	{
		name:       "SuccessfulInvocation",
		email:      "email",
		mockInput:  getUserMockInput("email"),
		mockOutput: getUserMockOutput(&User{Email: "email", Password: "password", Token: "token"}),
		wantUser:   &User{Email: "email", Password: "password", Token: "token"},
	},
}

func getUserMockInput(email string) *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}
}

func getUserMockOutput(user *User) *dynamodb.GetItemOutput {
	av, _ := dynamodbattribute.MarshalMap(user)
	return &dynamodb.GetItemOutput{Item: av}
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
			gotUser, gotErr := Dynamo.GetUser(test.email)

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

// ------------- UpdateUserToken Tests -------------------

var updateUserTokenTests = []struct {
	name      string
	email     string
	token     string
	mockInput *dynamodb.UpdateItemInput
	mockErr   error
	wantErr   error
}{
	{
		name:      "ServiceError",
		email:     "email",
		token:     "token",
		mockInput: updateUserTokenMockInput("email", "token"),
		mockErr:   errors.NewServer("DynamoDB failure"),
		wantErr:   errors.Wrap(errors.NewServer("DynamoDB failure"), "Failed DynamoDB UpdateItem call"),
	},
	{
		name:      "SuccessfulInvocation",
		email:     "email",
		token:     "token",
		mockInput: updateUserTokenMockInput("email", "token"),
	},
}

func updateUserTokenMockInput(email string, token string) *dynamodb.UpdateItemInput {
	return &dynamodb.UpdateItemInput{
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

// -------------- UpdateItem Tests -----------------

var updateTokenInput = &dynamodb.UpdateItemInput{
	ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
		":item": {
			S: aws.String("token"),
		},
	},
	Key: map[string]*dynamodb.AttributeValue{
		"Email": {
			S: aws.String("test@example.com"),
		},
	},
	TableName:        aws.String(os.Getenv("TABLE_NAME")),
	UpdateExpression: aws.String("SET path = :item"),
}

var updateObjectInput = &dynamodb.UpdateItemInput{
	ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
		":item": {
			M: map[string]*dynamodb.AttributeValue{
				"Name": {
					S: aws.String("Object Name"),
				},
				"Description": {
					S: aws.String("Test description"),
				},
				"Id": {
					S: aws.String("objectId"),
				},
			},
		},
	},
	Key: map[string]*dynamodb.AttributeValue{
		"Email": {
			S: aws.String("test@example.com"),
		},
	},
	TableName:        aws.String(os.Getenv("TABLE_NAME")),
	UpdateExpression: aws.String("SET path = :item"),
}

var updateItemTests = []struct {
	name string

	// Input
	item interface{}

	// Mock data
	mockInput *dynamodb.UpdateItemInput
	mockErr   error

	// Expected output
	wantErr error
}{
	{
		name:      "ServiceError",
		item:      "token",
		mockInput: updateTokenInput,
		mockErr:   errors.NewServer("DynamoDB failure"),
		wantErr:   errors.Wrap(errors.NewServer("DynamoDB failure"), "Failed DynamoDB UpdateItem call"),
	},
	{
		name:      "AuthTokenItem",
		item:      "token",
		mockInput: updateTokenInput,
		mockErr:   nil,
		wantErr:   nil,
	},
	{
		name:      "ObjectItem",
		item:      &Object{ID: "objectId", Name: "Object Name", Description: "Test description"},
		mockInput: updateObjectInput,
		mockErr:   nil,
		wantErr:   nil,
	},
}

func TestUpdateItem(t *testing.T) {
	for _, test := range updateItemTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			updateSvc = updateItemMock(test.mockInput, nil, test.mockErr)
			defer func() {
				updateSvc = defaultSvc
			}()

			// Execute
			gotErr := Dynamo.UpdateItem("test@example.com", "path", test.item)

			// Verify
			if !errors.Equal(gotErr, test.wantErr) {
				t.Errorf("Got error '%s'; want '%s'", gotErr, test.wantErr)
			}
		})
	}
}
