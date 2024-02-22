package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/Zentech-Development/go-template/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccountHandler struct {
	Adapters *domain.Adapters
}

func NewAccountHandler(adapters *domain.Adapters) AccountHandler {
	return AccountHandler{
		Adapters: adapters,
	}
}

func (h AccountHandler) GetByID(id string) (domain.Account, error) {
	ctx := context.Background()

	account, err := h.Adapters.Repos.Accounts.GetByID(ctx, id)
	if err != nil {
		return domain.Account{}, err
	}

	return account, nil
}

func (h AccountHandler) Add(account domain.AccountInput) (domain.Account, error) {
	ctx := context.Background()

	hashedPassword, err := hashPassword(account.Password, domain.GetConfig().HashCost)
	if err != nil {
		return domain.Account{}, errors.New("failed to generate password hash")
	}

	accountToSave := domain.Account{
		ID:       uuid.NewString(),
		Email:    account.Email,
		Locked:   false,
		Password: hashedPassword,
	}

	savedAccount, err := h.Adapters.Repos.Accounts.Add(ctx, accountToSave)
	if err != nil {
		return domain.Account{}, err
	}

	return savedAccount, nil
}

func (h AccountHandler) Login(credentials domain.LoginInput) (domain.Account, error) {
	ctx := context.Background()

	account, err := h.Adapters.Repos.Accounts.GetByEmail(ctx, credentials.Email)
	if err != nil {
		time.Sleep(time.Second)
		return domain.Account{}, errors.New("invalid credentials")
	}

	if !checkPassword(credentials.Password, account.Password) {
		time.Sleep(time.Second)
		return domain.Account{}, errors.New("invalid credentials")
	}

	return account, nil
}

func (h AccountHandler) ChangeStatus(id string, isLocked bool) error {
	ctx := context.Background()

	account, err := h.Adapters.Repos.Accounts.GetByID(ctx, id)
	if err != nil {
		return err
	}

	account.Locked = isLocked

	_, err = h.Adapters.Repos.Accounts.Update(ctx, account)
	return err
}

func (h AccountHandler) Remove(id string) error {
	ctx := context.Background()

	_, err := h.Adapters.Repos.Accounts.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return h.Adapters.Repos.Accounts.Remove(ctx, id)
}

func hashPassword(password string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func checkPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
