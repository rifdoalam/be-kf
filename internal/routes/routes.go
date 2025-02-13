package routes

import (
	"be-api/internal/handlers"
	"be-api/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the API",
		})
	})

	zoomGroup := r.Group("/api/zoom")
	{
		zoomGroup.POST("/create-meeting", func(c *gin.Context) {
			handlers.CreateZoomMeeting(c, db)
		})

		zoomGroup.GET("/meetings", func(c *gin.Context) {
			handlers.GetZoomMeeting(c, db)
		})

		zoomGroup.PUT("/meeting/:meetingId", func(c *gin.Context) {
			handlers.UpdateZoomMeeting(c, db)
		})

		zoomGroup.DELETE("/meeting/:meetingId", func(c *gin.Context) {
			handlers.DeleteZoomMeeting(c, db)
		})

	}

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
