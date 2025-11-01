package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListSentMessages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"messages": []interface{}{},
		"total":    0,
	})
}
