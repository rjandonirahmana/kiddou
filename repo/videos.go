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
	querry := `INSERT INTO videos (category_id, name, descriptions, price, url, subscribers) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err := tx.QueryRowContext(ctx, querry, video.CategoryID, video.Name, video.Descriptions, video.Price, video.Url, video.Subscribers).Scan(&video.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repoVideos) UpdateVideo(ctx context.Context, tx *sql.Tx, video *domain.Videos) error {
	querry := `UPDATE videos SET category_id = $1, name = $2, descriptions = $3, price = $4, url = $5, subscribers = $6 WHERE id = $7`

	err := tx.QueryRowContext(ctx, querry, video.CategoryID, video.Name, video.Descriptions, video.Price, video.Url, video.Subscribers, video.ID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *repoVideos) GetByCategory(ctx context.Context, categoryID int) (res []domain.Videos, err error) {
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
		res = append(res, video)

	}

	return
}

func (r *repoVideos) GetAllCategories(ctx context.Context) (res []domain.Categories, err error) {
	querry := `SELECT id, name FROM categories`

	row, err := r.db.QueryContext(ctx, querry)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var cat domain.Categories
		err = row.Scan(
			&cat.ID,
			&cat.Name,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, cat)

	}

	return
}

func (r *repoVideos) GetByID(ctx context.Context, ID int) (res domain.Videos, err error) {
	querry := `SELECT id, category_id, name, descriptions, price, url, subscribers FROM videos WHERE id = $1`

	err = r.db.QueryRowContext(ctx, querry, ID).Scan(
		&res.ID,
		&res.CategoryID,
		&res.Name,
		&res.Descriptions,
		&res.Price,
		&res.Url,
		&res.Subscribers,
	)
	if err != nil {
		return res, err
	}
	return

}
