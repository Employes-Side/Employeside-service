package employeside

import (
	"context"
	"time"
)

type Writer struct {
	ID         string    `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	UserName   string    `json:"user_name"`
	Password   string    `json:"password"`
	IsVerified bool      `json:"is_verified"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  time.Time `json:"created_by"`
}

type ReadWriterRequest struct {
	By    string
	Value string
}

type CreateWriterRequest struct {
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	UserName   string    `json:"user_name"`
	Password   string    `json:"password"`
	IsVerified bool      `json:"is_verified"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  time.Time `json:"created_by"`
}

type WriterManger interface {
	Read(ctx context.Context, req ReadWriterRequest) (*Writer, error)
	Create(ctx context.Context, params CreateWriterRequest) (*Writer, error)
}
