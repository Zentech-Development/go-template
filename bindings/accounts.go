package bindings

import (
	"github.com/Zentech-Development/go-template/domain"
)

type AccountsBinding struct {
	Handlers *domain.Handlers
}

func newAccountsBinding(handlers *domain.Handlers) *AccountsBinding {
	return &AccountsBinding{
		Handlers: handlers,
	}
}
