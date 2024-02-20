package main

import (
	"log"

	"github.com/Zentech-Development/go-template/bindings"
	"github.com/Zentech-Development/go-template/domain"
)

func main() {
	config := domain.GetConfig()

	// repos := mockdb.NewMockDB()

	// adapters := domain.Adapters{
	// 	Repos: repos,
	// }

	handlers := domain.Handlers{}

	server := bindings.NewServerBinding(&handlers, config)

	log.Fatal(server.Run(config.Host))
}
