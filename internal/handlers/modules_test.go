package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	modules "github.com/Employes-Side/employee-side"
	"github.com/Employes-Side/employee-side/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockS3Service struct {
	mock.Mock
}

func (m *MockS3Service) UploadFile(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	args := m.Called(ctx, key, data, contentType)
	return args.String(0), args.Error(1)
}

func (m *MockS3Service) DeleteFile(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

type MockModuleRepository struct {
	mock.Mock
}

func (m *MockModuleRepository) GetS3Keys(ctx context.Context, ids []string) ([]string, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockModuleRepository) BulkDelete(ctx context.Context, req modules.BulkDeleteRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// MockModuleEndpoints implements the endpoints.ModuleEndpointsInterface for testing
type MockModuleEndpoints struct {
	CreateFunc     func(ctx context.Context, req interface{}) (interface{}, error)
	ReadFunc       func(ctx context.Context, req interface{}) (interface{}, error)
	UpdateFunc     func(ctx context.Context, req interface{}) (interface{}, error)
	DeleteFunc     func(ctx context.Context, req interface{}) (interface{}, error)
	ListFunc       func(ctx context.Context, req interface{}) (interface{}, error)
	BulkCreateFunc func(ctx context.Context, req interface{}) (interface{}, error)
	BulkDeleteFunc func(ctx context.Context, req interface{}) (interface{}, error)
}

func (m *MockModuleEndpoints) Create(ctx context.Context, req interface{}) (interface{}, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockModuleEndpoints) Read(ctx context.Context, req interface{}) (interface{}, error) {
	if m.ReadFunc != nil {
		return m.ReadFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockModuleEndpoints) Update(ctx context.Context, req interface{}) (interface{}, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockModuleEndpoints) Delete(ctx context.Context, req interface{}) (interface{}, error) {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockModuleEndpoints) List(ctx context.Context, req interface{}) (interface{}, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockModuleEndpoints) BulkCreate(ctx context.Context, req interface{}) (interface{}, error) {
	if m.BulkCreateFunc != nil {
		return m.BulkCreateFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockModuleEndpoints) BulkDelete(ctx context.Context, req interface{}) (interface{}, error) {
	if m.BulkDeleteFunc != nil {
		return m.BulkDeleteFunc(ctx, req)
	}
	return nil, nil
}

// --- Tests ---

func TestBulkCreate_Rollback_S3Success_DBFail(t *testing.T) {
	// Setup
	mockS3 := new(MockS3Service)
	mockRepo := new(MockModuleRepository)

	// Mock Endpoint that fails (simulating DB error)
	mockEndpoints := &MockModuleEndpoints{
		BulkCreateFunc: func(ctx context.Context, request interface{}) (interface{}, error) {
			return nil, errors.New("db transaction failed")
		},
	}

	router := mux.NewRouter()
	handlers.NewModuleHandler(router, mockEndpoints, mockS3, mockRepo)

	// Request Data
	reqData := modules.BulkImportRequest{
		Format: "json",
		Modules: []modules.Modules{
			{ModuleName: "Test Module"},
		},
	}
	body, _ := json.Marshal(reqData)

	// Expectations
	// 1. UploadFile called. Verify Request-ID is passed.
	mockS3.On("UploadFile", mock.MatchedBy(func(ctx context.Context) bool {
		val, ok := ctx.Value("request-id").(string)
		return ok && val == "test-req-id"
	}), mock.Anything, mock.Anything, "application/json").Return("http://s3-url", nil)

	// 2. DeleteFile called (Rollback). Verify Request-ID is preserved in cleanup context.
	mockS3.On("DeleteFile", mock.MatchedBy(func(ctx context.Context) bool {
		val, ok := ctx.Value("request-id").(string)
		return ok && val == "test-req-id"
	}), mock.Anything).Return(nil)

	// Execute
	req := httptest.NewRequest("POST", "/modules/bulk", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", "test-req-id")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "db transaction failed", resp["error"])

	// Wait a bit for the async cleanup to complete
	time.Sleep(100 * time.Millisecond)

	mockS3.AssertExpectations(t)
}

func TestAsyncBulkDelete_Polling(t *testing.T) {
	mockS3 := new(MockS3Service)
	mockRepo := new(MockModuleRepository)
	mockEndpoints := &MockModuleEndpoints{}

	router := mux.NewRouter()
	handlers.NewModuleHandler(router, mockEndpoints, mockS3, mockRepo)

	reqData := modules.BulkDeleteRequest{
		IDs: []string{"id1"},
	}
	body, _ := json.Marshal(reqData)

	var jobCompleted atomic.Bool

	// Expectations
	// Verify context has request-id in async job
	mockRepo.On("GetS3Keys", mock.MatchedBy(func(ctx context.Context) bool {
		val, ok := ctx.Value("request-id").(string)
		return ok && val == "async-req-id"
	}), []string{"id1"}).Return([]string{"key1"}, nil)

	mockS3.On("DeleteFile", mock.Anything, "key1").Return(nil)

	mockRepo.On("BulkDelete", mock.Anything, reqData).Run(func(args mock.Arguments) {
		jobCompleted.Store(true)
	}).Return(nil)

	// Execute
	req := httptest.NewRequest("POST", "/modules/bulk_delete", bytes.NewReader(body))
	req.Header.Set("X-Request-ID", "async-req-id")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert 202 Accepted
	assert.Equal(t, http.StatusAccepted, w.Code)

	// Polling assertion
	assert.Eventually(t, func() bool {
		return jobCompleted.Load()
	}, 2*time.Second, 50*time.Millisecond, "Async delete job did not complete in time")

	mockRepo.AssertExpectations(t)
	mockS3.AssertExpectations(t)
}
