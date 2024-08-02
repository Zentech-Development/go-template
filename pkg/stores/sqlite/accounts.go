package sqlitestore

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"math"
	"math/big"

	"github.com/Zentech-Development/go-template/pkg/entities"
)

func (s SQLiteStore) createTables() error {
	createTableStatement := `
	CREATE TABLE IF NOT EXISTS accounts (
		id      INTEGER NOT NULL PRIMARY KEY, 
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
	var isAdmin bool

	if err = stmt.QueryRowContext(ctx, username).Scan(&id, &password, &isAdmin); err != nil {
		return entities.Account{}, err
	}

	return entities.Account{
		Username: username,
		Password: password,
		IsAdmin:  isAdmin,
	}, nil
}

// Create inserts a new account into the database. It returns the inserted account or an error.
func (s SQLiteStore) Create(ctx context.Context, account entities.Account) (entities.Account, error) {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return entities.Account{}, err
	}

	id, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return entities.Account{}, err
	}

	query := `
	INSERT INTO accounts(id, username, password, is_admin) VALUES (?, ?, ?, ?)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return entities.Account{}, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id, account.Username, account.Password, false)
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

	if err = tx.Commit(); err != nil {
		return entities.Account{}, err
	}

	return account, nil
}
