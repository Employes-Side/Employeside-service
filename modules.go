package employeside

import (
	"context"
	"time"
)

type Modules struct {
	ID              string     `json:"id"`
	ModuleName      string     `json:"module_name"`
	UserID          string     `json:"user_id"`
	ModuleType      string     `json:"module_type"`
	Module_Desc     string     `json:"module_desc"`
	ModuleShortName string     `json:"module_short_name"`
	ModulePrice     int64      `json:"module_price"`
	Purchased       bool       `json:"purchased"`
	S3Key           *string    `json:"s3_key,omitempty"`
	S3Url           *string    `json:"s3_url,omitempty"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type ReadModulesRequest struct {
	By    string
	Value string
}

type ModulesPage struct {
	TotalRecords int       `json:"total_records"`
	Users        []Modules `json:"modules"`
	Limit        int       `json:"limit"`
	Offset       int       `json:"offset"`
}

type CreateModulesParameters struct {
	ModuleName      string `json:"module_name"`
	ModuleType      string `json:"module_type"`
	Module_Desc     string `json:"module_desc"`
	ModuleShortName string `json:"module_short_name"`
	ModulePrice     int64  `json:"module_price"`
	Purchased       bool   `json:"purchased"`
	UserID          string `json:"user_id"`
}

type UpdateModulesParameters struct {
	ID              string `json:"id"`
	ModuleName      string `json:"module_name"`
	ModuleType      string `json:"module_type"`
	Module_Desc     string `json:"module_desc"`
	ModuleShortName string `json:"module_short_name"`
	ModulePrice     int64  `json:"module_price"`
	Purchased       bool   `json:"purchased"`
	UserID          string `json:"user_id"`
}

type ModulesManager interface {
	Read(ctx context.Context, req ReadModulesRequest) (*Modules, error)
	Create(ctx context.Context, params CreateModulesParameters) (*Modules, error)
	List(ctx context.Context, params ListParameters) (ModulesPage, error)
	Update(ctx context.Context, req ReadModulesRequest, params UpdateModulesParameters) (*Modules, error)
	Delete(ctx context.Context, req ReadModulesRequest) (*Modules, error)
	BulkAddModules(ctx context.Context, req BulkModuleRequest) ([]*Modules, error)
	BulkDelete(ctx context.Context, req BulkDeleteRequest) error
}

type BulkModuleRequest struct {
	Modules []Modules `json:"modules"`
}

type BulkDeleteRequest struct {
	IDs []string `json:"ids"`
}

type BulkImportRequest struct {
	Format  string    `json:"format"`            // "json" or "csv"
	Data    string    `json:"data"`              // Raw JSON or CSV data
	Modules []Modules `json:"modules,omitempty"` // For JSON format
	S3Key   string    `json:"s3_key,omitempty"`  // S3 location after upload
}

type BulkImportResponse struct {
	Modules      []*Modules `json:"modules"`
	S3Location   string     `json:"s3_location"`
	TotalRecords int        `json:"total_records"`
}
