package controllers

import (
	"context"
	"encoding/json"

	"fmt"

	"github.com/gin-gonic/gin"

	"gojudge/db"
	"gojudge/generator"
	"gojudge/models"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetBoilerplate(language string) (string, error) {
	switch strings.ToLower(language) {
	case "cpp", "c++":
		return generator.GenerateFullCPlusPlusBoilerplate(), nil
	case "java":
		return generator.GenerateFullJavaBoilerplate(), nil
	case "js", "javascript":
		return generator.GenerateFullJavaScriptBoilerplate(), nil
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}
}

// Submission handles the submission of a problem solution.
func CreateSubmission(c *gin.Context) {
	var submission models.Submission

	// Bind JSON input to the Submission struct
	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	boilerplate, err := GetBoilerplate(submission.Language)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var problem models.Problem
	if err := db.GetDB().Preload("Function").Preload("Function.Parameters").Preload("TestCases").Where("id = ?", submission.ProblemID).First(&problem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
		return
	}

	submission.Problem = problem

	finalCode := strings.Replace(boilerplate, "#USER CODE HERE#", submission.Code, 1)
	submission.Code = finalCode
	now := time.Now()
	submission.SubmissionTime = now
	submission.CreatedAt = now
	submission.UpdatedAt = now

	for _, testCase := range problem.TestCases {
		submission.TestResults = append(submission.TestResults, models.TestResult{
			TestCaseID:   testCase.ID,
			Status:       "pending", // Initial status can be set to "pending"
			Output:       "",        // Output can be left empty for now
			ErrorMessage: "",        // No error message initially
		})
	}

	// Instead of marshaling submission to JSON, just save it directly
	if result := db.GetDB().Create(&submission); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save submission: " + result.Error.Error()})
		return
	}

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "supersecretpassword",
		DB:       0,
	})

	ctx := context.Background()

	// Marshal submission to JSON to add to Redis queue
	submissionJSON, err := json.Marshal(submission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error serializing submission"})
		return
	}

	// Add the submission to the queue
	err = rdb.RPush(ctx, "submission_queue", submissionJSON).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding submission to queue"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "Submission queued Successfully!", "submission": submission})
}

func UpdateSubmission(c *gin.Context) {
	submissionID := c.Param("id") // Get the submission ID from the URL parameters

	// Create a struct to hold the fields that can be updated
	var updateData struct {
		Status       string              `json:"status"`
		ErrorMessage string              `json:"errorMessage"`
		TestResults  []models.TestResult `json:"testCaseResults"`
	}

	// Bind JSON input to the updateData struct
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the existing submission from the database
	var existingSubmission models.Submission
	if err := db.GetDB().Where("id = ?", submissionID).First(&existingSubmission).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	// Update the specified fields
	if updateData.Status != "" {
		existingSubmission.Status = updateData.Status
	}
	if updateData.ErrorMessage != "" {
		existingSubmission.ErrorMessage = updateData.ErrorMessage
	}

	// Update test case results
	if len(updateData.TestResults) > 0 {
		// Clear existing test results and update with new data
		existingSubmission.TestResults = updateData.TestResults
	}

	// Save the updated submission to the database
	if err := db.GetDB().Save(&existingSubmission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating submission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Submission updated successfully!", "submission": existingSubmission})
}
