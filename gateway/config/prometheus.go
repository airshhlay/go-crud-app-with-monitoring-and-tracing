package config

// PrometheusConfig holds config for the prometheus metrics scraping
type PrometheusConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Endpoint string `mapstructure:"endpoint"`
}
