package config

type PrometheusConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Endpoint string `mapstructure:"endpoint"`
}
