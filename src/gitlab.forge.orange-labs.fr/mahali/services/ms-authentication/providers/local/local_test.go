package local

import (
	"testing"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	micro "github.com/micro/go-micro"

	"gitlab.forge.orange-labs.fr/mahali/services/ms-authentication/model"
	userModel "gitlab.forge.orange-labs.fr/mahali/services/ms-user/model"
	userProto "gitlab.forge.orange-labs.fr/mahali/services/ms-user/proto"
)

const mongoUrl = "mongodb://localhost:27017/test"

func initUserClient() userProto.UserService {
	ms := micro.NewService(micro.Name("ms-user.client"))
	return userProto.NewUserService("ms-user", ms.Client())
}

func initDb() (*mgo.Session, error) {
	return mgo.Dial(mongoUrl)
}

func Test_provider_Authenticate(t *testing.T) {
	session, err := initDb()

	if err != nil {
		t.Errorf("local.Authenticate() - %v", err)
		t.FailNow()
	}
	defer session.DB("").DropDatabase()

	userClient := initUserClient()

	userId := bson.NewObjectId()
	session.DB("").C(userModel.UserNamespace).Insert(&userModel.User{
		ID:       userId,
		Username: "local",
	})

	if password, err := HashAndSaltPassword([]byte("$password$")); err != nil {
		t.Errorf("local.Authenticate() - %v", err)
		t.FailNow()
	} else {
		session.DB("").C(model.AccountNamespace).Insert(&model.Account{
			ID:       bson.NewObjectId(),
			Password: password,
			Provider: "local",
			UserID:   userId,
		})
	}

	if local, err := New(map[string]string{"secret": "mysecret"}, session, userClient); err != nil {
		t.Errorf("local.New() - %v", err)
		t.FailNow()
	} else {
		user, err := local.Authenticate(map[string]string{
			"username": "local",
			"password": "$password$",
		})
		if err != nil {
			t.Errorf("local.Authenticate() - %v", err)
		} else {
			if bson.ObjectIdHex(user.Id) != userId {
				t.Errorf("local.Authenticate() - bad user ID, wanted %s, got %s", userId.Hex(), user.Id)
			}
		}
	}
}
