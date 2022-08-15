package config

// JaegerConfig config for opentracing
type JaegerConfig struct {
	Host        string `mapstructure:"host"`
	ServiceName string `mapstructure:"serviceName"`
	LogSpans    bool   `mapstructure:"logSpans"`
}
