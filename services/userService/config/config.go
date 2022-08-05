package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
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

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %v\n", err)
		return nil, err
	}

	config := &Config{}
	err = viper.Unmarshal(config)

	if err != nil {
		fmt.Printf("Unable to unmarshal into config struct, %v\n", err)
		return nil, err
	}

	err = viper.UnmarshalKey("db", &config.DbConfig)

	if err != nil {
		fmt.Printf("Unable to unmarshal into dbconfig struct, %v\n", err)
		return nil, err
	}

	return config, nil
}
