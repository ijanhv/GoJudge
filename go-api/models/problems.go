package models

import "time"

type Problem struct {
    ID          uint       `gorm:"primaryKey" json:"id"`
    Title       string     `gorm:"not null" json:"title"`
    Description string     `gorm:"type:text;not null" json:"description"`
    Difficulty  string     `gorm:"not null" json:"difficulty"`
    TestCases   []TestCase `gorm:"foreignKey:ProblemID" json:"test_cases"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}