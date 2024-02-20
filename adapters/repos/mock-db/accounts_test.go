package mockdb_test

import (
	"context"
	"testing"

	mockAdapters "github.com/Zentech-Development/go-template/adapters/repos/mock-db"
	"github.com/Zentech-Development/go-template/domain"
	"github.com/google/uuid"
)

func TestAccountGetByID(t *testing.T) {
	id := uuid.NewString()
	locked := false
	email := "test@test.com"

	db := mockAdapters.NewMockDB()

	account := domain.Account{
		ID:     id,
		Email:  email,
		Locked: locked,
	}

	returnedAccount, err := db.Accounts.Add(context.Background(), account)
	if err != nil {
		t.Fatal("Unexpected error occurred: ", err)
	}

	if !isAccountCorrect(account, returnedAccount) {
		t.Fatal("Did not return same account as input")
	}

	savedAccount, err := db.Accounts.GetByID(context.Background(), account.ID)
	if err != nil {
		t.Fatal("Unexpected error occurred: ", err)
	}

	if !isAccountCorrect(account, savedAccount) {
		t.Fatal("User did not save properly")
	}

	if _, err := db.Accounts.Add(context.Background(), account); err == nil {
		t.Fatal("Expected account already exists error")
	}
}
func TestAccountGetByEmail(t *testing.T) {
	t.FailNow()
}
func TestAccountAdd(t *testing.T) {
	t.FailNow()
}
func TestAccountUpdate(t *testing.T) {
	t.FailNow()
}

func TestAccountRemove(t *testing.T) {
	t.FailNow()
}

func isAccountCorrect(expected domain.Account, actual domain.Account) bool {
	isMatch := true

	if expected.ID != actual.ID {
		isMatch = false
	}

	if expected.Email != actual.Email {
		isMatch = false
	}

	if expected.Locked != actual.Locked {
		isMatch = false
	}

	return isMatch
}
