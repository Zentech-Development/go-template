package sqlitestore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Zentech-Development/go-template/pkg/entities"
)

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

	account := &entities.Account{
		Username: username,
	}

	if err = stmt.QueryRowContext(ctx, username).Scan(&(account.ID), &(account.Password)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Account{}, &entities.ErrNotFound{}
		}
		return entities.Account{}, err
	}

	return *account, nil
}

// GetByID retrieves a stored acount by ID. If the ID is not found
// or the database query fails, GetByID will return an error.
func (s SQLiteStore) GetByID(ctx context.Context, id int64) (entities.Account, error) {
	query := `
	SELECT username, password FROM accounts WHERE id = ?;
	`

	stmt, err := s.DB.PrepareContext(ctx, query)
	if err != nil {
		return entities.Account{}, err
	}
	defer stmt.Close()

	account := &entities.Account{
		ID: id,
	}

	if err = stmt.QueryRowContext(ctx, id).Scan(&(account.Username), &(account.Password)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Account{}, &entities.ErrNotFound{}
		}
		return entities.Account{}, err
	}

	return *account, nil
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
