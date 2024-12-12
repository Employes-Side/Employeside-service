package employeside

import (
	"context"
)

type Blogs struct {
	ID          string `json:"id"`
	BlogTitle   string `json:"blog_title"`
	BlogContent string `json:"blog_content"`
	Status      string `json:"status"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type ReadBlogRequest struct {
	By    string
	Value string
}

type ListBlogsParameters struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	SortBy string `json:"sort_by"`
	Order  string `json:"order"`
}

type PageBlocks struct {
	TotalRecords int     `json:"total_records"`
	Blogs        []Blogs `json:"blogs"`
	Limit        int     `json:"limit"`
	Offset       int     `json:"offset"`
}

type CreatBlogParameters struct {
	BlogTitle   string `json:"blog_title"`
	BlogContent string `json:"blog_content"`
	Status      string `json:"status"`
}

type UpdateBlogParameters struct {
	BlogTitle   string `json:"blog_title"`
	BlogContent string `json:"blog_content"`
	Status      string `json:"status"`
}

type BlogManager interface {
	Read(ctx context.Context, req ReadBlogRequest) (*Blogs, error)
	Create(ctx context.Context, params CreatBlogParameters) (*Blogs, error)
	List(ctx context.Context, params ListBlogsParameters) (*Blogs, error)
	Update(ctx context.Context, req ReadBlogRequest, params UpdateBlogParameters) (*Blogs, error)
	Delete(ctx context.Context, req ReadBlogRequest) (*Blogs, error)
}
