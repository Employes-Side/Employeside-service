package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	models "github.com/Employes-Side/employee-side"
	modules "github.com/Employes-Side/employee-side"
	"github.com/Employes-Side/employee-side/internal/endpoints"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
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

func decodeCreateModuleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var params modules.CreateModulesParameters
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
