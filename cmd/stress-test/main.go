package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/noah-blockchain/Hiload_testing/internal/app"
	"github.com/noah-blockchain/Hiload_testing/internal/dal"
	"github.com/noah-blockchain/Hiload_testing/internal/env"
)

const (
	dbFolderPath = "db"
	dbPath       = dbFolderPath + "/db.sqlite"
)

const SqlCommand = `
	CREATE TABLE IF NOT EXISTS wallets (
		id INTEGER 	PRIMARY KEY AUTOINCREMENT,
		address 	TEXT NOT NULL,
		seed_phrase TEXT NOT NULL,
		mnemonic	TEXT NOT NULL,
		private_key TEXT NOT NULL,
		amount 		NUMERIC(70) DEFAULT 0,
		status 		BOOL
	)
`

func openAndCreateDB() (*sqlx.DB, error) {
	if err := os.MkdirAll(dbFolderPath, 0774); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	statement, _ := db.Prepare(SqlCommand)
	_, err = statement.Exec()
	if err != nil {
		return nil, err
	}

	return sqlx.NewDb(db, "sqlite3"), nil
}

func main() {
	db, err := openAndCreateDB()
	if err != nil {
		log.Panicln(err)
	}

	repo := dal.New(db)
	per := env.GetEnvAsInt("PER_SEC", 1)
	appl := app.New(repo, app.RateLimiter{
		Freq: env.GetEnvAsInt("FREQ", 150),
		Per:  time.Duration(per) * time.Second,
	})

	if env.GetEnvAsBool("CREATE_WALLETS", false) {
		if err := appl.CreateWallets(); err != nil {
			log.Panicln(err)
		}
	}

	if err := appl.Start(); err != nil {
		log.Panicln(err)
	}
}
