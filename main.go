package main

import (
	"errors"
	"flag"

	httpserver "github.com/Zentech-Development/go-template/internals/bindings/httpserver"
	"github.com/Zentech-Development/go-template/internals/entities"
	"github.com/Zentech-Development/go-template/internals/services"
	sqliteStore "github.com/Zentech-Development/go-template/internals/stores/sqlite"
	"github.com/Zentech-Development/go-template/pkg/config"
	"github.com/Zentech-Development/go-template/pkg/logger"
	"github.com/rs/zerolog"
)

const (
	STORE_TYPE_SQLITE = "sqlite"
	BINDING_TYPE_GIN  = "gin"
)

func main() {
	storeType := flag.String("store", STORE_TYPE_SQLITE, "the type of store to use, allowed values: sqlite")
	bindingType := flag.String("binding", BINDING_TYPE_GIN, "the type of binding to use, allowed values: gin")
	logLevel := flag.Int("log", int(zerolog.DebugLevel), "the zerolog log level to use, allowed values: 0, 1, 2, 3, 4, 5, default is 0 (debug)")
	configFile := flag.String("config", "./APPNAME-config.json", "the path to the config file, default is ./APPNAME-config.json")
	flag.Parse()

	logOpts := logger.Opts{
		Level:      zerolog.Level(*logLevel),
		TimeFormat: zerolog.TimeFormatUnix,
	}
	logger.InitLogger(logOpts)

	config.Init(*configFile)

	accountStore, err := getStores(*storeType, config.C)
	if err != nil {
		logger.L.Fatal().Err(err).Msg("Store initialization failed")
	}

	s := services.NewService(accountStore)

	run(*bindingType, config.C, s)
}

func getStores(storeType string, c *config.Config) (entities.AccountStore, error) {
	switch storeType {
	case STORE_TYPE_SQLITE:
		sqliteStore := sqliteStore.NewSQLiteStore(&sqliteStore.Opts{
			DBPath:         c.SQLiteOpts.DBPath,
			MigrationsPath: c.SQLiteOpts.MigrationsDir,
		})
		return sqliteStore, nil

	default:
		return nil, errors.New("invalid store type, try ./APPNAME --help")
	}
}

func run(bindingType string, c *config.Config, s *services.Services) {
	switch bindingType {
	case BINDING_TYPE_GIN:
		app := httpserver.NewBinding(s, httpserver.HTTPServerOpts{
			DebugMode:  c.Debug,
			SecretKey:  c.SecretKey,
			ListenAddr: c.Host,
		})
		logger.L.Fatal().Err(app.Run()).Msg("Application crashed")

	default:
		logger.L.Fatal().Msg("Invalid binding type, try ./APPNAME --help")
	}
}
