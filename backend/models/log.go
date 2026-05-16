package models

import "time"

// ActivityLog records user registration and deletion events
type ActivityLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"not null"`
	Action    string    `json:"action" gorm:"not null"` // e.g., "REGISTER", "DELETE_ACCOUNT", "DELETED_BY_ADMIN"
	Timestamp time.Time `json:"timestamp" gorm:"autoCreateTime"`
}
