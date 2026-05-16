package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetActivityLogs returns all user activity logs, ordered by timestamp descending
func GetActivityLogs(c *gin.Context) {
	var logs []models.ActivityLog
	if err := database.DB.Order("timestamp desc").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity logs"})
		return
	}
	c.JSON(http.StatusOK, logs)
}
