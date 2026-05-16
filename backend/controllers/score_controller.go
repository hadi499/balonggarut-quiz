// controllers/score_controller.go
package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// Sesuaikan path import ini dengan nama modul project Golang kamu
	// Asumsi config.DB adalah instance GORM kamu
	"backend/database"
	"backend/models"
)

// Struct khusus untuk menangkap body JSON dari Svelte
type ScoreInput struct {
	QuizID uint `json:"quiz_id" binding:"required"`
	Score  int  `json:"score"` // Bisa 0, jadi tidak perlu binding required ketat
}

func SubmitScore(c *gin.Context) {
	// 1. Ambil username dari context yang sudah di-set oleh Middleware JWT
	// Catatan: Pastikan key "username" ini sama dengan yang kamu set di fungsi Middleware kamu.
	usernameContext, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Anda belum login"})
		return
	}
	username := usernameContext.(string)

	// 2. Tangkap data dari Svelte (quiz_id dan score)
	var input ScoreInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// 3. Rakit data untuk disimpan ke Database
	newScore := models.Score{
		Username: username,
		QuizID:   input.QuizID,
		Score:    input.Score,
	}

	// 4. Simpan menggunakan GORM
	if err := database.DB.Create(&newScore).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan skor ke database"})
		return
	}

	// 5. Kirim balasan sukses ke Svelte
	c.JSON(http.StatusOK, gin.H{
		"message": "Skor berhasil disimpan!",
		"data":    newScore,
	})
}

func GetMyScores(c *gin.Context) {
	// 1. Ambil username dari token JWT yang sudah dilewati Middleware
	usernameContext, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak, Anda belum login"})
		return
	}

	// Konversi tipe data ke string
	username := usernameContext.(string)

	// 2. Siapkan variabel untuk menampung daftar skor
	var scores []models.Score

	// 3. Cari di database berdasarkan username
	// Order("created_at desc") berfungsi mengurutkan dari nilai yang paling baru dikerjakan
	if err := database.DB.Where("username = ?", username).Order("created_at desc").Find(&scores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan saat mengambil riwayat nilai"})
		return
	}

	// 4. Kembalikan data skor (berupa array/list) dalam format JSON
	// Jika belum ada nilai, pastikan mengembalikan array kosong [] bukan null
	if scores == nil {
		scores = []models.Score{}
	}
	c.JSON(http.StatusOK, scores)
}

// GetAllScores mengambil semua nilai dari semua user untuk Leaderboard
func GetAllScores(c *gin.Context) {
	var results []struct {
		ID        uint      `json:"id"`
		Username  string    `json:"username"`
		Score     int       `json:"score"`
		QuizTitle string    `json:"quiz_title"` // Kita ambil judulnya dari tabel kuis
		CreatedAt time.Time `json:"created_at"`
	}

	// Menggunakan Joins untuk mendapatkan judul kuis dari tabel quizzes
	err := database.DB.Table("scores").
		Select("scores.id, scores.username, scores.score, scores.created_at, quizzes.title as quiz_title").
		Joins("left join quizzes on quizzes.id = scores.quiz_id").
		Order("scores.score desc, scores.created_at asc").
		Limit(5).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peringkat"})
		return
	}

	if results == nil {
		results = make([]struct {
			ID        uint      `json:"id"`
			Username  string    `json:"username"`
			Score     int       `json:"score"`
			QuizTitle string    `json:"quiz_title"`
			CreatedAt time.Time `json:"created_at"`
		}, 0)
	}
	c.JSON(http.StatusOK, results)
}

// AdminGetAllScores mengambil SEMUA skor tanpa batas untuk manajemen admin
func AdminGetAllScores(c *gin.Context) {
	var results []struct {
		ID        uint      `json:"id"`
		Username  string    `json:"username"`
		Score     int       `json:"score"`
		QuizTitle string    `json:"quiz_title"`
		CreatedAt time.Time `json:"created_at"`
	}

	err := database.DB.Table("scores").
		Select("scores.id, scores.username, scores.score, scores.created_at, quizzes.title as quiz_title").
		Joins("left join quizzes on quizzes.id = scores.quiz_id").
		Order("scores.created_at desc").
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	if results == nil {
		results = make([]struct {
			ID        uint      `json:"id"`
			Username  string    `json:"username"`
			Score     int       `json:"score"`
			QuizTitle string    `json:"quiz_title"`
			CreatedAt time.Time `json:"created_at"`
		}, 0)
	}
	c.JSON(http.StatusOK, results)
}

// DeleteAllScores menghapus seluruh isi tabel scores (Gunakan dengan hati-hati!)
func DeleteAllScores(c *gin.Context) {
	// Menghapus semua baris di tabel scores
	// GORM membutuhkan AllowGlobalUpdate: true untuk menghapus tanpa WHERE clause
	if err := database.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Score{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengosongkan data nilai"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seluruh riwayat nilai berhasil dihapus (Refresh Sukses)"})
}
