package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Hostname string `mapstructure:hostname`
	Port     string `mapstructure:port`
}

func LoadConfig() *Config {
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("fatal error config file, %v\n", err)
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		fmt.Printf("unable to unmarshal into config struct, %v\n", err)
	}

	return config
}
