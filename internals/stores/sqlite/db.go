package sqlitestore

import (
	"database/sql"
	"errors"

	"github.com/Zentech-Development/go-template/pkg/logger"
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
	if opts.DBPath == "" || opts.MigrationsPath == "" {
		logger.L.Fatal().Msg("Missing required SQLITE_OPTS in config")
	}

	db, err := sql.Open("sqlite3", opts.DBPath)
	if err != nil {
		logger.L.Fatal().Err(err)
	}

	store := SQLiteStore{
		DB: db,
	}

	migrationDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		logger.L.Fatal().Err(err)
	}
	migrator, err := migrate.NewWithDatabaseInstance(opts.MigrationsPath, "sqlite3", migrationDriver)
	if err != nil {
		logger.L.Fatal().Err(err)
	}

	if err = migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.L.Info().Msg("SQLite migrations made no changes")
		} else {
			logger.L.Fatal().Err(err).Msg("SQLite migrations failed")
		}
	}

	if err := db.Ping(); err != nil {
		logger.L.Fatal().Err(err).Msg("SQLite connection failed")
	}

	logger.L.Info().Msg("SQLite database connected successfully")

	return store
}
