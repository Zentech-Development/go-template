package mockdb

import (
	"github.com/Zentech-Development/go-template/domain"
)

type MockDBData struct {
	Accounts []domain.Account
}

type MockDB struct {
	Data *MockDBData
}

func NewMockDB() domain.Repos {
	data := &MockDBData{
		Accounts: make([]domain.Account, 0),
	}

	return domain.Repos{
		Accounts: NewMockAccountRepo(data),
	}
}
