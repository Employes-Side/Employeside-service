package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	models "github.com/Employes-Side/employee-side"
	"github.com/Employes-Side/employee-side/internal/endpoints"

	"context"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewWriterHandler(router *mux.Router, writers *endpoints.WriterEndpoints) http.Handler {
	writerPath := router.PathPrefix("/writer").Subrouter()

	{
		writerPath.Methods(http.MethodPost).Path("").Handler(
			kithttp.NewServer(
				writers.Create,
				decodeCreateWriterRequest,
				kithttp.EncodeJSONResponse,
			),
		)

		writerPath.Methods(http.MethodGet).Path("/{id}").Handler(
			kithttp.NewServer(
				writers.Read,
				decodeReadWriterRequest,
				kithttp.EncodeJSONResponse,
			),
		)
		writerPath.Methods(http.MethodGet).Path("").Handler(
			kithttp.NewServer(
				writers.List,
				decodeListWriterRequest,
				kithttp.EncodeJSONResponse,
			),
		)

		writerPath.Methods(http.MethodPut).Path("/{id}").Handler(
			kithttp.NewServer(
				writers.Update,
				decodeUpdateWriterRequest,
				kithttp.EncodeJSONResponse,
			),
		)

		writerPath.Methods(http.MethodDelete).Path("/{id}").Handler(
			kithttp.NewServer(
				writers.Delete,
				decodeReadWriterRequest,
				kithttp.EncodeJSONResponse,
			),
		)
	}
	return router
}

func decodeCreateWriterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var params models.CreateWriterRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, errBadRequest
	}
	return params, nil
}

func decodeReadWriterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	if id == "" {
		return nil, errInvalidRequest
	}

	return models.ReadWriterRequest{
		By:    "id",
		Value: id,
	}, nil
}

func decodeListWriterRequest(_ context.Context, r *http.Request) (interface{}, error) {
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

	return models.ListParameters{
		Limit:  limit,
		Offset: offset,
		Order:  order,
	}, nil
}

func decodeUpdateWriterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	if id == "" {
		return nil, errInvalidRequest
	}

	var params models.UpdateWriterParameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, errBadRequest
	}
	params.UserName = id
	return params, nil
}
