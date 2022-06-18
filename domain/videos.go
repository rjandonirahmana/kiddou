package domain

import (
	"context"
	"database/sql"
	"time"
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

type Subscribers struct {
	ID               int
	UserID           string
	VideoID          int
	TypeSubscription string
	SubscribeAT      time.Time
	ExpiredAt        time.Time
	Status           string
}

type SubsriptionStatusDTO struct {
	VideoName        string `json:"video_name"`
	TypeSubscription string `json:"type_subscription"`
	ExpiredAt        string `json:"expired_at"`
	SubscribeAT      string `json:"subscribe_at"`
	Status           string `json:"status"`
}

type InserVideosRequest struct {
	CategoryID   int    `form:"category_id" validate:"required"`
	Name         string `form:"name" validate:"required"`
	Descriptions string `form:"desc" validate:"required"`
	Price        string `form:"price" validate:"required"`
}

type RepoVideos interface {
	InsertVideos(ctx context.Context, tx *sql.Tx, video *Videos) error
	GetByCategory(ctx context.Context, categoryID int) (res []Videos, err error)
	GetAllCategories(ctx context.Context) (res []Categories, err error)
	UpdateVideo(ctx context.Context, tx *sql.Tx, video *Videos) error
	GetByID(ctx context.Context, ID int) (res Videos, err error)
}

type UsecaseVideos interface {
	InsertVideos(ctx context.Context, videos []byte, input *InserVideosRequest) error
	SubscribtionVideo(ctx context.Context, userID string, videoID int) error
	SubsribesStatus(ctx context.Context, userID string, videoID int) (res SubsriptionStatusDTO, err error)
	GetByCategory(ctx context.Context, categoryID int) (res []Videos, err error)
}

type SubscribersRepo interface {
	InsertSub(ctx context.Context, tx *sql.Tx, input *Subscribers) error
	UpdateSub(ctx context.Context, tx *sql.Tx, input *Subscribers) error
	GetByVideoID(ctx context.Context, ID int, userID string) (res Subscribers, err error)
	GetByUserID(ctx context.Context, userID string) (res []Subscribers, err error)
}
