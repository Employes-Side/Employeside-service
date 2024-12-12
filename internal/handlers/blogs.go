package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	models "github.com/Employes-Side/employee-side"
	"github.com/Employes-Side/employee-side/internal/endpoints"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewBlogHandler(blogs *endpoints.BlogEndpoints) http.Handler {
	router := mux.NewRouter()

	blogPath := router.PathPrefix("/blogs").Subrouter()

	{

		blogPath.Methods(http.MethodPost).Path("").Handler(
			kithttp.NewServer(
				blogs.Create,
				decodeCreateBlogRequest,
				kithttp.EncodeJSONResponse,
			),
		)

		blogPath.Methods(http.MethodGet).Path("").Handler(
			kithttp.NewServer(
				blogs.List,
				decodeListBlogRequest,
				kithttp.EncodeJSONResponse,
			),
		)

		blogPath.Methods(http.MethodGet).Path("/{id}").Handler(
			kithttp.NewServer(
				blogs.Read,
				decodeReadBlogRequest,
				kithttp.EncodeJSONResponse,
			),
		)
		blogPath.Methods(http.MethodPut).Path("/{id}").Handler(
			kithttp.NewServer(
				blogs.Update,
				decodeUpdateBlogRequest,
				kithttp.EncodeJSONResponse,
			),
		)

		blogPath.Methods(http.MethodDelete).Path("/{id}").Handler(
			kithttp.NewServer(
				blogs.Delete,
				decodeReadBlogRequest,
				kithttp.EncodeJSONResponse,
			),
		)

	}
	return router
}

func decodeCreateBlogRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var params models.CreatBlogParameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, errBadRequest
	}
	return params, nil
}

func decodeReadBlogRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	if id == "" {
		return nil, errInvalidRequest
	}

	return models.ReadBlogRequest{
		By:    "id",
		Value: id,
	}, nil
}

func decodeListBlogRequest(_ context.Context, r *http.Request) (interface{}, error) {
	query := r.URL.Query()

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil {
		offset = 0
	}

	order := query.Get("order")
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	return models.ListBlogsParameters{
		Limit:  limit,
		Offset: offset,
		Order:  order,
	}, nil
}

func decodeUpdateBlogRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	if id == "" {
		return nil, errInvalidRequest
	}

	var params models.UpdateBlogParameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, errBadRequest
	}
	params.BlogTitle = id
	return params, nil
}
