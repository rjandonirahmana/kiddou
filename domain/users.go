package domain

import (
	"context"
	"database/sql"
	"time"
)

type Users struct {
	UserID      string
	Name        string
	Email       string
	Password    string
	Salt        string
	Avatar      sql.NullString
	PhoneNumber sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Admin struct {
	ID     int
	UserID string
}

type RegisterInput struct {
	Email           string `form:"email" validate:"required,email"`
	Password        string `form:"password" validate:"required,min=8,max=32,alphanum"`
	ConfirmPassword string `form:"confirm_password" validate:"required,eqfield=Password"`
}

type LoginInput struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
}

type RepositoryUser interface {
	Insertuser(ctx context.Context, tx *sql.Tx, user *Users) error
	GetByEmail(ctx context.Context, email string) (res Users, err error)
	IsUserAdmin(ctx context.Context, userID string) (admin Admin, err error)
}

type UsecaseUser interface {
	Register(ctx context.Context, input *RegisterInput) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
}
