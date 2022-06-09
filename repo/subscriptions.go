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

func (r *repoSub) InsertSub(ctx context.Context, tx *sql.Tx, input *domain.Subscribers) error {
	querry := `INSERT INTO subcriptions (user_id, video_id, type_subscription, subscribe_at, expired_at, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := tx.QueryRowContext(ctx, querry, input.UserID, input.VideoID, input.TypeSubscription, input.SubscribeAT, input.ExpiredAt, input.Status).Scan(&input.VideoID)
	return err
}

func (r *repoSub) UpdateSub(ctx context.Context, tx *sql.Tx, input *domain.Subscribers) error {
	querry := `UPDATE subcriptions SET user_id = $1, video_id = $2, type_subscription = $3, subscribe_at = $4, expired_at = $5, status = $6 WHERE id = $7`
	err := tx.QueryRowContext(ctx, querry, input.UserID, input.VideoID, input.TypeSubscription, input.SubscribeAT, input.ExpiredAt, input.Status, input.ID).Err()
	return err
}

func (r *repoSub) GetByVideoID(ctx context.Context, ID int, userID string) (res domain.Subscribers, err error) {
	querry := `SELECT id, user_id, video_id, type_subscription, subscribe_at, expired_at, status FROM subcriptions WHERE video_id = $1 AND user_id = $2`

	err = r.db.QueryRowContext(ctx, querry, ID, userID).Scan(&res.ID, &res.UserID, &res.VideoID, &res.TypeSubscription, &res.SubscribeAT, &res.ExpiredAt, &res.Status)
	return
}

func (r *repoSub) GetByUserID(ctx context.Context, userID string) (res []domain.Subscribers, err error) {
	querry := `SELECT id, user_id, video_id, type_subscription, subscribe_at, expired_at, status FROM subcriptions WHERE user_id = $1`

	row, err := r.db.QueryContext(ctx, querry, userID)
	if err != nil {
		return res, err
	}
	defer row.Close()
	for row.Next() {
		var sub domain.Subscribers
		row.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.VideoID,
			&sub.TypeSubscription,
			&sub.SubscribeAT,
			&sub.ExpiredAt,
			&sub.Status,
		)
		if err != nil {
			return res, err
		}

		res = append(res, sub)

	}

	return
}
