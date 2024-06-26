package main

import (
	"errors"
	"flag"
	"log"

	ginBinding "github.com/Zentech-Development/go-template/pkg/bindings/gin"
	"github.com/Zentech-Development/go-template/pkg/config"
	"github.com/Zentech-Development/go-template/pkg/entities"
	"github.com/Zentech-Development/go-template/pkg/service"
	sqliteStore "github.com/Zentech-Development/go-template/pkg/stores/sqlite"
)

const (
	STORE_TYPE_SQLITE = "sqlite"
	BINDING_TYPE_GIN  = "gin"
)

func main() {
	config := config.GetConfig()

	storeType := flag.String("store", STORE_TYPE_SQLITE, "the type of store to use, allowed values: sqlite")
	bindingType := flag.String("binding", BINDING_TYPE_GIN, "the type of binding to use, allowed values: gin")
	flag.Parse()

	accountStore, err := getStores(*storeType)
	if err != nil {
		log.Fatal(err)
	}

	services := service.NewService(accountStore)

	run(*bindingType, config, services)
}

func getStores(storeType string) (entities.AccountStore, error) {
	switch storeType {
	case STORE_TYPE_SQLITE:
		sqliteStore := sqliteStore.NewSQLiteStore("./APPNAME.db")
		return sqliteStore.Accounts, nil

	default:
		return nil, errors.New("invalid store type, try ./APPNAME --help")
	}
}

func run(bindingType string, config *config.Config, services *service.Service) {
	switch bindingType {
	case BINDING_TYPE_GIN:
		app := ginBinding.NewBinding(services, config.Host, config.Debug, config.SecretKey)
		log.Fatal(app.Run())

	default:
		log.Fatal("Invalid binding type, try ./APPNAME --help")
	}

}
