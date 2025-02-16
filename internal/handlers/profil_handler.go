package handlers

import (
	"be-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProfile (c *gin.Context, db *gorm.DB) {
	// Get the user ID from the context
	userID := c.Value("userID").(uint)

	// Fetch the user from the database
	var user models.User
	result := db.First(&user, userID)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	// Hide the user's email address
	user.Email = ""

	// Return the user
	c.JSON(200, gin.H{"data": user})
}