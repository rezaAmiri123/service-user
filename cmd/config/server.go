package config

type CommonConfig struct {
	JWTSecret string `arg:"env:JWT_SECRET"`
}

type ServerConfig struct {
	ServerAddress string `arg:"env:SERVER_ADDRESS"`
	ServerPort    string `arg:"env:SERVER_PORT"`
}

type GatewayConfig struct {
	GatewayAddress string `arg:"env:GATEWAY_ADDRESS"`
	GatewayPort    string `arg:"env:GATEWAY_PORT"`
}

type DatabaseConfig struct {
	DBType string `arg:"env:DB_TYPE"`
	DBUser string `arg:"env:DB_USER"`
	DBPass string `arg:"env:DB_PASSWORD"`
	DBName string `arg:"env:DB_NAME"`
	DBHost string `arg:"env:DB_HOST"`
	DBPort string `arg:"env:DB_PORT"`
}
type Config struct {
	ServerConfig
	GatewayConfig
	DatabaseConfig
	CommonConfig
}

func DefaultConfiguration() *Config {
	return &Config{
		CommonConfig: CommonConfig{
			JWTSecret: "1234%^&*ukfykjSCFAVARBTSDN",
		},
		ServerConfig: ServerConfig{
			ServerAddress: "localhost",
			ServerPort:    "6000",
		},
		GatewayConfig: GatewayConfig{
			GatewayAddress: "localhost",
			GatewayPort:    "8000",
		},
		DatabaseConfig: DatabaseConfig{
			DBType: "mysql",
			DBName: "go",
			DBUser: "go",
			DBPass: "go",
			DBHost: "localhost",
			DBPort: "3306",
		},
	}
}

func (cfg *ServerConfig) GetServerAddress() string {
	return cfg.ServerAddress + ":" + cfg.ServerPort
}

func (cfg *GatewayConfig) GetGatewayAddress() string {
	return cfg.GatewayAddress + ":" + cfg.GatewayPort
}
