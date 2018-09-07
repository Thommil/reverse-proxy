package service

import (
	"context"
	"encoding/json"

	"github.com/globalsign/mgo/bson"
	"github.com/micro/go-micro/errors"

	"github.com/globalsign/mgo"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/copier"
	"gitlab.forge.orange-labs.fr/mahali/services/ms-user/model"
	proto "gitlab.forge.orange-labs.fr/mahali/services/ms-user/proto"
)

// DefaultPaginationLimit indicates the default pagination limit
const DefaultPaginationLimit = 10

// MaxPaginationLimit indicates the max pagination limit
const MaxPaginationLimit = 50

type userService struct {
	db *mgo.Session
}

// New create a new instance of user service
func NewUserService(session *mgo.Session) *userService {
	return &userService{
		db: session,
	}
}

func (c *userService) Create(ctx context.Context, in *proto.User, out *proto.User) error {
	collection := c.db.DB("").C(model.UserNamespace)

	user := model.User{}
	copier.Copy(&user, in)
	user.ID = bson.NewObjectId()

	err := collection.Insert(&user)
	if err != nil {
		errors.InternalServerError("user_internal_error", err.Error())
	}

	copier.Copy(out, &user)
	out.Id = user.ID.Hex()

	return nil
}

func (c *userService) Get(ctx context.Context, in *proto.Id, out *proto.User) error {
	collection := c.db.DB("").C(model.UserNamespace)

	user := model.User{}

	if !bson.IsObjectIdHex(in.Value) {
		return errors.BadRequest("user_bad_request", "invalid id")
	}

	err := collection.FindId(bson.ObjectIdHex(in.Value)).One(&user)

	if err != nil {
		if err.Error() == "not found" {
			return errors.NotFound("user_not_found", err.Error())
		}
		return errors.InternalServerError("user_internal_error", err.Error())
	}

	copier.Copy(out, &user)
	out.Id = user.ID.Hex()

	return nil
}

func (c *userService) Find(ctx context.Context, in *proto.Query, out *proto.Result) error {
	collection := c.db.DB("").C(model.UserNamespace)

	filter := make(map[string]interface{})
	err := json.Unmarshal([]byte(in.Filter), &filter)

	if err != nil {
		errors.InternalServerError("user_internal_error", err.Error())
	}

	q := collection.Find(filter)

	//Total
	total, err := q.Count()

	if err != nil {
		errors.InternalServerError("user_internal_error", err.Error())
	}

	//sort
	if in.Sort != "" {
		q.Sort(in.Sort)
	}

	//offset
	if in.Offset > 0 {
		q.Skip(int(in.Offset))
	}

	//limit
	if in.Limit <= 0 {
		q.Limit(DefaultPaginationLimit)
	} else if in.Limit > MaxPaginationLimit {
		q.Limit(MaxPaginationLimit)
	} else {
		q.Limit(int(in.Limit))
	}

	users := []model.User{}
	err = q.All(&users)

	if err != nil {
		errors.InternalServerError("user_internal_error", err.Error())
	}
	for index := range users {
		user := &proto.User{}
		copier.Copy(user, users[index])
		user.Id = users[index].ID.Hex()
		out.Items = append(out.Items, user)
	}

	out.Total = int32(total)
	out.Offset = in.Offset
	out.Limit = in.Limit
	out.Sort = in.Sort

	return nil
}

func (c *userService) Update(ctx context.Context, in *proto.User, out *proto.User) error {
	collection := c.db.DB("").C(model.UserNamespace)

	user := model.User{}
	copier.Copy(&user, in)
	user.ID = bson.ObjectIdHex(in.Id)

	err := collection.UpdateId(user.ID, &user)
	if err != nil {
		if err.Error() == "not found" {
			return errors.NotFound("user_not_found", err.Error())
		}
		return errors.InternalServerError("user_internal_error", err.Error())
	}

	copier.Copy(out, in)

	return nil
}

func (c *userService) Delete(ctx context.Context, in *proto.Id, out *google_protobuf.Empty) error {
	collection := c.db.DB("").C(model.UserNamespace)

	if !bson.IsObjectIdHex(in.Value) {
		return errors.BadRequest("user_bad_request", "invalid id")
	}

	err := collection.RemoveId(bson.ObjectIdHex(in.Value))

	if err != nil {
		if err.Error() == "not found" {
			return errors.NotFound("user_not_found", err.Error())
		}
		return errors.InternalServerError("user_internal_error", err.Error())
	}

	return nil
}
