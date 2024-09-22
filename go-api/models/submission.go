package models

import "time"

type Submission struct {
	Language  string    `json:"language"`
	Code      string    `json:"code"`
	Input     string    `json:"input,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
