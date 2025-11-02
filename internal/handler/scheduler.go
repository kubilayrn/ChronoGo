package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ToggleScheduler godoc
// @Summary      Toggle scheduler on/off
// @Description  Start or stop the automatic message sending scheduler
// @Tags         scheduler
// @Accept       json
// @Produce      json
// @Success      200  {object}  ToggleSchedulerResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /scheduler/toggle [post]
func (h *Handler) ToggleScheduler(c *gin.Context) {
	isRunning := h.scheduler.IsRunning()

	if isRunning {
		h.scheduler.Stop()
		c.JSON(http.StatusOK, ToggleSchedulerResponse{
			Message: "Scheduler stopped",
			Status:  "stopped",
		})
	} else {
		if err := h.scheduler.Start(); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to start scheduler",
			})
			return
		}
		c.JSON(http.StatusOK, ToggleSchedulerResponse{
			Message: "Scheduler started",
			Status:  "running",
		})
	}
}
