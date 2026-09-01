package config

import (
	"log"
	"os"
	"strconv"
	"time"
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName string
	DBPassword string
	DBUser string
	DBHost string
	DBPort int
	ServerPort string
	Domain string
	Issuer string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	Secure bool
	JWTSecret string
	Version string
}

func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println("Файл .env не найден, используются стандартные настройки")
	}

	version := getEnv("version", "dev")

	dbport, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	dbmax, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "10"))
	maxIdle, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_LIFETIME", "10"))
	lifetime, _ := time.ParseDuration(getEnv("DB_CONN_MAX_LIFETIME", "5m"))

	cfg := &Config{
		DBName: getEnv("DB_NAME", "test"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPort: dbport,
		DBHost: getEnv("DB_HOST", "localhost"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBMaxOpenConns: dbmax,
		DBMaxIdleConns: maxIdle,
		DBConnMaxLifetime: lifetime,
		JWTSecret: getEnv("SECRET", "pass@w0rd"),
		Version: version,
	}

	switch version {
	case "dev":
		cfg.Domain = "localhost"
		cfg.Issuer = "localhost"
		cfg.Secure = false
	case "prod":
		cfg.Domain = getEnv("DOMAIN", "localhost")
		cfg.Issuer = getEnv("ISSUER", "localhost")
		cfg.Secure = true
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}