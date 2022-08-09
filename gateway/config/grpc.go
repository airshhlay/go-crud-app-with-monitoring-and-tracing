package config

type GrpcConfig struct {
	UserService GrpcServiceConfig `mapstructure:userService`
	ItemService GrpcServiceConfig `mapstructure:itemService`
}
type GrpcServiceConfig struct {
	Label string `mapstructure:label`
	Host  string `mapstructure:host`
	Port  string `mapstructure:port`
}
