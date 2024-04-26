package routes

import (
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDashboardStatsRoute(c *gin.Context) {
	registerStats, err := services.GetRegisterStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"register": registerStats,
	})
}
