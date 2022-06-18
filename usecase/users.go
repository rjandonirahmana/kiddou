package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"kiddou/base"
	"kiddou/domain"
	"log"
	"strings"
	"time"
)

type usecaseUser struct {
	repoUser  domain.RepositoryUser
	authRedis base.AuthRedis
	db        *sql.DB
	secret    string
}

func NewUsecaseUser(repoUsers domain.RepositoryUser, secret string, db *sql.DB, authRedis base.AuthRedis) *usecaseUser {
	return &usecaseUser{repoUser: repoUsers, secret: secret, db: db, authRedis: authRedis}
}

func RandStringBytes(s int) (string, error) {
	b := make([]byte, s)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func (u *usecaseUser) Register(ctx context.Context, input *domain.RegisterInput) (string, error) {
	user, err := u.repoUser.GetByEmail(ctx, input.Email)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if user.UserID != "" {
		return "", errors.New("email has been used")
	}
	tx, err := u.db.Begin()
	if err != nil {
		return "", err
	}
	salt, err := RandStringBytes(10)
	if err != nil {
		return "", err
	}
	password := input.Password + salt
	h := sha256.New()
	h.Write([]byte(password))
	hashpassword := fmt.Sprintf("%x", h.Sum([]byte(u.secret)))

	user1 := &domain.Users{
		UserID:    base.GenerateUserID(),
		Name:      strings.Split(input.Email, "@")[0],
		Password:  hashpassword,
		Email:     input.Email,
		Salt:      salt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = SendEMail(ctx, "template/registration.html", user1)
	if err != nil {
		return "", err
	}

	err = u.repoUser.Insertuser(ctx, tx, user1)
	if err != nil {
		return "", err
	}

	token, err := u.authRedis.GenerateTokenRedis(ctx, user1.UserID, input.Email, "user", user1.Name)
	if err != nil {
		return "", err
	}
	tx.Commit()

	return token, nil

}

func (u *usecaseUser) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	password += user.Salt
	h := sha256.New()
	h.Write([]byte(password))
	hashpassword := fmt.Sprintf("%x", h.Sum([]byte(u.secret)))

	if hashpassword != user.Password {
		return "", errors.New("wrong password")
	}

	admin, err := u.repoUser.IsUserAdmin(ctx, user.UserID)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	role := "user"
	if admin.ID > 0 {
		role = "admin"
	}
	token, err := u.authRedis.GenerateTokenRedis(ctx, user.UserID, user.Email, role, user.Name)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *usecaseUser) LoginGoogle(ctx context.Context, users *domain.Users, googleID string) (token string, err error) {
	sosmed, err := u.repoUser.GetSosmedID(ctx, googleID, "google")
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	log.Println(googleID)
	log.Println(sosmed.ID)
	if sosmed.ID == 0 {
		user, err := u.repoUser.GetByEmail(ctx, users.Email)
		if err != nil && err != sql.ErrNoRows {
			return "", err
		}
		if user.Email != "" {
			return "", errors.New("you are login using another type of login")
		}
		salt, err := RandStringBytes(10)
		if err != nil {
			return "", err
		}

		tx, err := u.db.Begin()
		if err != nil {
			return "", err
		}

		users.UserID = base.GenerateUserID()
		users.Salt = salt
		users.CreatedAt = time.Now()
		users.UpdatedAt = time.Now()
		users.Password = "google"
		users.Name = strings.Split(users.Email, "@")[0]
		err = u.repoUser.Insertuser(ctx, tx, users)
		if err != nil {
			tx.Rollback()
			return "", err
		}

		sosmed.UserID = users.UserID
		sosmed.GoogleID = sql.NullString{String: googleID, Valid: true}
		err = u.repoUser.InsertSosmed(ctx, tx, &sosmed)
		if err != nil {
			tx.Rollback()
			return "", err
		}

		token, err := u.authRedis.GenerateTokenRedis(ctx, users.UserID, users.Email, "user", users.Name)
		if err != nil {
			return "", err
		}

		tx.Commit()
		return token, nil

	} else {
		user, err := u.repoUser.GetByUserID(ctx, sosmed.UserID)
		if err != nil {
			return "", err
		}

		admin, err := u.repoUser.IsUserAdmin(ctx, user.UserID)
		if err != nil && err != sql.ErrNoRows {
			return "", err
		}

		role := "user"
		if admin.ID > 0 {
			role = "admin"
		}

		token, err := u.authRedis.GenerateTokenRedis(ctx, user.UserID, user.Email, role, user.Name)
		if err != nil {
			return "", err
		}

		return token, nil

	}
}
