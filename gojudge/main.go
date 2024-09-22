package main

import (
	"context"
	"encoding/json"
	"gojudge/docker"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

type Submission struct {
	Language       string    `json:"language"`
	Code           string    `json:"code"`
	Input          string    `json:"input,omitempty"`
	ExpectedOutput string    `json:"expected_output"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func main() {
	// initialize redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "supersecretpassword",
		DB:       0,
	})

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to listen for shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Start the worker
	go worker(ctx, rdb)

	// Wait for shutdown signal
	<-shutdown
	log.Println("Shutting down gracefully...")

}

func worker(ctx context.Context, rdb *redis.Client) {
	for {
		select {
		case <-ctx.Done():
			return

		default:
			result, err := rdb.BLPop(ctx, 0, "submission_queue").Result()

			if err != nil {
				log.Printf("Error popping from queue: %v", err)
				continue
			}

			submissionJSON := result[1]

			var submission Submission

			err = json.Unmarshal([]byte(submissionJSON), &submission)

			if err != nil {
				log.Printf("Error unmarshaling submission: %v", err)
				continue
			}

			processSubmission(submission)

		}
	}
}

func processSubmission(submission Submission) {

	success, result, err := docker.RunCodeInContainer(submission.Language, submission.Code, submission.Input, submission.ExpectedOutput)
	if err != nil {
		log.Printf("Error running submission: %v", err)
		return
	}

	log.Printf("Code executed successfully. Result: %s Success: %t", result, success)

}
