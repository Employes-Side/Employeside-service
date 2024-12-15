package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	models "github.com/Employes-Side/employee-side"
	"github.com/Employes-Side/employee-side/generated/employeside/model"
	"github.com/Employes-Side/employee-side/generated/employeside/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

func NewWriterManager(db *sql.DB) *WriterRepository {
	return &WriterRepository{db}
}

type WriterRepository struct {
	db *sql.DB
}

func (mgr *WriterRepository) List(ctx context.Context, params models.ListParameters) (*models.WriterPage, error) {
	var orderExpr mysql.OrderByClause
	if params.Order == "desc" {
		orderExpr = table.Writer.UserName.DESC()
	} else {
		orderExpr = table.Writer.UserName.ASC()
	}

	statement := table.Writer.
		SELECT(table.Writer.AllColumns).
		ORDER_BY(orderExpr).
		LIMIT(int64(params.Limit)).
		OFFSET(int64(params.Offset))

	sqlQuery, args := statement.Sql()
	log.Printf("Executing SQL query: %s, with args: %v", sqlQuery, args)

	var writer []model.Writer

	if err := statement.QueryContext(ctx, mgr.db, &writer); err != nil {
		if err == qrm.ErrNoRows {
			return &models.WriterPage{}, nil
		}
		return &models.WriterPage{}, err
	}

	writerModels := make([]models.Writer, len(writer))
	for i, u := range writer {
		writerModel, err := convertoWriterDBModel(u)
		if err != nil {
			return &models.WriterPage{}, err
		}
		writerModels[i] = *writerModel
	}

	countStatement := table.Writer.SELECT(mysql.COUNT(table.Writer.ID))

	countSqlQuery, countArgs := countStatement.Sql()
	log.Printf("Executing COUNT SQL query: %s, with args: %v", countSqlQuery, countArgs)

	var totalRecordsSlice []int

	if err := countStatement.QueryContext(ctx, mgr.db, &totalRecordsSlice); err != nil {
		return &models.WriterPage{}, err
	}
	totalRecords := 0
	if len(totalRecordsSlice) > 0 {
		totalRecords = totalRecordsSlice[0]
	}

	return &models.WriterPage{
		TotalRecords: totalRecords,
		Users:        writerModels,
		Limit:        params.Limit,
		Offset:       params.Offset,
	}, nil
}

func (mgr *WriterRepository) Read(ctx context.Context, req models.ReadWriterRequest) (*models.Writer, error) {
	conditions, err := mgr.buildReadClause(req)
	if err != nil {
		return nil, err
	}

	statement := table.Writer.SELECT(table.Writer.AllColumns).WHERE(conditions)

	var writer model.Writer
	if err := statement.QueryContext(ctx, mgr.db, &writer); err != nil {
		if err == qrm.ErrNoRows {
			return nil, errors.New("user not found")

		}
		return nil, err
	}

	return convertoWriterDBModel(writer)
}

func (mgr *WriterRepository) Create(ctx context.Context, params models.CreateWriterRequest) (*models.Writer, error) {
	id := uuid.New()

	now := time.Now()

	realm := model.Writer{
		ID:         id.String(),
		FirstName:  params.FirstName,
		LastName:   params.LastName,
		IsVerified: params.IsVerified,
		IsActive:   params.IsActive,
		UserName:   params.UserName,
		Email:      params.Email,
		Password:   params.Password,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	statement := table.Writer.INSERT(table.Writer.AllColumns).MODEL(realm)

	_, err := statement.ExecContext(ctx, mgr.db)
	if err != nil {
		return nil, err
	}
	return mgr.Read(ctx, models.ReadWriterRequest{By: "id", Value: id.String()})

}

func (mgr *WriterRepository) Delete(ctx context.Context, req models.ReadWriterRequest) (*models.Writer, error) {
	writer, err := mgr.Read(ctx, req)
	if err != nil {
		return nil, err
	}

	conditions := table.Writer.ID.EQ(mysql.String(writer.ID))
	statement := table.Writer.DELETE().WHERE(conditions)
	if _, err := statement.ExecContext(ctx, mgr.db); err != nil {
		return nil, err
	}

	return writer, nil
}

func (mgr *WriterRepository) Update(ctx context.Context, req models.ReadWriterRequest, params models.UpdateWriterParameters) (*models.Writer, error) {
	writer, err := mgr.Read(ctx, req)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	updateModel := model.Writer{
		UserName:   params.UserName,
		Email:      params.Email,
		Password:   params.Password,
		FirstName:  params.FirstName,
		LastName:   params.LastName,
		IsVerified: params.IsVerified,
		IsActive:   params.IsActive,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	updateStatement := table.Writer.UPDATE(
		table.Writer.UserName,
		table.Writer.Email,
		table.Writer.Password,
		table.Writer.FirstName,
		table.Writer.LastName,
		table.Writer.IsActive,
		table.Writer.IsVerified,
		table.Writer.CreatedAt,
		table.Writer.UpdatedAt,
	).MODEL(updateModel).WHERE(table.Writer.ID.EQ(mysql.String(writer.ID)))

	if _, err := updateStatement.ExecContext(ctx, mgr.db); err != nil {
		return nil, err
	}

	return mgr.Read(ctx, req)
}

func (mgr *WriterRepository) buildReadClause(req models.ReadWriterRequest) (mysql.BoolExpression, error) {
	switch req.By {
	case "id":
		return table.Writer.ID.EQ(mysql.String(req.Value)), nil
	case "user_name":
		return table.Writer.UserName.EQ(mysql.String(req.Value)), nil
	default:
		return nil, errors.New("by should be one of 'id' or 'username'")
	}
}

func convertoWriterDBModel(writer model.Writer) (*models.Writer, error) {

	return &models.Writer{
		ID:         writer.ID,
		FirstName:  writer.FirstName,
		LastName:   writer.LastName,
		Email:      writer.Email,
		UserName:   writer.UserName,
		Password:   writer.Password,
		IsVerified: writer.IsVerified,
		IsActive:   writer.IsActive,
		CreatedAt:  writer.CreatedAt,
		UpdatedAt:  writer.UpdatedAt,
	}, nil
}
