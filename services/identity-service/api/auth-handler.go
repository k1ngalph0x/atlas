package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/k1ngalph0x/atlas/services/identity-service/config"
	"github.com/k1ngalph0x/atlas/services/identity-service/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
	Config *config.Config
}

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type SignInRequest struct{
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type Claims struct{
	UserId string `json:"user_id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthHandler(db *gorm.DB, config *config.Config) *AuthHandler {
	return &AuthHandler{
		DB: db,
		Config: config,
	}
}

func(a *AuthHandler) GenerateJWT(userId, email string)(string, error){
	expiration := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserId: userId,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer:"atlas-identity",
			Subject: userId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.Config.TOKEN.JwtKey))
	if err != nil{
		return "", err
	}
	return tokenString, nil
}

func(a *AuthHandler) SignUp(c *gin.Context){
	var newUser SignUpRequest
	var existingUser models.User

	err := c.ShouldBindJSON(&newUser)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
	}

	if newUser.Email == "" || newUser.Password == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
        return
	}

	email := strings.ToLower(strings.TrimSpace(newUser.Email))
	password := strings.TrimSpace(newUser.Password)

	result := a.DB.Where("email = ?", email).First(&existingUser)
	if result.Error == nil{
		c.JSON(http.StatusConflict, gin.H{"error":"Email already exists"})
        return 
	}else if !errors.Is(result.Error, gorm.ErrRecordNotFound){
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Database error"})
        return 
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err!= nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create a new user"})
		return 
	}

	user := models.User{
		Email: email,
		Password: string(hashedPassword),
	}

	result = a.DB.Create(&user)
	if result.Error	!= nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create a new user"})
		return 
	}

	token, err := a.GenerateJWT(user.UserID, user.Email)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Something went wrong"})
		return 
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"token":   token,
		"user": gin.H{
			"id":    user.UserID,
			"email": user.Email,
		},
	})
}

func(a *AuthHandler) SignIn(c *gin.Context){
	var user SignUpRequest
	var existingUser models.User

	err := c.ShouldBindJSON(&user)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
	}

	if user.Email == "" || user.Password == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Email and password are required"})
		return 
	} 

	
	email := strings.ToLower(strings.TrimSpace(user.Email))
	password := strings.TrimSpace(user.Password)


	result := a.DB.Where("email = ?", email).First(&existingUser)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err  = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password))
	if err != nil{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := a.GenerateJWT(existingUser.UserID, existingUser.Email)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":"Login successful", 
		"token":token, 
		"email": existingUser.Email,
	})	
}