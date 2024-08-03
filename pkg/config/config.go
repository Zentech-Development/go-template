package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Host      string `mapstructure:"HOST" json:"HOST"`
	Debug     bool   `mapstructure:"DEBUG" json:"DEBUG"`
	SecretKey string `mapstructure:"SECRET_KEY" json:"-"`
	Lifecycle string `mapstructure:"LIFECYCLE" json:"LIFECYCLE"`
}

var (
	activeConfig *Config
	lock         = &sync.Mutex{}
)

const (
	LIFECYCLE_PRODUCTION = "prod"
	LIFECYCLE_DEVELOP    = "dev"
	LIFECYCLE_LOCAL      = "local"
)

func GetConfig() *Config {
	// extra check here to avoid using the (very expensive) lock whenever possible
	if activeConfig == nil {
		lock.Lock()
		defer lock.Unlock()

		if activeConfig == nil {
			activeConfig = newConfig()
		}
	}

	return activeConfig
}

func newConfig() *Config {
	conf := &Config{}

	v := viper.New()

	v.SetEnvPrefix("APPNAME")

	v.SetDefault("HOST", "localhost:8000")
	v.SetDefault("LIFECYCLE", LIFECYCLE_PRODUCTION)
	v.SetDefault("DEBUG", false)

	v.SetConfigName("APPNAME-config")

	projectDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Somehow failed to get the working directory")
	}
	v.AddConfigPath(projectDir)

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatal("Failed to load the configuration: ", err)
	}

	if err := v.Unmarshal(conf); err != nil {
		log.Fatal("Failed to parse configuration: ", err)
	}

	if err = conf.validate(); err != nil {
		log.Fatal("Invalid configuration: ", err)
	}

	log.Default().Println("APPNAME configuration initialized")

	vals, _ := json.MarshalIndent(conf, "", "\t")

	if conf.Debug {
		log.Default().Println(string(vals))
	}

	return conf
}

func (c *Config) validate() error {
	if c.SecretKey == "" {
		return errors.New("must provide a SECRET_KEY")
	}

	return nil
}
