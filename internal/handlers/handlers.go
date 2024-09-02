package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/Employes-Side/employee-side/internal/endpoints"
	"github.com/Employes-Side/employee-side/internal/models"
	"github.com/gorilla/mux"
)

var (
	errInvalidRequest = errors.New("invalid request")
	errBadRequest     = errors.New("bad request")
)

func NewHandler(users *endpoints.UserEndpoints) http.Handler {
	router := mux.NewRouter()

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

func decodeUpdateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["userID"]
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
