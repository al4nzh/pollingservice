package config

import (
	"log"
	"os"
	"strconv"
	"time"
)
type Config struct {
	Port         string
	DBDSN        string
	PollInterval time.Duration
	WorkerCount  int
	HTTPTimeout  time.Duration
	CSFloatAPIKey string
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Load() Config {
	port := getEnv("PORT", "8080")
	dbdsn := getEnv("DB_DSN", "postgres://postgres:postgres@localhost:5432/pollingservice")

	pollIntervalStr := getEnv("POLL_INTERVAL", "60") // seconds
	pollIntervalSec, err := strconv.Atoi(pollIntervalStr)
	if err != nil {
		log.Fatal("invalid POLL_INTERVAL")
	}

	workerCountStr := getEnv("WORKER_COUNT", "5")
	workerCount, err := strconv.Atoi(workerCountStr)
	if err != nil {
		log.Fatal("invalid WORKER_COUNT")
	}

	httpTimeoutStr := getEnv("HTTP_TIMEOUT", "5") // seconds
	httpTimeoutSec, err := strconv.Atoi(httpTimeoutStr)
	if err != nil {
		log.Fatal("invalid HTTP_TIMEOUT")
	}

	return Config{
		Port:         port,
		DBDSN:        dbdsn,
		PollInterval: time.Duration(pollIntervalSec) * time.Second,
		WorkerCount:  workerCount,
		HTTPTimeout:  time.Duration(httpTimeoutSec) * time.Second,
		CSFloatAPIKey: getEnv("CSFLOAT_API_KEY", ""),
	}
}