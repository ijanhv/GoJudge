package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gojudge/db"
	"gojudge/generator"
	"gojudge/models"
	"net/http"
	"strings"
	"time"
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

func GetSubmission(c* gin.Context) {
	submissionId := c.Param("id")

	var submission models.Submission
	if err := db.GetDB().Preload("TestResults").Where("id = ?", submissionId).First(&submission).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return

	}


	c.JSON(http.StatusAccepted, gin.H{"status": "Submission fetched Successfully!", "submission": submission})


}

// Submission handles the submission of a problem solution.
func CreateSubmission(c *gin.Context) {
	var submission models.Submission

	user, _ := c.Get("currentUser")

	submission.UserID = user.(models.User).ID

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
			Status:       "pending",
			Output:       "",
			ErrorMessage: "",
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
	submissionID := c.Param("id")

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

func GetResults(c *gin.Context) {
    submissionID := c.Param("id")
    fmt.Printf("Submission ID: %s\n", submissionID)

    // Bind the incoming JSON array to a slice of TestResult structs
    var results []models.TestResult
    if err := c.ShouldBindJSON(&results); err != nil {
        fmt.Printf("ERROR BIND: %s\n", err.Error())
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    fmt.Printf("Received Results: %+v\n", results)

    // Fetch the existing submission from the database
    var existingSubmission models.Submission
    if err := db.GetDB().Where("id = ?", submissionID).Preload("TestResults").First(&existingSubmission).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
        return
    }

    fmt.Printf("Existing Submission: %+v\n", existingSubmission)

    allSuccess := true // Variable to keep track of whether all test results are successful

    for _, result := range results {
        var existingTestResult models.TestResult
        if err := db.GetDB().Where("submission_id = ? AND test_case_id = ?", result.SubmissionID, result.TestCaseID).First(&existingTestResult).Error; err != nil {
            fmt.Printf("Test result not found for submissionID: %d, testCaseId: %d\n", result.SubmissionID, result.TestCaseID)
            continue
        }

        // Update the fields of the existing test result
        existingTestResult.Status = result.Status
        existingTestResult.Output = result.Output
        existingTestResult.ErrorMessage = result.ErrorMessage

        // Check if the status is not 'Success', set allSuccess to false
        if result.Status != "Success" {
            allSuccess = false
        }

        // Save the updated test result back to the database
        if err := db.GetDB().Save(&existingTestResult).Error; err != nil {
            fmt.Printf("Error updating test result: %s\n", err.Error())
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating test result"})
            return
        }
    }

    // Update the submission status based on the test results
    if allSuccess {
        existingSubmission.Status = "Accepted"
    } else {
        existingSubmission.Status = "Rejected"
    }

    // Save the updated submission status back to the database
    if err := db.GetDB().Save(&existingSubmission).Error; err != nil {
        fmt.Printf("Error updating submission status: %s\n", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating submission status"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "Test results updated successfully!"})
}
