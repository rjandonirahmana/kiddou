package repo

import (
	"context"
	"database/sql"
	"kiddou/domain"
)

type repoVideos struct {
	db *sql.DB
}

func NewRepositoryVideos(db *sql.DB) *repoVideos {
	return &repoVideos{db: db}
}

func (r *repoVideos) InsertVideos(ctx context.Context, tx *sql.Tx, video *domain.Videos) error {
	querry := `INSERT INTO videos (category_id, name, descriptions, price, url, subscribers) VALUES ($1, $2, $3, $4, $5, $6)`
	stmt, err := tx.PrepareContext(ctx, querry)
	if err != nil {
		return err
	}
	rslt, err := stmt.ExecContext(ctx, video.CategoryID, video.Name, video.Descriptions, video.Price, video.Url, video.Subscribers)
	if err != nil {
		return err
	}
	last, err := rslt.LastInsertId()
	if err != nil {
		return err
	}

	video.ID = int(last)
	return nil
}

func (r *repoVideos) GetByCategory(ctx context.Context, categoryID int) (res []*domain.Videos, err error) {
	querry := `SELECT id, category_id, name, descriptions, price, url, subscribers FROM videos WHERE category_id = $1`

	row, err := r.db.QueryContext(ctx, querry, categoryID)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var video domain.Videos
		err = row.Scan(
			&video.ID,
			&video.CategoryID,
			&video.Name,
			&video.Descriptions,
			&video.Price,
			&video.Url,
			&video.Subscribers,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, &video)

	}

	return
}
