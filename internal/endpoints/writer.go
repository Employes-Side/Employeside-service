package endpoints

import (
	"context"
	"errors"

	models "github.com/Employes-Side/employee-side"
	"github.com/Employes-Side/employee-side/internal/repositories"
)

func NewWriterEndpoint(manager repositories.WriterRepository) *WriterEndpoints {
	return &WriterEndpoints{manager}
}

type WriterEndpoints struct {
	manager repositories.WriterRepository
}

func (ep *WriterEndpoints) Create(ctx context.Context, req interface{}) (interface{}, error) {
	params, ok := req.(models.CreateWriterRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.Create(ctx, params)
}

func (ep *WriterEndpoints) Read(ctx context.Context, req interface{}) (interface{}, error) {
	readReq, ok := req.(models.ReadWriterRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.Read(ctx, readReq)
}

func (ep *WriterEndpoints) Update(ctx context.Context, req interface{}) (interface{}, error) {
	updateReq, ok := req.(models.UpdateWriterParameters)
	if !ok {
		return nil, errors.New("invalid request")
	}
	readReq := models.ReadWriterRequest{
		By:    "id",
		Value: updateReq.UserName,
	}
	return ep.manager.Update(ctx, readReq, updateReq)
}

func (ep *WriterEndpoints) Delete(ctx context.Context, req interface{}) (interface{}, error) {
	deleteReq, ok := req.(models.ReadWriterRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.Delete(ctx, deleteReq)
}

func (ep *WriterEndpoints) List(ctx context.Context, req interface{}) (interface{}, error) {
	listReq, ok := req.(models.ListParameters)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.List(ctx, listReq)
}
