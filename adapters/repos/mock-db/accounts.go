package mockdb

import (
	"context"

	"github.com/Zentech-Development/go-template/domain"
)

type MockAccountRepo struct {
	Data *MockDBData
}

func NewMockAccountRepo(data *MockDBData) MockAccountRepo {
	return MockAccountRepo{
		Data: data,
	}
}

func (r MockAccountRepo) GetByID(ctx context.Context, id string) (domain.Account, error) {
	for _, account := range r.Data.Accounts {
		if account.ID == id {
			return account, nil
		}
	}

	return domain.Account{}, &domain.AccountNotFoundError{ID: id}
}

func (r MockAccountRepo) GetByEmail(ctx context.Context, email string) (domain.Account, error) {
	for _, account := range r.Data.Accounts {
		if account.Email == email {
			return account, nil
		}
	}

	return domain.Account{}, &domain.AccountNotFoundError{Email: email}
}

func (r MockAccountRepo) Add(ctx context.Context, account domain.Account) (domain.Account, error) {
	if _, err := r.GetByEmail(ctx, account.Email); err == nil {
		return domain.Account{}, &domain.AccountAlreadyExistsError{Email: account.Email}
	}

	if _, err := r.GetByID(ctx, account.ID); err == nil {
		return domain.Account{}, &domain.AccountAlreadyExistsError{ID: account.ID}
	}

	r.Data.Accounts = append(r.Data.Accounts, account)
	return account, nil
}

func (r MockAccountRepo) Update(ctx context.Context, account domain.Account) (domain.Account, error) {
	if _, err := r.GetByID(ctx, account.ID); err != nil {
		return domain.Account{}, err
	}

	for i, currAccount := range r.Data.Accounts {
		if currAccount.ID == account.ID {
			r.Data.Accounts[i] = account
		}
	}

	return account, nil
}

func (r MockAccountRepo) Remove(ctx context.Context, id string) error {
	newAccounts := make([]domain.Account, 0)

	for _, account := range r.Data.Accounts {
		if account.ID != id {
			newAccounts = append(newAccounts, account)
		}
	}

	if len(newAccounts) == len(r.Data.Accounts) {
		return &domain.AccountNotFoundError{ID: id}
	}

	r.Data.Accounts = newAccounts

	return nil
}
