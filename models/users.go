package models

import (
	"gopkg.in/mgo.v2/bson"
)

// User model
type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Email    string        `json:"email"`
	Password string        `json:"password"`
	Role     string        `json:"role"`
}

// PasswordChange model
type PasswordChange struct {
	Email       string
	OldPassword string
	NewPassword string
}
