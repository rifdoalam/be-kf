package main

import (
	"be-api/internal/db"
	"be-api/internal/routes"
	"be-api/pkg/config"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	config.Load()
	

	
	database := db.Init()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // Allow methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allow headers
		AllowCredentials: true,
	}))
	routes.SetupRoutes(r, database)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start the server: ", err)
	}
}
