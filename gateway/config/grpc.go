package config

// GrpcConfig holds config for the different grpc clients
type GrpcConfig struct {
	UserService GrpcServiceConfig `mapstructure:userService`
	ItemService GrpcServiceConfig `mapstructure:itemService`
}

// GrpcServiceConfig defines the config for a grpc client
type GrpcServiceConfig struct {
	Label string `mapstructure:label`
	Host  string `mapstructure:host`
	Port  string `mapstructure:port`
}
