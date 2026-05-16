package middleware

import (
	"backend/database"
	"backend/models"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	jwtKey []byte
	once   sync.Once
)

func loadEnv() {
	once.Do(func() {
		// Load file .env (abaikan jika environment variable sudah diset dari sistem)
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found, using system environment variables")
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			log.Fatal("JWT_SECRET is not set")
		}

		jwtKey = []byte(secret)
	})
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		loadEnv()

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format. Use 'Bearer <token>'"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		hash := sha256.Sum256([]byte(tokenString))
		tokenHash := hex.EncodeToString(hash[:])

		database.DB.Where("expires_at < ?", time.Now()).Delete(&models.BlacklistedToken{})

		var blacklisted models.BlacklistedToken
		if err := database.DB.Where("token_hash = ?", tokenHash).First(&blacklisted).Error; err == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		}

		// simpan ke context
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func TeacherOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")

		if !exists || role != "teacher" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: teacher only"})
			c.Abort()
			return
		}

		c.Next()
	}
}
