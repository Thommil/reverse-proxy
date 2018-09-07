// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/authentication.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	proto/authentication.proto

It has these top-level messages:
	AuthenticateRequest
	JWT
	Token
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import proto2 "gitlab.forge.orange-labs.fr/mahali/services/ms-user/proto"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = proto2.User{}

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for AuthenticationService service

type AuthenticationService interface {
	Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...client.CallOption) (*JWT, error)
	Validate(ctx context.Context, in *Token, opts ...client.CallOption) (*proto2.User, error)
}

type authenticationService struct {
	c    client.Client
	name string
}

func NewAuthenticationService(name string, c client.Client) AuthenticationService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "proto"
	}
	return &authenticationService{
		c:    c,
		name: name,
	}
}

func (c *authenticationService) Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...client.CallOption) (*JWT, error) {
	req := c.c.NewRequest(c.name, "AuthenticationService.Authenticate", in)
	out := new(JWT)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationService) Validate(ctx context.Context, in *Token, opts ...client.CallOption) (*proto2.User, error) {
	req := c.c.NewRequest(c.name, "AuthenticationService.Validate", in)
	out := new(proto2.User)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AuthenticationService service

type AuthenticationServiceHandler interface {
	Authenticate(context.Context, *AuthenticateRequest, *JWT) error
	Validate(context.Context, *Token, *proto2.User) error
}

func RegisterAuthenticationServiceHandler(s server.Server, hdlr AuthenticationServiceHandler, opts ...server.HandlerOption) error {
	type authenticationService interface {
		Authenticate(ctx context.Context, in *AuthenticateRequest, out *JWT) error
		Validate(ctx context.Context, in *Token, out *proto2.User) error
	}
	type AuthenticationService struct {
		authenticationService
	}
	h := &authenticationServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&AuthenticationService{h}, opts...))
}

type authenticationServiceHandler struct {
	AuthenticationServiceHandler
}

func (h *authenticationServiceHandler) Authenticate(ctx context.Context, in *AuthenticateRequest, out *JWT) error {
	return h.AuthenticationServiceHandler.Authenticate(ctx, in, out)
}

func (h *authenticationServiceHandler) Validate(ctx context.Context, in *Token, out *proto2.User) error {
	return h.AuthenticationServiceHandler.Validate(ctx, in, out)
}