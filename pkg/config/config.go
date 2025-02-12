package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)
var (
	DbUser     string
	DbPassword string
	DbName     string
	DbHost     string
	DbPort     string
)
func Load() {
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbName = os.Getenv("DB_NAME")
	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
}
