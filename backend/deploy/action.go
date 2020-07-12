package deploy

import (
	"github.com/jackstenglein/rest_api_creator/backend/auth"
	"github.com/jackstenglein/rest_api_creator/backend/dao"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
	"github.com/jackstenglein/rest_api_creator/backend/log"
)

// deployDatabase wraps the database functions used by the deployProject action in order to allow dependency injection.
type deployDatabase interface {
	auth.UserGetter
	GetProject(string, string) (*dao.Project, error)
	UpdateDeployment(string, string, string, string) error
}

// deployer wraps the EC2 functions used by the deployProject action in order to allow dependency injection.
type deployer interface {
	GetPublicURL(string) (string, error)
	LaunchInstance(string) (string, string, error)
}

// deployProject launches an EC2 instance to run the given project. If the deployment is successful and the EC2 instance launches
// within 5 seconds, the public DNS name of the instance is returned.
func deployProject(cookie string, projectID string, deployRequest deployRequest, verifyCookie auth.VerifyCookieFunc, db deployDatabase, ec2 deployer) (string, string, error) {

	if projectID == "" {
		return "", "", errors.NewClient("Parameter `pid` is required")
	}

	if deployRequest.URL == "" {
		return "", "", errors.NewClient("Parameter `url` is required")
	}

	email, err := verifyCookie(cookie, db)
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to verify cookie")
	}

	project, err := db.GetProject(email, projectID)
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to get project")
	}
	log.Info("Got project:", project)

	if project.InstanceID != "" {
		// TODO: terminate old instance
	}

	// Launch new instance
	log.Info("Launching instance")
	instanceID, url, err := ec2.LaunchInstance(deployRequest.URL)
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to launch EC2 instance")
	}

	log.Info("Updating deployment")
	err = db.UpdateDeployment(email, projectID, instanceID, url)
	if err != nil {
		// TODO: terminate instance in order to not leak EC2 instances
		return "", "", errors.Wrap(err, "Failed to update deployment info")
	}

	return instanceID, url, nil
}
