package handlers

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	models "github.com/Employes-Side/employee-side"
	modules "github.com/Employes-Side/employee-side"
	"github.com/Employes-Side/employee-side/internal/endpoints"
	"github.com/Employes-Side/employee-side/internal/services"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	errInvalidRequest = errors.New("invalid request")
	errBadRequest     = errors.New("bad request")
)

type ModuleHandler struct {
	endpoints *endpoints.ModulesEndpoints
	s3Service *services.S3Service
}

func NewModuleHandler(router *mux.Router, moduleEndpoints *endpoints.ModulesEndpoints, s3Svc *services.S3Service) http.Handler {
	handler := &ModuleHandler{
		endpoints: moduleEndpoints,
		s3Service: s3Svc,
	}

	modulePath := router.PathPrefix("/modules").Subrouter()

	{
		modulePath.Methods(http.MethodPost).Path("").Handler(
			kithttp.NewServer(
				handler.endpoints.Create,
				decodeCreateModuleRequest,
				kithttp.EncodeJSONResponse,
			),
		)

		modulePath.Methods(http.MethodPost).Path("/bulk").Handler(
			kithttp.NewServer(
				handler.endpoints.BulkCreate,
				handler.decodeBulkImportRequest,
				kithttp.EncodeJSONResponse,
				kithttp.ServerErrorEncoder(customErrorEncoder),
			),
		)

		modulePath.Methods(http.MethodPost).Path("/bulk_delete").Handler(
			kithttp.NewServer(
				handler.endpoints.BulkDelete,
				decodeBulkDeleteRequest,
				kithttp.EncodeJSONResponse,
				kithttp.ServerErrorEncoder(customErrorEncoder),
			),
		)

		modulePath.Methods(http.MethodGet).Path("/{id}").Handler(
			kithttp.NewServer(
				handler.endpoints.Read,
				decodeReadModuleRequest,
				kithttp.EncodeJSONResponse,
			),
		)
		modulePath.Methods(http.MethodGet).Path("").Handler(
			kithttp.NewServer(
				handler.endpoints.List,
				decodeListModuleRequest,
				kithttp.EncodeJSONResponse,
			),
		)

		modulePath.Methods(http.MethodPut).Path("/{id}").Handler(
			kithttp.NewServer(
				handler.endpoints.Update,
				decodeUpdateModuleRequest,
				kithttp.EncodeJSONResponse,
			),
		)

		modulePath.Methods(http.MethodDelete).Path("/{id}").Handler(
			kithttp.NewServer(
				handler.endpoints.Delete,
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

func (h *ModuleHandler) decodeBulkImportRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req modules.BulkImportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest
	}

	// Determine format
	format := strings.ToLower(req.Format)
	if format == "" {
		format = "json"
	}

	var moduleData []byte
	var contentType string

	if format == "csv" {
		moduleData = []byte(req.Data)
		contentType = "text/csv"
	} else {
		
		if len(req.Modules) > 0 {
			var err error
			moduleData, err = json.Marshal(req.Modules)
			if err != nil {
				return nil, errBadRequest
			}
		} else if req.Data != "" {
			moduleData = []byte(req.Data)
		}
		contentType = "application/json"
	}

	if len(moduleData) == 0 {
		return nil, errBadRequest
	}

	
	s3Key := services.GenerateS3Key("modules/bulk", format)
	s3Location, err := h.s3Service.UploadFile(context.Background(), s3Key, moduleData, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %w", err)
	}

	
	var parsedModules []modules.Modules
	if format == "csv" {
		parsedModules, err = parseCSVModules(string(moduleData))
		if err != nil {
			return nil, errBadRequest
		}
	} else {
		if err := json.Unmarshal(moduleData, &parsedModules); err != nil {
			return nil, errBadRequest
		}
	}

	return modules.BulkImportRequest{
		Format:  format,
		Modules: parsedModules,
		S3Key:   s3Location,
	}, nil
}

func parseCSVModules(csvData string) ([]modules.Modules, error) {
	reader := csv.NewReader(strings.NewReader(csvData))

	
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	
	fieldIndex := make(map[string]int)
	for i, h := range header {
		fieldIndex[strings.TrimSpace(h)] = i
	}

	var modulesList []modules.Modules
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		module := modules.Modules{}

		if idx, ok := fieldIndex["module_name"]; ok && idx < len(record) {
			module.ModuleName = record[idx]
		}
		if idx, ok := fieldIndex["user_id"]; ok && idx < len(record) {
			module.UserID = record[idx]
		}
		if idx, ok := fieldIndex["module_type"]; ok && idx < len(record) {
			module.ModuleType = record[idx]
		}
		if idx, ok := fieldIndex["module_desc"]; ok && idx < len(record) {
			module.Module_Desc = record[idx]
		}
		if idx, ok := fieldIndex["module_short_name"]; ok && idx < len(record) {
			module.ModuleShortName = record[idx]
		}
		if idx, ok := fieldIndex["module_price"]; ok && idx < len(record) {
			module.ModulePrice, _ = strconv.ParseInt(record[idx], 10, 64)
		}
		if idx, ok := fieldIndex["purchased"]; ok && idx < len(record) {
			module.Purchased = strings.ToLower(record[idx]) == "true"
		}

		modulesList = append(modulesList, module)
	}

	return modulesList, nil
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
