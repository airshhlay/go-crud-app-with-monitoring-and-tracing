package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Hostname       string     `mapstructure:hostname`
	Port           string     `mapstructure:port`
	GinMode        string     `mapstructure: ginMode`
	TrustedProxies []string   `mapstructure:trustedProxies`
	AllowedOrigins []string   `mapstructure:allowedOrigins`
	HttpConfig     HttpConfig `mapstructure:http`
	GrpcConfig     GrpcConfig `mapstructure:grpc`
}

func LoadConfig(logger *zap.Logger) (*Config, error) {
	viper.AddConfigPath("/app/gateway/config")
	viper.AddConfigPath("/app/config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal(
			"error_config_file",
			zap.Error(err),
		)
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		logger.Fatal(
			"error_config_unmarshal",
			zap.Error(err),
		)
		return nil, err
	}

	err = viper.UnmarshalKey("grpc", &config.GrpcConfig)
	if err != nil {
		logger.Fatal(
			"Unable to unmarshal into GrpcConfig struct",
			zap.Error(err),
		)
		return nil, err
	}

	err = viper.UnmarshalKey("http", &config.HttpConfig)
	if err != nil {
		logger.Fatal(
			"Unable to unmarshal into HttpConfig struct",
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info(
		"info_config_loaded",
		zap.Any("config", config),
	)

	return config, nil
}
