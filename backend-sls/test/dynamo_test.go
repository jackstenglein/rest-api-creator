package test

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	gomock "github.com/golang/mock/gomock"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/mock"
)

const email = "test@example.com"
const projectId = "asdf"
const projectName = "default"

var project = dao.Project{Id: projectId, Name: projectName}

var getProjectTests = []struct {
	name                 string
	email                string
	projectId            string
	projectName          string
	projectionExpression string
	mockOutput           *dynamodb.GetItemOutput
	mockErr              error
	wantProject          *dao.Project
	wantErr              string
}{
	{
		name:                 "TestError",
		email:                email,
		projectId:            projectId,
		projectName:          projectName,
		projectionExpression: fmt.Sprintf("Projects.%s", projectId),
		mockOutput:           getMockOutput(project),
		mockErr:              errors.New("DynamoDB error"),
		wantErr:              "Failed DynamoDB GetItem call",
	},
	{
		name:                 "TestSuccess",
		email:                email,
		projectId:            projectId,
		projectName:          projectName,
		projectionExpression: fmt.Sprintf("Projects.%s", projectId),
		mockOutput:           getMockOutput(project),
		wantProject:          &project,
	},
	{
		name:                 "NonexistentProject",
		email:                email,
		projectId:            projectId,
		projectName:          projectName,
		projectionExpression: fmt.Sprintf("Projects.%s", projectId),
		mockOutput:           &dynamodb.GetItemOutput{},
		wantProject:          &dao.Project{},
	},
}

func getMockOutput(project dao.Project) *dynamodb.GetItemOutput {
	av, _ := dynamodbattribute.MarshalMap(project)
	return &dynamodb.GetItemOutput{Item: av}
}

func TestGetProject(t *testing.T) {
	for _, test := range getProjectTests {
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockService := mock.NewMockDynamoService(mockCtrl)
			dynamo := dao.NewDynamoStore(mockService)

			input := &dynamodb.GetItemInput{
				Key: map[string]*dynamodb.AttributeValue{
					"Email": {
						S: aws.String(test.email),
					},
				},
				ProjectionExpression: aws.String(test.projectionExpression),
				TableName:            aws.String(os.Getenv("TABLE_NAME")),
			}

			mockService.EXPECT().GetItem(input).Return(test.mockOutput, test.mockErr).Times(1)
			project, err := dynamo.GetProject(test.email, test.projectId)

			if !reflect.DeepEqual(project, test.wantProject) {
				t.Errorf("Got project %v; want %v", project, test.wantProject)
			}

			if err == nil {
				if test.wantErr != "" {
					t.Errorf("Got error nil; want '%s'", test.wantErr)
				}
			} else if !strings.Contains(err.Error(), test.wantErr) {
				t.Errorf("Got error %s; want '%s'", err, test.wantErr)
			}
		})
	}
}
