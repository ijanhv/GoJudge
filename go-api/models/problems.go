package models

import (
	"time"
)

type BaseModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type Problem struct {
	BaseModel
	Title       string            `gorm:"type:varchar(255);not null" json:"title"` // Problem title.
	Slug        string            `gorm:"type:varchar(255);not null" json:"slug"`
	Description string            `gorm:"type:text;not null" json:"description"`                              // Problem description.
	Difficulty  string            `gorm:"type:varchar(50);not null" json:"difficulty"`                        // Difficulty level (e.g., Easy, Medium, Hard).
	Function    FunctionSignature `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE;" json:"function"`  // Function signature.
	TestCases   []TestCase        `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE;" json:"testCases"` // Test cases for the problem.
	Submissions []Submission      `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE;" json:"submissions"` // List of submissions related to the problem.

}

// FunctionSignature represents the function signature for a problem.
type FunctionSignature struct {
	BaseModel
	ProblemID    uint        `gorm:"not null" json:"problemId"`                                             // Reference to the problem.
	FunctionName string      `gorm:"type:varchar(100);not null" json:"functionName"`                        // Name of the function.
	Parameters   []Parameter `gorm:"foreignKey:SignatureID;constraint:OnDelete:CASCADE;" json:"parameters"` // List of function parameters.
	ReturnType   string      `gorm:"type:varchar(50);not null" json:"returnType"`                           // Expected return type of the function.
}

// Parameter represents a parameter of the function signature.
type Parameter struct {
	BaseModel
	SignatureID uint   `gorm:"not null" json:"signatureId"`           // Reference to the function signature.
	Name        string `gorm:"type:varchar(50);not null" json:"name"` // Parameter name.
	Type        string `gorm:"type:varchar(50);not null" json:"type"` // Parameter type (e.g., "int[]", "TreeNode").
}
