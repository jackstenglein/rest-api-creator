package actions

import (
	"github.com/rest_api_creator/backend-sls/authentication"
	"github.com/rest_api_creator/backend-sls/dao"
	"github.com/rest_api_creator/backend-sls/errors"
)

type GetProjectRequest struct {
	ProjectId string `json:"project"`
	Cookie    string `json:"-"`
}

type GetProjectResponse struct {
	Project dao.Project `json:"project,omitempty"`
	Error   error       `json:"error,omitempty"`
}

type GetProjectAction struct {
	store dao.DataStore
	auth  authentication.Authenticator
}

func NewGetProjectAction(store dao.DataStore, auth authentication.Authenticator) *GetProjectAction {
	return &GetProjectAction{store: store, auth: auth}
}

func (action *GetProjectAction) GetProject(request GetProjectRequest) GetProjectResponse {
	if request.ProjectId == "" {
		return GetProjectResponse{Error: errors.NewUserError("project is required")}
	}

	email, token, hmac, err := authentication.SplitCookie(request.Cookie)
	if err != nil {
		return GetProjectResponse{Error: errors.NewUserError("Not authenticated")}
	}

	ok, err := action.auth.VerifyCookie(email, token, hmac, action.store)
	if err != nil {
		return GetProjectResponse{Error: err}
	}
	if !ok {
		return GetProjectResponse{Error: errors.NewUserError("Not authenticated")}
	}

	project, err := action.store.GetProject(email, request.ProjectId)
	return GetProjectResponse{Project: project, Error: err}
}
