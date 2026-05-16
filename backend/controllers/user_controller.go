package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUsers returns all users
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// UpdateUserRole allows an admin to change a user's role
func UpdateUserRole(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Role != "teacher" && input.Role != "student" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update role
	user.Role = input.Role
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully", "user": user})
}

// DeleteUser deletes a user and all their associated scores
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Mulai transaksi untuk menghapus skor dan pengguna agar aman
	tx := database.DB.Begin()

	// 1. Hapus semua skor pengguna berdasarkan username
	if err := tx.Where("username = ?", user.Username).Delete(&models.Score{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user's scores"})
		return
	}

	// 2. Hapus pengguna
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// 3. Catat aktivitas di log
	database.DB.Create(&models.ActivityLog{
		Username: user.Username,
		Action:   "DELETED_BY_ADMIN",
	})

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "User and their scores deleted successfully"})
}

// DeleteMe allows a logged-in user to delete their own account and scores
func DeleteMe(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Mulai transaksi
	tx := database.DB.Begin()

	// 1. Hapus semua skor pengguna
	if err := tx.Where("username = ?", user.Username).Delete(&models.Score{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete scores"})
		return
	}

	// 2. Hapus pengguna
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	// 3. Catat aktivitas di log
	database.DB.Create(&models.ActivityLog{
		Username: user.Username,
		Action:   "DELETE_ACCOUNT",
	})

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
