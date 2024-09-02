package repositories

import (
	"context"
	"database/sql"
	"errors"

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

// func NewUserRepository(db *sql.DB) *UserRepository {
// 	return &UserRepository{db: db}
// }

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
	idBytes, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}

	//	now := time.Now()

	realm := model.Users{
		ID:       idBytes,
		UserName: params.UserName,
		Email:    params.Email,
		Password: params.Password,
	}

	statement := table.Users.INSERT(table.Users.AllColumns).MODEL(realm)

	_, err = statement.ExecContext(ctx, mgr.db)
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

	conditions := table.Users.ID.EQ(mysql.String(string(user.ID)))
	statement := table.Users.DELETE().WHERE(conditions)
	if _, err := statement.ExecContext(ctx, mgr.db); err != nil {
		return nil, err
	}

	return user, nil
}

func (mgr *UserRepository) Update(ctx context.Context, req models.ReadUserRequest, params models.UpdateUserParameters) (*models.User, error) {
	realm, err := mgr.Read(ctx, req)
	if err != nil {
		return nil, err
	}

	updateModel := model.Users{
		UserName: params.UserName,
		Email:    params.Email,
		Password: params.Password,
	}

	updateStatement := table.Users.UPDATE(
		table.Users.UserName,
		table.Users.Email,
		table.Users.Password,
	).MODEL(updateModel).WHERE(table.Users.ID.EQ(mysql.String(string(realm.ID))))

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
	id, err := uuid.FromBytes(user.ID)
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:       []byte(id.String()),
		UserName: user.UserName,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}
