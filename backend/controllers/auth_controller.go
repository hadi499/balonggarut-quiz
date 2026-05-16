package controllers

import (
	"backend/config"
	"backend/database"
	"backend/models"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}



func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username exists
	var count int64
	database.DB.Model(&models.User{}).Where("username = ?", input.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	// Validate password
	if len(input.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}
	var hasUpper, hasLower, hasSpecial bool
	for _, char := range input.Password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:',.<>?/`~\"\\", char):
			hasSpecial = true
		}
	}
	if !hasUpper || !hasLower || !hasSpecial {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one uppercase letter, one lowercase letter, and one special character"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     "student", // Hardcoded default role
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&models.ActivityLog{
		Username: user.Username,
		Action:   "REGISTER",
	})

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful", "user": user})
}

func Login(c *gin.Context) {

	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// ⬇️ INI BAGIAN PENTING
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Username: user.Username,
		Role:     user.Role, // ⬅️ WAJIB ADA
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(config.JWTKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}

func Me(c *gin.Context) {
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"role":     role,
	})
}

func Logout(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is missing"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization format"})
		return
	}

	tokenString := parts[1]

	claims := &Claims{}
	parsedToken, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JWTKey, nil
	})

	hash := sha256.Sum256([]byte(tokenString))
	tokenHash := hex.EncodeToString(hash[:])

	expiresAt := time.Now().Add(24 * time.Hour)
	if parsedToken != nil && parsedToken.Valid {
		if claims.ExpiresAt != nil {
			expiresAt = claims.ExpiresAt.Time
		}
	}

	blacklisted := models.BlacklistedToken{
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	database.DB.Create(&blacklisted)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
