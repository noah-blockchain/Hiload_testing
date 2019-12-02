package dal

import (
	"github.com/jmoiron/sqlx"
	"github.com/noah-blockchain/Hiload_testing/internal/app"
)

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) app.Repo {
	return &repo{
		db: db,
	}
}
