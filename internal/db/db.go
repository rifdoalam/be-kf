package db

import (
	"be-api/internal/models"
	"be-api/pkg/config"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Init initializes and returns a MySQL database connection
func Init() *gorm.DB {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", 
		config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{})
	return db
}

