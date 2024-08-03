package service

import (
	"context"
	"errors"

	"github.com/Zentech-Development/go-template/pkg/entities"
	"github.com/alexedwards/argon2id"
)

// GetAccountByUsername retrieves a stored user by username. If the username is not
// found or the account store call fails, an error will be returned.
func (s *Service) GetAccountByUsername(ctx context.Context, username string) (entities.Account, error) {
	account, err := s.accountStore.GetByUsername(ctx, username)
	if err != nil {
		return entities.Account{}, err
	}

	return account, nil
}

// GetAccountByID retrieves a stored user by ID. If the ID is not
// found or the account store call fails, an error will be returned.
func (s *Service) GetAccountByID(ctx context.Context, id int64) (entities.Account, error) {
	account, err := s.accountStore.GetByID(ctx, id)
	if err != nil {
		return entities.Account{}, err
	}

	return account, nil
}

// Create adds a new user to the database. If the username already exists an error will be returned.
func (s *Service) CreateAccount(ctx context.Context, input entities.AccountCreateInput) (entities.Account, error) {
	_, err := s.accountStore.GetByUsername(ctx, input.Username)
	if err == nil {
		return entities.Account{}, &entities.ErrAlreadyExists{}
	}
	if !errors.Is(err, &entities.ErrNotFound{}) {
		return entities.Account{}, err
	}

	hashedPassword, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		return entities.Account{}, errors.New("password hash failed")
	}

	input.Password = hashedPassword

	savedAccount, err := s.accountStore.Create(ctx, input)
	if err != nil {
		return entities.Account{}, err
	}

	return savedAccount, nil
}

// Login checks the provided input against accounts and returns an error if the credentials
// are not valid.
func (s *Service) Login(ctx context.Context, input entities.AccountLoginInput) (entities.Account, error) {
	savedAccount, err := s.accountStore.GetByUsername(ctx, input.Username)
	if err != nil {
		return entities.Account{}, &entities.ErrBadCredentials{}
	}

	match, err := argon2id.ComparePasswordAndHash(input.Password, savedAccount.Password)
	if err != nil {
		return entities.Account{}, err
	}
	if !match {
		return entities.Account{}, &entities.ErrBadCredentials{}
	}

	return savedAccount, nil
}
