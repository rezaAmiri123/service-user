package config

type ServerConfig struct {
	Port    string `arg:"env:SERVER_PORT"`
	Address string `arg:"env:SERVER_ADDRESS"`
}

type Database struct {
	DBType string `arg:"env:DB_TYPE"`
	DBUser string `arg:"env:DB_USER"`
	DBPass string `arg:"env:DB_PASSWORD"`
	DBName string `arg:"env:DB_NAME"`
	DBHost string `arg:"env:DB_HOST"`
	DBPort string `arg:"env:DB_PORT"`
}
type Config struct {
	ServerConfig
	Database
}

func DefaultConfiguration() *Config {
	return &Config{
		ServerConfig: ServerConfig{
			Address: "localhost",
			Port:    "6000",
		},
		Database: Database{
			DBType: "mysql",
			DBName: "go",
			DBUser: "go",
			DBPass: "go",
			DBHost: "localhost",
			DBPort: "3306",
		},
	}
}

func (cfg *ServerConfig) GetAddress() string {
	return cfg.Address + ":" + cfg.Port
}
