package employeside

import (
	"context"
	"time"
)

type WriterPage struct {
	TotalRecords int      `json:"total_records"`
	Users        []Writer `json:"writer"`
	Limit        int      `json:"limit"`
	Offset       int      `json:"offset"`
}

type Writer struct {
	ID         string     `json:"id"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Email      string     `json:"email"`
	UserName   string     `json:"user_name"`
	Password   string     `json:"password"`
	IsVerified bool       `json:"is_Verified"`
	IsActive   bool       `json:"is_Active"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type ReadWriterRequest struct {
	By    string
	Value string
}

type CreateWriterRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
	IsVerified bool   `json:"is_Verified"`
	IsActive   bool   `json:"is_Active"`
}
type UpdateWriterParameters struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
	IsVerified bool   `json:"is_Verified"`
	IsActive   bool   `json:"is_Active"`
}

type WriterManger interface {
	Read(ctx context.Context, req ReadWriterRequest) (*Writer, error)
	Create(ctx context.Context, params CreateWriterRequest) (*Writer, error)
	List(ctx context.Context, params ListParameters) (WriterPage, error)
	Delete(ctx context.Context, req ReadWriterRequest) (*Writer, error)
	Update(ctx context.Context, req ReadWriterRequest, params UpdateWriterParameters) (*Writer, error)
}
