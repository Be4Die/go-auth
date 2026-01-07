package config

import (
	"os"
	"strconv"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Security SecurityConfig
}

type AppConfig struct {
	Name        string
	Environment string
	LogLevel    string
}

type HTTPConfig struct {
	Port string
}

type PostgresConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
}

type SecurityConfig struct {
	BcryptCost int
}

func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "go-auth"),
			Environment: getEnv("APP_ENV", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
		HTTP: HTTPConfig{
			Port: getEnv("HTTP_PORT", "8080"),
		},
		Postgres: PostgresConfig{
			DSN: getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/auth?sslmode=disable"),
		},
		Redis: RedisConfig{
			Addr: getEnv("REDIS_ADDR", "localhost:6379"),
		},
		JWT: JWTConfig{
			AccessSecret:  getEnv("JWT_ACCESS_SECRET", "super-secret-access-key"),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", "super-secret-refresh-key"),
		},
		Security: SecurityConfig{},
	}

	if v := os.Getenv("BCRYPT_COST"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Security.BcryptCost = n
		}
	}

	if cfg.App.Environment == "production" {
		if os.Getenv("JWT_ACCESS_SECRET") == "" || os.Getenv("JWT_REFRESH_SECRET") == "" || os.Getenv("DATABASE_URL") == "" {
			return nil, ErrMissingProdEnv
		}
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var ErrMissingProdEnv = Err("missing required environment variables for production")

type Err string

func (e Err) Error() string { return string(e) }
