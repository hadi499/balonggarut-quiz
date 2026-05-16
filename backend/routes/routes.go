package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/api/auth/register", controllers.Register)
	r.POST("/api/auth/login", controllers.Login)
	r.GET("/api/quizzes", controllers.GetQuizzes)
	r.GET("/api/quizzes/:id", controllers.GetQuiz)
	r.GET("/api/questions", controllers.GetQuestions)
	r.GET("/api/questions/:id", controllers.GetQuestion)
	r.GET("/me", middleware.AuthMiddleware(), controllers.Me)

	// Grup Siswa
	student := r.Group("/api")
	student.Use(middleware.AuthMiddleware())
	{

		student.POST("/scores", controllers.SubmitScore)
		student.GET("/scores", controllers.GetMyScores)
		student.GET("/leaderboard", controllers.GetAllScores)
		student.DELETE("/users/me", controllers.DeleteMe)
	}

	r.POST("/api/auth/logout", middleware.AuthMiddleware(), controllers.Logout)

	// --- 3. ROUTE GURU / ADMIN (Butuh Token & Role "teacher") ---
	admin := r.Group("/api")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.TeacherOnly())
	{
		// Pindahkan semua aksi berbahaya (CRUD Kuis & Soal) ke sini!
		admin.POST("/quizzes", controllers.CreateQuiz)
		admin.PUT("/quizzes/:id", controllers.UpdateQuiz)
		admin.DELETE("/quizzes/:id", controllers.DeleteQuiz)

		admin.POST("/questions", controllers.CreateQuestion)
		admin.PUT("/questions/:id", controllers.UpdateQuestion)
		admin.DELETE("/questions/:id", controllers.DeleteQuestion)

		// Fitur Manajemen Nilai (Baru)
		admin.GET("/admin/scores", controllers.AdminGetAllScores)
		admin.DELETE("/admin/scores/reset", controllers.DeleteAllScores)

		// Fitur Manajemen User
		admin.GET("/admin/users", controllers.GetUsers)
		admin.PUT("/admin/users/:id/role", controllers.UpdateUserRole)
		admin.DELETE("/admin/users/:id", controllers.DeleteUser)

		// Fitur Log Aktivitas
		admin.GET("/admin/logs", controllers.GetActivityLogs)
	}
}
