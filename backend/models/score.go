package models

import "time"

type Score struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username"`
	QuizID    uint      `json:"quiz_id"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}
