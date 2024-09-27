package controllers

import (
	"gojudge/db"
	"gojudge/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateProblem(c *gin.Context) {
	var problem models.Problem
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the problem in the database
	if result := db.GetDB().Create(&problem); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, problem)
}

func GetProblem(c *gin.Context) {
	id := c.Param("id")

	var problem models.Problem

	if result := db.GetDB().Preload("TestCases").First(&problem, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve problem"})
		return
	}

	c.JSON(http.StatusOK, problem)

}
