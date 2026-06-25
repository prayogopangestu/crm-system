package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Env      string `yaml:"env"`
		BaseURL  string `yaml:"base_url"`
		Timezone string `yaml:"timezone"`
		LogLevel string `yaml:"log_level"`
	} `yaml:"app"`
	HTTP struct {
		Addr           string   `yaml:"addr"`
		AllowedOrigins []string `yaml:"allowed_origins"`
	} `yaml:"http"`
	GRPC struct {
		Addr string `yaml:"addr"`
	} `yaml:"grpc"`
	Database struct {
		URL      string `yaml:"url"`
		MaxConns int32  `yaml:"max_conns"`
		MinConns int32  `yaml:"min_conns"`
	} `yaml:"database"`
	Redis struct {
		URL string `yaml:"url"`
	} `yaml:"redis"`
	Auth struct {
		JWTSecret  string        `yaml:"jwt_secret"`
		JWTTTL     time.Duration `yaml:"-"`
		JWTTTLText string        `yaml:"jwt_ttl"`
		BcryptCost int           `yaml:"bcrypt_cost"`
	} `yaml:"auth"`
	Security struct {
		EncryptionKey string `yaml:"encryption_key"`
	} `yaml:"security"`
	Telegram struct {
		WorkerInterval     time.Duration `yaml:"-"`
		WorkerIntervalText string        `yaml:"worker_interval"`
		WorkerBatchSize    int           `yaml:"worker_batch_size"`
	} `yaml:"telegram"`
}

func Load(path string) (Config, error) {
	var cfg Config
	raw, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return cfg, fmt.Errorf("decode config: %w", err)
	}

	override(&cfg)
	if cfg.Auth.JWTTTL, err = time.ParseDuration(cfg.Auth.JWTTTLText); err != nil {
		return cfg, fmt.Errorf("invalid jwt ttl: %w", err)
	}
	if cfg.Telegram.WorkerInterval, err = time.ParseDuration(cfg.Telegram.WorkerIntervalText); err != nil {
		return cfg, fmt.Errorf("invalid telegram worker interval: %w", err)
	}
	if len(cfg.Auth.JWTSecret) < 32 {
		return cfg, errors.New("JWT_SECRET must be at least 32 characters")
	}
	key, err := base64.StdEncoding.DecodeString(cfg.Security.EncryptionKey)
	if err != nil || len(key) != 32 {
		return cfg, errors.New("APP_ENCRYPTION_KEY must be base64 encoded 32 bytes")
	}
	if _, err := time.LoadLocation(cfg.App.Timezone); err != nil {
		return cfg, fmt.Errorf("invalid timezone: %w", err)
	}
	return cfg, nil
}

func override(cfg *Config) {
	setString("APP_ENV", &cfg.App.Env)
	setString("APP_BASE_URL", &cfg.App.BaseURL)
	setString("APP_TIMEZONE", &cfg.App.Timezone)
	setString("LOG_LEVEL", &cfg.App.LogLevel)
	setString("HTTP_ADDR", &cfg.HTTP.Addr)
	setString("GRPC_ADDR", &cfg.GRPC.Addr)
	setString("DATABASE_URL", &cfg.Database.URL)
	setString("REDIS_URL", &cfg.Redis.URL)
	setString("JWT_SECRET", &cfg.Auth.JWTSecret)
	setString("JWT_TTL", &cfg.Auth.JWTTTLText)
	setString("APP_ENCRYPTION_KEY", &cfg.Security.EncryptionKey)
	setString("TELEGRAM_WORKER_INTERVAL", &cfg.Telegram.WorkerIntervalText)
	if v := os.Getenv("CORS_ALLOWED_ORIGINS"); v != "" {
		cfg.HTTP.AllowedOrigins = strings.Split(v, ",")
	}
	setInt32("DATABASE_MAX_CONNS", &cfg.Database.MaxConns)
	setInt32("DATABASE_MIN_CONNS", &cfg.Database.MinConns)
	setInt("BCRYPT_COST", &cfg.Auth.BcryptCost)
	setInt("TELEGRAM_WORKER_BATCH_SIZE", &cfg.Telegram.WorkerBatchSize)
}

func setString(key string, dst *string) {
	if value := os.Getenv(key); value != "" {
		*dst = value
	}
}

func setInt(key string, dst *int) {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			*dst = parsed
		}
	}
}

func setInt32(key string, dst *int32) {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 32); err == nil {
			*dst = int32(parsed)
		}
	}
}
