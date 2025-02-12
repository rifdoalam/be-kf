package handlers

import (
	"be-api/internal/dtos"
	"be-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUsers retrieves all users from the database
func GetUsers(c *gin.Context, db *gorm.DB) {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	var usersResponse []dtos.GetUsers
	for _, user := range users {
		usersResponse = append(usersResponse, dtos.GetUsers{
			Name:  user.Name,
			Email: user.Email,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": usersResponse, "status": "success", "message": "Users fetched successfully"})
}

// CreateUser creates a new user in the database
func CreateUser(c *gin.Context, db *gorm.DB) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := db.Where("email = ?", user.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user, "message": "User created successfully"} )
}
