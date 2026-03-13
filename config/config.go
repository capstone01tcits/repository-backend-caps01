package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Config struct {
	AppPort               string
	AppEnv                string
	DBHost                string
	DBPort                string
	DBUser                string
	DBPassword            string
	DBName                string
	JWTSecret             string
	JWTExpireHours        string
	JWTRefreshSecret      string
	JWTRefreshExpireHours string
	AIServiceURL          string
}

var Cfg Config

func Load() {
	_ = godotenv.Load()

	appEnv := getEnv("APP_ENV", "development")
	jwtSecret := getEnv("JWT_SECRET", "")
	jwtRefreshSecret := getEnv("JWT_REFRESH_SECRET", "")

	if jwtSecret == "" {
		if appEnv == "production" {
			log.Fatal("ERROR: JWT_SECRET must be set in production environment. Please set JWT_SECRET in .env")
		}
		jwtSecret = "dev-secret-change-in-production"
		log.Println("WARNING: Using default JWT secret. Set JWT_SECRET in .env for production")
	}

	if jwtRefreshSecret == "" {
		if appEnv == "production" {
			log.Fatal("ERROR: JWT_REFRESH_SECRET must be set in production environment. Please set JWT_REFRESH_SECRET in .env")
		}
		jwtRefreshSecret = "dev-refresh-secret-change-in-production"
		if appEnv != "production" {
			log.Println("WARNING: Using default JWT refresh secret. Set JWT_REFRESH_SECRET in .env for production")
		}
	}

	Cfg = Config{
		AppPort:               getEnv("APP_PORT", "3000"),
		AppEnv:                appEnv,
		DBHost:                getEnv("DB_HOST", "localhost"),
		DBPort:                getEnv("DB_PORT", "5432"),
		DBUser:                getEnv("DB_USER", "postgres"),
		DBPassword:            getEnv("DB_PASSWORD", ""),
		DBName:                getEnv("DB_NAME", "go_auth"),
		JWTSecret:             jwtSecret,
		JWTExpireHours:        getEnv("JWT_EXPIRE_HOURS", "24"),
		JWTRefreshSecret:      jwtRefreshSecret,
		JWTRefreshExpireHours: getEnv("JWT_REFRESH_EXPIRE_HOURS", "168"),
		AIServiceURL:          getEnv("AI_SERVICE_URL", "http://localhost:8000"),
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

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
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
