package service

import (
	"context"

	"github.com/Zentech-Development/go-template/pkg/entities"
)

// Ping handles a ping request. This should always return nil.
func (s *Service) Ping(ctx context.Context) error {
	return nil
}

// GetUserByUsername retrieves a stored user by username. If the username is not
// found or the account store call fails, an error will be returned,
func (s *Service) GetUserByUsername(ctx context.Context, username string) (entities.Account, error) {
	return entities.Account{}, nil
}
