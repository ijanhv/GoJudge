package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gojudge/docker"
	"log"
	"net/http"
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
type Problem struct {
	BaseModel
	Title       string            `gorm:"type:varchar(255);not null" json:"title"` // Problem title.
	Slug        string            `gorm:"type:varchar(255);not null" json:"slug"`
	Description string            `gorm:"type:text;not null" json:"description"`                                // Problem description.
	Difficulty  string            `gorm:"type:varchar(50);not null" json:"difficulty"`                          // Difficulty level (e.g., Easy, Medium, Hard).
	Function    FunctionSignature `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE;" json:"function"`    // Function signature.
	TestCases   []TestCase        `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE;" json:"testCases"`   // Test cases for the problem.
	Submissions []Submission      `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE;" json:"submissions"` // List of submissions related to the problem.

}

// FunctionSignature represents the function signature for a problem.
type FunctionSignature struct {
	BaseModel
	ProblemID    uint        `gorm:"not null" json:"problemId"`                                             // Reference to the problem.
	FunctionName string      `gorm:"type:varchar(100);not null" json:"functionName"`                        // Name of the function.
	Parameters   []Parameter `gorm:"foreignKey:SignatureID;constraint:OnDelete:CASCADE;" json:"parameters"` // List of function parameters.
	ReturnType   string      `gorm:"type:varchar(50);not null" json:"returnType"`                           // Expected return type of the function.
}

// Parameter represents a parameter of the function signature.
type Parameter struct {
	BaseModel
	SignatureID uint   `gorm:"not null" json:"signatureId"`           // Reference to the function signature.
	Name        string `gorm:"type:varchar(50);not null" json:"name"` // Parameter name.
	Type        string `gorm:"type:varchar(50);not null" json:"type"` // Parameter type (e.g., "int[]", "TreeNode").
}

type TestCase struct {
	BaseModel
	ProblemID uint   `gorm:"not null;constraint:OnDelete:CASCADE;" json:"problemId"`
	Input     string `gorm:"type:jsonb;not null" json:"input"`  // Store as string
	Output    string `gorm:"type:jsonb;not null" json:"output"` // Store as string
}

type Submission struct {
	BaseModel
	ProblemID      uint         `gorm:"not null;constraint:OnDelete:CASCADE;" json:"problemId"`
	Problem        Problem      `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE;" json:"problem"`
	UserID         uint         `gorm:"not null" json:"userId"`
	TestResults    []TestResult `gorm:"foreignKey:SubmissionID;constraint:OnDelete:CASCADE;" json:"testResults"`
	SubmissionTime time.Time    `json:"submissionTime"`
	Status         string       `gorm:"type:varchar(50);default:'pending';not null" json:"status"`
	ErrorMessage   string       `gorm:"type:text" json:"errorMessage"`
	Language       string       `gorm:"type:text" json:"language"`
	Code           string       `gorm:"type:text" json:"code"`
}

// TestResult represents the result of an individual test case in a submission.
type TestResult struct {
	BaseModel
	SubmissionID uint   `gorm:"not null;constraint:OnDelete:CASCADE;" json:"submissionId"`
	TestCaseID   uint   `gorm:"not null;constraint:OnDelete:CASCADE;" json:"testCaseId"`
	Status       string `gorm:"type:varchar(50);default:'pending';not null" json:"status"`
	Output       string `gorm:"type:text;not null" json:"output"`
	ErrorMessage string `gorm:"type:text" json:"errorMessage"`
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
	var results []TestResult
	// Loop through each test case and inject code one by one
	for i, testCase := range submission.Problem.TestCases {
		testCaseInput, err := testCase.UnmarshalInput()
		if err != nil {
			log.Printf("Error unmarshalling input for test case %d: %v", i, err)
			continue
		}

		// Prepare the submission code with a placeholder for test case inputs
		testCaseCode := generateTestCases(submission.Language, submission.Problem.Function.FunctionName, testCaseInput, submission.Problem.Function.ReturnType)
		injectedCode := strings.Replace(submission.Code, "#TEST CASE INPUT#", testCaseCode, 1)

		fmt.Printf("TEST CASE CODE: %s INJECTED CODE: %s", testCaseCode, injectedCode)

		expectedOutput, err := testCase.UnmarshalOutput()
		if err != nil {
			log.Printf("Error formatting expected output for test case %d: %v", i, err)
			continue
		}

		// Run the code in the Docker container
		success, result, err := docker.RunCodeInContainer(submission.Language, injectedCode, formatInput(testCaseInput), formatValue(submission.Language, expectedOutput))
		if err != nil {
			log.Printf("Error running submission for test case %d: %v", i, err)
			continue
		}

		// Log result
		log.Printf("Test Case id %d - SUCCESS: %t, RESULT: %s", i, success, result)

		testResult := TestResult{
			SubmissionID: submission.ID,
			TestCaseID:   testCase.ID,
			Status:       "pending",
			Output:       result,
			ErrorMessage: "",
		}
		if success {
			testResult.Status = "success"
		} else {
			testResult.Status = "failure"
			testResult.ErrorMessage = "Output did not match expected results" // Set an appropriate message
		}

		results = append(results, testResult)

	}

	if err := SendResults(submission.ID, results); err != nil {
		log.Printf("Error sending results for submission %d: %v", submission.ID, err)
	}

}

func SendResults(submissionID uint, results []TestResult) error {
	url := fmt.Sprintf("http://localhost:8001/api/submission/%d/results", submissionID)

	// Convert results to JSON
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		fmt.Printf("ERROR : %s", err.Error())

		return fmt.Errorf("failed to marshal results: %w", err)
	}

	// Make a POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(resultsJSON))
	if err != nil {
		fmt.Printf("ERROR : %s", err.Error())
		return fmt.Errorf("failed to send results: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	log.Printf("Successfully sent results for submission %d", submissionID)
	return nil
}

func generateTestCases(lang, functionName string, input map[string]interface{}, returnType string) string {
	var testCode strings.Builder

	// Generate input variable declarations
	for paramName, paramValue := range input {
		testCode.WriteString(generateVariableDeclaration(lang, getType(lang, paramValue), paramName, paramValue))
		testCode.WriteString("\n")
	}

	// Generate function call
	functionCall := generateFunctionCall(lang, functionName, strings.Join(getInputVariableNames(input), ", "), returnType)
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
		return "{}"
	}
	return string(data)
}

func formatExpectedOutput(output interface{}) (string, error) {
	data, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal output: %v", err)
	}
	fmt.Println(data)
	return string(data), nil
}
func generateVariableDeclaration(lang, typeName, name string, value interface{}) string {
	switch lang {
	case "java":
		return generateJavaVariableDeclaration(typeName, name, value)
	case "cpp":
		return fmt.Sprintf("%s %s = %s;", typeName, name, formatValue(lang, value))
	case "javascript":
		return fmt.Sprintf("const %s = %s;", name, formatValue(lang, value))
	case "typescript":
		return fmt.Sprintf("const %s: %s = %s;", name, typeName, formatValue(lang, value))
	default:
		return fmt.Sprintf("let %s = %s;", name, formatValue(lang, value))
	}
}

func generateJavaVariableDeclaration(typeName, name string, value interface{}) string {
	switch typeName {
	case "String[]", "int[]", "long[]", "double[]", "boolean[]", "char[]":
		return fmt.Sprintf("%s %s = {%s};", typeName, name, formatValue("java", value))
	default:
		return fmt.Sprintf("%s %s = %s;", typeName, name, formatValue("java", value))
	}
}

func generateFunctionCall(lang, name, args, returnType string) string {
	switch lang {
	case "java":
		if returnType == "arr<int>" {
			return fmt.Sprintf("int[] result = solution.%s(%s);\nSystem.out.println(Arrays.toString(result));", name, args)
		}
		return fmt.Sprintf("int result = solution.%s(%s);\nSystem.out.println(result);", name, args)
	case "cpp":
		return fmt.Sprintf("auto result = %s(%s);\nstd::cout << result;", name, args)
	case "javascript":
		if returnType == "number[]" {
			return fmt.Sprintf("const result = %s(%s);\nconsole.log(JSON.stringify(result));", name, args)
		}
		return fmt.Sprintf("const result = %s(%s);\nconsole.log(result);", name, args)

	case "typescript":
		if returnType == "arr<int>" {
			return fmt.Sprintf("const result = %s(%s);\nconsole.log(result);", name, args)
		}
		return fmt.Sprintf("const result = %s(%s);\nconsole.log(result);", name, args)
	default:
		return fmt.Sprintf("%s(%s);", name, args)
	}
}
func formatValue(lang string, value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case []interface{}:
		return formatSlice(lang, v)
	case bool:
		if lang == "java" {
			return fmt.Sprintf("%t", v)
		}
		return strings.ToLower(fmt.Sprintf("%t", v))
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatSlice(lang string, slice []interface{}) string {
	var formatted []string
	for _, val := range slice {
		formatted = append(formatted, formatValue(lang, val))
	}
	if lang == "typescript" || lang == "javascript" {
		return "[" + strings.Join(formatted, ", ") + "]"
	}
	return strings.Join(formatted, ", ")
}
func getType(lang string, value interface{}) string {
	switch lang {
	case "java":
		return getJavaType(value)
	case "cpp":
		return getCppType(value)
	case "javascript":
		return getJSType(value)
	case "typescript":
		return getTSType(value)
	default:
		return "unknown"
	}
}
func getJavaType(value interface{}) string {
	switch v := value.(type) {
	case []interface{}:
		if len(v) > 0 {
			switch v[0].(type) {
			case string:
				return "String[]"
			case float64, int:
				return "int[]"
			case bool:
				return "boolean[]"
			}
		}
		return "Object[]"
	case float64, int:
		return "int"
	case string:
		return "String"
	case bool:
		return "boolean"
	default:
		return "Object"
	}
}
func getCppType(value interface{}) string {
	switch v := value.(type) {
	case []interface{}:
		if len(v) > 0 {
			switch v[0].(type) {
			case float64, int:
				return "std::vector<int>"
			case string:
				return "std::vector<std::string>"
			case bool:
				return "std::vector<bool>"
			}
		}
		return "std::vector<void*>"
	case float64, int:
		return "int"
	case string:
		return "std::string"
	case bool:
		return "bool"
	default:
		return "auto"
	}
}

func getJSType(value interface{}) string {
	switch v := value.(type) {
	case []interface{}:
		if len(v) > 0 {
			switch v[0].(type) {
			case float64, int:
				return "number[]"
			case string:
				return "string[]"
			case bool:
				return "boolean[]"
			}
		}
		return "any[]"
	case float64, int:
		return "number"
	case string:
		return "string"
	case bool:
		return "boolean"
	default:
		return "any"
	}
}

func getTSType(value interface{}) string {
	switch v := value.(type) {
	case []interface{}:
		if len(v) > 0 {
			switch v[0].(type) {
			case float64, int:
				return "number[]"
			case string:
				return "string[]"
			case bool:
				return "boolean[]"
			}
		}
		return "any[]"
	case float64, int:
		return "number"
	case string:
		return "string"
	case bool:
		return "boolean"
	default:
		return "any"
	}
}
