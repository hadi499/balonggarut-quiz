package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET ALL: Mengambil semua kuis
func GetQuizzes(c *gin.Context) {
	var quizzes []models.Quiz

	// Preload pertanyaan agar data questions ikut terambil
	if err := database.DB.Preload("Questions").Order("id desc").Find(&quizzes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quizzes)
}

// GET ONE: Mengambil satu kuis beserta pertanyaannya
func GetQuiz(c *gin.Context) {
	var quiz models.Quiz
	id := c.Param("id")
	if err := database.DB.Preload("Questions").First(&quiz, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kuis tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, quiz)
}

// CREATE: Menambah kuis baru
func CreateQuiz(c *gin.Context) {
	var newQuiz models.Quiz
	if err := c.ShouldBindJSON(&newQuiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&newQuiz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newQuiz)
}

// UPDATE: Mengubah kuis yang sudah ada
func UpdateQuiz(c *gin.Context) {
	var quiz models.Quiz
	id := c.Param("id")

	// Cek apakah data ada
	if err := database.DB.First(&quiz, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kuis tidak ditemukan"})
		return
	}

	// Bind data baru dari JSON
	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&quiz)
	c.JSON(http.StatusOK, quiz)
}

// DELETE: Menghapus kuis
func DeleteQuiz(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Quiz{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kuis berhasil dihapus"})
}
