package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Environment variables
const (
	AppEnv          = "APP_ENV"
	AppName         = "APP_NAME"
	AppPort         = "APP_PORT"

	DbHost       = "DB_HOST"
	DbPort       = "DB_PORT"
	DbUser       = "DB_USER"
	DbName       = "DB_NAME"
	DbPassword   = "DB_PASSWORD"
	DbSslMode    = "DB_SSL_MODE"

	JWTSecret = "JWT_SECRET"
)

type EnvConfig struct {
	AppEnv          string
	AppName         string
	AppPort         string

	DbDetails    string
	DbName       string
	DbSslMode	string

	JWTSecret string
}

// Load environment variables with godotenv and initialize configuration
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func GetConfig() *EnvConfig {
	dbDetails, dbName := getDbDetails()


	return &EnvConfig{
		AppEnv:          os.Getenv(AppEnv),
		AppName:         os.Getenv(AppName),
		AppPort:         os.Getenv(AppPort),

		DbDetails: dbDetails,
		DbName:    dbName,

		JWTSecret: os.Getenv(JWTSecret),
	}
}

func getDbDetails() (string, string) {
	host := os.Getenv(DbHost)
	port := os.Getenv(DbPort)
	user := os.Getenv(DbUser)
	password := os.Getenv(DbPassword)
	name := os.Getenv(DbName)
	sslMode := os.Getenv(DbSslMode)

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		host, user, password, name, port, sslMode), name
}

// Parse environment variable as int, with default value
func getEnvAsInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		val, err := strconv.Atoi(value)
		if err == nil {
			return val
		}
		log.Println("Invalid integer value for environment variable", key)
	}
	return defaultVal
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Parse environment variable as bool, with default value
func getEnvAsBool(key string, defaultVal bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		val, err := strconv.ParseBool(value)
		if err == nil {
			return val
		}
		log.Println("Invalid boolean value for environment variable", key)
	}
	return defaultVal
}
