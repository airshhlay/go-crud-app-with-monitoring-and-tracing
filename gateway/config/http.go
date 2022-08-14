package config

// HTTPConfig defines config for the gateway's incoming HTTP requests
type HTTPConfig struct {
	UserService UserServiceConfig `mapstructure:userService`
	ItemService ItemServiceConfig `mapstructure:itemService`
}

// ItemServiceConfig holds config for routes to item service
type ItemServiceConfig struct {
	Label    string          `mapstructure:label`
	Host     string          `mapstructure:host`
	Port     string          `mapstructure:port`
	URLGroup string          `mapstructure:urlGroup`
	APIs     ItemServiceAPIs `mapstructure:apis`
}

// UserServiceConfig holds config for routes to user service
type UserServiceConfig struct {
	Label    string          `mapstructure:label`
	Host     string          `mapstructure:host`
	Port     string          `mapstructure:port`
	Secret   string          `mapstructure:secret`
	URLGroup string          `mapstructure:urlGroup`
	APIs     UserServiceAPIs `mapstructure:apis`
	Expiry   int             `mapstructure:"expiry"`
}

// UserServiceAPIs defines the public APIs to the user service
type UserServiceAPIs struct {
	Signup API `mapstructure:signup`
	Login  API `mapstructure:login`
}

// ItemServiceAPIs defines the public APIs to the item service
type ItemServiceAPIs struct {
	AddFav     API `mapstructure:addFav`
	DeleteFav  API `mapstructure:deleteFav`
	GetFavList API `mapstructure:getFavList`
}

// API config for a public API
type API struct {
	Endpoint string `mapstructure:endpoint`
	Method   string `mapstructure:method`
}
