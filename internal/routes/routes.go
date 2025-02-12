package routes

import (
	"be-api/internal/handlers"
	"be-api/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/login", func(c *gin.Context) {
			handlers.Login(c, db)
		})
		authGroup.POST("/register", func(c *gin.Context) {
			handlers.Register(c, db)
		})
	}	

	userGroup := r.Group("/api/users")
	{
		userGroup.Use(middlewares.JWTMiddleware)
		userGroup.GET("/", func(c *gin.Context) {
			handlers.GetUsers(c, db)
		})
	}
}
