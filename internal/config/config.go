package config

import (
	"os"
	"sync"
)

var (
	config *Config
	once   sync.Once
)

// Config is global config.
type Config struct {
	*Server
	*PostgreSQL
	*JWT
}

// Server is server config.
type Server struct {
	Addr string
}

// PostgreSQL is PostgreSQL config.
type PostgreSQL struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWT is Json Web Token config.
type JWT struct {
	Secret string
}

// Get reads config once and returns it.
func Get() *Config {
	once.Do(func() {
		config = &Config{
			Server: &Server{
				Addr: getEnv("SERVER_ADDR", ":8080"),
			},
			PostgreSQL: &PostgreSQL{
				Host:     getEnv("POSTGRES_HOST", "locahost"),
				Port:     getEnv("POSTGRES_PORT", "5432"),
				User:     getEnv("POSTGRES_USER", "postgres"),
				Password: getEnv("POSTGRES_PASSWORD", "123"),
				DBName:   getEnv("POSTGRES_DBNAME", "jwt"),
				SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
			},
			JWT: &JWT{
				Secret: getEnv("JWT_SECRET", "jwt_secret"),
			},
		}
	})

	return config
}

// getEnv is the os.Getenv but with default value if environment variable wasn't set.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}

	return value
}
