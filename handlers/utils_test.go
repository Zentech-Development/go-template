package handlers

import (
	adapters "github.com/Zentech-Development/go-template/adapters/repos/mock-db"
	"github.com/Zentech-Development/go-template/domain"
)

func newHandlers() domain.Handlers {
	mockDB := adapters.NewMockDB()

	adapts := domain.Adapters{
		Repos:  mockDB,
		Logger: domain.Logger{},
	}

	handlers := domain.Handlers{
		Accounts: NewAccountHandler(&adapts),
	}

	return handlers
}
