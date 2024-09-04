package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/Employes-Side/employee-side/internal/models"
	"github.com/Employes-Side/employee-side/internal/users/model"
	"github.com/Employes-Side/employee-side/internal/users/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

func NewManager(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

type UserRepository struct {
	db *sql.DB
}

func (mgr *UserRepository) List(ctx context.Context, params models.ListUserParameters) (*models.Page, error) {
	var orderExpr mysql.OrderByClause
	if params.Order == "desc" {
		orderExpr = table.Users.UserName.DESC()
	} else {
		orderExpr = table.Users.UserName.ASC()
	}

	statement := table.Users.
		SELECT(table.Users.AllColumns).
		ORDER_BY(orderExpr).
		LIMIT(int64(params.Limit)).
		OFFSET(int64(params.Offset))

	sqlQuery, args := statement.Sql()
	log.Printf("Executing SQL query: %s, with args: %v", sqlQuery, args)

	var users []model.Users

	if err := statement.QueryContext(ctx, mgr.db, &users); err != nil {
		if err == qrm.ErrNoRows {
			return &models.Page{}, nil
		}
		return &models.Page{}, err
	}

	userModels := make([]models.User, len(users))
	for i, u := range users {
		userModel, err := convertoDBModel(u)
		if err != nil {
			return &models.Page{}, err
		}
		userModels[i] = *userModel
	}

	countStatement := table.Users.SELECT(mysql.COUNT(table.Users.ID))

	countSqlQuery, countArgs := countStatement.Sql()
	log.Printf("Executing COUNT SQL query: %s, with args: %v", countSqlQuery, countArgs)

	var totalRecordsSlice []int

	if err := countStatement.QueryContext(ctx, mgr.db, &totalRecordsSlice); err != nil {
		return &models.Page{}, err
	}
	totalRecords := 0
	if len(totalRecordsSlice) > 0 {
		totalRecords = totalRecordsSlice[0]
	}

	return &models.Page{
		TotalRecords: totalRecords,
		Users:        userModels,
		Limit:        params.Limit,
		Offset:       params.Offset,
	}, nil
}

func (mgr *UserRepository) Read(ctx context.Context, req models.ReadUserRequest) (*models.User, error) {
	conditions, err := mgr.buildReadClause(req)
	if err != nil {
		return nil, err
	}

	statement := table.Users.SELECT(table.Users.AllColumns).WHERE(conditions)

	var user model.Users
	if err := statement.QueryContext(ctx, mgr.db, &user); err != nil {
		if err == qrm.ErrNoRows {
			return nil, errors.New("user not found")

		}
		return nil, err
	}

	return convertoDBModel(user)

}
func (mgr *UserRepository) Create(ctx context.Context, params models.CreateUserParameters) (*models.User, error) {
	id := uuid.New()

	now := time.Now()

	realm := model.Users{
		ID:        id.String(),
		UserName:  params.UserName,
		Email:     params.Email,
		Password:  params.Password,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	statement := table.Users.INSERT(table.Users.AllColumns).MODEL(realm)

	_, err := statement.ExecContext(ctx, mgr.db)
	if err != nil {
		return nil, err
	}
	return mgr.Read(ctx, models.ReadUserRequest{By: "id", Value: id.String()})

}

func (mgr *UserRepository) Delete(ctx context.Context, req models.ReadUserRequest) (*models.User, error) {
	user, err := mgr.Read(ctx, req)
	if err != nil {
		return nil, err
	}

	conditions := table.Users.ID.EQ(mysql.String(user.ID))
	statement := table.Users.DELETE().WHERE(conditions)
	if _, err := statement.ExecContext(ctx, mgr.db); err != nil {
		return nil, err
	}

	return user, nil
}

func (mgr *UserRepository) Update(ctx context.Context, req models.ReadUserRequest, params models.UpdateUserParameters) (*models.User, error) {
	user, err := mgr.Read(ctx, req)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	updateModel := model.Users{
		UserName:  params.UserName,
		Email:     params.Email,
		Password:  params.Password,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	updateStatement := table.Users.UPDATE(
		table.Users.UserName,
		table.Users.Email,
		table.Users.Password,
		table.Users.CreatedAt,
		table.Users.UpdatedAt,
	).MODEL(updateModel).WHERE(table.Users.ID.EQ(mysql.String(user.ID)))

	if _, err := updateStatement.ExecContext(ctx, mgr.db); err != nil {
		return nil, err
	}

	return mgr.Read(ctx, req)
}

func (mgr *UserRepository) buildReadClause(req models.ReadUserRequest) (mysql.BoolExpression, error) {
	switch req.By {
	case "id":
		return table.Users.ID.EQ(mysql.String(req.Value)), nil
	case "name":
		return table.Users.UserName.EQ(mysql.String(req.Value)), nil
	default:
		return nil, errors.New("by should be one of 'id' or 'name'")
	}

}

func convertoDBModel(user model.Users) (*models.User, error) {

	return &models.User{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.UnixMilli(),
		UpdatedAt: user.UpdatedAt.UnixMilli(),
	}, nil
}
