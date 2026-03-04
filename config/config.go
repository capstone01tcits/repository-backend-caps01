package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	AppPort              string
	AppEnv               string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	JWTSecret            string
	JWTExpireHours       string
	JWTRefreshSecret     string
	JWTRefreshExpireHours string
}

var Cfg Config

func Load() {
	_ = godotenv.Load()

	Cfg = Config{
		AppPort:              getEnv("APP_PORT", "3000"),
		AppEnv:               getEnv("APP_ENV", "development"),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", ""),
		DBName:               getEnv("DB_NAME", "go_auth"),
		JWTSecret:            getEnv("JWT_SECRET", "secret"),
		JWTExpireHours:       getEnv("JWT_EXPIRE_HOURS", "24"),
		JWTRefreshSecret:     getEnv("JWT_REFRESH_SECRET", "refresh_secret"),
		JWTRefreshExpireHours: getEnv("JWT_REFRESH_EXPIRE_HOURS", "168"),
	}
}

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		Cfg.DBHost, Cfg.DBUser, Cfg.DBPassword, Cfg.DBName, Cfg.DBPort,
	)

	logLevel := logger.Silent
	if Cfg.AppEnv == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	fmt.Println("✓ Database connected!")
	return db
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
