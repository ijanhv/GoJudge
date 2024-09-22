package models

import "time"

type TestCase struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProblemID uint      `gorm:"not null" json:"problem_id"`
	Input     string    `gorm:"type:text;not null" json:"input"`
	Expected  string    `gorm:"type:text;not null" json:"expected"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
