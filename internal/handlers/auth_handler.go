package handlers

import (
	"be-api/internal/dtos"
	"be-api/internal/models"
	"be-api/internal/utils"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Login(c *gin.Context, db *gorm.DB) {
	var userDto *dtos.LoginRequest;
	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	user := models.User{
		Email:    userDto.Email,
		Password: userDto.Password,
	}
	if err := db.Where("email = ?", user.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	if !user.CheckPassword(userDto.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"token": token, "email": user.Email},
		"success": true,
		"message": "User logged in successfully",
	})
}

func Register(c *gin.Context, db *gorm.DB) {
	var userDto *dtos.RegisterRequest;
	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := db.Where("email = ?", userDto.Email).First( &models.User{	}).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		return
	}
	user := models.User{
		Name:     userDto.Name,
		Email:    userDto.Email,
		Password: userDto.Password,
	}
	if err := user.HashedPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}


	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User registered successfully"})
}