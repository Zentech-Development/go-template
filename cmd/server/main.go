package main

import (
	"log"

	ginBinding "github.com/Zentech-Development/go-template/pkg/bindings/gin"
	"github.com/Zentech-Development/go-template/pkg/config"
	"github.com/Zentech-Development/go-template/pkg/service"
	sqliteStore "github.com/Zentech-Development/go-template/pkg/stores/sqlite"
)

func main() {
	config := config.GetConfig()

	sqliteStore := sqliteStore.NewSQLiteStore("./APPNAME.db")
	services := service.NewService(sqliteStore.Accounts)
	app := ginBinding.NewBinding(services, config.Host, config.Debug, config.SecretKey)

	log.Fatal(app.Run())
}
