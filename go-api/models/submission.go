package models

import "time"

type Submission struct {
    ID             uint      `gorm:"primaryKey" json:"id"`
    UserID         uint      `json:"user_id"`
    ProblemID      uint      `json:"problem_id"`
    Code           string    `gorm:"type:text;not null" json:"code"`
    Language       string    `json:"language"`
    Input          string    `json:"input,omitempty"`            // Optional
    ExpectedOutput string    `json:"expected_output,omitempty"`  // Optional
    Status         string    `json:"status"`                     // "Accepted", "Wrong Answer", etc.
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}