package utils

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func GetInt64(key string, fallback int64) int64 {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
	}
	return fallback
}

func GetDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if i, err := time.ParseDuration(value); err == nil {
			return i
		}
	}
	return fallback
}

func ParseList(value string) []string {
	var result []string
	for _, item := range strings.Split(value, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
