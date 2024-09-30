package models

import (
	"time"
)

type Submission struct {
	BaseModel
	ProblemID      uint         `gorm:"not null;constraint:OnDelete:CASCADE;" json:"problemId"`
	Problem        Problem      `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE;" json:"problem"`
	UserID         uint         `gorm:"not null" json:"userId"`
	TestResults    []TestResult `gorm:"foreignKey:SubmissionID;constraint:OnDelete:CASCADE;" json:"testResults"`
	SubmissionTime time.Time    `json:"submissionTime"`
	Status         string       `gorm:"type:varchar(50);default:'pending';not null" json:"status"`
	ErrorMessage   string       `gorm:"type:text" json:"errorMessage"`
	Language       string       `gorm:"type:text" json:"language"`
	Code           string       `gorm:"type:text" json:"code"`
}

// TestResult represents the result of an individual test case in a submission.
type TestResult struct {
	BaseModel
	SubmissionID uint   `gorm:"not null;constraint:OnDelete:CASCADE;" json:"submissionId"`
	TestCaseID   uint   `gorm:"not null;constraint:OnDelete:CASCADE;" json:"testCaseId"`
	Status       string `gorm:"type:varchar(50);default:'pending';not null" json:"status"`
	Output       string `gorm:"type:text;not null" json:"output"`
	ErrorMessage string `gorm:"type:text" json:"errorMessage"`
}
