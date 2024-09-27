package models

import (
	"time"

	"gorm.io/gorm"
)


type TestCase struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProblemID uint      `json:"problem_id"` 
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
