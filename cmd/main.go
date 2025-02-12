package main

import (
	"be-api/internal/db"
	"be-api/internal/routes"
	"be-api/pkg/config"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Initializing MySQL")
	config.Load()

	
	database := db.Init()
	r := gin.Default()
	routes.SetupRoutes(r, database)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start the server: ", err)
	}
}
