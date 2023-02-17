package utils

import (
	"os"
	"strings"
)

func GetEnvOrFallback(key string, fallback string) string {
	if env, set := os.LookupEnv(key); set {
		return env
	}
	return fallback
}

func FlagToEnv(s string) string {
	return strings.Replace(strings.ToUpper(s), "-", "_", -1)
}
