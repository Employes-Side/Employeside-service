package endpoints

import (
	"context"
	"errors"

	models "github.com/Employes-Side/employee-side"
	"github.com/Employes-Side/employee-side/internal/repositories"
)

func NewBlogEndpoint(manager repositories.BlogRepository) *BlogEndpoints {
	return &BlogEndpoints{manager}
}

type BlogEndpoints struct {
	manager repositories.BlogRepository
}

func (ep *BlogEndpoints) Create(ctx context.Context, req interface{}) (interface{}, error) {
	params, ok := req.(models.CreatBlogParameters)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.Create(ctx, params)
}

func (ep *BlogEndpoints) Read(ctx context.Context, req interface{}) (interface{}, error) {
	readReq, ok := req.(models.ReadBlogRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.Read(ctx, readReq)
}

func (ep *BlogEndpoints) Update(ctx context.Context, req interface{}) (interface{}, error) {
	updateReq, ok := req.(models.UpdateBlogParameters)
	if !ok {
		return nil, errors.New("invalid request")
	}
	readReq := models.ReadBlogRequest{
		By:    "id",
		Value: updateReq.BlogTitle,
	}
	return ep.manager.Update(ctx, readReq, updateReq)
}

func (ep *BlogEndpoints) Delete(ctx context.Context, req interface{}) (interface{}, error) {
	deleteReq, ok := req.(models.ReadBlogRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.Delete(ctx, deleteReq)
}

func (ep *BlogEndpoints) List(ctx context.Context, req interface{}) (interface{}, error) {
	listReq, ok := req.(models.ListBlogsParameters)
	if !ok {
		return nil, errors.New("invalid request")
	}
	return ep.manager.List(ctx, listReq)
}
