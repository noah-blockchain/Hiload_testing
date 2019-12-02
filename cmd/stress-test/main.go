package main

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/noah-blockchain/Hiload_testing/internal/app"
	"github.com/noah-blockchain/Hiload_testing/internal/dal"
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
	_, _ = statement.Exec()

	return sqlx.NewDb(db, "sqlite3"), nil
}

func main() {
	db, err := openAndCreateDB()
	if err != nil {
		log.Panicln(err)
	}

	repo := dal.New(db)
	appl := app.New(repo)
	if err := appl.CheckAndCreateWallets(); err != nil {
		log.Panicln(err)
	}
}
