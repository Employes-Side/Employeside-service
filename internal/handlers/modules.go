package handlers

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	modules "github.com/Employes-Side/employee-side"
	"github.com/Employes-Side/employee-side/internal/endpoints"
	"github.com/Employes-Side/employee-side/internal/services"
	gokitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	errInvalidRequest = errors.New("invalid request")
	errBadRequest     = errors.New("bad request")
)

type S3Service interface {
	UploadFile(ctx context.Context, key string, data []byte, contentType string) (string, error)
	DeleteFile(ctx context.Context, key string) error
}

type ModuleRepository interface {
	GetS3Keys(ctx context.Context, ids []string) ([]string, error)
	BulkDelete(ctx context.Context, req modules.BulkDeleteRequest) error
}

type ModuleHandler struct {
	endpoints endpoints.ModuleEndpointsInterface
	s3Service S3Service
	repo      ModuleRepository
}

func NewModuleHandler(router *mux.Router, moduleEndpoints endpoints.ModuleEndpointsInterface, s3Svc S3Service, repo ModuleRepository) http.Handler {
	handler := &ModuleHandler{
		endpoints: moduleEndpoints,
		s3Service: s3Svc,
		repo:      repo,
	}

	modulePath := router.PathPrefix("/modules").Subrouter()
	modulePath.Use(requestIDMiddleware)
	modulePath.Use(latencyMiddleware)

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
				handler.bulkCreatePipeline(handler.endpoints.BulkCreate),
				handler.decodeBulkImportRequest,
				kithttp.EncodeJSONResponse,
				kithttp.ServerErrorEncoder(customErrorEncoder),
			),
		)

		modulePath.Methods(http.MethodPost).Path("/bulk_delete").HandlerFunc(handler.handleAsyncBulkDelete)

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

func latencyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID, _ := r.Context().Value("request-id").(string)
		
		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapped, r.WithContext(r.Context()))
		
		duration := time.Since(start)
		if duration > 2*time.Second {
			log.Printf("WARN [%s] High latency detected: %v", reqID, duration)
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), "request-id", reqID)
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *ModuleHandler) bulkCreatePipeline(next gokitendpoint.Endpoint) gokitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(modules.BulkImportRequest)
		reqID, _ := ctx.Value("request-id").(string)

		start := time.Now()
		log.Printf("[%s] Starting Bulk Create Operation", reqID)

		s3Start := time.Now()

		var moduleData []byte
		var contentType string
		if req.Format == "csv" {
			moduleData = []byte(req.Data)
			contentType = "text/csv"
		} else {
			moduleData, _ = json.Marshal(req.Modules)
			contentType = "application/json"
		}

		s3Key := services.GenerateS3Key("modules/bulk", req.Format)
		s3Url, err := h.s3Service.UploadFile(ctx, s3Key, moduleData, contentType)
		s3Latency := time.Since(s3Start)

		if err != nil {
			log.Printf("[%s] S3 Upload Failed: %v", reqID, err)
			return nil, fmt.Errorf("failed to upload to S3: %w", err)
		}

		for i := range req.Modules {
			req.Modules[i].S3Key = &s3Key
			req.Modules[i].S3Url = &s3Url
		}

		bulkReq := modules.BulkModuleRequest{
			Modules: req.Modules,
		}

		dbStart := time.Now()
		resp, err := next(ctx, bulkReq)
		dbLatency := time.Since(dbStart)

		if err != nil {
			log.Printf("[%s] DB Transaction Failed: %v. Cleaning up S3...", reqID, err)
			cleanupCtx := context.WithValue(context.WithoutCancel(ctx), "request-id", reqID)
			go func() {
				if delErr := h.s3Service.DeleteFile(cleanupCtx, s3Key); delErr != nil {
					log.Printf("[%s] CRITICAL: Failed to cleanup S3 file %s: %v", reqID, s3Key, delErr)
				} else {
					log.Printf("[%s] S3 Cleanup Successful", reqID)
				}
			}()
			return nil, err
		}

		totalDuration := time.Since(start)
		log.Printf("[%s] Bulk Create Completed. Total: %v (S3: %v, DB: %v)", reqID, totalDuration, s3Latency, dbLatency)

		if totalDuration > 2*time.Second {
			log.Printf("WARN [%s] High Latency Detected! Total: %v (S3: %v, DB: %v)", reqID, totalDuration, s3Latency, dbLatency)
		}

		return resp, nil
	}
}

func (h *ModuleHandler) handleAsyncBulkDelete(w http.ResponseWriter, r *http.Request) {
	var req modules.BulkDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	jobID := uuid.New().String()
	reqID, _ := r.Context().Value("request-id").(string)

	go func(jobID, reqID string, deleteReq modules.BulkDeleteRequest) {
		// Use detached context for async processing
		ctx := context.WithValue(context.Background(), "request-id", reqID)
		log.Printf("[%s] Starting Async Delete Job %s for %d items", reqID, jobID, len(deleteReq.IDs))

		keys, err := h.repo.GetS3Keys(ctx, deleteReq.IDs)
		if err != nil {
			log.Printf("[%s] Job %s Failed to fetch S3 keys: %v", reqID, jobID, err)
			return
		}

		for _, key := range keys {
			if key != "" {
				if err := h.s3Service.DeleteFile(ctx, key); err != nil {
					log.Printf("[%s] Job %s Failed to delete S3 key %s: %v", reqID, jobID, key, err)
				}
			}
		}

		if err := h.repo.BulkDelete(ctx, deleteReq); err != nil {
			log.Printf("[%s] Job %s Failed to delete from DB: %v", reqID, jobID, err)
			return
		}

		log.Printf("[%s] Async Delete Job %s Completed", reqID, jobID)
	}(jobID, reqID, req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"job_id": jobID,
		"status": "accepted",
	})
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

	format := strings.ToLower(req.Format)
	if format == "" {
		format = "json"
	}

	var parsedModules []modules.Modules
	if format == "csv" {
		var err error
		parsedModules, err = parseCSVModules(req.Data)
		if err != nil {
			return nil, errBadRequest
		}
	} else {
		if len(req.Modules) == 0 && req.Data != "" {
			if err := json.Unmarshal([]byte(req.Data), &parsedModules); err != nil {
				return nil, errBadRequest
			}
		} else {
			parsedModules = req.Modules
		}
	}

	return modules.BulkImportRequest{
		Format:  format,
		Modules: parsedModules,
		Data:    req.Data,
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

	return modules.ListParameters{
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

	var params modules.UpdateModulesParameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, errBadRequest
	}
	params.ID = id
	return params, nil
}
