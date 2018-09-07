package service

import (
	"context"
	"strconv"
	"testing"

	"gitlab.forge.orange-labs.fr/mahali/services/ms-user/model"

	"github.com/globalsign/mgo"

	"github.com/globalsign/mgo/bson"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/errors"
	proto "gitlab.forge.orange-labs.fr/mahali/services/ms-user/proto"
)

const mongoUrl = "mongodb://localhost:27017/test"

func initClient() proto.UserService {
	service := micro.NewService(micro.Name("user.client"))
	return proto.NewUserService("user", service.Client())
}

func initDb() (*mgo.Session, error) {
	return mgo.Dial(mongoUrl)
}

func Test_userService_Create(t *testing.T) {
	client := initClient()
	inUser := &proto.User{Username: "create"}
	user, err := client.Create(context.TODO(), inUser)

	if err != nil {
		t.Errorf("user.Create() - %v", err)
		t.FailNow()
	}

	if user.Id == "" {
		t.Errorf("user.Create() - No ID found in returned user")
	}

	if user.Username != inUser.Username {
		t.Errorf("user.Create() - Username not set, wanted %s, got %s", inUser.Username, user.Username)
	}
}

func Test_userService_Get(t *testing.T) {
	client := initClient()
	id := &proto.Id{Value: bson.NewObjectId().Hex()}

	session, err := initDb()

	if err != nil {
		t.Errorf("user.Get() - %v", err)
		t.FailNow()
	}
	defer session.DB("").DropDatabase()

	var user *proto.User
	user, err = client.Get(context.TODO(), id)

	if err == nil {
		t.Errorf("user.Get() - Waiting not found error, got no error")
	} else {
		code := errors.Parse(err.Error()).Code
		if code != 404 {
			t.Errorf("user.Get() - Waiting not found bad error, got %d", code)
		}
	}

	session.DB("").C(model.UserNamespace).Insert(&model.User{
		ID:       bson.ObjectIdHex(id.Value),
		Username: "get",
	})

	user, err = client.Get(context.TODO(), id)

	if err != nil {
		t.Errorf("user.Get() - %v", err)
	} else {
		if user.Username != "get" {
			t.Errorf("user.Get() - bad user, wanted 'get', got %s", user.Username)
		}
	}

}

func Test_userService_Find(t *testing.T) {
	client := initClient()
	session, err := initDb()

	if err != nil {
		t.Errorf("user.Find() - %v", err)
		t.FailNow()
	}
	defer session.DB("").DropDatabase()

	for i := 0; i < 111; i++ {
		err := session.DB("").C(model.UserNamespace).Insert(&model.User{
			Username: strconv.Itoa(i),
			Role:     "user",
		})
		if err != nil {
			t.Errorf("user.Find() - Mongo dataset creation failed : %v", err)
			t.FailNow()
		}
	}

	for i := 0; i < 10; i++ {
		err := session.DB("").C(model.UserNamespace).Insert(&model.User{
			Username: strconv.Itoa(i),
			Role:     "admin",
		})
		if err != nil {
			t.Errorf("user.Find() - Mongo dataset creation failed : %v", err)
			t.FailNow()
		}
	}

	//Test Find Mongo Query OK
	query := proto.Query{
		Offset: 0,
		Limit:  DefaultPaginationLimit,
		Sort:   "",
		Filter: "{\"role\": \"user\"}",
	}

	if result, err := client.Find(context.TODO(), &query); err != nil {
		t.Errorf("user.Find() Test Find Mongo Query OK, error = %v", err)
	} else {
		if result.Total != 111 {
			t.Errorf("user.Find() Test Find Mongo Query OK, bad total : wanted 111, got %d", result.Total)
		}
		if len(result.Items) != DefaultPaginationLimit {
			t.Errorf("user.Find() Test Find Mongo Query OK, bad number of items : wanted %d, got %d", DefaultPaginationLimit, len(result.Items))
		}
	}

	//Test Find Mongo Query Empty
	query = proto.Query{
		Offset: 0,
		Limit:  DefaultPaginationLimit,
		Sort:   "",
		Filter: "{\"role\": \"fake\"}",
	}

	if result, err := client.Find(context.TODO(), &query); err != nil {
		t.Errorf("user.Find() Test Find Mongo Query Empty, error = %v", err)
	} else {
		if len(result.Items) != 0 {
			t.Errorf("uesr.Find() Test Find Mongo Query Empty, bad number of items : wanted %d, got %d", 0, len(result.Items))
		}
	}

	//Test Find Mongo Query Sorted OK
	query = proto.Query{
		Offset: 0,
		Limit:  DefaultPaginationLimit,
		Sort:   "-username",
		Filter: "{\"role\": \"user\"}",
	}

	if result, err := client.Find(context.TODO(), &query); err != nil {
		t.Errorf("user.Find() Test Find Mongo Query Sorted OK, error = %v", err)
	} else {
		if result.Items[0].Username != "99" {
			t.Errorf("user.Find() Test Find Mongo Query Sorted OK, first item KO, wanted 99, got %s", result.Items[0].Username)
		}
	}

	//Test Find Mongo Query Max OK
	query = proto.Query{
		Offset: 0,
		Limit:  MaxPaginationLimit,
		Sort:   "",
		Filter: "{\"role\": \"user\"}",
	}

	if result, err := client.Find(context.TODO(), &query); err != nil {
		t.Errorf("user.Find() Test Find Mongo Query Max OK, error = %v", err)
	} else {
		if len(result.Items) != MaxPaginationLimit {
			t.Errorf("user.Find() Test Find Mongo Query Max OK, bad number of items : wanted %d, got %d", MaxPaginationLimit, len(result.Items))
		}
	}

	//Test Find Mongo Query Page OK
	query = proto.Query{
		Offset: 10,
		Limit:  DefaultPaginationLimit,
		Sort:   "",
		Filter: "{\"role\": \"user\"}",
	}

	if result, err := client.Find(context.TODO(), &query); err != nil {
		t.Errorf("user.Find() Test Find Mongo Query Page OK, error = %v", err)
	} else {
		if result.Items[0].Username != strconv.Itoa(DefaultPaginationLimit) {
			t.Errorf("user.Find() Test Find Mongo Query Page OK, first item KO, wanted %d, got %s", DefaultPaginationLimit, result.Items[0].Username)
		}
	}

	//Test Find Mongo Query Bool OK
	query = proto.Query{
		Offset: 0,
		Limit:  DefaultPaginationLimit,
		Sort:   "",
		Filter: "{ \"$or\": [ { \"username\": \"1\" }, { \"username\": \"2\" } ] }",
	}
	if result, err := client.Find(context.TODO(), &query); err != nil {
		t.Errorf("user.Find() Test Find Mongo Query Bool OK, error = %v", err)
	} else {
		if result.Total != 4 {
			t.Errorf("user.Find() Test Find Bool OK, bad total : wanted 4, got %d", result.Total)
		}
		if len(result.Items) != 4 {
			t.Errorf("user.Find() Test Find Bool OK, bad number of items, wanted 4, got %d", len(result.Items))
		}
	}
}

func Test_userService_Update(t *testing.T) {
	client := initClient()
	inUser := &proto.User{Id: bson.NewObjectId().Hex(), Username: "update"}
	session, err := initDb()

	if err != nil {
		t.Errorf("user.Update() - %v", err)
		t.FailNow()
	}
	defer session.DB("").DropDatabase()

	session.DB("").C(model.UserNamespace).Insert(&model.User{
		ID:       bson.ObjectIdHex(inUser.Id),
		Username: "create",
	})

	inUser, err = client.Update(context.TODO(), inUser)

	if err != nil {
		t.Errorf("user.Update() - %v", err)
	} else {
		if inUser.Username != "update" {
			t.Errorf("user.Update() - bad user, wanted 'update', got %s", inUser.Username)
		}
	}

}

func Test_userService_Delete(t *testing.T) {
	client := initClient()
	id := &proto.Id{Value: bson.NewObjectId().Hex()}

	session, err := initDb()

	if err != nil {
		t.Errorf("user.Delete() - %v", err)
		t.FailNow()
	}
	defer session.DB("").DropDatabase()

	_, err = client.Delete(context.TODO(), id)

	if err == nil {
		t.Errorf("user.Delete() - Waiting not found error, got no error")
	} else {
		code := errors.Parse(err.Error()).Code
		if code != 404 {
			t.Errorf("user.Delete() - Waiting not found bad error, got %d", code)
		}
	}

	session.DB("").C(model.UserNamespace).Insert(&model.User{
		ID:       bson.ObjectIdHex(id.Value),
		Username: "delete",
	})

	_, err = client.Delete(context.TODO(), id)

	if err != nil {
		t.Errorf("user.Delete() - %v", err)
	}
}
