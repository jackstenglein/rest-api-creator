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
	ID      string             `dynamodbav:"Id" json:"id"`
	Name    string             `dynamodbav:"Name" json:"name"`
	Objects map[string]*Object `dynamodbav:"Objects" json:"objects,omitempty"`
}

// Object represents an instance of the Object model in the database.
type Object struct {
	ID          string `dynamodbav:"Id" json:"id"`
	Name        string `dynamodbav:"Name" json:"name"`
	Description string `dynamodbav:"Description" json:"description"`
}
