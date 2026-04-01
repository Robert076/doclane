package utils

import (
	"log"
	"os"
)

func RequireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s env var is not set", key)
	}
	return val
}
