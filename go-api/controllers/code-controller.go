package controllers

import (
	"context"
	"encoding/json"
	"gojudge/models"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func Submission(w http.ResponseWriter, r *http.Request) {

	var submission models.Submission

	err := json.NewDecoder(r.Body).Decode(&submission)

	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	now := time.Now()

	submission.CreatedAt = now
	submission.UpdatedAt = now

	submissionJSON, err := json.Marshal(submission)

	if err != nil {
		http.Error(w, "Error serializing", http.StatusInternalServerError)
		return
	}

	// create redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "supersecretpassword",
		DB:       0, // use default DB
	})

	ctx := context.Background()

	// Add the submission to the queue
	err = rdb.RPush(ctx, "submission_queue", submissionJSON).Err()
	if err != nil {
		http.Error(w, "Error adding submission to queue", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "Submission queued Successfully!"})

}
