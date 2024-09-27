package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gojudge/docker"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

type BaseModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Submission struct {
	Problem   Problem    `json:"problem"`
	Code      string     `json:"code"`
	Language  string     `json:"language"`
	TestCases []TestCase `json:"testCases"`
}

type Problem struct {
	BaseModel
	Title       string            `gorm:"type:varchar(255);not null" json:"title"`
	Description string            `gorm:"type:text;not null" json:"description"`
	Difficulty  string            `gorm:"type:varchar(50);not null" json:"difficulty"`
	Tags        []string          `gorm:"type:varchar(255);" json:"tags"`
	Author      string            `gorm:"type:varchar(255);not null" json:"author"`
	Function    FunctionSignature `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE;" json:"function"`
	TestCases   []TestCase        `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE;" json:"testCases"`
}

type FunctionSignature struct {
	BaseModel
	ProblemID    uint        `gorm:"not null" json:"problemId"`
	FunctionName string      `gorm:"type:varchar(100);not null" json:"functionName"`
	Parameters   []Parameter `gorm:"foreignKey:SignatureID;constraint:OnDelete:CASCADE;" json:"parameters"`
	ReturnType   string      `gorm:"type:varchar(50);not null" json:"returnType"`
}

type Parameter struct {
	BaseModel
	SignatureID uint   `gorm:"not null" json:"signatureId"`
	Name        string `gorm:"type:varchar(50);not null" json:"name"`
	Type        string `gorm:"type:varchar(50);not null" json:"type"`
}

type TestCase struct {
	BaseModel
	ProblemID uint                   `gorm:"not null" json:"problemId"`
	Input     map[string]interface{} `gorm:"type:text;not null" json:"input"`
	Output    interface{}            `gorm:"type:text;not null" json:"output"`
}

func main() {
	// Initialize Redis client
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

			if err = json.Unmarshal([]byte(submissionJSON), &submission); err != nil {
				log.Printf("Error unmarshaling submission: %v", err)
				continue
			}

			processSubmission(submission)
		}
	}
}

func processSubmission(submission Submission) {
	fmt.Println(submission.Code, submission.Language)
	// Loop through each test case and inject code one by one
	for i, testCase := range submission.Problem.TestCases {
		// Prepare the submission code with a placeholder for test case inputs
		testCaseCode := generateTestCases(submission.Language, submission.Problem.Function.FunctionName, testCase.Input)
		injectedCode := strings.Replace(submission.Code, "#TEST CASE INPUT#", testCaseCode, 1)
		fmt.Println(testCaseCode, injectedCode)
		// Replace input variable placeholders in the injected code
		for paramName, paramValue := range testCase.Input {
			inputVarName := fmt.Sprintf("input%d_%s", i, paramName)
			injectedCode = strings.Replace(injectedCode, inputVarName, formatValue(submission.Language, paramValue), -1)
		}

		log.Printf("Language: %s, Injected Code for Test Case %d:\n%s", submission.Language, i, injectedCode)

		// Format expected output
		expectedOutput, err := formatExpectedOutput(testCase.Output)
		if err != nil {
			log.Printf("Error formatting expected output for test case %d: %v", i, err)
			continue
		}

		// Run the code in the Docker container
		success, result, err := docker.RunCodeInContainer(submission.Language, injectedCode, formatInput(testCase.Input), expectedOutput)
		if err != nil {
			log.Printf("Error running submission for test case %d: %v", i, err)
			continue
		}

		// Log only the result without additional text
		fmt.Println(result, success)
	}
}

func generateTestCases(lang, functionName string, input map[string]interface{}) string {
	fmt.Println(input)
	var testCode strings.Builder

	// Generate input variable declarations
	for paramName, paramValue := range input {
		testCode.WriteString(generateVariableDeclaration(lang, getType(lang, paramValue), paramName, paramValue))
		testCode.WriteString("\n")
	}

	// Generate function call
	functionCall := generateFunctionCall(lang, functionName, strings.Join(getInputVariableNames(input), ", "))
	testCode.WriteString(functionCall)
	testCode.WriteString("\n")

	return testCode.String()
}

func getInputVariableNames(input map[string]interface{}) []string {
	var names []string
	for paramName := range input {
		names = append(names, paramName)
	}
	return names
}

func formatInput(input map[string]interface{}) string {
	data, err := json.Marshal(input)
	if err != nil {
		log.Printf("Error marshalling input: %v", err)
		return "{}" // Return an empty JSON object on error
	}
	return string(data)
}

func formatExpectedOutput(output interface{}) (string, error) {
	data, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal output: %v", err)
	}
	return string(data), nil
}

func generateVariableDeclaration(lang, type_, name string, value interface{}) string {
	switch lang {
	case "java":
		return fmt.Sprintf("%s %s = %s;", type_, name, formatValue(lang, value))
	case "cpp":
		return fmt.Sprintf("%s %s = %s;", type_, name, formatValue(lang, value))
	case "javascript":
		return fmt.Sprintf("const %s = %s;", name, formatValue(lang, value))
	case "typescript":
		return fmt.Sprintf("const %s: %s = %s;", name, type_, formatValue(lang, value))
	default:
		return fmt.Sprintf("let %s = %s;", name, formatValue(lang, value))
	}
}

func generateFunctionCall(lang, name, args string) string {
	return fmt.Sprintf("%s(%s);", name, args)
}

func formatValue(lang string, value interface{}) string {
	switch v := value.(type) {
	case []interface{}:
		elements := make([]string, len(v))
		for i, elem := range v {
			elements[i] = fmt.Sprint(elem)
		}
		switch lang {
		case "java":
			return fmt.Sprintf("new int[]{%s}", strings.Join(elements, ", "))
		case "cpp":
			return fmt.Sprintf("{%s}", strings.Join(elements, ", "))
		case "javascript", "typescript":
			return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
		}
	}
	return fmt.Sprint(value)
}

func getType(lang string, value interface{}) string {
	switch lang {
	case "java":
		switch value.(type) {
		case []interface{}:
			return "int[]"
		case float64:
			return "int"
		case string:
			return "String"
		case bool:
			return "boolean"
		}
	case "cpp":
		switch value.(type) {
		case []interface{}:
			return "std::vector<int>"
		case float64:
			return "int"
		case string:
			return "std::string"
		case bool:
			return "bool"
		}
	case "javascript", "typescript":
		switch value.(type) {
		case []interface{}:
			return "number[]"
		case float64:
			return "number"
		case string:
			return "string"
		case bool:
			return "boolean"
		}
	}
	return "unknown"
}
