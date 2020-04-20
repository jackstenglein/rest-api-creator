package actions

import (
	"github.com/jackstenglein/rest_api_creator/backend-sls/authentication"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	apierrors "github.com/jackstenglein/rest_api_creator/backend-sls/errors"
	"github.com/pkg/errors"
)

type GetProjectRequest struct {
	Id     string `json:"id"`
	Cookie string `json:"-"`
}

type GetProjectAction struct {
	store dao.DataStore
	auth  authentication.Authenticator
}

func DefaultGetProjectAction() *GetProjectAction {
	return &GetProjectAction{store: dao.DefaultDynamoStore(), auth: authentication.NewSessionAuthenticator()}
}

func NewGetProjectAction(store dao.DataStore, auth authentication.Authenticator) *GetProjectAction {
	return &GetProjectAction{store: store, auth: auth}
}

func (action *GetProjectAction) GetProject(request GetProjectRequest) (*dao.Project, error) {
	if request.Id == "" {
		return nil, apierrors.NewUserError("Parameter `id` is required")
	}

	email, token, hmac, err := authentication.SplitCookie(request.Cookie)
	if err != nil {
		return nil, apierrors.NewUserError("Not authenticated")
	}

	ok, err := action.auth.VerifyCookie(email, token, hmac, action.store)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to verify cookie")
	}
	if !ok {
		return nil, apierrors.NewUserError("Not authenticated")
	}

	project, err := action.store.GetProject(email, request.Id)
	return project, errors.Wrap(err, "Failed to get project")
}
