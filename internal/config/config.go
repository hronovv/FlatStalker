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
	CORS     CORS
}

type Server struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Database struct {
	URL string
}

type Telegram struct {
	BotToken string
}

type CORS struct {
	Origins []string
}


func MustLoad() *Config {
	_ = godotenv.Load(".env")

	cfg := &Config{
		Server: Server{
			Addr:         getEnv("HTTP_ADDR", ":8080"),
			ReadTimeout:  getDuration("HTTP_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDuration("HTTP_WRITE_TIMEOUT", 10*time.Second),
		},
		Database: Database{
			URL: strings.TrimSpace(os.Getenv("DATABASE_URL")),
		},
		Telegram: Telegram{
			BotToken: strings.TrimSpace(os.Getenv("BOT_TOKEN")),
		},
		CORS: CORS{
			Origins: splitCSV(getEnv(
				"CORS_ORIGINS",
				"http://localhost:5173,https://hronovv.github.io",
			)),
		},
	}

	if cfg.Telegram.BotToken == "" {
		log.Fatal("config: BOT_TOKEN is required")
	}
	if cfg.Database.URL == "" {
		log.Fatal("config: DATABASE_URL is required")
	}
	if cfg.Server.Addr == "" {
		log.Fatal("config: HTTP_ADDR is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	if secs, err := strconv.Atoi(raw); err == nil {
		return time.Duration(secs) * time.Second
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		log.Fatalf("config: invalid %s=%q", key, raw)
	}
	return d
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
