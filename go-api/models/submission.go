package models

import "time"

type Submission struct {
	Language       string    `json:"language"`
	Code           string    `json:"code"`
	Input          string    `json:"input,omitempty"`
	ExpectedOutput string    `json:"expected_output"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
