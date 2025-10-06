package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	HTTP_PORT string

	POSTGRES_HOST     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string
	POSTGRES_PORT     string

	DB *gorm.DB
)

func InitialENV() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	HTTP_PORT = requiredEnv("HTTP_PORT")
	POSTGRES_HOST = requiredEnv("POSTGRES_HOST")
	POSTGRES_USER = requiredEnv("POSTGRES_USER")
	POSTGRES_PASSWORD = requiredEnv("POSTGRES_PASSWORD")
	POSTGRES_DB = requiredEnv("POSTGRES_DB")
	POSTGRES_PORT = requiredEnv("POSTGRES_PORT")
}

func ConnectDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_PORT)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
}

func requiredEnv(key string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("required env %s not set", key)
	}
	return env
}
