package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Database `mapstructure:"database"`
	}
	Database struct {
		User     string `env-required:"true" mapstructure:"user" env:"DATABASE_USER"`
		Password string `env-required:"true" mapstructure:"password" env:"DATABASE_PASSWORD"`
		Port     string `env-required:"true" mapstructure:"port" env:"DATABASE_PORT"`
		Schema   string `env-required:"true" mapstructure:"schema" env:"DATABASE_SCHEMA"`
		Database string `env-required:"true" mapstructure:"db" env:"DATABASE_DB"`
		Host     string `env-required:"true" mapstructure:"host" env:"DATABASE_HOST"`
	}
)

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config") // optionally look for config in the working directory
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	err = viper.Unmarshal(&cfg)

	if err != nil {
		return nil, fmt.Errorf("Config error %s", err)
	}

	return cfg, nil
}
