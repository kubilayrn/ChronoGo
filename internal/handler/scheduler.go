package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ToggleScheduler(c *gin.Context) {
	isRunning := h.scheduler.IsRunning()

	if isRunning {
		h.scheduler.Stop()
		c.JSON(http.StatusOK, gin.H{
			"message": "Scheduler stopped",
			"status":  "stopped",
		})
	} else {
		if err := h.scheduler.Start(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to start scheduler",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Scheduler started",
			"status":  "running",
		})
	}
}
