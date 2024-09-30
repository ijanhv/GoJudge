package controllers

import (
	"bytes"
	"fmt"
	"gojudge/db"
	"gojudge/generator"
	"gojudge/models"
	"log"

	"net/http"
	"os"
	"path/filepath"

	"github.com/gosimple/slug"

	storage_go "github.com/supabase-community/storage-go"

	"github.com/gin-gonic/gin"
)

func CreateProblem(c *gin.Context) {
	// Check for Supabase credentials at the beginning
	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseAnonKey := os.Getenv("ANON_KEY")
	bucketName := "gojudge"

	if supabaseUrl == "" || supabaseAnonKey == "" || bucketName == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Supabase configuration not set"})
		return
	}

	var problem models.Problem

	// Bind JSON input to Problem struct
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	problem.Slug = slug.Make(problem.Title)

	// Save the problem in the database
	if result := db.GetDB().Create(&problem); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save problem: " + result.Error.Error()})
		return
	}

	// Initialize Supabase storage client
	storageClient := storage_go.NewClient(supabaseUrl, supabaseAnonKey, nil)

	// Test Supabase connection
	_, listErr := storageClient.ListBuckets()
	if listErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Supabase storage: " + listErr.Error()})
		return
	}

	// Generate boilerplate code for various languages
	boilerplateCodes := map[string]string{
		"problem.cpp":  generator.GenerateCPlusPlusBoilerplate(problem),
		"problem.js":   generator.GenerateJavaScriptBoilerplate(problem),
		"problem.java": generator.GenerateJavaBoilerplate(problem),
	}

	uploadErrors := make([]string, 0)

	// Upload boilerplate code to Supabase storage
	for fileName, content := range boilerplateCodes {
		folderName := problem.Slug
		fullPath := filepath.Join(folderName, fileName)

		contentReader := bytes.NewReader([]byte(content))

		// Upload file
		result, uploadErr := storageClient.UploadFile(bucketName, fullPath, contentReader, storage_go.FileOptions{
			ContentType: getContentType(fileName),
		})

		if uploadErr != nil {
			errorMsg := fmt.Sprintf("Failed to upload %s: %s", fullPath, uploadErr.Error())
			uploadErrors = append(uploadErrors, errorMsg)
			log.Printf("Upload error: %s", errorMsg)
		} else {
			log.Printf("Successfully uploaded file %s: %s\n", fullPath, result.Key)
		}
	}

	if len(uploadErrors) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{
			"message": "Problem created but some files failed to upload",
			"problem": problem,
			"errors":  uploadErrors,
		})
		return
	}

	c.JSON(http.StatusOK, problem)
}

// Helper function to determine content type based on file extension
func getContentType(fileName string) *string {
	switch filepath.Ext(fileName) {
	case ".cpp":
		return strPtr("text/x-c++")
	case ".js":
		return strPtr("application/javascript")
	case ".java":
		return strPtr("text/x-java-source")
	default:
		return nil
	}
}

// Helper function to convert string to *string
func strPtr(s string) *string {
	return &s
}

func GetAllProblems(c *gin.Context) {
	var problems []models.Problem

	result := db.GetDB().Find(&problems)

	if result.Error != nil {
		log.Printf("Database error while fetching problems: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve problems: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Problems Retrieved successfully!",
		"problems": problems,
	})
}

func GetProblem(c *gin.Context) {
	slug := c.Param("slug")

	fmt.Printf("SLUG: %s", slug)
	var problem models.Problem

	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseAnonKey := os.Getenv("ANON_KEY")
	bucketName := "gojudge"

	storageClient := storage_go.NewClient(supabaseUrl, supabaseAnonKey, nil)

    result := db.GetDB().Preload("Function").Preload("Function.Parameters").Preload("TestCases").Where("slug = ?", slug).First(&problem)

	if result.Error != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve problem: " + result.Error.Error()})
		return
	}

	codeFiles := []string{"problem.cpp", "problem.js", "problem.java"}
	fileContents := make(map[string]string)

	for _, fileName := range codeFiles {
		filePath := filepath.Join(problem.Slug, fileName)

		fileBytes, err := storageClient.DownloadFile(bucketName, filePath, storage_go.UrlOptions{})
		if err != nil {
			log.Printf("Error downloading %s: %v", filePath, err)
			fileContents[fileName] = fmt.Sprintf("Error: Unable to retrieve %s", fileName)
			continue
		}

		fileContents[fileName] = string(fileBytes)
	}

	response := gin.H{
		"message":     "Problem Retrieved successfully!",
		"problem":     problem,
		"boilerplate": fileContents,
	}

	c.JSON(http.StatusOK, response)
}
