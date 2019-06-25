package project

import (
	"time"
)

const (
	UserUrl = "/api/v1/users"
)

/*
|--------------------------------------------------------------------------
| Other Entities
|--------------------------------------------------------------------------
|
*/
type Timestamp struct {
	CreatedAt 	time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	CreatedBy 	AuthUser  `json:"created_by,omitempty" bson:"created_by,omitempty"`
	UpdatedAt 	time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	UpdatedBy 	AuthUser  `json:"updated_by,omitempty" bson:"updated_by,omitempty"`
	DeletedAt 	time.Time `json:"-,omitempty" bson:"deleted_at,omitempty"`
}

/*
|--------------------------------------------------------------------------
| In Source Entities
|--------------------------------------------------------------------------
|
*/
type Auth struct {
	User 	AuthUser 	`json: "auth_user, omitempty" bson: "auth_user,omitempty"`
	Token   string 		`json:"token"`
}

type AuthUser struct {
	ID 			string `json:"id,omitempty" bson:"id,omitempty"`
	FullName 	string `json:"full_name,omitempty" bson:"full_name,omitempty"`
}

type User struct {
	ID 			interface{}	`json:"id,omitempty" bson:"_id,omitempty"`
	FullName	string		`json:"full_name,omitempty" bson:"full_name,omitempty"`

	Timestamp
}