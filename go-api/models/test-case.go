package models


type TestCase struct {
    BaseModel
    SubmissionID uint                   `gorm:"not null" json:"submissionId"` // Reference to the submission.
    ProblemID    uint                   `gorm:"not null" json:"problemId"`     // Reference to the problem.
    Input        map[string]interface{} `gorm:"type:text;not null" json:"input"`  // Input as a JSON object.
    Output       interface{}            `gorm:"type:text;not null" json:"output"` // Expected output as a JSON object or array.
}
