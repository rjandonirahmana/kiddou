package base

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

func randStringBytes(s int) (string, error) {
	b := make([]byte, s)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

type RedisAuth struct {
	db *redis.Client
}

func NewRedisAuth(db *redis.Client) *RedisAuth {
	return &RedisAuth{db: db}
}

type AuthRedis interface {
	GenerateTokenRedis(ctx context.Context, userID, email, role, name string) (string, error)
	DeleteToken(ctx context.Context, email string) error
	Authentication(ctx context.Context, token string) (AuthToken, error)
	RegenerateToken(ctx context.Context, email, name string) error
}

type AuthToken struct {
	UserID  string `redis:"user_id"`
	Role    string `redis:"role"`
	Token   string `redis:"token"`
	Expired string `redis:"expired"`
	Name    string `redis:"name"`
}

func (a *RedisAuth) GenerateTokenRedis(ctx context.Context, userID, email, role, name string) (string, error) {
	date := time.Now()

	token, err := randStringBytes(25)
	if err != nil {
		return token, err
	}
	expired := date.Add(24 * time.Hour)

	exist, err := a.db.Exists(ctx, email).Result()
	if err != nil {
		return "", err
	}
	if exist == 1 {
		oldtoken, err := a.db.HGet(ctx, email, "token").Result()
		if err != nil {
			return "", err
		}
		_, err = a.db.Del(ctx, oldtoken).Result()
		if err != nil {
			return "", err
		}
		pipe := a.db.TxPipeline()
		pipe.HSet(ctx, email, "token", token)
		pipe.HSet(ctx, email, "role", role)
		pipe.HSet(ctx, email, "expired", expired.Format(time.RFC3339))
		pipe.Set(ctx, token, email, time.Hour*24)

		_, err = pipe.Exec(ctx)
		if err != nil {
			return "", err
		}
	} else {
		pipe := a.db.TxPipeline()
		pipe.HSet(ctx, email, "user_id", userID, "role", role, "token", token, "expired", expired.Format(time.RFC3339), "name", name)
		pipe.Set(ctx, token, email, time.Hour*24)
		_, err := pipe.Exec(ctx)

		if err != nil {
			return "", err
		}
	}

	return token, nil

}

func (r *RedisAuth) DeleteToken(ctx context.Context, email string) error {
	exits, err := r.db.Exists(ctx, email).Result()
	if err != nil {
		return err
	}
	if exits == 1 {
		oldToken, err := r.db.HGet(ctx, email, "token").Result()
		if err != nil {
			return err
		}
		_, err = r.db.Del(ctx, oldToken).Result()
		if err != nil {
			return err
		}
		_, err = r.db.Del(ctx, email).Result()
		if err != nil {
			return err
		}

		return nil

	} else {
		return errors.New("token by email not found")
	}
}

func (r *RedisAuth) Authentication(ctx context.Context, token string) (res AuthToken, err error) {
	email, err := r.db.Get(ctx, token).Result()
	if err != nil {
		return res, err
	}

	var auth AuthToken
	err = r.db.HGetAll(ctx, email).Scan(&auth)

	if err != nil {
		return res, err
	}

	return auth, nil
}

func (r *RedisAuth) RegenerateToken(ctx context.Context, email, name string) error {
	exits, err := r.db.Exists(ctx, email).Result()
	if err != nil {
		return err
	}
	if exits == 1 {
		_, err := r.db.HSet(ctx, email, "name", name).Result()
		if err != nil {
			return err
		}

		return nil

	} else {
		return errors.New("failed to generate token")
	}
}
