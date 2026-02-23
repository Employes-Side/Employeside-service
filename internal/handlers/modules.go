package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	models "github.com/Employes-Side/employee-side"
	modules "github.com/Employes-Side/employee-side"
	"github.com/Employes-Side/employee-side/internal/endpoints"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	errInvalidRequest = errors.New("invalid request")
	errBadRequest     = errors.New("bad request")
)

func NewModuleHandler(router *mux.Router, modules *endpoints.ModulesEndpoints) http.Handler {

	modulePath := router.PathPrefix("/modules").Subrouter()

	{
		modulePath.Methods(http.MethodPost).Path("").Handler(
			kithttp.NewServer(
				modules.Create,
				decodeCreateModuleRequest,
				kithttp.EncodeJSONResponse,
			),
		)

		modulePath.Methods(http.MethodPost).Path("/bulk").Handler(
			kithttp.NewServer(
				modules.BulkCreate,
				decodeBulkCreateModuleRequest,
				kithttp.EncodeJSONResponse,
				kithttp.ServerErrorEncoder(customErrorEncoder),
			),
		)

		modulePath.Methods(http.MethodPost).Path("/bulk_delete").Handler(
			kithttp.NewServer(
				modules.BulkDelete,
				decodeBulkDeleteRequest,
				kithttp.EncodeJSONResponse,
				kithttp.ServerErrorEncoder(customErrorEncoder),
			),
		)

		modulePath.Methods(http.MethodGet).Path("/{id}").Handler(
			kithttp.NewServer(
				modules.Read,
				decodeReadModuleRequest,
				kithttp.EncodeJSONResponse,
			),
		)
		modulePath.Methods(http.MethodGet).Path("").Handler(
			kithttp.NewServer(
				modules.List,
				decodeListModuleRequest,
				kithttp.EncodeJSONResponse,
			),
		)

		modulePath.Methods(http.MethodPut).Path("/{id}").Handler(
			kithttp.NewServer(
				modules.Update,
				decodeUpdateModuleRequest,
				kithttp.EncodeJSONResponse,
			),
		)

		modulePath.Methods(http.MethodDelete).Path("/{id}").Handler(
			kithttp.NewServer(
				modules.Delete,
				decodeReadModuleRequest,
				kithttp.EncodeJSONResponse,
			),
		)

	}
	return router
}

func customErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if errors.Is(err, errBadRequest) {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func decodeCreateModuleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var params modules.CreateModulesParameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, errBadRequest
	}
	return params, nil
}

func decodeBulkCreateModuleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var params modules.BulkModuleRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, errBadRequest
	}
	return params, nil
}

func decodeBulkDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var params modules.BulkDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, errBadRequest
	}
	return params, nil
}

func decodeReadModuleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	if id == "" {
		return nil, errInvalidRequest
	}

	return modules.ReadModulesRequest{
		By:    "id",
		Value: id,
	}, nil
}

func decodeListModuleRequest(_ context.Context, r *http.Request) (interface{}, error) {
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

func decodeUpdateModuleRequest(_ context.Context, r *http.Request) (interface{}, error) {
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
