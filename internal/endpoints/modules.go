package endpoints

import (
	"context"
	"errors"

	modules "github.com/Employes-Side/employee-side"

	"github.com/Employes-Side/employee-side/internal/repositories"
)

func NewModuleEndpoint(manager repositories.ModulesRepository) *ModulesEndpoints {
	return &ModulesEndpoints{manager}
}

type ModulesEndpoints struct {
	manager repositories.ModulesRepository
}

func (ep *ModulesEndpoints) Create(ctx context.Context, req interface{}) (interface{}, error) {
	params, ok := req.(modules.CreateModulesParameters)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.Create(ctx, params)
}

func (ep *ModulesEndpoints) Read(ctx context.Context, req interface{}) (interface{}, error) {
	readReq, ok := req.(modules.ReadModulesRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.Read(ctx, readReq)
}

func (ep *ModulesEndpoints) Update(ctx context.Context, req interface{}) (interface{}, error) {
	updateReq, ok := req.(modules.UpdateModulesParameters)
	if !ok {
		return nil, errors.New("invalid request")
	}
	readReq := modules.ReadModulesRequest{
		By:    "id",
		Value: updateReq.ModuleName,
	}
	return ep.manager.Update(ctx, readReq, updateReq)
}

func (ep *ModulesEndpoints) Delete(ctx context.Context, req interface{}) (interface{}, error) {
	deleteReq, ok := req.(modules.ReadModulesRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.Delete(ctx, deleteReq)
}

func (ep *ModulesEndpoints) List(ctx context.Context, req interface{}) (interface{}, error) {
	listReq, ok := req.(modules.ListParameters)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.List(ctx, listReq)
}
