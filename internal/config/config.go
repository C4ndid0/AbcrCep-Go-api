package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerAddr string
	Timeout    time.Duration
}

func LoadConfig() Config {
	timeoutStr := getEnv("TIMEOUT_MS", "30000")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeout = 30000
	}

	return Config{
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
		Timeout:    time.Duration(timeout) * time.Millisecond,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
