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
	ModulePrice     string     `json:"module_price"`
	Purchased       bool       `json:"purchased"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
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
	ModulePrice     string `json:"module_price"`
	Purchased       bool   `json:"purchased"`
	UserID          string `json:"user_id"`
}

type UpdateModulesParameters struct {
	ModuleName      string `json:"module_name"`
	ModuleType      string `json:"module_type"`
	Module_Desc     string `json:"module_desc"`
	ModuleShortName string `json:"module_short_name"`
	ModulePrice     string `json:"module_price"`
	Purchased       bool   `json:"purchased"`
	UserID          string `json:"user_id"`
}

type ModulesManager interface {
	Read(ctx context.Context, req ReadModulesRequest) (*Modules, error)
	Create(ctx context.Context, params CreateModulesParameters) (*Modules, error)
	List(ctx context.Context, params ListParameters) (ModulesPage, error)
	Update(ctx context.Context, req ReadModulesRequest, params UpdateModulesParameters) (*Modules, error)
	Delete(ctx context.Context, req ReadModulesRequest) (*Modules, error)
}
