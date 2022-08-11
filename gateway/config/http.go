package config

type HttpConfig struct {
	UserService UserServiceConfig `mapstructure:userService`
	ItemService ItemServiceConfig `mapstructure:itemService`
}
type ItemServiceConfig struct {
	Label    string          `mapstructure:label`
	Host     string          `mapstructure:host`
	Port     string          `mapstructure:port`
	UrlGroup string          `mapstructure:urlGroup`
	Apis     ItemServiceApis `mapstructure:apis`
}

type UserServiceConfig struct {
	Label    string          `mapstructure:label`
	Host     string          `mapstructure:host`
	Port     string          `mapstructure:port`
	Secret   string          `mapstructure:secret`
	UrlGroup string          `mapstructure:urlGroup`
	Apis     UserServiceApis `mapstructure:apis`
}

type UserServiceApis struct {
	Signup Api `mapstructure:signup`
	Login  Api `mapstructure:login`
}

type ItemServiceApis struct {
	AddFav     Api `mapstructure:addFav`
	DeleteFav  Api `mapstructure:deleteFav`
	GetFavList Api `mapstructure:getFavList`
}

type Api struct {
	Endpoint string `mapstructure:endpoint`
	Method   string `mapstructure:method`
}
