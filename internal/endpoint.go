package internal

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"testgenerate_backend_user/internal/app"
)

type Endpoints struct {
	GetUserEndpoint      endpoint.Endpoint
	GetUsersRoleEndpoint endpoint.Endpoint
	PostUserEndpoint     endpoint.Endpoint
	PutUserEndpoint      endpoint.Endpoint
	DeleteUserEndpoint   endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetUserEndpoint:      MakeGetUserEndpoint(s),
		GetUsersRoleEndpoint: MakeGetUsersRoleEndpoint(s),
		PostUserEndpoint:     MakePostUserEndpoint(s),
		PutUserEndpoint:      MakePutUserEndpoint(s),
		DeleteUserEndpoint:   MakeDeleteUserEndpoint(s),
	}
}

func (e Endpoints) GetUser(ctx context.Context, user, role string) (app.User, error) {
	request := getUserRequest{user, role}
	response, err := e.GetUserEndpoint(ctx, request)
	if err != nil {
		return app.User{}, err
	}
	resp := response.(getUserResponse)
	return resp.User, resp.Err
}

func (e Endpoints) GetUsersRole(ctx context.Context) ([]app.User, error) {
	request := getUsersRoleRequest{}
	response, err := e.GetUsersRoleEndpoint(ctx, request)
	if err != nil {
		return []app.User{}, err
	}
	resp := response.(getUsersRoleResponse)
	return resp.Users, resp.Err
}

func (e Endpoints) PostUser(ctx context.Context, user app.User) error {
	request := postUserRequest{user}
	response, err := e.PostUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(postUserResponse)
	return resp.Err
}

func (e Endpoints) PutUser(ctx context.Context, user app.User) error {
	request := putUserRequest{user}
	response, err := e.PutUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(putUserResponse)
	return resp.Err
}

func (e Endpoints) DeleteUser(ctx context.Context, userName string) error {
	request := deleteUserRequest{userName}
	response, err := e.DeleteUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteUserResponse)
	return resp.Err
}

// ----------------------------------------------------------------------------------------------------------------------
type getUserRequest struct {
	User string
	Role string
}

type getUserResponse struct {
	User app.User `json:"user,omitempty"`
	Err  error    `json:"err,omitempty"`
}

type getUsersRoleRequest struct {
}

type getUsersRoleResponse struct {
	Users []app.User `json:"users,omitempty"`
	Err   error      `json:"err,omitempty"`
}

type postUserRequest struct {
	User app.User
}

type postUserResponse struct {
	Err error `json:"err,omitempty"`
}

type putUserRequest struct {
	User app.User
}

type putUserResponse struct {
	Err error `json:"err,omitempty"`
}

type deleteUserRequest struct {
	UserName string
}

type deleteUserResponse struct {
	Err error `json:"err,omitempty"`
}

// ----------------------------------------------------------------------------------------------------------------------
func MakeGetUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getUserRequest)
		t, e := s.GetUser(ctx, req.User, req.Role)
		return getUserResponse{t, e}, nil
	}
}

func MakeGetUsersRoleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//req := request.(getUsersRoleRequest)
		t, e := s.GetUsersRole(ctx)
		return getUsersRoleResponse{t, e}, nil
	}
}

func MakePostUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postUserRequest)
		e := s.AddUser(ctx, req.User)
		return postUserResponse{e}, nil
	}
}

func MakePutUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(putUserRequest)
		e := s.UpdateUser(ctx, req.User)
		return putUserResponse{e}, nil
	}
}

func MakeDeleteUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteUserRequest)
		e := s.DeleteUser(ctx, req.UserName)
		return deleteUserResponse{e}, nil
	}
}
