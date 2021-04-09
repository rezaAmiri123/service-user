package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type CommonConfig struct {
	JWTSecret string
}

type GatewayConfig struct {
	Port string
	ServerAddress string
	ServerPort string
}

func (c *GatewayConfig)GetServerAddress()string{
	return c.ServerAddress + c.ServerPort
}
type DatabaseConfig struct {
	DBType string
	DBUser string
	DBPass string
	DBName string
	DBHost string
	DBPort string
}
type Config struct {
	Server   ServerConfig
	Gateway  GatewayConfig
	Database DatabaseConfig
	Common   CommonConfig
	Redis    RedisConfig
	Logger   LoggerConfig
	Jaeger   JaegerConfig
	Metrics  MetricsConfig
}

// Server config struct
type ServerConfig struct {
	AppVersion        string
	Port              string
	PprofPort         string
	Mode              string
	JwtSecretKey      string
	CookieName        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
	MaxConnectionIdle time.Duration
	Timeout           time.Duration
	MaxConnectionAge  time.Duration
	Time              time.Duration
}

// Metrics config
type MetricsConfig struct {
	URL         string
	ServiceName string
}

// Logger config
type LoggerConfig struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// Redis config
type RedisConfig struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultdb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

// Jaeger
type JaegerConfig struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

// Get config
func GetConfig(configPath string) (*Config, error) {
	cfgFile, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	cfg, err := ParseConfig(cfgFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

//func DefaultConfiguration() *Config {
//	return &Config{
//		CommonConfig: CommonConfig{
//			JWTSecret: "1234%^&*ukfykjSCFAVARBTSDN",
//		},
//		ServerConfig: ServerConfig{
//			ServerAddress: "localhost",
//			ServerPort:    "6000",
//		},
//		GatewayConfig: GatewayConfig{
//			GatewayAddress: "localhost",
//			GatewayPort:    "8000",
//		},
//		DatabaseConfig: DatabaseConfig{
//			DBType: "mysql",
//			DBName: "go",
//			DBUser: "go",
//			DBPass: "go",
//			DBHost: "localhost",
//			DBPort: "3306",
//		},
//	}
//}

//func (cfg *ServerConfig) GetServerAddress() string {
//	return cfg.ServerAddress + ":" + cfg.ServerPort
//}
//
//func (cfg *GatewayConfig) GetGatewayAddress() string {
//	return cfg.GatewayAddress + ":" + cfg.GatewayPort
//}
