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
	Kafka    KafkaConfig
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

type KafkaConfig struct {
	Broker            string
	AuditTopic        string
	ReplicationFactor int16
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
		"postgres://flagr:flagr@localhost:5432/flagr?sslmode=disable")

	viper.SetDefault("redis.addr", "localhost:6380")
	viper.SetDefault("redis.password", "")

	viper.SetDefault("kafka.broker", "localhost:9092")
	viper.SetDefault("kafka.audit_topic", "flag.audit")
	viper.SetDefault("kafka.replication_factor", 1)

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
		Kafka: KafkaConfig{
			Broker:            viper.GetString("kafka.broker"),
			AuditTopic:        viper.GetString("kafka.audit_topic"),
			ReplicationFactor: int16(viper.GetInt("kafka.replication_factor")),
		},
		Auth: AuthConfig{
			JWTSecret:            viper.GetString("auth.jwt_secret"),
			AccessTokenDuration:  viper.GetDuration("auth.access_token_duration"),
			RefreshTokenDuration: viper.GetDuration("auth.refresh_token_duration"),
		},
	}
}
