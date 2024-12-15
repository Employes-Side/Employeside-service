package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	models "github.com/Employes-Side/employee-side"
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

func (mgr *ModulesRepository) List(ctx context.Context, params models.ListParameters) (*models.ModulesPage, error) {
	var orderExpr mysql.OrderByClause
	if params.Order == "desc" {
		orderExpr = table.Modules.ModuleName.DESC()
	} else {
		orderExpr = table.Modules.ModuleName.ASC()
	}

	statement := table.Modules.
		SELECT(table.Modules.AllColumns).
		ORDER_BY(orderExpr).
		LIMIT(int64(params.Limit)).
		OFFSET(int64(params.Offset))

	sqlQuery, args := statement.Sql()
	log.Printf("Executing SQL query: %s, with args: %v", sqlQuery, args)

	var modules []model.Modules

	if err := statement.QueryContext(ctx, mgr.db, &modules); err != nil {
		if err == qrm.ErrNoRows {
			return &models.ModulesPage{}, nil
		}
		return &models.ModulesPage{}, err
	}

	userModels := make([]models.Modules, len(modules))
	for i, u := range modules {
		userModel, err := convertToModulesDBMode(u)
		if err != nil {
			return &models.ModulesPage{}, err
		}
		userModels[i] = *userModel
	}

	countStatement := table.Modules.SELECT(mysql.COUNT(table.Modules.ID))

	countSqlQuery, countArgs := countStatement.Sql()
	log.Printf("Executing COUNT SQL query: %s, with args: %v", countSqlQuery, countArgs)

	var totalRecordsSlice []int

	if err := countStatement.QueryContext(ctx, mgr.db, &totalRecordsSlice); err != nil {
		return &models.ModulesPage{}, err
	}
	totalRecords := 0
	if len(totalRecordsSlice) > 0 {
		totalRecords = totalRecordsSlice[0]
	}

	return &models.ModulesPage{
		TotalRecords: totalRecords,
		Users:        userModels,
		Limit:        params.Limit,
		Offset:       params.Offset,
	}, nil
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

func (mgr *ModulesRepository) Delete(ctx context.Context, req models.ReadModulesRequest) (*models.Modules, error) {
	module, err := mgr.Read(ctx, req)
	if err != nil {
		return nil, err
	}

	conditions := table.Modules.ID.EQ(mysql.String(module.ID))
	statement := table.Modules.DELETE().WHERE(conditions)
	if _, err := statement.ExecContext(ctx, mgr.db); err != nil {
		return nil, err
	}

	return module, nil
}

func (mgr *ModulesRepository) Update(ctx context.Context, req models.ReadModulesRequest, params models.UpdateModulesParameters) (*models.Modules, error) {
	module, err := mgr.Read(ctx, req)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	updateModel := model.Modules{
		UserID:          params.UserID,
		ModuleName:      params.ModuleName,
		ModuleType:      params.ModuleType,
		ModuleDesc:      &params.Module_Desc,
		ModuleShortName: &params.ModuleShortName,
		ModulePrice:     &params.ModulePrice,
		Purchased:       &params.Purchased,
		UpdatedAt:       &now,
	}

	updateStatement := table.Modules.UPDATE(
		table.Users.UserName,
		table.Users.Email,
		table.Users.Password,
		table.Users.FirstName,
		table.Users.LastName,
		table.Users.UpdatedAt,
	).MODEL(updateModel).WHERE(table.Modules.ID.EQ(mysql.String(module.ID)))

	if _, err := updateStatement.ExecContext(ctx, mgr.db); err != nil {
		return nil, err
	}

	return mgr.Read(ctx, req)
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
