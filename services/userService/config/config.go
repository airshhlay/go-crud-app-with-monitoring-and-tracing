package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Hostname string   `mapstructure:hostname`
	Port     string   `mapstructure:port`
	DbConfig DbConfig `mapstructure:db`
}

type DbConfig struct {
	Driver   string `mapstructure:driver`
	Host     string `mapstructure:host`
	Port     string `mapstructure:port`
	User     string `mapstructure:user`
	Net      string `mapstructure:net`
	DbName   string `mapstructure:dbName`
	Password string `mapstructure:password`
}

func LoadConfig(logger *zap.Logger) (*Config, error) {
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal(
			"Error reading config file",
			zap.Error(err),
		)
		return nil, err
	}

	config := &Config{}
	err = viper.Unmarshal(config)

	if err != nil {
		logger.Fatal(
			"Unable to unmarshal into config struct",
			zap.Error(err),
		)
		return nil, err
	}

	err = viper.UnmarshalKey("db", &config.DbConfig)
	if err != nil {
		logger.Fatal(
			"Unable to unmarshal into dbconfig struct",
			zap.Error(err),
		)
		return nil, err
	}

	return config, nil
}
