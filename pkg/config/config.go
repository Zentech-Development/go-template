package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Zentech-Development/go-template/pkg/logger"
	"github.com/spf13/viper"
)

var C *Config

type Config struct {
	Host          string     `mapstructure:"HOST" json:"HOST"`
	Debug         bool       `mapstructure:"DEBUG" json:"DEBUG"`
	SecretKey     string     `mapstructure:"SECRET_KEY" json:"-"`
	SecretKeyFile string     `mapstructure:"SECRET_KEY_FILE" json:"SECRET_KEY_FILE"`
	SQLiteOpts    SQLiteOpts `mapstructure:"SQLITE_OPTS" json:"SQLITE_OPTS"`
}

type SQLiteOpts struct {
	DBPath        string `mapstructure:"DB_PATH" json:"DB_PATH"`
	MigrationsDir string `mapstructure:"MIGRATIONS_DIR" json:"MIGRATIONS_DIR"`
}

func Init(configFilePath string) {
	logger.L.Info().Msgf("Initializing config from file %s", configFilePath)

	C = &Config{}

	v := viper.New()

	v.SetDefault("HOST", "localhost:8000")
	v.SetDefault("DEBUG", false)

	v.SetConfigFile(configFilePath)

	if err := v.ReadInConfig(); err != nil {
		logger.L.Fatal().Err(err).Msgf("Failed to load config file: %s", configFilePath)
	}

	if err := v.Unmarshal(C); err != nil {
		logger.L.Fatal().Err(err).Msgf("Failed to parse configuration: %s", configFilePath)
	}

	if err := C.loadSecretsFromFiles(); err != nil {
		logger.L.Fatal().Err(err).Msgf("Failed to read secrets from files: %s", configFilePath)
	}

	if err := C.validate(); err != nil {
		logger.L.Fatal().Err(err).Msgf("Invalid configuration: %s", configFilePath)
	}

	vals, _ := json.MarshalIndent(C, "", "\t")
	logger.L.Info().RawJSON("config", vals).Msg("APPNAME configuration initialized successfully")
}

func (c *Config) loadSecretsFromFiles() error {
	if c.SecretKeyFile != "" {
		if c.SecretKey != "" {
			logger.L.Warn().Msg("WARNING: overwriting SECRET_KEY value from SECRET_KEY_FILE")
		}
		secretKey, err := os.ReadFile(c.SecretKeyFile)
		if err != nil {
			return err
		}

		c.SecretKey = string(secretKey)
	}

	return nil
}

func (c *Config) validate() error {
	if c.SecretKey == "" {
		return errors.New("must provide a SECRET_KEY or SECRET_KEY_FILE")
	}

	return nil
}
