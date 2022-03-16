package config

import (
	"context"
	"github.com/brunomdev/digital-account/infra/log"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const fileConfig = ".env"

type Config struct {
	AppDebug           bool   `mapstructure:"APP_DEBUG"`
	HTTPPort           string `mapstructure:"HTTP_PORT"`
	DBHost             string `mapstructure:"DB_HOST"`
	DBPort             string `mapstructure:"DB_PORT"`
	DBDatabase         string `mapstructure:"DB_DATABASE"`
	DBUser             string `mapstructure:"DB_USER"`
	DBPass             string `mapstructure:"DB_PASS"`
	NewRelicAppName    string `mapstructure:"NEW_RELIC_APP_NAME"`
	NewRelicLicenseKey string `mapstructure:"NEW_RELIC_LICENSE_KEY"`
}

// Load the config from file or env to the Config struct
func Load() (*Config, error) {
	viper.SetConfigFile(fileConfig)
	viper.AutomaticEnv()

	viper.SetDefault("HTTP_PORT", "8080")

	var cfg Config

	if err := loadMappedEnvVariables(&cfg); err != nil {
		return nil, errors.WithMessage(err, "unable to load environment variables")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Info(context.TODO(), "unable find or read configuration file")
	}

	if err := viper.UnmarshalExact(&cfg); err != nil {
		return nil, errors.WithMessage(err, "error while unmarshal config")
	}

	return &cfg, nil
}

func loadMappedEnvVariables(cfg *Config) error {
	envKeysMap := &map[string]interface{}{}
	if err := mapstructure.Decode(cfg, &envKeysMap); err != nil {
		return err
	}

	for k := range *envKeysMap {
		if err := viper.BindEnv(k); err != nil {
			return err
		}
	}

	return nil
}
