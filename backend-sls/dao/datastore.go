package dao

type User struct {
	Email    string `dynamodbav:"Email" json:"email"`
	Password string `dynamodbav:"Password" json:"-"`
	Token    string `dynamodbav:"SessionToken" json:"-"`
}

type Project struct {
	Id      string   `dynamodbav:"Id" json:"id"`
	Name    string   `dynamodbav:"Name" json:"name"`
	Objects []Object `dynamodbav:"Objects" json:"objects,omitempty"`
}

type Object struct {
	Id          string `dynamodbav:"Id" json:"id"`
	Name        string `dynamodbav:"Name" json:"name"`
	Description string `dynamodbav:"Description" json:"description"`
}

type DataStore interface {
	CreateUser(string, string, string) error
	GetUser(string) (User, error)
	GetProject(string, string) (*Project, error)
	UpdateUserToken(string, string) error
}
