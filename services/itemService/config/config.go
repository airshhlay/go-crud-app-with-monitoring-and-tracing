package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Hostname       string         `mapstructure:hostname`
	Port           string         `mapstructure:port`
	MaxPerPage     int            `mapstructure:maxPerPage`
	DbConfig       DbConfig       `mapstructure:db`
	RedisConfig    RedisConfig    `mapstructure:redis`
	ExternalConfig ExternalConfig `mapstructure:external`
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

type RedisConfig struct {
	Host     string `mapstructure:host`
	Port     string `mapstructure:port`
	Password string `mapstructure:password`
	Db       int    `mapstructure:db`
	Expire   int    `mapstructure:expire`
}

type ExternalConfig struct {
	Shopee Shopee `mapstructure:shopee`
}

type Shopee struct {
	GetItem Api `mapstructure:getItem`
}

type Api struct {
	Endpoint string `mapstructure:endpoint`
	Method   string `mapstructure:method`
}

func LoadConfig(logger *zap.Logger) (*Config, error) {
	viper.AddConfigPath("/app/config")
	viper.AddConfigPath("/app/itemService/config")
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

	err = viper.UnmarshalKey("redis", &config.RedisConfig)
	if err != nil {
		logger.Fatal(
			"Unable to unmarshal into redisconfig struct",
			zap.Error(err),
		)
		return nil, err
	}

	err = viper.UnmarshalKey("external", &config.ExternalConfig)
	if err != nil {
		logger.Fatal(
			"Unable to unmarshal into externalconfig struct",
			zap.Error(err),
		)
		return nil, err
	}

	return config, nil
}
