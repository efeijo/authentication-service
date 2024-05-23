package config

import (
	"fmt"
	"os"
)

type Config struct {
	RedisConfig       RedisConfig
	AuthServiceConfig AuthServiceConfig
}

// RedisConfig holds the configuration for the redis server.
type RedisConfig struct {
	// The address of the redis server
	Host string `json:"address,omitempty"`
	// the port of the redis server
	Port string `json:"port,omitempty"`
	// ttl for the redis keys in seconds
	TTL int `json:"ttl,omitempty"`
}

// AuthServiceConfig holds the configuration for the auth service.
type AuthServiceConfig struct {
	// The secret to use for JWT
	Secret string `json:"secret,omitempty"`
	// The Port to bind to
	Port string `json:"port,omitempty"`
}

// LoadConfigFromEnv loads the configuration from environment variables.
func LoadConfigFromEnv() *Config {
	config := &Config{
		RedisConfig: RedisConfig{
			Host: os.Getenv("REDIS_HOST"),
			Port: os.Getenv("REDIS_PORT"),
			TTL:  3600,
		},
		AuthServiceConfig: AuthServiceConfig{
			Secret: os.Getenv("JWT_SECRET"),
			Port:   os.Getenv("AUTH_PORT"),
		},
	}
	fmt.Printf("config: %+v\n", config)

	return config
}
