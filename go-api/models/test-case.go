package models

import (
	"encoding/json"
)

type TestCase struct {
	BaseModel
	ProblemID uint   `gorm:"not null;constraint:OnDelete:CASCADE;" json:"problemId"`
	Input     string `gorm:"type:jsonb;not null" json:"input"`  // Store as string
	Output    string `gorm:"type:jsonb;not null" json:"output"` // Store as string
}

// UnmarshalInput method to decode JSON
func (t *TestCase) UnmarshalInput() (map[string]interface{}, error) {
	var input map[string]interface{}
	if err := json.Unmarshal([]byte(t.Input), &input); err != nil {
		return nil, err
	}
	return input, nil
}

// UnmarshalOutput method to decode JSON
func (t *TestCase) UnmarshalOutput() (interface{}, error) {
	var output interface{}
	if err := json.Unmarshal([]byte(t.Output), &output); err != nil {
		return nil, err
	}
	return output, nil
}
