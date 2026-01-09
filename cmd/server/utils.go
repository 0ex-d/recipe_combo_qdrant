package main

import (
	"log"
	"os"
	"strconv"
)

func envOrDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func mustEnv(key string) string {
	val := envOrDefault(key, "")
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return val
}

func envIntOrDefault(key string, def int) int {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return def
}

func mustIntEnv(key string) int {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("invalid integer for %s: %v", key, err)
	}
	return parsed
}
