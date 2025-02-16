package routes

import (
	"be-api/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the API",
		})
	})

	// Auth routes
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", func(c *gin.Context) {
			handlers.RegistrationHandler(c, db)
		})
		authRoutes.POST("/login", func(c *gin.Context) {
			handlers.LoginHandler(c, db)
		})
		
	}

	// Profile routes
	profileRoutes := r.Group("/api/profile")
	{
		
		profileRoutes.GET("/", func(c *gin.Context) {
			handlers.GetProfile(c, db)
		})
		
	}
	
}