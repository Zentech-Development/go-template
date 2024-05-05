package service

import "github.com/Zentech-Development/go-template/pkg/entities"

// Service represents the services available to compose applications with.
// The Service struct provides implementations of the actions available
// for each entity in the application. External dependencies (databases, caching,
// logging, etc.) are provided to the Service struct via dependency injection.
type Service struct {
	accountStore entities.AccountStore
}

// NewService creates an instance of the Service struct with the provided
// dependencies injected.
func NewService(accountStore entities.AccountStore) *Service {
	return &Service{
		accountStore: accountStore,
	}
}
