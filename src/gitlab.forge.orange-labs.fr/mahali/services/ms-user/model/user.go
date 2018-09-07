package model

import "github.com/globalsign/mgo/bson"

//TODO valid tags in proto and controls

// UserNamespace used in DAO for storage of User
const UserNamespace = "user"

// User model definition
type User struct {
	// Id of the user
	ID bson.ObjectId `bson:"_id,omitempty"`
	// Username of the user
	Username string `bson:"username,omitempty"`
	// Phone Number of the user
	PhoneNumber string `bson:"phone_number,omitempty"`
	// Email of the user
	EmailAddress string `bson:"email_address,omitempty"`
	// Firstname of the user
	Firstname string `bson:"firstname,omitempty"`
	// Lastname of the user
	Lastname string `bson:"lastname,omitempty"`
	// Locale of the user in format fr, en ...
	Locale string `bson:"locale,omitempty"`
	// Role of the user
	Role string `bson:"role,omitempty"`
	// Picture url of the user
	Picture string `bson:"picture,omitempty"`
}
