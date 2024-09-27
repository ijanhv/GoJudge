package models

import (
	"time"

	"gorm.io/gorm"
)

type Problem struct {
	gorm.Model
	ID              uint       `gorm:"primaryKey" json:"id"`
	Name            string     `json:"name"`
	FunctionName    string     `json:"function_name"`
	InputStructure  string     `json:"input_structure"`
	OutputStructure string     `json:"output_structure"`
	Description     string     `json:"description"`
	TestCases       []TestCase `gorm:"foreignKey:ProblemID" json:"test_cases"` // Establishing relationship
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
