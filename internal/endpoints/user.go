package endpoints

import (
	"context"
	"errors"

	"github.com/Employes-Side/employee-side/internal/models"
	"github.com/Employes-Side/employee-side/internal/repositories"
)

func NewUserEndpoint(manager repositories.UserRepository) *UserEndpoints {
	return &UserEndpoints{manager}
}

type UserEndpoints struct {
	manager repositories.UserRepository
}

func (ep *UserEndpoints) Create(ctx context.Context, req interface{}) (interface{}, error) {
	params, ok := req.(models.CreateUserParameters)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.Create(ctx, params)
}

func (ep *UserEndpoints) Read(ctx context.Context, req interface{}) (interface{}, error) {
	readReq, ok := req.(models.ReadUserRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.Read(ctx, readReq)
}

func (ep *UserEndpoints) Update(ctx context.Context, req interface{}) (interface{}, error) {
	updateReq, ok := req.(models.UpdateUserParameters)
	if !ok {
		return nil, errors.New("invalid request")
	}
	readReq := models.ReadUserRequest{
		By:    "id",
		Value: updateReq.UserName,
	}
	return ep.manager.Update(ctx, readReq, updateReq)
}

func (ep *UserEndpoints) Delete(ctx context.Context, req interface{}) (interface{}, error) {
	deleteReq, ok := req.(models.ReadUserRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.Delete(ctx, deleteReq)
}

func (ep *UserEndpoints) List(ctx context.Context, req interface{}) (interface{}, error) {
	listReq, ok := req.(models.ListUserParameters)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.List(ctx, listReq)
}
