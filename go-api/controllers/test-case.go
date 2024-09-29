package controllers

import (
	"gojudge/db"
	"gojudge/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTestCase (c *gin.Context) {
	id := c.Param("id")

	var testCase models.TestCase
	result := db.GetDB().First(&testCase, id)

	if result.Error != nil {
		log.Printf("Database error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve problem: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, testCase)

}