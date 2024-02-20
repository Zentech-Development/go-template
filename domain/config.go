package domain

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

var lock = &sync.Mutex{}

var appConfig *ApplicationConfig

type ApplicationConfig struct {
	Lifecycle              string `mapstructure:"lifecycle"`
	Host                   string `mapstructure:"host"`
	SecretKey              string `mapstructure:"secret_key" json:"-"`
	AppName                string `mapstructure:"app_name" json:"appName"`
	LoginExpirationSeconds int    `mapstructure:"login_expiration_seconds" json:"loginExpirationSeconds"`
	TokenExpirationSeconds int    `mapstructure:"token_expiration_seconds" json:"tokenExpirationSeconds"`
	AutoRefreshSeconds     int    `mapstructure:"auto_refresh_seconds" json:"autoRefreshSeconds"`
	TokenName              string `mapstructure:"token_name" json:"tokenName"`
}

const (
	LIFECYCLE_PRODUCTION = "prod"
	LIFECYCLE_DEVELOP    = "dev"
	LIFECYCLE_LOCAL      = "local"
)

func GetConfig() *ApplicationConfig {
	if appConfig == nil {
		lock.Lock()
		defer lock.Unlock()

		if appConfig == nil {
			appConfig = newConfig()
		}
	}

	return appConfig
}

func newConfig() *ApplicationConfig {
	config := &ApplicationConfig{}

	v := viper.New()

	v.SetDefault("lifecycle", LIFECYCLE_PRODUCTION)
	v.SetDefault("host", "localhost:8000")
	v.SetDefault("secret_key", "")
	v.SetDefault("app_name", "go-template")
	v.SetDefault("token_name", "X-API-KEY")

	one_day := time.Second * 60 * 60 * 24
	v.SetDefault("login_expiration_seconds", one_day)

	four_hours := time.Second * 60 * 60 * 4
	v.SetDefault("token_expiration_seconds", four_hours)

	fifteen_minutes := time.Second * 60 * 15
	v.SetDefault("auto_refresh_seconds", fifteen_minutes)

	v.SetConfigFile("go-template-config.json")

	v.SetEnvPrefix("go-template")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatal("Failed to load the configuration file: ", err)
	}

	if err := v.Unmarshal(config); err != nil {
		log.Fatal("Invalid configuration: ", err)
	}

	var validLifecycles = []string{LIFECYCLE_DEVELOP, LIFECYCLE_LOCAL, LIFECYCLE_PRODUCTION}
	if !slices.Contains(validLifecycles, config.Lifecycle) {
		log.Fatalf("Invalid configuration: lifecycle must be one of (%s, %s, %s)", LIFECYCLE_DEVELOP, LIFECYCLE_LOCAL, LIFECYCLE_PRODUCTION)
	}

	if config.SecretKey == "" {
		log.Fatal("Invalid configuration: must provide a secret_key")
	}

	log.Default().Println("Application configuration initialized")
	vals, _ := json.MarshalIndent(config, "", "\t")
	log.Default().Println(string(vals))

	return config
}
