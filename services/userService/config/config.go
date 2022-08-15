package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config struct to hold main configuration from config.yaml
type Config struct {
	Hostname         string           `mapstructure:hostname`
	Port             string           `mapstructure:port`
	ServiceLabel     string           `mapstructure:serviceLabel`
	DbConfig         DbConfig         `mapstructure:db`
	PrometheusConfig PrometheusConfig `mapstructure:prometheus`
}

// DbConfig holds configurations for the database.
type DbConfig struct {
	ServiceLabel string `mapstructure:"serviceLabel"`
	Driver       string `mapstructure:driver`
	Host         string `mapstructure:host`
	Port         string `mapstructure:port`
	User         string `mapstructure:user`
	Net          string `mapstructure:net`
	DbName       string `mapstructure:dbName`
	Password     string `mapstructure:password`
}

// LoadConfig is called in main.go to load all config
func LoadConfig(logger *zap.Logger) (*Config, error) {
	// ex, error := os.Executable()
	// if error != nil {
	// 	panic(error)
	// }
	// exPath := filepath.Dir(ex)
	// logger.Info(exPath)
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/app/config")
	viper.AddConfigPath("/app/userService/config")
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

	err = viper.UnmarshalKey("prometheus", &config.PrometheusConfig)
	if err != nil {
		logger.Fatal(
			"Unable to unmarshal into prometheusconfig struct",
			zap.Error(err),
		)
		return nil, err
	}

	return config, nil
}
