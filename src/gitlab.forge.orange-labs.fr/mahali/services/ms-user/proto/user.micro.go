// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/user.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	proto/user.proto

It has these top-level messages:
	Id
	User
	Query
	Result
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/empty"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = google_protobuf.Empty{}

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for UserService service

type UserService interface {
	Create(ctx context.Context, in *User, opts ...client.CallOption) (*User, error)
	Get(ctx context.Context, in *Id, opts ...client.CallOption) (*User, error)
	Find(ctx context.Context, in *Query, opts ...client.CallOption) (*Result, error)
	Update(ctx context.Context, in *User, opts ...client.CallOption) (*User, error)
	Delete(ctx context.Context, in *Id, opts ...client.CallOption) (*google_protobuf.Empty, error)
}

type userService struct {
	c    client.Client
	name string
}

func NewUserService(name string, c client.Client) UserService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "proto"
	}
	return &userService{
		c:    c,
		name: name,
	}
}

func (c *userService) Create(ctx context.Context, in *User, opts ...client.CallOption) (*User, error) {
	req := c.c.NewRequest(c.name, "UserService.Create", in)
	out := new(User)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) Get(ctx context.Context, in *Id, opts ...client.CallOption) (*User, error) {
	req := c.c.NewRequest(c.name, "UserService.Get", in)
	out := new(User)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) Find(ctx context.Context, in *Query, opts ...client.CallOption) (*Result, error) {
	req := c.c.NewRequest(c.name, "UserService.Find", in)
	out := new(Result)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) Update(ctx context.Context, in *User, opts ...client.CallOption) (*User, error) {
	req := c.c.NewRequest(c.name, "UserService.Update", in)
	out := new(User)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) Delete(ctx context.Context, in *Id, opts ...client.CallOption) (*google_protobuf.Empty, error) {
	req := c.c.NewRequest(c.name, "UserService.Delete", in)
	out := new(google_protobuf.Empty)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserService service

type UserServiceHandler interface {
	Create(context.Context, *User, *User) error
	Get(context.Context, *Id, *User) error
	Find(context.Context, *Query, *Result) error
	Update(context.Context, *User, *User) error
	Delete(context.Context, *Id, *google_protobuf.Empty) error
}

func RegisterUserServiceHandler(s server.Server, hdlr UserServiceHandler, opts ...server.HandlerOption) error {
	type userService interface {
		Create(ctx context.Context, in *User, out *User) error
		Get(ctx context.Context, in *Id, out *User) error
		Find(ctx context.Context, in *Query, out *Result) error
		Update(ctx context.Context, in *User, out *User) error
		Delete(ctx context.Context, in *Id, out *google_protobuf.Empty) error
	}
	type UserService struct {
		userService
	}
	h := &userServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&UserService{h}, opts...))
}

type userServiceHandler struct {
	UserServiceHandler
}

func (h *userServiceHandler) Create(ctx context.Context, in *User, out *User) error {
	return h.UserServiceHandler.Create(ctx, in, out)
}

func (h *userServiceHandler) Get(ctx context.Context, in *Id, out *User) error {
	return h.UserServiceHandler.Get(ctx, in, out)
}

func (h *userServiceHandler) Find(ctx context.Context, in *Query, out *Result) error {
	return h.UserServiceHandler.Find(ctx, in, out)
}

func (h *userServiceHandler) Update(ctx context.Context, in *User, out *User) error {
	return h.UserServiceHandler.Update(ctx, in, out)
}

func (h *userServiceHandler) Delete(ctx context.Context, in *Id, out *google_protobuf.Empty) error {
	return h.UserServiceHandler.Delete(ctx, in, out)
}
