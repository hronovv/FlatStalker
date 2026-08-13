package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   Server
	Database Database
	Telegram Telegram
	Worker   Worker
	CORS     CORS
	App      App
}

type Worker struct {
	Free time.Duration
	Plus time.Duration
	Pro  time.Duration
}

type App struct {
	GinMode string
}

type Server struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Database struct {
	URL             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	PingTimeout     time.Duration
}

type Telegram struct {
	BotToken             string
	SupportContact       string
	PaymentProviderToken string
}

type CORS struct {
	Origins []string
}

// MustLoad reads .env and environment. All keys are required — no hardcoded defaults.
func MustLoad() *Config {
	_ = godotenv.Load(".env")

	cfg := &Config{
		App: App{
			GinMode: mustEnv("GIN_MODE"),
		},
		Server: Server{
			Addr:         mustEnv("HTTP_ADDR"),
			ReadTimeout:  mustDuration("HTTP_READ_TIMEOUT"),
			WriteTimeout: mustDuration("HTTP_WRITE_TIMEOUT"),
		},
		Database: Database{
			URL:             mustEnv("DATABASE_URL"),
			MaxConns:        mustInt32("DB_MAX_CONNS"),
			MinConns:        mustInt32("DB_MIN_CONNS"),
			MaxConnLifetime: mustDuration("DB_MAX_CONN_LIFETIME"),
			MaxConnIdleTime: mustDuration("DB_MAX_CONN_IDLE_TIME"),
			PingTimeout:     mustDuration("DB_PING_TIMEOUT"),
		},
		Telegram: Telegram{
			BotToken:             mustEnv("BOT_TOKEN"),
			SupportContact:       mustEnv("TELEGRAM_SUPPORT"),
			PaymentProviderToken: mustEnv("PAYMENT_PROVIDER_TOKEN"),
		},
		Worker: Worker{
			Free: mustDuration("WORKER_INTERVAL_FREE"),
			Plus: mustDuration("WORKER_INTERVAL_PLUS"),
			Pro:  mustDuration("WORKER_INTERVAL_PRO"),
		},
		CORS: CORS{
			Origins: splitCSV(mustEnv("CORS_ORIGINS")),
		},
	}

	if cfg.Database.MaxConns < 1 {
		log.Fatal("config: DB_MAX_CONNS must be >= 1")
	}
	if cfg.Database.MinConns < 0 || cfg.Database.MinConns > cfg.Database.MaxConns {
		log.Fatal("config: DB_MIN_CONNS must be between 0 and DB_MAX_CONNS")
	}
	if len(cfg.CORS.Origins) == 0 {
		log.Fatal("config: CORS_ORIGINS must contain at least one origin")
	}

	return cfg
}

func mustEnv(key string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		log.Fatalf("config: %s is required", key)
	}
	return v
}

func mustDuration(key string) time.Duration {
	raw := mustEnv(key)
	if secs, err := strconv.Atoi(raw); err == nil {
		return time.Duration(secs) * time.Second
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		log.Fatalf("config: invalid %s=%q", key, raw)
	}
	return d
}

func mustInt32(key string) int32 {
	raw := mustEnv(key)
	n, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		log.Fatalf("config: invalid %s=%q", key, raw)
	}
	return int32(n)
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
