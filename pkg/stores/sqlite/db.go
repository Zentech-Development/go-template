package sqlitestore

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	DB *sql.DB
}

func NewSQLiteStore(dbPath string) SQLiteStore {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	store := SQLiteStore{
		DB: db,
	}

	if err := store.createTables(); err != nil {
		log.Fatalf("Failed to create accounts table %w\n", err)
	}

	return store
}
