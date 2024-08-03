package sqlitestore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Zentech-Development/go-template/pkg/entities"
)

func (s SQLiteStore) createTables() error {
	createTableStatement := `
	CREATE TABLE IF NOT EXISTS accounts (
		id       INTEGER NOT NULL PRIMARY KEY, 
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	`
	_, err := s.DB.Exec(createTableStatement)
	return err
}

// GetUserByUsername retrieves a stored user by username. If the username is not found
// or the database query fails, GetUserByUsername will return an error.
func (s SQLiteStore) GetByUsername(ctx context.Context, username string) (entities.Account, error) {
	query := `
	SELECT id, password FROM accounts WHERE username = ?;
	`

	stmt, err := s.DB.PrepareContext(ctx, query)
	if err != nil {
		return entities.Account{}, err
	}
	defer stmt.Close()

	var id int64
	var password string

	if err = stmt.QueryRowContext(ctx, username).Scan(&id, &password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Account{}, &entities.ErrNotFound{}
		}
		return entities.Account{}, err
	}

	return entities.Account{
		Username: username,
		Password: password,
	}, nil
}

// Create inserts a new account into the database. It returns the inserted account or an error.
func (s SQLiteStore) Create(ctx context.Context, account entities.AccountCreateInput) (entities.Account, error) {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return entities.Account{}, err
	}

	query := `
	INSERT INTO accounts(username, password) VALUES (?, ?)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return entities.Account{}, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, account.Username, account.Password)
	if err != nil {
		return entities.Account{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return entities.Account{}, err
	}

	if rowsAffected != 1 {
		return entities.Account{}, errors.New("failed to insert account")
	}

	newAccountID, err := result.LastInsertId()
	if err != nil {
		return entities.Account{}, err
	}

	if err = tx.Commit(); err != nil {
		return entities.Account{}, err
	}

	savedAccount := entities.Account{
		ID:       newAccountID,
		Username: account.Username,
		Password: account.Password,
	}

	return savedAccount, nil
}
