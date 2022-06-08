package repo

import "database/sql"

type repoSub struct {
	db *sql.DB
}

func NewRepositorySSub(db *sql.DB) *repoSub {
	return &repoSub{db: db}
}
