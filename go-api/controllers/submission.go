package controllers

import (
	"context"
	"encoding/json"
	"gojudge/models"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Submission handles the submission of a problem solution.
func Submission(c *gin.Context) {
	var submission models.Submission

	// Bind JSON input to the Submission struct
	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	submission.CreatedAt = now
	submission.UpdatedAt = now

	submissionJSON, err := json.Marshal(submission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error serializing submission"})
		return
	}

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "supersecretpassword",
		DB:       0, // use default DB
	})

	ctx := context.Background()

	// Add the submission to the queue
	err = rdb.RPush(ctx, "submission_queue", submissionJSON).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding submission to queue"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "Submission queued Successfully!"})
}
