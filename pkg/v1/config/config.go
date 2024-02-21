package config

import (
	"time"
)

func CreateConfig() Config {
	return Config{
		Port:          8080,
		CheckDuration: 10 * time.Second,
	}
}

type Config struct {
	Port          int
	CheckDuration time.Duration
}
