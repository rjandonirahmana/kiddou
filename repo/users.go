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
