package handlers

import (
	"strings"
	"testing"

	"github.com/Zentech-Development/go-template/domain"
)

func TestHashFunctions(t *testing.T) {
	testPassword := "test123"
	hash, err := hashPassword(testPassword, 12)
	if err != nil {
		t.Fatal("Failed to hash password")
	}

	if !strings.HasPrefix(hash, "$2a") {
		t.Fatal("Token is bad")
	}

	if !checkPassword(testPassword, hash) {
		t.Fatal("Password should have been correct")
	}

	if checkPassword("wrong", hash) {
		t.Fatal("Password should have been incorrect")
	}
}

func TestGetByID(t *testing.T) {
	handlers := newHandlers()

	email := "test@test.com"
	password := "password123"
	account := domain.AccountInput{
		Email:    email,
		Password: password,
	}

	originalAccount, err := handlers.Accounts.Add(account)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	savedAccount, err := handlers.Accounts.GetByID(originalAccount.ID)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedAccount.ID != originalAccount.ID || savedAccount.Email != originalAccount.Email {
		t.Fatal("Return the wrong account")
	}

	if _, err = handlers.Accounts.GetByID("bad-id"); err == nil {
		t.Fatal("Expected not found error")
	}
}

func TestAccountAdd(t *testing.T) {
	handlers := newHandlers()

	email := "test@test.com"
	password := "password123"
	account := domain.AccountInput{
		Email:    email,
		Password: password,
	}

	savedAccount, err := handlers.Accounts.Add(account)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedAccount.ID == "" || savedAccount.Email != email {
		t.Fatal("Account saved improperly")
	}

	if savedAccount.Password == password {
		t.Fatal("Failed to hash password... yikes")
	}
}

func TestLogin(t *testing.T) {
	handlers := newHandlers()

	email := "test@test.com"
	password := "password123"
	account := domain.AccountInput{
		Email:    email,
		Password: password,
	}

	_, err := handlers.Accounts.Add(account)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	validCreds := domain.LoginInput{
		Email:    email,
		Password: password,
	}

	result, err := handlers.Accounts.Login(validCreds)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if result.Email != email {
		t.Fatal("Login returned wrong account")
	}

	badEmail := domain.LoginInput{
		Email:    "bad",
		Password: password,
	}
	if _, err = handlers.Accounts.Login(badEmail); err == nil {
		t.Fatal("Expected bad credentials error")
	}

	badPassword := domain.LoginInput{
		Email:    email,
		Password: "bad",
	}
	if _, err = handlers.Accounts.Login(badPassword); err == nil {
		t.Fatal("Expected bad credentials error")
	}

	badUsernameAndPasskey := domain.LoginInput{
		Email:    "bad",
		Password: "also-bad",
	}
	if _, err = handlers.Accounts.Login(badUsernameAndPasskey); err == nil {
		t.Fatal("Expected bad credentials error")
	}
}

func TestUpdateAccountStatus(t *testing.T) {
	handlers := newHandlers()

	email := "test@test.com"
	password := "password123"
	account := domain.AccountInput{
		Email:    email,
		Password: password,
	}

	savedAccount, err := handlers.Accounts.Add(account)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if err = handlers.Accounts.ChangeStatus(savedAccount.ID, true); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	savedAccount, err = handlers.Accounts.GetByID(savedAccount.ID)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if !savedAccount.Locked {
		t.Fatal("Failed to set account to locked")
	}

	if err = handlers.Accounts.ChangeStatus(savedAccount.ID, false); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	savedAccount, err = handlers.Accounts.GetByID(savedAccount.ID)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if savedAccount.Locked {
		t.Fatal("Failed to set account to unlocked")
	}

	if err = handlers.Accounts.ChangeStatus("bad-id", true); err == nil {
		t.Fatal("Expected not found error")
	}
}

func TestRemoveAccount(t *testing.T) {
	handlers := newHandlers()

	email := "test@test.com"
	password := "password123"
	account := domain.AccountInput{
		Email:    email,
		Password: password,
	}

	savedAccount, err := handlers.Accounts.Add(account)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if err = handlers.Accounts.Remove(savedAccount.ID); err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if _, err = handlers.Accounts.GetByID(savedAccount.ID); err == nil {
		t.Fatal("Expected not found error, failed to remove account")
	}

	if err = handlers.Accounts.Remove("bad-id"); err == nil {
		t.Fatal("Expected not found error")
	}
}
