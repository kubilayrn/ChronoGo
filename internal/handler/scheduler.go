package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartScheduler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Scheduler started",
		"status":  "running",
	})
}

func StopScheduler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Scheduler stopped",
		"status":  "stopped",
	})
}
