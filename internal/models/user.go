package models

import (
	"context"
)

type User struct {
	ID       []byte `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type ReadUserRequest struct {
	By    string
	Value string
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
	Create(ctx context.Context, params CreateUserParameters) (*User, error)
	Delete(ctx context.Context, req ReadUserRequest) (*User, error)
	Update(ctx context.Context, req ReadUserRequest, params UpdateUserParameters) (*User, error)
}
