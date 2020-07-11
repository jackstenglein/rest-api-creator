// Package dao contains structs for representing database objects in memory as well as functions
// to manipulate those objects in the database.
package dao

// User represents an instance of the User model in the database.
type User struct {
	Email    string              `dynamodbav:"Email" json:"email"`
	Password string              `dynamodbav:"Password" json:"-"`
	Token    string              `dynamodbav:"SessionToken" json:"-"`
	Projects map[string]*Project `dynamodbav:"Projects" json:"projects,omitempty"`
}

// Project represents an instance of the Project model in the database.
type Project struct {
	ID          string             `dynamodbav:"Id" json:"id"`
	Name        string             `dynamodbav:"Name" json:"name"`
	Description string             `dynamodbav:"Description" json:"description"`
	InstanceID  string             `dynamodbav:"InstanceId" json:"-"`
	DeployURL   string             `dynamodbav:"DeployUrl" json:"url"`
	Objects     map[string]*Object `dynamodbav:"Objects" json:"objects"`
}

// Object represents an instance of the Object model in the database.
type Object struct {
	ID          string       `dynamodbav:"Id" json:"id"`
	Name        string       `dynamodbav:"Name" json:"name"`
	CodeName    string       `dynamodbav:"CodeName" json:"-"`
	Description string       `dynamodbav:"Description" json:"description"`
	Attributes  []*Attribute `dynamodbav:"Attributes,omitempty" json:"attributes,omitempty"`
}

// Attribute represents an instance of the Attribute model in the database.
type Attribute struct {
	Name        string `dynamodbav:"Name" json:"name"`
	CodeName    string `dynamodbav:"CodeName" json:"-"`
	Type        string `dynamodbav:"Type" json:"type"`
	Required    bool   `dynamodbav:"Required" json:"required"`
	Description string `dynamodbav:"Description" json:"description"`
}
