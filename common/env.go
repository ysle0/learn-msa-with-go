package common

import (
	"fmt"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
)

func EnvString(key, fallback string) string {
	if v, ok := syscall.Getenv(key); ok {
		fmt.Println("key", key, "value", v)
		return v
	}
	return fallback
}
