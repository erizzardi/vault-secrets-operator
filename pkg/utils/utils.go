package utils

import (
	"math/rand"
	"os"
	"strings"
	"time"
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

func RandSeq(n int, letters []rune) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
