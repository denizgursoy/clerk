package config

import (
	"time"
)

func CreateConfig() Config {
	return Config{
		Port:             8080,
		CheckDuration:    10 * time.Second,
		LifeSpanDuration: 30 * time.Second,
		ETCDEndpoint:     "localhost:2379",
	}
}

type Config struct {
	Port             int
	CheckDuration    time.Duration
	LifeSpanDuration time.Duration
	ETCDEndpoint     string
}
