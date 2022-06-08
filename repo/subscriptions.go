package repo

import (
	"context"
	"database/sql"
	"kiddou/domain"
)

type repoSub struct {
	db *sql.DB
}

func NewRepositorySSub(db *sql.DB) *repoSub {
	return &repoSub{db: db}
}

func (r *repoSub) InsertSub(ctx context.Context, tx *sql.Tx, input domain.Subscribers) error {
	querry := `INSERT INTO subcriptions (user_id, video_id, type_subscription, expired_at, status) VALUES ($1, $2, $3, $4, $5)`
	err := tx.QueryRowContext(ctx, querry, input.UserID, input.VideoID, input.TypeSubscription, input.ExpiredAt, input.Status).Err()
	return err
}
