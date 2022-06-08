package domain

import (
	"context"
	"database/sql"
)

type Videos struct {
	ID           int
	CategoryID   int
	Name         string
	Descriptions string
	Price        string
	Url          string
	Subscribers  int
}

type InserVideosRequest struct {
	CategoryID   int    `form:"category_id" validate:"required"`
	Name         string `form:"name" validate:"required"`
	Descriptions string `form:"desc" validate:"required"`
	Price        string `form:"price" validate:"required"`
}

type RepoVideos interface {
	InsertVideos(ctx context.Context, tx *sql.Tx, video *Videos) error
	GetByCategory(ctx context.Context, categoryID int) (res []*Videos, err error)
}
