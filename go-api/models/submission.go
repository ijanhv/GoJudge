package models

import "time"

type Submission struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `json:"user_id"`
	ProblemID uint       `json:"problem_id"`
	Problem   Problem    `json:"problem"`
	Code      string     `json:"code"`
	Language  string     `json:"language"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
