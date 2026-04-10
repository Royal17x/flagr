package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP     HTTPConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Auth     AuthConfig
}

type HTTPConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type PostgresConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr     string
	Password string
}

type AuthConfig struct {
	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func Load() *Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("http.port", "8080")
	viper.SetDefault("http.read_timeout", 15*time.Second)
	viper.SetDefault("http.write_timeout", 15*time.Second)

	viper.SetDefault("postgres.dsn",
		"postgres://flagr:flagr@localhost:5433/flagr?sslmode=disable")

	viper.SetDefault("redis.addr", "localhost:6380")
	viper.SetDefault("redis.password", "")

	viper.SetDefault("auth.jwt_secret", "dev-secret-change-in-prod")
	viper.SetDefault("auth.access_token_duration", 15*time.Minute)
	viper.SetDefault("auth.refresh_token_duration", 7*24*time.Hour)

	return &Config{
		HTTP: HTTPConfig{
			Port:         viper.GetString("http.port"),
			ReadTimeout:  viper.GetDuration("http.read_timeout"),
			WriteTimeout: viper.GetDuration("http.write_timeout"),
		},
		Postgres: PostgresConfig{
			DSN: viper.GetString("postgres.dsn"),
		},
		Redis: RedisConfig{
			Addr:     viper.GetString("redis.addr"),
			Password: viper.GetString("redis.password"),
		},
		Auth: AuthConfig{
			JWTSecret:            viper.GetString("auth.jwt_secret"),
			AccessTokenDuration:  viper.GetDuration("auth.access_token_duration"),
			RefreshTokenDuration: viper.GetDuration("auth.refresh_token_duration"),
		},
	}
}
