package handlers

import (
	"be-api/internal/dtos"
	"be-api/internal/models"
	"be-api/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


func RegistrationHandler(c *gin.Context, db *gorm.DB) {
	var register dtos.Register
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check email address
	var user models.User
	result := db.Where("email = ?", register.Email).First(&user)
	if result.Error == nil {
		c.JSON(400, gin.H{"error": "Email address already in use"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err) // Handle the error accordingly
	}
	// Create a new user
	user = models.User{
		Name:     register.Name,
		Email:    register.Email,
		Password: string(passwordHash),
		Phone:    register.Phone,
	}
	db.Create(&user)
	c.JSON(201, gin.H{"message": "User created successfully", "data": user})
}

func LoginHandler(c *gin.Context, db *gorm.DB) {
	var login dtos.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check email address
	var user models.User
	result := db.Where("email = ?", login.Email).First(&user)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}
	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		c.JSON(400, gin.H{"error": "Invalid password"})
		return
	}
	// Generate JWT token
	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}
	var dataUser dtos.GetDataUser
	dataUser = dtos.GetDataUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}

	c.JSON(200, gin.H{"message": "Login successful", "data": gin.H{"token": token, "user": dataUser}})
}

