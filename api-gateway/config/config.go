package config

type Config struct {
	Server   ServerConfig
	Services ServicesConfig
}

type ServerConfig struct {
	Port string
}

type ServicesConfig struct {
	UserService    string
	ProductService string
	OrderService   string
}

var DefaultConfig = Config{
	Server: ServerConfig{
		Port: "8080",
	},
	Services: ServicesConfig{
		UserService:    "localhost:50051",
		ProductService: "localhost:50052",
		OrderService:   "localhost:50053",
	},
}
