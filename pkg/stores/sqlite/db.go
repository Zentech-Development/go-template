package sqlitestore

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	DB       *sql.DB
	Accounts *AccountStore
}

func NewSQLiteStore(dbPath string) SQLiteStore {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	accountStore := newAccountStore(db)

	store := SQLiteStore{
		DB:       db,
		Accounts: accountStore,
	}

	return store
}
