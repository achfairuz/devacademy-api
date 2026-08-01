package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Enabled  bool
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret string
	Expiry int
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_MODE", "debug")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "devacademy")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("REDIS_ENABLED", false)
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("JWT_SECRET", "supersecret")
	viper.SetDefault("JWT_EXPIRY", 24)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("error reading config file: %v", err)
		}
	}

	return &Config{
		Server: ServerConfig{
			Port: viper.GetString("APP_PORT"),
			Mode: viper.GetString("APP_MODE"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Enabled:  viper.GetBool("REDIS_ENABLED"),
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
			Expiry: viper.GetInt("JWT_EXPIRY"),
		},
	}
}

func (c *DatabaseConfig) DSN() string {
	return "host=" + c.Host +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.Name +
		" port=" + c.Port +
		" sslmode=" + c.SSLMode
}
