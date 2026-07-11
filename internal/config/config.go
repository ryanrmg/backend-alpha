package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port             string
	DBConnString     string
	ProjectXHttps    string
	ProjectXSocket   string
	ProjectXUsername string
	ProjectXApiKey   string
}

func Load() Config {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	db := os.Getenv("DB_NAME")

	projectx_username := os.Getenv("PROJECTX_USERNAME")
	projectx_api_key := os.Getenv("PROJECTX_API_KEY")

	return Config{
		Port: getEnv("PORT", "8080"),
		DBConnString: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			user,
			password,
			host,
			port,
			db,
		),
		ProjectXHttps:    "https://api.topstepx.com/api",
		ProjectXSocket:   "rtc.topstepx.com/hubs/",
		ProjectXUsername: projectx_username,
		ProjectXApiKey:   projectx_api_key,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
