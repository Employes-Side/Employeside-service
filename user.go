package employeside

import (
	"context"
)

type User struct {
	ID        string `json:"id"`
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
type ReadUserRequest struct {
	By    string
	Value string
}

type ListUserParameters struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	SortBy string `json:"sort_by"`
	Order  string `json:"order"`
}

type Page struct {
	TotalRecords int    `json:"total_records"`
	Users        []User `json:"users"`
	Limit        int    `json:"limit"`
	Offset       int    `json:"offset"`
}

type UpdateUserParameters struct {
	UserName string `json:"user_name" validate:"required,max=150"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserParameters struct {
	UserName string `json:"user_name" validate:"required,max=150"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserManager interface {
	Read(ctx context.Context, req ReadUserRequest) (*User, error)
	List(ctx context.Context, parms ListUserParameters) (Page, error)
	Create(ctx context.Context, params CreateUserParameters) (*User, error)
	Delete(ctx context.Context, req ReadUserRequest) (*User, error)
	Update(ctx context.Context, req ReadUserRequest, params UpdateUserParameters) (*User, error)
}
