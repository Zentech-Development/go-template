package entities

import (
	"context"
)

// Account represents a stored account.
type Account struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

// AccountCreateInput represents valid input data for creating an account.
type AccountCreateInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AccountLoginInput represents valid input data for a login request.
type AccountLoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AccountStore interface {
	GetByUsername(ctx context.Context, username string) (Account, error)
	Create(ctx context.Context, account AccountCreateInput) (Account, error)
}
