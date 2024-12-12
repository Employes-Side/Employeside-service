package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	models "github.com/Employes-Side/employee-side"
	"github.com/Employes-Side/employee-side/generated/users/model"
	"github.com/Employes-Side/employee-side/generated/users/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

func NewBlogManager(db *sql.DB) *BlogRepository {
	return &BlogRepository{db}
}

type BlogRepository struct {
	db *sql.DB
}

func (mgr *BlogRepository) List(ctx context.Context, params models.ListBlogsParameters) (*models.PageBlocks, error) {

	var orderExpr mysql.OrderByClause
	if params.Order == "desc" {
		orderExpr = table.Blogs.BlogTitle.DESC()
	} else {
		orderExpr = table.Blogs.BlogTitle.ASC()
	}
	statement := table.Blogs.SELECT(table.Blogs.AllColumns).ORDER_BY(orderExpr).LIMIT(int64(params.Limit)).OFFSET(int64(params.Offset))

	sqlQuery, args := statement.Sql()
	log.Printf("Executing SQL query: %s, with args: %v", sqlQuery, args)

	var blogs []model.Blogs

	if err := statement.QueryContext(ctx, mgr.db, &blogs); err != nil {
		if err == qrm.ErrNoRows {
			return &models.PageBlocks{}, nil
		}
	}

	blogModels := make([]models.Blogs, len(blogs))
	for i, b := range blogs {
		blogModel, err := convertToBlogsDBModel(b)
		if err != nil {
			return &models.PageBlocks{}, err
		}
		blogModels[i] = *blogModel
	}

	countStatement := table.Blogs.SELECT(mysql.COUNT(table.Blogs.ID))

	countSqlQuery, countArgs := countStatement.Sql()
	log.Printf("Executing count sql query with: %s, with args %v", countSqlQuery, countArgs)

	var totalRecordsSlice []int

	if err := countStatement.QueryContext(ctx, mgr.db, &totalRecordsSlice); err != nil {
		return &models.PageBlocks{}, err
	}
	totalRecords := 0
	if len(totalRecordsSlice) > 0 {
		totalRecords = totalRecordsSlice[0]
	}

	return &models.PageBlocks{
		TotalRecords: totalRecords,
		Blogs:        blogModels,
		Limit:        params.Limit,
		Offset:       params.Offset,
	}, nil
}

func (mgr *BlogRepository) Read(ctx context.Context, req models.ReadBlogRequest) (*models.Blogs, error) {
	conditions, err := mgr.buildReadClause(req)
	if err != nil {
		return nil, err
	}

	statement := table.Blogs.SELECT(table.Blogs.AllColumns).WHERE(conditions)

	var blog model.Blogs
	if err := statement.QueryContext(ctx, mgr.db, &blog); err != nil {
		if err == qrm.ErrNoRows {
			return nil, errors.New("blog not found")
		}
		return nil, err
	}
	return convertToBlogsDBModel(blog)
}

func (mgr *BlogRepository) Create(ctx context.Context, params models.CreatBlogParameters) (*models.Blogs, error) {
	id := uuid.New()

	now := time.Now()

	realm := model.Blogs{
		ID:          stringPtr(id.String()),
		BlogTitle:   &params.BlogTitle,
		BlogContent: &params.BlogContent,
		Status:      &params.Status,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	statement := table.Blogs.INSERT(table.Blogs.AllColumns).MODEL(realm)

	_, err := statement.ExecContext(ctx, mgr.db)
	if err != nil {
		return nil, err
	}

	return mgr.Read(ctx, models.ReadBlogRequest{By: "id", Value: id.String()})
}

func (mgr *BlogRepository) Delete(ctx context.Context, req models.ReadBlogRequest) (*models.Blogs, error) {
	blog, err := mgr.Read(ctx, req)
	if err != nil {
		return nil, err
	}

	conditions := table.Blogs.ID.EQ(mysql.String(blog.ID))
	statement := table.Blogs.DELETE().WHERE(conditions)
	if _, err := statement.ExecContext(ctx, mgr.db); err != nil {
		return nil, err
	}

	return blog, nil
}

func (mgr *BlogRepository) Update(ctx context.Context, req models.ReadBlogRequest, params models.UpdateBlogParameters) (*models.Blogs, error) {
	blog, err := mgr.Read(ctx, req)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	updateModel := model.Blogs{
		BlogTitle:   &params.BlogTitle,
		BlogContent: &params.BlogContent,
		Status:      &params.Status,
		UpdatedAt:   &now,
	}

	updateStatement := table.Blogs.UPDATE(
		table.Blogs.BlogContent,
		table.Blogs.BlogTitle,
		table.Blogs.Status,
		table.Users.UpdatedAt,
	).MODEL(updateModel).WHERE(table.Blogs.ID.EQ(mysql.String(blog.ID)))

	if _, err := updateStatement.ExecContext(ctx, mgr.db); err != nil {
		return nil, err
	}

	return mgr.Read(ctx, req)
}

func (mgr *BlogRepository) buildReadClause(req models.ReadBlogRequest) (mysql.BoolExpression, error) {
	switch req.By {
	case "id":
		return table.Blogs.ID.EQ(mysql.String(req.Value)), nil
	case "blog_title":
		return table.Blogs.BlogTitle.EQ(mysql.String(req.Value)), nil
	default:
		return nil, errors.New("by should be one of ID or TITLe")
	}
}
func convertToBlogsDBModel(blog model.Blogs) (*models.Blogs, error) {
	return &models.Blogs{
		ID:          *blog.ID,
		BlogTitle:   *blog.BlogTitle,
		BlogContent: *blog.BlogContent,
		Status:      *blog.Status,
		CreatedAt:   blog.CreatedAt.UnixMilli(),
		UpdatedAt:   blog.UpdatedAt.UnixMilli(),
	}, nil
}

func stringPtr(s string) *string {
	return &s
}
