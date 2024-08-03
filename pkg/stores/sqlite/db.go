package sqlitestore

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	DB *sql.DB
}

type Opts struct {
	DBPath         string
	MigrationsPath string
}

func NewSQLiteStore(opts *Opts) SQLiteStore {
	db, err := sql.Open("sqlite3", opts.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	store := SQLiteStore{
		DB: db,
	}

	migrationDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}
	migrator, err := migrate.NewWithDatabaseInstance(opts.MigrationsPath, "sqlite3", migrationDriver)
	if err != nil {
		log.Fatal(err)
	}

	if err = migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("SQLite migrations made no changes")
		} else {
			log.Fatal("SQLite migrations failed: ", err)
		}
	}

	log.Println("SQLite database connected and migrated successfully")

	return store
}
