package config

type Config struct {
	DBPath        string
	Port          int
}

func LoadConfig() *Config {
	return &Config{
		DBPath:        "./employees.db",
		Port:          8080,
	}
} 