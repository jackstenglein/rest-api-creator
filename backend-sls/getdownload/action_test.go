package getdownload

import (
	"reflect"
	"testing"

	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type databaseMock struct {
	project *dao.Project
	err     error
}

func (mock *databaseMock) GetProject(email string, projectID string) (*dao.Project, error) {
	return mock.project, mock.err
}

func (mock *databaseMock) GetUserInfo(email string) (*dao.User, error) {
	return nil, nil
}

func verifyCookieMock(mockCookie string, mockDB auth.UserGetter, mockEmail string, mockErr error) cookieVerifier {
	return func(cookie string, db auth.UserGetter) (string, error) {
		if cookie != mockCookie || !reflect.DeepEqual(db, mockDB) {
			return "", errors.NewServer("Incorrect input to VerifyCookie mock")
		}
		return mockEmail, mockErr
	}
}

func createMock(mock1 string, mock2 string, mockErr error) func(string, string) error {
	return func(in1 string, in2 string) error {
		if in1 != mock1 || in2 != mock2 {
			return errors.NewServer("Incorrect input to mock")
		}
		return mockErr
	}
}

func generatorMock(mockProject *dao.Project, mockPath string, mockErr error) func(*dao.Project, string) error {
	return func(project *dao.Project, dirPath string) error {
		if !reflect.DeepEqual(project, mockProject) || dirPath != mockPath {
			return errors.NewServer("Incorrect input to generate mock")
		}
		return mockErr
	}
}

func presignMock(mockKey string, mockURL string, mockErr error) func(string) (string, error) {
	return func(key string) (string, error) {
		if mockKey != key {
			return "", errors.NewServer("Incorrect input to presign mock")
		}
		return mockURL, mockErr
	}
}

var generateCodeTests = []struct {
	name string

	// Input
	projectID string
	cookie    string

	// Mock data
	db         *databaseMock
	email      string
	verifyErr  error
	downloader func(string, string) error
	unzipper   func(string, string) error
	generator  func(*dao.Project, string) error
	zipper     func(string, string) error
	uploader   func(string, string) error
	presigner  func(string) (string, error)

	// Expected output
	wantURL string
	wantErr error
}{
	{
		name:    "EmptyCookie",
		wantErr: errors.NewClient("Not authenticated"),
	},
	{
		name:    "EmptyProjectID",
		cookie:  "TestCookie",
		wantErr: errors.NewClient("Parameter `pid` is required"),
	},
	{
		name:      "InvalidCookie",
		cookie:    "invalidcookie",
		projectID: "projectID",
		verifyErr: errors.NewClient("Not authenticated"),
		wantErr:   errors.Wrap(errors.NewClient("Not authenticated"), "Failed to verify cookie"),
	},
	{
		name:      "DatabaseError",
		cookie:    "validcookie",
		projectID: "projectID",
		email:     "test@example.com",
		db:        &databaseMock{nil, errors.NewServer("DynamoDB failure")},
		wantErr:   errors.Wrap(errors.NewServer("DynamoDB failure"), "Failed to get project from database"),
	},
	{
		name:       "DownloadError",
		cookie:     "validcookie",
		projectID:  "projectID",
		email:      "test@example.com",
		db:         &databaseMock{&dao.Project{Name: "Default Project", ID: "defaultProject"}, nil},
		downloader: createMock("/tmp/blank-sails.zip", "templates/sails.zip", errors.NewServer("S3 failure")),
		wantErr:    errors.Wrap(errors.NewServer("S3 failure"), "Failed to get project template from S3"),
	},
	{
		name:       "UnzipError",
		cookie:     "validcookie",
		projectID:  "projectID",
		email:      "test@example.com",
		db:         &databaseMock{&dao.Project{Name: "Default Project", ID: "defaultProject"}, nil},
		downloader: createMock("/tmp/blank-sails.zip", "templates/sails.zip", nil),
		unzipper:   createMock("/tmp/blank-sails.zip", "/tmp", errors.NewServer("Unzip failure")),
		wantErr:    errors.Wrap(errors.NewServer("Unzip failure"), "Failed to unzip project template"),
	},
	{
		name:       "GenerateError",
		cookie:     "validcookie",
		projectID:  "projectID",
		email:      "test@example.com",
		db:         &databaseMock{&dao.Project{Name: "Default Project", ID: "defaultProject"}, nil},
		downloader: createMock("/tmp/blank-sails.zip", "templates/sails.zip", nil),
		unzipper:   createMock("/tmp/blank-sails.zip", "/tmp", nil),
		generator:  generatorMock(&dao.Project{Name: "Default Project", ID: "defaultProject"}, "/tmp/blank-sails-project", errors.NewServer("Generate failure")),
		wantErr:    errors.Wrap(errors.NewServer("Generate failure"), "Failed to generate project code"),
	},
	{
		name:       "ZipError",
		cookie:     "validcookie",
		projectID:  "projectID",
		email:      "test@example.com",
		db:         &databaseMock{&dao.Project{Name: "Default Project", ID: "defaultProject"}, nil},
		downloader: createMock("/tmp/blank-sails.zip", "templates/sails.zip", nil),
		unzipper:   createMock("/tmp/blank-sails.zip", "/tmp", nil),
		generator:  generatorMock(&dao.Project{Name: "Default Project", ID: "defaultProject"}, "/tmp/blank-sails-project", nil),
		zipper:     createMock("/tmp/generated-project.zip", "/tmp/blank-sails-project", errors.NewServer("Zip failure")),
		wantErr:    errors.Wrap(errors.NewServer("Zip failure"), "Failed to zip generated code"),
	},
	{
		name:       "UploadError",
		cookie:     "validcookie",
		projectID:  "projectID",
		email:      "test@example.com",
		db:         &databaseMock{&dao.Project{Name: "Default Project", ID: "defaultProject"}, nil},
		downloader: createMock("/tmp/blank-sails.zip", "templates/sails.zip", nil),
		unzipper:   createMock("/tmp/blank-sails.zip", "/tmp", nil),
		generator:  generatorMock(&dao.Project{Name: "Default Project", ID: "defaultProject"}, "/tmp/blank-sails-project", nil),
		zipper:     createMock("/tmp/generated-project.zip", "/tmp/blank-sails-project", nil),
		uploader:   createMock("/tmp/generated-project.zip", "test@example.com/defaultProject.zip", errors.NewServer("Upload failure")),
		wantErr:    errors.Wrap(errors.NewServer("Upload failure"), "Failed to upload generated zip to S3"),
	},
	{
		name:       "PresignError",
		cookie:     "validcookie",
		projectID:  "projectID",
		email:      "test@example.com",
		db:         &databaseMock{&dao.Project{Name: "Default Project", ID: "defaultProject"}, nil},
		downloader: createMock("/tmp/blank-sails.zip", "templates/sails.zip", nil),
		unzipper:   createMock("/tmp/blank-sails.zip", "/tmp", nil),
		generator:  generatorMock(&dao.Project{Name: "Default Project", ID: "defaultProject"}, "/tmp/blank-sails-project", nil),
		zipper:     createMock("/tmp/generated-project.zip", "/tmp/blank-sails-project", nil),
		uploader:   createMock("/tmp/generated-project.zip", "test@example.com/defaultProject.zip", nil),
		presigner:  presignMock("test@example.com/defaultProject.zip", "", errors.NewServer("Presign failure")),
		wantErr:    errors.Wrap(errors.NewServer("Presign failure"), "Failed to generate pre-signed URL"),
	},
	{
		name:       "SuccessfulInvocation",
		cookie:     "validcookie",
		projectID:  "projectID",
		email:      "test@example.com",
		db:         &databaseMock{&dao.Project{Name: "Default Project", ID: "defaultProject"}, nil},
		downloader: createMock("/tmp/blank-sails.zip", "templates/sails.zip", nil),
		unzipper:   createMock("/tmp/blank-sails.zip", "/tmp", nil),
		generator:  generatorMock(&dao.Project{Name: "Default Project", ID: "defaultProject"}, "/tmp/blank-sails-project", nil),
		zipper:     createMock("/tmp/generated-project.zip", "/tmp/blank-sails-project", nil),
		uploader:   createMock("/tmp/generated-project.zip", "test@example.com/defaultProject.zip", nil),
		presigner:  presignMock("test@example.com/defaultProject.zip", "example.com", nil),
		wantURL:    "example.com",
	},
}

func TestGenerateCode(t *testing.T) {
	for _, test := range generateCodeTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			verifyCookie := verifyCookieMock(test.cookie, test.db, test.email, test.verifyErr)
			download = test.downloader
			unzip = test.unzipper
			generate = test.generator
			zipper = test.zipper
			upload = test.uploader
			presign = test.presigner

			// Execute
			url, err := generateCode(test.projectID, test.cookie, verifyCookie, test.db)

			// Verify
			if url != test.wantURL {
				t.Errorf("Got url %v; want %v", url, test.wantURL)
			}
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got err %v; want %v", err, test.wantErr)
			}
		})
	}
}
