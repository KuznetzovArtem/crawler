package env

import "syscall"

// GetEnvString reads environment variable
func GetEnvString(key string, fallback string) string {
	value, ok := syscall.Getenv(key)
	if !ok {
		return fallback
	}
	return value
}
