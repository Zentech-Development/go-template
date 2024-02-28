package domain

import (
	"context"
	"fmt"
)

type Account struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Locked   bool   `json:"-"`
	Password string `json:"-"`
}

type AccountInput struct {
	Email           string `json:"email" form:"email" binding:"required"`
	Password        string `json:"password" form:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type AccountRepo interface {
	GetByID(ctx context.Context, id string) (Account, error)
	GetByEmail(ctx context.Context, email string) (Account, error)
	Add(ctx context.Context, account Account) (Account, error)
	Update(ctx context.Context, account Account) (Account, error)
	Remove(ctx context.Context, id string) error
}

type AccountHandlers interface {
	GetByID(id string) (Account, error)
	Add(account AccountInput) (Account, error)
	ChangeStatus(id string, isLocked bool) error
	Remove(id string) error
	Login(credentials LoginInput) (Account, error)
}

type AccountNotFoundError struct {
	ID    string
	Email string
}

func (e *AccountNotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("account with id '%s' not found", e.ID)
	}
	return fmt.Sprintf("account with email '%s' not found", e.Email)
}

type AccountAlreadyExistsError struct {
	ID    string
	Email string
}

func (e *AccountAlreadyExistsError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("account with id '%s' already exists", e.ID)
	}
	return fmt.Sprintf("account with email '%s' already exists", e.Email)
}
