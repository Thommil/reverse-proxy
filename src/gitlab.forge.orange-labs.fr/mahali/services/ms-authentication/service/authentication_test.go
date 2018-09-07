package service

import (
	"context"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	micro "github.com/micro/go-micro"
	"gitlab.forge.orange-labs.fr/mahali/services/ms-authentication/model"
	proto "gitlab.forge.orange-labs.fr/mahali/services/ms-authentication/proto"
	"gitlab.forge.orange-labs.fr/mahali/services/ms-common/crypto"
	userModel "gitlab.forge.orange-labs.fr/mahali/services/ms-user/model"
	userProto "gitlab.forge.orange-labs.fr/mahali/services/ms-user/proto"
)

const mongoUrl = "mongodb://localhost:27017/test"

func initClient() proto.AuthenticationService {
	service := micro.NewService(micro.Name("authentication.client"))
	return proto.NewAuthenticationService("authentication", service.Client())
}

func initUserClient() userProto.UserService {
	service := micro.NewService(micro.Name("user.client"))
	return userProto.NewUserService("user", service.Client())
}

func initDb() (*mgo.Session, error) {
	return mgo.Dial(mongoUrl)
}

func TestNewAuthenticationService(t *testing.T) {
	session, err := initDb()

	if err != nil {
		t.Errorf("user.Get() - %v", err)
		t.FailNow()
	}
	defer session.DB("").DropDatabase()

	userService := initUserClient()

	providersSettings := make(map[string]map[string]string)
	providersSettings["local"] = make(map[string]string)
	providersSettings["local"]["secret"] = "secret"

	jwtSettings := JWTSettings{
		Secret: "secret",
	}

	if _, err := NewAuthenticationService(providersSettings, session, jwtSettings, userService); err != nil {
		t.Errorf("NewAuthenticationService() - %v", err)
	}
}

func Test_authenticationService_AuthenticateValidate(t *testing.T) {
	client := initClient()
	session, err := initDb()

	if err != nil {
		t.Errorf("user.Get() - %v", err)
		t.FailNow()
	}
	defer session.DB("").DropDatabase()

	userId := bson.NewObjectId()
	session.DB("").C(userModel.UserNamespace).Insert(&userModel.User{
		ID:       userId,
		Username: "local",
	})

	if password, err := crypto.HashAndSaltPassword([]byte("$password$")); err != nil {
		t.Errorf("authentication.Authenticate() - Failed to create test dataset : %v", err)
		t.FailNow()
	} else {
		session.DB("").C(model.AccountNamespace).Insert(&model.Account{
			ID:       bson.NewObjectId(),
			Password: password,
			Provider: "local",
			UserID:   userId,
		})
	}

	authenticateRequest := &proto.AuthenticateRequest{
		Provider: "local",
		Credentials: map[string]string{
			"username": "local",
			"password": "$password$",
		},
	}

	token, err := client.Authenticate(context.TODO(), authenticateRequest)
	if err != nil {
		t.Errorf("authentication.Authenticate() - %v", err)
		t.FailNow()
	}

	user, err := client.Validate(context.TODO(), &proto.Token{Value: token.Token})

	if err != nil {
		t.Errorf("authentication.Validate() - %v", err)
	} else {
		if user.Id != userId.Hex() {
			t.Errorf("authentication.Validate() - Bad user returned %v", user)
		}
	}
}
