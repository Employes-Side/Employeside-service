package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	kithttp "github.com/go-kit/kit/transport/http"

	models "github.com/Employes-Side/employee-side"
	"github.com/Employes-Side/employee-side/internal/endpoints"
	"github.com/gorilla/mux"
)

var (
	errInvalidRequest = errors.New("invalid request")
	errBadRequest     = errors.New("bad request")
)

func NewHandler(router *mux.Router, users *endpoints.UserEndpoints) http.Handler {

	usersPath := router.PathPrefix("/users").Subrouter()

	{
		usersPath.Methods(http.MethodPost).Path("").Handler(
			kithttp.NewServer(
				users.Create,
				decodeCreateUserRequest,
				kithttp.EncodeJSONResponse,
			),
		)

		usersPath.Methods(http.MethodGet).Path("/{id}").Handler(
			kithttp.NewServer(
				users.Read,
				decodeReadUserRequest,
				kithttp.EncodeJSONResponse,
			),
		)
		usersPath.Methods(http.MethodGet).Path("").Handler(
			kithttp.NewServer(
				users.List,
				decodeListUserRequest,
				kithttp.EncodeJSONResponse,
			),
		)

		usersPath.Methods(http.MethodPut).Path("/{id}").Handler(
			kithttp.NewServer(
				users.Update,
				decodeUpdateUserRequest,
				kithttp.EncodeJSONResponse,
			),
		)

		usersPath.Methods(http.MethodDelete).Path("/{id}").Handler(
			kithttp.NewServer(
				users.Delete,
				decodeReadUserRequest,
				kithttp.EncodeJSONResponse,
			),
		)

	}
	return router
}

func decodeCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var params models.CreateUserParameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, errBadRequest
	}
	return params, nil
}

func decodeReadUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	if id == "" {
		return nil, errInvalidRequest
	}

	return models.ReadUserRequest{
		By:    "id",
		Value: id,
	}, nil
}

func decodeListUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
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

	return models.ListUserParameters{
		Limit:  limit,
		Offset: offset,
		Order:  order,
	}, nil
}

func decodeUpdateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	if id == "" {
		return nil, errInvalidRequest
	}

	var params models.UpdateUserParameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, errBadRequest
	}
	params.UserName = id
	return params, nil
}
