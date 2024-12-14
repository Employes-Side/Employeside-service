package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	modules "github.com/Employes-Side/employee-side"
	"github.com/Employes-Side/employee-side/generated/employeside/model"
	"github.com/Employes-Side/employee-side/generated/employeside/table"
	"github.com/google/uuid"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func NewModulesManger(db *sql.DB) *ModulesRepository {
	return &ModulesRepository{db}
}

type ModulesRepository struct {
	db *sql.DB
}

func (mgr *ModulesRepository) Read(ctx context.Context, req modules.ReadModulesRequest) (*modules.Modules, error) {
	conditions, err := mgr.buildReadClause(req)
	if err != nil {
		return nil, err

	}

	statement := table.Modules.SELECT(table.Modules.AllColumns).WHERE(conditions)

	var module model.Modules
	if err := statement.QueryContext(ctx, mgr.db, &module); err != nil {
		if err == qrm.ErrNoRows {
			return nil, errors.New("module not found")
		}
		return nil, err
	}

	return convertToModulesDBMode(module)
}

func (mgr *ModulesRepository) Create(ctx context.Context, params modules.CreateModulesParameters) (*modules.Modules, error) {
	id := uuid.New()

	now := time.Now()

	module := model.Modules{
		ID:              id.String(),
		ModuleName:      params.ModuleName,
		ModuleType:      params.ModuleType,
		ModuleDesc:      &params.Module_Desc,
		ModuleShortName: &params.ModuleShortName,
		ModulePrice:     &params.ModulePrice,
		Purchased:       &params.Purchased,
		UserID:          params.UserID,
		CreatedAt:       &now,
		UpdatedAt:       &now,
	}

	statement := table.Modules.INSERT(table.Modules.AllColumns).MODEL(module)

	_, err := statement.ExecContext(ctx, mgr.db)
	if err != nil {
		return nil, err
	}

	return mgr.Read(ctx, modules.ReadModulesRequest{By: "id", Value: id.String()})
}

func (mgr *ModulesRepository) buildReadClause(req modules.ReadModulesRequest) (mysql.BoolExpression, error) {
	switch req.By {
	case "id":
		return table.Modules.ID.EQ(mysql.String(req.Value)), nil

	default:
		return nil, errors.New("by should be one of id ") //TODO: add other case to read for
	}
}

func convertToModulesDBMode(module model.Modules) (*modules.Modules, error) {
	return &modules.Modules{
		ID:              module.ID,
		ModuleName:      module.ModuleName,
		ModuleType:      module.ModuleType,
		Module_Desc:     *module.ModuleDesc,
		ModuleShortName: *module.ModuleShortName,
		ModulePrice:     *module.ModulePrice,
		Purchased:       *module.Purchased,
		UserID:          module.UserID,
		CreatedAt:       module.CreatedAt,
		UpdatedAt:       module.UpdatedAt,
	}, nil
}
