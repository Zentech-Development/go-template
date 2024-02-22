package main

import (
	"log"

	mockdb "github.com/Zentech-Development/go-template/adapters/repos/mock-db"
	"github.com/Zentech-Development/go-template/bindings"
	"github.com/Zentech-Development/go-template/domain"
	"github.com/Zentech-Development/go-template/handlers"
)

func main() {
	config := domain.GetConfig()

	repos := mockdb.NewMockDB()

	adapters := domain.Adapters{
		Repos: repos,
	}

	handlers := &domain.Handlers{
		Accounts: handlers.NewAccountHandler(&adapters),
	}

	server := bindings.NewServerBinding(handlers, config)

	log.Fatal(server.Run(config.Host))
}
