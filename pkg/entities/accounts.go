package entities

import (
	"context"
)

// Account represents a stored account.
type Account struct {
	Username string `json:"username"`
	Password string `json:"-"`
	IsAdmin  bool   `json:"isAdmin"`
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
	Create(ctx context.Context, account Account) (Account, error)
}
