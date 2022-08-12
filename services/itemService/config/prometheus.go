package config

// PrometheusConfig stores config for enabling prometheus scraping
type PrometheusConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Endpoint string `mapstructure:"endpoint"`
}
