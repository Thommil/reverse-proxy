package local

import (
	"context"

	"github.com/globalsign/mgo/bson"

	"gitlab.forge.orange-labs.fr/mahali/services/ms-authentication/model"
	"gitlab.forge.orange-labs.fr/mahali/services/ms-common/crypto"
	userProto "gitlab.forge.orange-labs.fr/mahali/services/ms-user/proto"

	"github.com/globalsign/mgo"
	"github.com/micro/go-micro/errors"
)

// Allows to check user entry against local DB
type provider struct {
	configuration map[string]string
	db            *mgo.Session
	userService   userProto.UserService
}

// New instanciates a new Local provider
func New(configuration map[string]string, db *mgo.Session, userService userProto.UserService) (*provider, error) {
	return &provider{
		configuration: configuration,
		db:            db,
		userService:   userService,
	}, nil
}

func (p provider) Authenticate(credentials map[string]string) (*userProto.User, error) {
	if username, ok := credentials["username"]; ok == true {
		if password, ok := credentials["password"]; ok == true {
			query := &userProto.Query{Filter: "{\"username\": \"" + username + "\"}"}
			result, err := p.userService.Find(context.TODO(), query)
			if err != nil {
				return nil, err
			}

			if result.Total == 0 || len(result.Items) == 0 {
				return nil, errors.NotFound("account_not_found", "invalid username or password")
			}

			user := result.Items[0]

			account := &model.Account{
				Provider: "local",
				UserID:   bson.ObjectIdHex(user.Id),
			}
			err = p.db.DB("").C(model.AccountNamespace).Find(account).One(account)
			if err != nil {
				return nil, errors.NotFound("account_not_found", err.Error())
			}

			if !crypto.ComparePasswords(account.Password, password) {
				return nil, errors.NotFound("account_not_found", "invalid username or password")
			}

			return user, nil
		}
	}

	return nil, errors.BadRequest("authentication_bad_request", "malformed credentials")
}
