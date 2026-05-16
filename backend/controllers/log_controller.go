package controllers

import (
	"backend/database"
	"backend/models"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetActivityLogs returns all user activity logs, ordered by timestamp descending
func GetActivityLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}

	offset := (page - 1) * limit

	query := database.DB.Model(&models.ActivityLog{})

	if search != "" {
		// Use ILIKE for case-insensitive search in PostgreSQL
		query = query.Where("username ILIKE ? OR action ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var logs []models.ActivityLog
	if err := query.Order("timestamp desc").Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity logs"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data":       logs,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}
