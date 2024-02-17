package config

func CreateConfig() Config {
	return Config{
		Port: 8080,
	}
}

type Config struct {
	Port int
}
