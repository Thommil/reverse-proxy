package model

import "github.com/globalsign/mgo/bson"

//TODO valid tags in proto and controls

// AccountNamespace used in DAO for storage of Account
const AccountNamespace = "account"

// Account model definition
type Account struct {
	// ID of the account
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`
	// Provider of the account
	Provider string `json:"provider,omitempty" bson:"provider,omitempty"`
	// ExternalID of the account
	ExternalID string `json:"external_id,omitempty" bson:"external_id,omitempty"`
	// User ID linked to this account
	UserID bson.ObjectId `json:"user_id,omitempty" bson:"user_id,omitempty"`
	// Password of the account (depends on provider)
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}
