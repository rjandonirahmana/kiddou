package repo

import (
	"context"
	"database/sql"
	"errors"
	"kiddou/domain"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

func (r userRepo) Insertuser(ctx context.Context, tx *sql.Tx, user *domain.Users) error {
	querry := `INSERT INTO users (user_id, name, email, password, salt, avatar, phone_number, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	stmt, err := tx.PrepareContext(ctx, querry)
	if err != nil {
		return err
	}
	result, err := stmt.ExecContext(ctx, user.UserID, user.Name, user.Email, user.Password, user.Salt, user.Avatar, user.PhoneNumber, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if id == 0 {
		return errors.New("not inserted")
	}
	return nil

}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (res domain.Users, err error) {
	querry := `SELECT user_id, name, email, password, salt, avatar, phone_number, created_at, updated_at FROM users WHERE email = $1`

	err = r.db.QueryRowContext(ctx, querry, email).Scan(&res.UserID, &res.Name, &res.Email, &res.Password, &res.Salt, &res.Avatar, &res.PhoneNumber, &res.CreatedAt, &res.UpdatedAt)
	return
}

func (r *userRepo) GetByUserID(ctx context.Context, userID string) (res domain.Users, err error) {
	querry := `SELECT user_id, name, email, password, salt, avatar, phone_number, created_at, updated_at FROM users WHERE user_id = $1`

	err = r.db.QueryRowContext(ctx, querry, userID).Scan(&res.UserID, &res.Name, &res.Email, &res.Password, &res.Salt, &res.Avatar, &res.PhoneNumber, &res.CreatedAt, &res.UpdatedAt)
	return
}

func (r *userRepo) IsUserAdmin(ctx context.Context, userID string) (admin domain.Admin, err error) {
	querry := `SELECT id, user_id FROM admin WHERE user_id = $1`
	err = r.db.QueryRowContext(ctx, querry, userID).Scan(&admin.ID, &admin.UserID)
	return

}

func (r *userRepo) GetSosmedID(ctx context.Context, sosmedID string, sosmed string) (res domain.SocialMedia, err error) {
	queeryWHERE := ""
	if sosmed == "google" {
		queeryWHERE = "WHERE google_id = $1"
	} else if sosmed == "github" {
		queeryWHERE = "WHERE github_id = $1"
	} else {
		queeryWHERE = "WHERE facebook_id = $1"
	}
	querry := `SELECT id, user_id, google_id, facebook_id, github_id FROM sosial_media `
	querry += queeryWHERE

	err = r.db.QueryRowContext(ctx, querry, sosmed).Scan(&res.ID, &res.UserID, &res.GoogleID, &res.FacebookID, &res.GithubID)
	return
}

func (r *userRepo) InsertSosmed(ctx context.Context, tx *sql.Tx, sosmed *domain.SocialMedia) error {
	querry := `INSERT INTO sosial_media (user_id, google_id, facebook_id, github_id) VALUES ($1, $2, $3, $4) RETURNING id`

	err := tx.QueryRowContext(ctx, querry, sosmed.UserID, sosmed.GoogleID, sosmed.FacebookID, sosmed.GithubID).Scan(&sosmed.ID)
	return err
}
