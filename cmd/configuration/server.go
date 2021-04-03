package configuration

type ServerConfig struct {
	Port string `arg:"env:SERVER_PORT"`
	Address string `arg:"env:SERVER_ADDRESS"`
}

type Config struct {
	ServerConfig
}

func DefaultConfiguration() *Config{
	return &Config{
		ServerConfig: ServerConfig{
			Address: "localhost",
			Port: "6000",
		},
	}
}
