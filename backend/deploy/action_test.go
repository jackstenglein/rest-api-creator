package deploy

import (
	"reflect"
	"testing"

	"github.com/jackstenglein/rest_api_creator/backend/auth"
	"github.com/jackstenglein/rest_api_creator/backend/dao"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
)

func verifyCookieMock(mockCookie string, mockDB auth.UserGetter, mockEmail string, mockErr error) auth.VerifyCookieFunc {
	return func(cookie string, db auth.UserGetter) (string, error) {
		if cookie != mockCookie || !reflect.DeepEqual(db, mockDB) {
			return "", errors.NewServer("Incorrect input to VerifyCookie mock")
		}
		return mockEmail, mockErr
	}
}

type databaseMock struct {
	email      string
	projectID  string
	project    *dao.Project
	getErr     error
	instanceID string
	url        string
	updateErr  error
}

func (mock *databaseMock) GetUserInfo(email string) (*dao.User, error) {
	return nil, nil
}

func (mock *databaseMock) GetProject(email string, projectID string) (*dao.Project, error) {
	if email != mock.email || projectID != mock.projectID {
		return nil, errors.NewServer("Incorrect input to GetProject mock.")
	}
	return mock.project, mock.getErr
}

func (mock *databaseMock) UpdateDeployment(email string, projectID string, instanceID string, url string) error {
	if email != mock.email || projectID != mock.projectID || instanceID != mock.instanceID || url != mock.url {
		return errors.NewServer("Incorrect input to UpdateDeployment mock")
	}
	return mock.updateErr
}

type ec2Mock struct {
	instanceID string
	launchURL  string
	launchErr  error
	getURL     string
	getErr     error
}

func (mock *ec2Mock) LaunchInstance() (string, string, error) {
	return mock.instanceID, mock.launchURL, mock.launchErr
}

func (mock *ec2Mock) GetPublicURL(instanceID string) (string, error) {
	if instanceID != mock.instanceID {
		return "", errors.NewServer("Incorrect input to GetPublicURL mock")
	}
	return mock.getURL, mock.getErr
}

var deployProjectTests = []struct {
	name      string
	cookie    string
	projectID string

	// Mock data
	db        *databaseMock
	email     string
	verifyErr error
	ec2       *ec2Mock

	wantURL string
	wantErr error
}{
	{
		name:    "EmptyProjectID",
		wantErr: errors.NewClient("Parameter `pid` is required"),
	},
	{
		name:      "InvalidCookie",
		cookie:    "cookievalue",
		projectID: "project",
		verifyErr: errors.NewClient("Not authenticated"),
		wantErr:   errors.Wrap(errors.NewClient("Not authenticated"), "Failed to verify cookie"),
	},
	{
		name:      "GetProjectFailure",
		cookie:    "cookievalue",
		projectID: "project",
		db:        &databaseMock{email: "test@example.com", projectID: "project", project: nil, getErr: errors.NewServer("Database failure")},
		email:     "test@example.com",
		wantErr:   errors.Wrap(errors.NewServer("Database failure"), "Failed to get project"),
	},
	{
		name:      "LaunchInstanceFailure",
		cookie:    "cookievalue",
		projectID: "project",
		db:        &databaseMock{email: "test@example.com", projectID: "project", project: &dao.Project{ID: "project", Name: "Project", InstanceID: "instance"}},
		email:     "test@example.com",
		ec2:       &ec2Mock{launchErr: errors.NewServer("EC2 failure")},
		wantErr:   errors.Wrap(errors.NewServer("EC2 failure"), "Failed to launch EC2 instance"),
	},
	{
		name:      "UpdateDeploymentFailure",
		cookie:    "cookievalue",
		projectID: "project",
		db: &databaseMock{
			email:      "test@example.com",
			projectID:  "project",
			project:    &dao.Project{ID: "project", Name: "Project", InstanceID: "instance"},
			instanceID: "newinstance",
			url:        "instanceurl",
			updateErr:  errors.NewServer("Database failure"),
		},
		email:   "test@example.com",
		ec2:     &ec2Mock{instanceID: "newinstance", launchURL: "instanceurl", launchErr: nil},
		wantErr: errors.Wrap(errors.NewServer("Database failure"), "Failed to update deployment info"),
	},
	{
		name:      "SuccessfulInvocation",
		cookie:    "cookievalue",
		projectID: "project",
		db: &databaseMock{
			email:      "test@example.com",
			projectID:  "project",
			project:    &dao.Project{ID: "project", Name: "Project", InstanceID: "instance"},
			instanceID: "newinstance",
			url:        "instanceurl",
		},
		email:   "test@example.com",
		ec2:     &ec2Mock{instanceID: "newinstance", launchURL: "instanceurl", launchErr: nil},
		wantURL: "instanceurl",
	},
	{
		name:      "GetURLFailure",
		cookie:    "cookievalue",
		projectID: "project",
		db: &databaseMock{
			email:      "test@example.com",
			projectID:  "project",
			project:    &dao.Project{ID: "project", Name: "Project", InstanceID: "instance"},
			instanceID: "newinstance",
			url:        "instanceurl",
		},
		email:   "test@example.com",
		ec2:     &ec2Mock{instanceID: "newinstance", launchURL: "", launchErr: nil, getErr: errors.NewServer("DescribeInstance failure")},
		wantErr: errors.Wrap(errors.NewServer("DescribeInstance failure"), "Failed to get instance DNS name"),
	},
	{
		name:      "GetURLSuccess",
		cookie:    "cookievalue",
		projectID: "project",
		db: &databaseMock{
			email:      "test@example.com",
			projectID:  "project",
			project:    &dao.Project{ID: "project", Name: "Project", InstanceID: "instance"},
			instanceID: "newinstance",
			url:        "instanceurl",
		},
		email:   "test@example.com",
		ec2:     &ec2Mock{instanceID: "newinstance", launchURL: "", launchErr: nil, getURL: "instanceurl"},
		wantURL: "instanceurl",
	},
}

func TestDeployProject(t *testing.T) {
	for _, test := range deployProjectTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			verifyCookie := verifyCookieMock(test.cookie, test.db, test.email, test.verifyErr)
			delay = func() {}

			// Execute
			url, err := deployProject(test.cookie, test.projectID, verifyCookie, test.db, test.ec2)

			// Verify
			if url != test.wantURL {
				t.Errorf("Got url `%s`; want `%s`", url, test.wantURL)
			}
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got err %v; want %v", err, test.wantErr)
			}
		})
	}
}
