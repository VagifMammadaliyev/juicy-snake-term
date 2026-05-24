package config

import (
	"os"
	"strconv"
)

const (
	configEnvKeyPrefix = "JST_"
	serverHostEnvKey   = "SERVER_HOST"
	serverPortEnvKey   = "SERVER_PORT"
)

type ServerConfig struct {
	Host string
	Port int
}

type GameConfig struct {
	Server ServerConfig
}

const (
	defaultHost = "localhost"
	defaultPort = 21000
)

func NewGameConfig() GameConfig {
	return GameConfig{
		Server: newServerConfig(),
	}
}

func newServerConfig() ServerConfig {
	return ServerConfig{
		Host: parseEnvString(envKey(serverHostEnvKey), defaultHost),
		Port: parseEnvInt(envKey(serverPortEnvKey), defaultPort),
	}
}

func envKey(key string) string {
	return configEnvKeyPrefix + key
}

func parseEnvInt(key string, defaultValue int) int {
	valueStr, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func parseEnvString(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return value
}
